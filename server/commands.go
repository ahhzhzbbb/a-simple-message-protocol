package main

import "net"

func (s *Server) handleCommand(conn net.Conn, parts []string) bool {
	instruction := parts[0]

	switch instruction {
	case "/user":
		s.handleUser(conn, parts)
	case "/users":
		s.handleUsers(conn, parts)
	case "/msg":
		s.handleBroadcast(conn, parts)
	case "/mtu":
		s.handleMessageToUser(conn, parts)
	case "/group":
		s.handleCreateGroup(conn, parts)
	case "/groups":
		s.handleListGroups(conn, parts)
	case "/getout":
		s.handleLeaveGroup(conn, parts)
	case "/join":
		s.handleJoinGroup(conn, parts)
	case "/mtg":
		s.handleMessageToGroup(conn, parts)
	case "/quit":
		return s.handleQuit(conn, parts)
	case "/help":
		_, _ = conn.Write([]byte(s.printManual()))
	default:
		_, _ = conn.Write([]byte("unknown command, use /help\n"))
	}

	return false
}
