package server

import (
	"net"
)

type User struct {
	Name    string
	Address string
	Chanel  chan string
	Conn    net.Conn
	server  *Server
}

// 创建用户
func NewUser(conn net.Conn, server *Server) *User {
	addr := conn.RemoteAddr().String()

	user := &User{
		Name:    addr,
		Address: addr,
		Chanel:  make(chan string),
		Conn:    conn,
		server:  server,
	}

	// 启动消息监听
	go user.ListenMessage()

	return user
}

// 监听当前user的channel，监听到消息发送至客户端
func (u *User) ListenMessage() {
	for {
		msg := <-u.Chanel
		u.Conn.Write([]byte(msg + "\n"))
	}
}

// 用户上线
func (u *User) Online() {

	u.server.mapLock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.mapLock.Unlock()

	// 广播用户上线
	u.server.Broadcast(u, "已上线")

}

// 用户下线
func (u *User) Offline() {
	u.server.mapLock.Lock()
	delete(u.server.OnlineMap, u.Name)
	u.server.mapLock.Unlock()

	// 广播用户下线
	u.server.Broadcast(u, "下线")

}

// 处理用户发送来的消息
func (u *User) HandleMessage(msg string) {
	u.server.Broadcast(u, msg)
}
