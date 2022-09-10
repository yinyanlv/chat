package server

import (
	"net"
	"strings"
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

// 处理用户发送过来的消息
func (u *User) HandleMessage(msg string) {
	if msg == "query" {
		// 查询所有在线用户
		var res string
		u.server.mapLock.Lock()
		for _, user := range u.server.OnlineMap {
			res += "[" + user.Address + "]" + user.Name + ":在线...\n"
		}
		u.server.mapLock.Unlock()
		u.SendMessage(res)
	} else if len(msg) > 5 && msg[:5] == "name|" {
		newName := msg[5:]
		u.server.mapLock.Lock()
		_, ok := u.server.OnlineMap[newName]
		// 当前用户名已存在
		if ok {
			u.SendMessage("当前用户名已被使用\n")
			return
		}
		delete(u.server.OnlineMap, u.Name)
		u.Name = newName
		u.server.OnlineMap[newName] = u
		u.server.mapLock.Unlock()

		u.SendMessage("您已设置用户名为：" + u.Name + "\n")
	} else if len(msg) > 3 && msg[:3] == "to|" {
		// 消息格式：to|张三|你好
		list := strings.Split(msg, "|")

		remoteName := list[1]
		if remoteName == "" {
			u.SendMessage("请输入正确的消息格式，如：\"to|张三|你好\"\n")
			return
		}

		remoteUser, ok := u.server.OnlineMap[remoteName]
		if !ok {
			u.SendMessage("用户" + remoteName + "不存在！\n")
			return
		}

		content := list[2]
		if content == "" {
			u.SendMessage("请输入消息内容\n")
			return
		}

		remoteUser.SendMessage(u.Name + "对你说：" + content + "\n")

	} else {
		u.server.Broadcast(u, msg)
	}
}

// 向当前用户的客户端发送消息
func (u *User) SendMessage(msg string) {
	u.Conn.Write([]byte(msg))
}
