package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
)

type Server struct {
	userToConn map[string]net.Conn
	connToUser map[net.Conn]string
	group      map[string][]net.Conn
	mu         sync.Mutex
}

func NewServer() *Server {
	return &Server{
		userToConn: make(map[string]net.Conn),
		connToUser: make(map[net.Conn]string),
		group:      make(map[string][]net.Conn),
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Println("failed to read data, err:", err)
			} else {
				fmt.Printf("client %s disconnected\n", conn.RemoteAddr())
			}
			s.removeConnection(conn)
			return
		}

		fmt.Printf("request: %s", input)

		request := strings.TrimSpace(input)
		parts := strings.Fields(request)
		if len(parts) == 0 {
			continue
		}
		if s.handleCommand(conn, parts) {
			return
		}
	}
}
