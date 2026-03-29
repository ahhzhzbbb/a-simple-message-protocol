package main

import (
	"net"
	"slices"
)

func (s *Server) usernameOf(conn net.Conn) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.connToUser[conn]
}

func (s *Server) removeConnection(conn net.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()

	username := s.connToUser[conn]
	delete(s.userToConn, username)
	delete(s.connToUser, conn)

	for groupName, members := range s.group {
		idx := slices.Index(members, conn)
		if idx == -1 {
			continue
		}

		updated := slices.Delete(members, idx, idx+1)
		if len(updated) == 0 {
			delete(s.group, groupName)
		} else {
			s.group[groupName] = updated
		}
	}
}

func (s *Server) changeYourUserName(newName string, conn net.Conn) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.userToConn[newName]; exists {
		return false
	}

	oldName, ok := s.connToUser[conn]
	if ok {
		delete(s.userToConn, oldName)
	}

	s.connToUser[conn] = newName
	s.userToConn[newName] = conn
	return true
}
