package server

import (
	"fmt"
	"log"
	"net"
	"sync"
)

// Server represents our TCP server
type Server struct {
	listener net.Listener
	clients  map[net.Conn]bool
	mutex    sync.RWMutex
	port     string
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
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		
		// Handle each connection in a separate goroutine
		go s.handleConnection(conn)
	}
}

// handleConnection handles individual client connections
func (s *Server) handleConnection(conn net.Conn) {
	fmt.Printf("New client connected: %s\n", conn.RemoteAddr().String())
	
	
	
	// Send welcome message
	conn.Write([]byte("HTTP/1.1 200 OK \r\n\r\n Engineer Pro\r\n"))
	conn.Close()
}

// Stop gracefully shuts down the server
func (s *Server) Stop() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
} 