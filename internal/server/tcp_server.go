package server

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"syscall"

	"github.com/lyxuansang91/redis-crash-course/internal/config"
	"github.com/lyxuansang91/redis-crash-course/internal/core"
	"github.com/lyxuansang91/redis-crash-course/internal/core/io_multiplexing"
	"github.com/lyxuansang91/redis-crash-course/internal/data_structure"
	"github.com/lyxuansang91/redis-crash-course/threadpool"
)

// Server represents our TCP server
type Server struct {
	config *config.Config
	listener net.Listener
	port     string
	executor core.CommandExecutor
}


// NewServer creates a new TCP server instance
func NewServer(config *config.Config) *Server {
	return &Server{
		config: config,
		port: config.Port,
		executor: core.NewCommandExecutor(data_structure.CreateDict()),
	}
}

func (s *Server) readCommand(fd int) (*core.Command, error) {
	var buf = make([]byte, 512)
	n, err := syscall.Read(fd, buf)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, io.EOF
	}
	return core.ParseCmd(buf)
}

func (s *Server) RunIoMultiplexingServer() error {
	log.Println("starting an I/O Multiplexing TCP server on", s.config.Port)
	listener, err := net.Listen(s.config.Protocol, s.config.Port)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %v", s.config.Port, err)
	}
	defer listener.Close()

	// Get the file descriptor from the listener
	tcpListener, ok := listener.(*net.TCPListener)
	if !ok {
		return fmt.Errorf("listener is not a TCPListener")
	}
	listenerFile, err := tcpListener.File()

	if err != nil {
		return fmt.Errorf("failed to get file descriptor from listener: %v", err)
	}
	defer listenerFile.Close()

	serverFd := int(listenerFile.Fd())

	// Create an ioMultiplexer instance (epoll in Linux, kqueue in MacOS)
	ioMultiplexer, err := io_multiplexing.CreateIOMultiplexer(s.config)
	if err != nil {
		return fmt.Errorf("failed to create io multiplexer: %v", err)
	}
	defer ioMultiplexer.Close()

	// Monitor "read" events on the Server FD
	if err = ioMultiplexer.Monitor(io_multiplexing.Event{
		Fd: serverFd,
		Op: io_multiplexing.OpRead,
	}); err != nil {
		return fmt.Errorf("failed to monitor server fd: %v", err)
	}

	var events = make([]io_multiplexing.Event, config.MaxConnections)
	for {
		// wait for file descriptors in the monitoring list to be ready for I/O
		// it is a blocking call.
		events, err = ioMultiplexer.Wait()
		if err != nil {
			continue
		}

		for i := 0; i < len(events); i++ {
			if events[i].Fd == serverFd {
				log.Println("new client is trying to connect")
				// set up new connection
				connFd, _, err := syscall.Accept(serverFd)
				if err != nil {
					log.Printf("err accept: %v\n", err)
					continue
				}
				log.Println("set up a new connection")
				// ask epoll to monitor this connection
				if err = ioMultiplexer.Monitor(io_multiplexing.Event{
					Fd: connFd,
					Op: io_multiplexing.OpRead,
				}); err != nil {
					return fmt.Errorf("failed to monitor connection fd: %v", err)
				}
			} else {
				cmd, err := s.readCommand(events[i].Fd)
				if err != nil {
					if err == io.EOF || err == syscall.ECONNRESET {
						log.Println("client disconnected")
						_ = syscall.Close(events[i].Fd)
						continue
					}
					log.Printf("read error: %v\n", err)
					continue
				}
				if err = s.executor.ExecuteAndResponse(cmd, events[i].Fd); err != nil {
					log.Printf("err write: %v\n", err)
				}
			}
		}
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