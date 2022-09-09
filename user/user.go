package user

import "net"

type User struct {
	Name    string
	Address string
	Chanel  chan string
	Conn    net.Conn
}

// 创建用户
func NewUser(conn net.Conn) *User {
	addr := conn.RemoteAddr().String()

	user := &User{
		Name:    addr,
		Address: addr,
		Chanel:  make(chan string),
		Conn:    conn,
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
