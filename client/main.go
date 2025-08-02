package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	fmt.Println("Connected to TCP server!")
	// Create a scanner to read from stdin
	serverScanner := bufio.NewScanner(conn)
	for serverScanner.Scan() {
		fmt.Printf("Server: %s\n", serverScanner.Text())
	}
	if err := serverScanner.Err(); err != nil {
		log.Printf("Error reading from server: %v", err)
	}
} 