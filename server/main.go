package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	server := NewServer()

	if len(os.Args) < 2 {
		fmt.Printf("Usage: go run %s <port>\n", os.Args[0])
		os.Exit(1)
	}
	port := fmt.Sprintf(":%s", os.Args[1])
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("failed to create listener, err:", err)
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Printf("listening on %s\n", listener.Addr())

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("failed to accept connection, err:", err)
			continue
		}

		fmt.Println("connection estalished!")

		go server.handleConnection(conn)
	}
}
