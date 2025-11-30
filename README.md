# Go Realtime Chat

This is a simple realtime chat system built with Go.
It supports multiple clients, broadcasting messages, and concurrency using goroutines and channels.

## Features

* Real-time broadcasting to all clients
* No self-echo: clients do not see their own messages
* Join/Leave notifications for other clients
* Each client has its own send/receive goroutines
* Clients list protected with sync.Mutex
* JSON-based protocol over TCP

## How it Works

1. Client connects to the server.
2. When a client joins, all other clients are notified:

   ```json
   {"type": "join", "from": "User1"}
   ```
3. When a client sends a message, it is broadcasted to all other clients (no self-echo):

   ```json
   {"type": "message", "from": "User1", "body": "Hello everyone!"}
   ```
4. When a client leaves, all other clients are notified:

   ```json
   {"type": "leave", "from": "User1"}
   ```

## Running the Project(Broadcast Chat Server)

### Server

```bash
go run server.go
```

### Client 1

```bash
go run client.go
```

### Client 2

```bash
go run client.go
```

### Client 3 (optional)

```bash
go run client.go
```

Start typing messages in the client terminals. Messages will be broadcasted to all other clients in real-time.

## Files

* `server.go`: TCP chat server with broadcasting and concurrency
* `client.go`: TCP chat client that sends/receives messages
* `README.md`: Project description and instructions
