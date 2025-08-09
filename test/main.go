package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

// TestClient represents a test client for the TCP server
type TestClient struct {
	conn net.Conn
	name string
}

// NewTestClient creates a new test client
func NewTestClient(name string) (*TestClient, error) {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		return nil, err
	}
	
	// Read welcome message
	scanner := bufio.NewScanner(conn)
	if scanner.Scan() {
		fmt.Printf("[%s] %s\n", name, scanner.Text())
	}
	
	return &TestClient{
		conn: conn,
		name: name,
	}, nil
}

// SendMessage sends a message to the server
func (tc *TestClient) SendMessage(message string) error {
	_, err := fmt.Fprintln(tc.conn, message)
	if err != nil {
		return err
	}
	
	// Read response
	scanner := bufio.NewScanner(tc.conn)
	if scanner.Scan() {
		fmt.Printf("[%s] Server: %s\n", tc.name, scanner.Text())
	}
	
	return nil
}

// Close closes the connection
func (tc *TestClient) Close() {
	tc.conn.Close()
}

func main() {
	fmt.Println("Starting TCP Server Test...")
	fmt.Println("Make sure the server is running with: go run test/main.go")
	
	clients := make([]*TestClient, 0)
	for i := 0; i < 10; i++ {
		client, err := NewTestClient(fmt.Sprintf("Client-%d", i))
		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
		}
		fmt.Printf("Created %s\n", client.name)
		clients = append(clients, client)
	}
	// Send messages from each client
	messages := []string{
		"Hello from client!",
		"How are you doing?",
		"This is a test message",
		"Testing concurrent connections",
	}

	// synchronize the clients
	wg := &sync.WaitGroup{}
	for _, client := range clients {
		wg.Add(1)
		go func(c *TestClient) {
			defer wg.Done()
			for _, message := range messages {
				fmt.Printf("\n[%s] Sending: %s\n", c.name, message)
				if err := c.SendMessage(message); err != nil {
					log.Printf("Error sending message from %s: %v", c.name, err)
				}
				time.Sleep(500 * time.Millisecond) // Small delay between messages	
			}
			fmt.Println("\nDisconnecting client...")
			c.SendMessage("quit")
		}(client)
	}
	wg.Wait()
	fmt.Println("\nTest completed successfully!")
} 