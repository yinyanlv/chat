package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIP   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

var serverIP string
var serverPort int

func init() {
	flag.StringVar(&serverIP, "ip", "127.0.0.1", "设置服务器ip")
	flag.IntVar(&serverPort, "port", 5000, "设置服务器端口")
}

func NewClient(ip string, port int) *Client {
	client := &Client{
		ServerIP:   ip,
		ServerPort: port,
		flag:       10000,
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", client.ServerIP, client.ServerPort))

	if err != nil {
		fmt.Println("net.Dial error: ", err)
		return nil
	}

	client.conn = conn

	return client
}

func (c *Client) menu() bool {
	var flag int
	fmt.Println("1: 公聊模式")
	fmt.Println("2: 私聊模式")
	fmt.Println("3: 修改用户名")
	fmt.Println("0: 退出")
	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		c.flag = flag
		return true
	} else {
		fmt.Println("请输入合法范围内的数字")
		return false
	}
}

func (c *Client) UpdateName() bool {
	fmt.Println(">>> 请输入用户名：")
	fmt.Scanln(&c.Name)

	sendMsg := "name|" + c.Name + "\n"
	_, err := c.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write error: ", err)
		return false
	}
	return true
}

// 处理服务器发来的消息
func (c *Client) HandleMessage() {
	// 相当于io.Copy(os.Stdout, client.conn)
	for {
		buf := make([]byte, 4096)
		c.conn.Read(buf) // 阻塞等待
		fmt.Println(string(buf))
	}
}

func (c *Client) PublicChat() {
	var chatMsg string

	fmt.Println(">>> 请输入聊天信息，输入exit退出聊天模式：")
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := c.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn.Write error: ", err)
				break
			}
		}

		// 重置
		chatMsg = ""
		fmt.Println(">>> 请输入聊天信息，输入exit退出聊天模式：")
		fmt.Scanln(&chatMsg)
	}
}

func (c *Client) QueryUsers() {
	sendMsg := "query\n"
	_, err := c.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write error: ", err)
		return
	}
}

func (c *Client) PrivateChat() {
	var remoteName string
	var chatMsg string

	c.QueryUsers()

	fmt.Println(">>> 请输入聊天对象的用户名，输入exit退出聊天模式：")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println(">>> 请输入聊天信息，输入exit退出聊天模式：")
		fmt.Scanln(&chatMsg)

		for chatMsg != "exit" {

			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n"
				_, err := c.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn.Write error: ", err)
					break
				}
			}

			// 重置
			chatMsg = ""
			fmt.Println(">>> 请输入聊天信息，输入exit退出聊天模式：")
			fmt.Scanln(&chatMsg)
		}

		// 重置
		remoteName = ""
		fmt.Println(">>> 请输入聊天对象的用户名，输入exit退出聊天模式：")
		fmt.Scanln(&remoteName)
	}
}

func (c *Client) Run() {
	for c.flag != 0 {
		// 阻塞
		if c.menu() != true {
		}
		switch c.flag {
		case 1:
			c.PublicChat()
		case 2:
			c.PrivateChat()
		case 3:
			c.UpdateName()
		case 0:
			fmt.Println("退出")
		}
	}
}

func main() {
	flag.Parse()

	client := NewClient(serverIP, serverPort)

	if client == nil {
		return
	}
	fmt.Println("连接服务器成功")

	go client.HandleMessage()
	client.Run()
}
