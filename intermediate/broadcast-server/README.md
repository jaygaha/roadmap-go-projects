# Broadcast Server

A simple WebSocket-based broadcast server built in Go that allows multiple clients to connect and exchange messages in real-time.

## Features

- Real-time message broadcasting to all connected clients
- WebSocket-based communication
- CLI-based server and client
- User join/leave notifications
- Live user count
- Graceful shutdown handling
- Clean, organized code structure

## Installation

```bash
git clone https://github.com/jaygaha/roadmap-go-projects.git
cd intermediate/broadcast-server
go mod tidy
go build
```

## Usage

The application has two main commands:

1. Start the server:
```bash
./broadcast-server start [--port PORT]
```
Default port is 8080 if not specified.

2. Connect as a client:
```bash
./broadcast-server connect [--host HOST] [--port PORT] [--username USERNAME]
```

If username is not provided, you will be prompted to enter one.

Default values:
- Host: localhost
- Port: 8080

## Commands

- `start` - Start the broadcast server
- `connect` - Connect to the server as a client

## Client Commands

Once connected as a client:
- Type any message and press Enter to broadcast
- Type `quit` or `exit` to disconnect
- Press `Ctrl+C` for graceful shutdown

The client will display:
- Join notifications when users connect
- Leave notifications when users disconnect
- User count updates
- Messages from other users

## Architecture

```
├── cmd/
│   └── root.go              # CLI entry point
├── internal/
│   ├── server/
│   │   ├── server.go        # HTTP server and routes
│   │   ├── hub.go           # Client management and broadcasting
│   │   └── client.go        # WebSocket client connection handling
│   ├── client/
│   │   └── client.go        # CLI client implementation
│   ├── config/
│   │   └── config.go        # Configuration management
│   └── handler/
│       ├── client_handler.go # Client command handler
│       ├── help_handler.go   # Help/usage command handler
│       └── start_handler.go  # Server start command handler
├── pkg/
│   └── message/
│       └── sys-message.go   # Message structures and utilities
├── go.mod
├── go.sum
├── README.md
└── main.go                  # Main application entry point
```

## Message Types

- `message` - Regular chat messages
- `join` - User joined notification
- `leave` - User left notification  
- `user_count` - Current user count update

## Example

Terminal 1 (Server):
```bash
$ ./broadcast-server start
Starting broadcast server on port 8080...
Server starting on :8080
Client connected: Alice (Total: 1)
Client connected: Bob (Total: 2)
```

Terminal 2 (Client Alice):
```bash
$ ./broadcast-server connect --username Alice
Connecting to localhost:8080 as Alice...
Users online: 1
[15:04:05] Bob joined the chat
Users online: 2
> Hello everyone!
[15:04:10] Bob: Hi Alice!
> 
```

Terminal 3 (Client Bob):
```bash
$ ./broadcast-server connect --username Bob  
Connecting to localhost:8080 as Bob...
[15:04:05] Alice: Hello everyone!
> Hi Alice!
> 
```

## Technical Details

- Uses `gorilla/websocket` for WebSocket functionality
- Implements hub pattern for client management
- Concurrent handling of multiple clients using goroutines
- JSON-based message protocol with support for handling multiple JSON objects in a single message
- Robust message parsing with the SplitJSONMessages function to prevent JSON parsing errors
- Graceful shutdown with signal handling
- Clean separation of concerns with layered architecture


## Key Features:

1. **Simple & Clean**: Minimal external dependencies, clear code structure
2. **Real-time Communication**: WebSocket-based broadcasting
3. **CLI Interface**: Easy-to-use command-line interface
4. **Concurrent**: Handles multiple clients simultaneously
5. **Robust Message Handling**: SplitJSONMessages function to handle multiple JSON objects in a single message
6. **Graceful Shutdown**: Proper cleanup on exit
7. **Organized**: Clear separation of server, client, and shared code
8. **Extensible**: Easy to add new message types and features

To run:
```bash
go mod tidy
go build
./broadcast-server start  # In one terminal
./broadcast-server connect # In another terminal
```

## Project Link

- [Broadcast Server](https://roadmap.sh/projects/broadcast-server)

## Acknowledgments

- Part of the Go programming language learning roadmap projects
- Created by [jaygaha](https://github.com/jaygaha)
