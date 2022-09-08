package server

import (
	"fmt"
	"net"
)

type Server struct {
	IP   string
	Port int
}

// 创建服务器
func NewServer(ip string, port int) *Server {
	return &Server{
		IP:   ip,
		Port: port,
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

	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept error: ", err)
			continue
		}

		// handle
		go s.Handle(conn)

	}

}

func (s *Server) Handle(conn net.Conn) {
	// nc 127.0.0.1 5000
	fmt.Println("连接建立成功")
}
