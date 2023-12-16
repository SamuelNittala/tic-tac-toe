package main

import (
	"fmt"
	"net"
)

func handleSession(conn net.Conn) string {
	return "hello"
}

func main() {
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		fmt.Println(err)
	}
	for {
		conn, _ := listener.Accept()
		go handleSession(conn)
	}
}
