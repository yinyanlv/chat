package main

import (
	"github.com/yinyanlv/chat/server"
)

func main() {
	server := server.NewServer("127.0.0.1", 5000)
	server.Start()
}
