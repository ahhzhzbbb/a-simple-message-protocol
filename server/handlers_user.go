package main

import (
	"fmt"
	"net"
	"strings"
)

func (s *Server) handleUser(conn net.Conn, parts []string) {
	if len(parts) < 2 {
		_, _ = conn.Write([]byte("usage: /user <name>\n"))
		return
	}

	newName := parts[1]
	if !s.changeYourUserName(newName, conn) {
		_, _ = conn.Write([]byte("username already taken\n"))
		return
	}

	response := fmt.Sprintf("your username was changed to be %s\n", newName)
	fmt.Printf("response: %s", response)
	_, _ = conn.Write([]byte(response))
}

func (s *Server) handleUsers(conn net.Conn, parts []string) {
	if len(parts) > 1 {
		_, _ = conn.Write([]byte("too many args\n"))
		return
	}

	var response strings.Builder
	response.WriteString("users list:\n")

	s.mu.Lock()
	for _, u := range s.connToUser {
		response.WriteString(u)
		response.WriteByte('\n')
	}
	s.mu.Unlock()

	fmt.Printf("response: {\n%s}", response.String())
	_, _ = conn.Write([]byte(response.String()))
}

func (s *Server) handleBroadcast(conn net.Conn, parts []string) {
	if len(parts) < 2 {
		_, _ = conn.Write([]byte("usage: /msg <message>\n"))
		return
	}

	msg := strings.Join(parts[1:], " ")

	s.mu.Lock()
	sender := s.connToUser[conn]
	receivers := make([]net.Conn, 0, len(s.connToUser))
	for c := range s.connToUser {
		if c != conn {
			receivers = append(receivers, c)
		}
	}
	s.mu.Unlock()

	response := fmt.Sprintf("\033[34m[%s]\033[0m %s\n", sender, msg)
	for _, c := range receivers {
		_, _ = c.Write([]byte(response))
	}
}

func (s *Server) handleMessageToUser(conn net.Conn, parts []string) {
	if len(parts) < 3 {
		_, _ = conn.Write([]byte("usage: /mtu <username> <message>\n"))
		return
	}

	username := parts[1]
	message := strings.Join(parts[2:], " ")

	s.mu.Lock()
	receiver, ok := s.userToConn[username]
	sender := s.connToUser[conn]
	s.mu.Unlock()

	if !ok {
		_, _ = conn.Write([]byte("user not found\n"))
		return
	}

	response := fmt.Sprintf("\033[31m[%s]\033[0m %s\n", sender, message)
	_, _ = receiver.Write([]byte(response))
}

func (s *Server) handleQuit(conn net.Conn, parts []string) bool {
	if len(parts) != 1 {
		_, _ = conn.Write([]byte("too many args, try again\n"))
		return false
	}

	username := s.usernameOf(conn)
	fmt.Printf("closing conn of %s\n", username)
	_, _ = conn.Write([]byte("Goodbye\n"))
	s.removeConnection(conn)
	return true
}
