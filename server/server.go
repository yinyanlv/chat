package server

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	IP        string
	Port      int
	OnlineMap map[string]*User
	mapLock   sync.RWMutex
	Message   chan string // 消息广播的channel
}

// 创建服务器
func NewServer(ip string, port int) *Server {
	return &Server{
		IP:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
}

func (s *Server) Start() {
	// listen socket
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("net.Listen error: ", err)
		return
	}
	fmt.Printf("服务器已启动！IP：%s，端口：%d\n", s.IP, s.Port)

	// close listen socket
	defer listener.Close()

	// 启动监听message管道
	go s.ListenMessage()

	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept error: ", err)
			continue
		}
		// 处理当前连接
		go s.Handle(conn)
	}
}

// 处理连接
func (s *Server) Handle(conn net.Conn) {

	// nc 127.0.0.1 5000
	user := NewUser(conn, s)

	user.Online()

	isLive := make(chan bool)

	// 处理客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf) // Read方法会阻塞等待消息
			// 客户端关闭
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("conn read error: ", err)
				return
			}
			// 提取信息，去除\n
			msg := string(buf[:n-1])
			// 广播消息
			user.HandleMessage(msg)
			// 标记用户为活跃用户
			isLive <- true
		}
	}()

	// 阻塞当前handler
	for {
		select {
		case <-isLive: // 当前用户活跃
			// 不做任何操作
		case <-time.After(time.Second * 300): // 超时时间5分钟，每次for循环，如果执行到该case，则重置定时器
			user.SendMessage("你被踢了")
			close(user.Chanel)
			user.Conn.Close()
			s.mapLock.Lock()
			delete(s.OnlineMap, user.Name)
			s.mapLock.Unlock()
			return // 或者runtime.Goexit()
		}
	}
}

// 广播
func (s *Server) Broadcast(user *User, msg string) {
	sendMsg := "[" + user.Address + "]" + user.Name + ":" + msg

	s.Message <- sendMsg
}

// 监听message管道
func (s *Server) ListenMessage() {
	for {
		msg := <-s.Message
		s.mapLock.Lock()
		for _, cli := range s.OnlineMap {
			cli.Chanel <- msg
		}
		s.mapLock.Unlock()
	}
}
