package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/lyxuansang91/redis-crash-course/server"
)

func main() {
	// Create and start the TCP server on port 8080
	tcpServer := server.NewServer("8080")
	
	// Start the server in a goroutine
	go func() {
		if err := tcpServer.Start(); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()
	
	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	fmt.Println("TCP Server is running...")
	fmt.Println("Press Ctrl+C to stop the server")
	
	// Wait for shutdown signal
	<-sigChan
	fmt.Println("\nShutting down server...")
	
	// Gracefully stop the server
	if err := tcpServer.Stop(); err != nil {
		log.Printf("Error stopping server: %v", err)
	}
	
	fmt.Println("Server stopped successfully")
}
