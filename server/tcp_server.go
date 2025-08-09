package server

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
)

// Server represents our TCP server
type Server struct {
	listener net.Listener
	clients  map[net.Conn]bool
	port     string
}

type Client struct {
	conn net.Conn
}

// NewServer creates a new TCP server instance
func NewServer(port string) *Server {
	return &Server{
		clients: make(map[net.Conn]bool),
		port:    port,
	}
}

// Start initializes and starts the TCP server
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}
	s.listener = listener
	
	fmt.Printf("TCP Server started on port %s\n", s.port)
	fmt.Println("Waiting for connections...")
	
	// Accept connections in a loop
	for {
		conn, err := s.listener.Accept()
		client := &Client{conn: conn}
		if err != nil {
            if errors.Is(err, net.ErrClosed) {
                // Listener closed via Stop(); exit accept loop
                return nil
            }
            log.Printf("Error accepting connection: %v", err)
            continue
		}
		
		// Handle each connection in a separate goroutine
		go client.handleConnection()
	}
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

// Stop gracefully shuts down the server
func (s *Server) Stop() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
} 