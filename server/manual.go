package main

func (s *Server) printManual() string {
	msg := `=============================================
/user + <username> : change your username
/users : list all usernames
/msg : send message to broadcast
/mtu : send message to user
/group: create new group
/groups : list all group
/join : join in the group
/getout : leave the group
/mtg : send message to the group members
/quit : disconnect from server
/help : print system manual
=============================================
`
	return msg
}
