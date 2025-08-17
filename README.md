# Redis Crash Course - TCP Server Implementation

A Redis-like multi-threaded TCP server implemented in Go, featuring concurrent connection handling, I/O multiplexing, and RESP (Redis Serialization Protocol) support.

## Features

- **Multi-threaded TCP Server**: Handles multiple concurrent connections efficiently
- **I/O Multiplexing**: Uses epoll (Linux) and kqueue (macOS) for optimal performance
- **RESP Protocol Support**: Implements Redis Serialization Protocol for data encoding/decoding
- **Thread Pool**: Configurable thread pool for connection handling
- **Graceful Shutdown**: Proper signal handling and cleanup
- **Cross-platform**: Works on Linux and macOS

## Prerequisites

- Go 1.21.6 or later
- Linux (for epoll) or macOS (for kqueue)

## Installation

1. **Clone the repository**:

   ```bash
   git clone https://github.com/lyxuansang91/redis-crash-course.git
   cd redis-crash-course
   ```

2. **Install dependencies**:

   ```bash
   go mod download
   ```

3. **Build the project**:
   ```bash
   go build -o redis-server
   ```

## Usage

### Running the Server

1. **Start the TCP server**:

   ```bash
   ./redis-server
   ```

   Or run directly with Go:

   ```bash
   go run cmd/main.go
   ```

2. **The server will start on port 8080** and display:

   ```
   TCP Server is running...
   Press Ctrl+C to stop the server
   ```

3. **Stop the server** by pressing `Ctrl+C` for graceful shutdown.

### Testing the Server

1. **In a new terminal, run the test client**:

   ```bash
   go run test/main.go
   ```

2. **The test client will**:
   - Create 10 concurrent client connections
   - Send test messages to the server
   - Demonstrate concurrent handling capabilities

### Manual Testing

You can also test manually using netcat or telnet:

```bash
# Connect to the server
nc localhost 8080

# Send RESP protocol commands
*2\r\n$3\r\nGET\r\n$4\r\nkey1\r\n
```

## Project Structure

```
redis-crash-course/
├── cmd/
|   ├── main.go                # Main application entry point
├── go.mod                     # Go module dependencies
├── go.sum                     # Go module checksums
├── internal/                  # Internal packages
│   ├── config/               # Configuration management
│   │   └── config.go
│   ├── core/                 # Core RESP protocol implementation
│   │   ├── resp.go           # RESP encoding/decoding
│   │   ├── resp_test.go      # RESP protocol tests
│   │   └── io_multiplexing/  # I/O multiplexing implementations
│   │       ├── epoll_linux.go
│   │       ├── kqueue_macos.go
│   │       └── io_multiplexing.go
│   └── server/               # TCP server implementation
│       ├── tcp_server.go     # Main server logic
│       └── client.go         # Client connection handling
├── threadpool/               # Thread pool implementation
│   └── pool.go
├── test/                     # Test utilities
│   └── main.go              # TCP client for testing
└── README.md                 # This file
```

## Configuration

The server configuration can be modified in `internal/config/config.go`. Default settings include:

- Port: 8080
- Thread pool size: Configurable
- I/O multiplexing strategy: Auto-detected based on OS

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./internal/core
go test ./internal/server
```

### Adding New Features

1. **RESP Protocol**: Extend `internal/core/resp.go` for new data types
2. **Server Commands**: Add new commands in `internal/server/tcp_server.go`
3. **I/O Multiplexing**: Extend platform-specific implementations in `internal/core/io_multiplexing/`

## Architecture

The project follows a clean architecture pattern:

- **Core Layer**: RESP protocol implementation
- **Server Layer**: TCP server and connection management
- **I/O Layer**: Platform-specific I/O multiplexing
- **Config Layer**: Configuration management
- **Thread Pool**: Concurrent connection handling

## Performance

- **I/O Multiplexing**: Uses epoll/kqueue for efficient event handling
- **Thread Pool**: Prevents thread explosion under high load
- **Concurrent Connections**: Handles multiple clients simultaneously
- **Memory Efficient**: Minimal memory allocation per connection

## Troubleshooting

### Common Issues

1. **Port already in use**: Change the port in configuration or kill existing processes
2. **Permission denied**: Ensure you have permission to bind to the specified port
3. **Platform not supported**: Currently supports Linux (epoll) and macOS (kqueue)

### Debug Mode

Enable debug logging by modifying the configuration or adding debug flags.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## License

This project is open source and available under the [MIT License](LICENSE).

## Acknowledgments

- Redis for the RESP protocol specification
- Go community for excellent I/O multiplexing examples
- Contributors and maintainers
