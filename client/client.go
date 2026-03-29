package main

import (
	"bufio"
	"fmt"
	"net"
)

type Client struct {
	username string
	conn     net.Conn
	requests chan []byte
}

func newClient(name string) *Client {
	return &Client{
		username: name,
		requests: make(chan []byte, 10),
	}
}

func (c *Client) handleRequest(conn net.Conn) {
	for r := range c.requests {
		_, err := conn.Write(r)
		if err != nil {
			fmt.Println("server disconnected")
			return
		}
	}
}
func (c *Client) handleResponse(conn net.Conn) {
	reader := bufio.NewReader(conn)

	for {
		r, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println("server disconnected")
			return
		}
		fmt.Print(string(r))
	}
}
