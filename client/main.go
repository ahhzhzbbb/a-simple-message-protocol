package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("type your username: ")
	input, err := reader.ReadBytes(byte('\n'))
	if err != nil {
		if err != io.EOF {
			fmt.Printf("disconecting...")
		}
		return
	}
	username := string(input)
	strings.TrimSpace(username)
	client := newClient(string(username))

	// =====================Connect to Server======================
	if len(os.Args) < 2 {
		fmt.Printf("Usage: go run %s <port>\n", os.Args[0])
		os.Exit(1)
	}
	port := fmt.Sprintf(":%s", os.Args[1])
	conn, err := net.Dial("tcp", port)
	if err != nil {
		fmt.Println("failed to connect server, err:", err)
		os.Exit(1)
	}
	// ============================================================

	register := fmt.Sprintf("/user %s", username)
	go client.handleRequest(conn)
	client.requests <- []byte(register)

	go client.handleResponse(conn)

	for {
		request, err := reader.ReadBytes(byte('\n'))
		if err != nil {
			if err != io.EOF {
				fmt.Printf("disconecting...")
			}
			return
		}
		client.requests <- request
	}
}
