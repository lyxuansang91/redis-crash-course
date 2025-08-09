package server

import (
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/lyxuansang91/redis-crash-course/threadpool"
)

// Server represents our TCP server
type Server struct {
	listener net.Listener
	port     string
}


// NewServer creates a new TCP server instance
func NewServer(port string) *Server {
	return &Server{
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

    pool := threadpool.NewPool(10)
    pool.Start()
	// Accept connections in a loop
	for {
        conn, err := s.listener.Accept()
		
		if err != nil {
            if errors.Is(err, net.ErrClosed) {
                // Listener closed via Stop(); exit accept loop
                return nil
            }
            log.Printf("Error accepting connection: %v", err)
            continue
		}
		
        client := NewClient(conn)

        // Queue a function task for the pool instead of calling directly
        pool.AddJob(func() {
            client.handleConnection()
        })
	}
}

// Stop gracefully shuts down the server
func (s *Server) Stop() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
} 