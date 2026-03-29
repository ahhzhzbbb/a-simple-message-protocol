package main

import (
	"fmt"
	"net"
	"slices"
	"strings"
)

func (s *Server) handleCreateGroup(conn net.Conn, parts []string) {
	if len(parts) < 2 {
		_, _ = conn.Write([]byte("usage: /group <nameGR>\n"))
		return
	}

	nameGR := parts[1]

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.group[nameGR]; exists {
		_, _ = conn.Write([]byte("this name was existed, try another one\n"))
		return
	}

	s.group[nameGR] = []net.Conn{conn}
	_, _ = fmt.Fprintf(conn, "Created a new group <%s>\n", nameGR)
}

func (s *Server) handleListGroups(conn net.Conn, parts []string) {
	if len(parts) > 1 {
		_, _ = conn.Write([]byte("\033[31mToo many args, try again\033[0m\n"))
		return
	}

	var groupIn strings.Builder
	var groupNIn strings.Builder

	s.mu.Lock()
	for g, users := range s.group {
		if slices.Contains(users, conn) {
			groupIn.WriteString(fmt.Sprintf("\033[34m[Joined] <%s>\033[0m\n", g))
		} else {
			groupNIn.WriteString(fmt.Sprintf("\033[31m[Public] <%s>\033[0m\n", g))
		}
	}
	s.mu.Unlock()

	if groupIn.Len() == 0 && groupNIn.Len() == 0 {
		_, _ = conn.Write([]byte("No groups available.\n"))
		return
	}

	var response strings.Builder
	response.WriteString("--- List of Groups ---\n")
	response.WriteString(groupIn.String())
	response.WriteString(groupNIn.String())
	_, _ = conn.Write([]byte(response.String()))
}

func (s *Server) handleLeaveGroup(conn net.Conn, parts []string) {
	if len(parts) < 2 {
		_, _ = conn.Write([]byte("usage: /getout <group>\n"))
		return
	}

	nameGR := parts[1]

	s.mu.Lock()
	defer s.mu.Unlock()

	groupConn, exists := s.group[nameGR]
	if !exists {
		_, _ = conn.Write([]byte("can not find this group, retry\n"))
		return
	}

	idx := slices.Index(groupConn, conn)
	if idx == -1 {
		_, _ = conn.Write([]byte("Youre not in this group\n"))
		return
	}

	updated := slices.Delete(groupConn, idx, idx+1)
	if len(updated) == 0 {
		delete(s.group, nameGR)
	} else {
		s.group[nameGR] = updated
	}

	_, _ = fmt.Fprintf(conn, "You are leaving the group %s\n", nameGR)
}

func (s *Server) handleJoinGroup(conn net.Conn, parts []string) {
	if len(parts) < 2 {
		_, _ = conn.Write([]byte("usage: /join <nameGR>\n"))
		return
	}

	nameGR := parts[1]

	s.mu.Lock()
	defer s.mu.Unlock()

	groupConn, exists := s.group[nameGR]
	if !exists {
		_, _ = conn.Write([]byte(fmt.Sprintf("The group <%s> is not exist\n", nameGR)))
		return
	}

	if slices.Contains(groupConn, conn) {
		_, _ = conn.Write([]byte(fmt.Sprintf("You already joined <%s>\n", nameGR)))
		return
	}

	s.group[nameGR] = append(groupConn, conn)
	_, _ = conn.Write([]byte(fmt.Sprintf("You joined in the group <%s>\n", nameGR)))
}

func (s *Server) handleMessageToGroup(conn net.Conn, parts []string) {
	if len(parts) < 3 {
		_, _ = conn.Write([]byte("usage: /mtg <nameGR> <message>\n"))
		return
	}

	nameGR := parts[1]
	msg := strings.Join(parts[2:], " ")

	s.mu.Lock()
	members, exists := s.group[nameGR]
	if !exists {
		s.mu.Unlock()
		_, _ = conn.Write([]byte(fmt.Sprintf("The group <%s> is not exist\n", nameGR)))
		return
	}
	if !slices.Contains(members, conn) {
		s.mu.Unlock()
		_, _ = conn.Write([]byte("You are not in this group\n"))
		return
	}

	sender := s.connToUser[conn]
	receivers := make([]net.Conn, 0, len(members))
	for _, c := range members {
		if c != conn {
			receivers = append(receivers, c)
		}
	}
	s.mu.Unlock()

	response := fmt.Sprintf("\033[31m<%s>\033[0m \033[34m[%s]\033[0m %s\n", nameGR, sender, msg)
	for _, c := range receivers {
		_, _ = c.Write([]byte(response))
	}
}
