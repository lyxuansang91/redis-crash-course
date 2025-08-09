package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type Client struct {
	conn net.Conn
}

func NewClient(conn net.Conn) *Client {
	return &Client{conn: conn}
}


// handleConnection handles individual client connections
func (c *Client) handleConnection() {
	conn := c.conn
	fmt.Printf("New client connected: %s\n", conn.RemoteAddr().String())
    defer func() {
        _ = conn.Close()
        fmt.Printf("Client disconnected: %s\n", conn.RemoteAddr().String())
    }()

    // Send welcome message (line-based so client scanner can read it)
    fmt.Fprintln(conn, "Welcome to the TCP Server! Send 'quit' to disconnect.")
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		message := scanner.Text()
		fmt.Println("Received message:", message)
		if message == "quit" || message == "exit" {
			fmt.Fprintln(conn, "Goodbye!")
			return
		}
		fmt.Fprintf(conn, "You said: %s\n", message)
	}

    if err := scanner.Err(); err != nil {
        log.Printf("Connection error from %s: %v", conn.RemoteAddr().String(), err)
    }
}