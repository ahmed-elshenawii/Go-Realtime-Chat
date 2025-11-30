package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"sync"
)

type Client struct {
	id   int
	name string
	conn net.Conn
	ch   chan Message
}

type Message struct {
	Type string `json:"type"` // "join", "leave", "message"
	From string `json:"from,omitempty"`
	Body string `json:"body,omitempty"`
}

var (
	clients      = make(map[int]*Client)
	clientsMutex sync.Mutex
	broadcastCh  = make(chan Message)
	nextID       = 1
)

func main() {
	go broadcaster()
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}
	fmt.Println("Server running on port 9000")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}

		clientsMutex.Lock()
		id := nextID
		nextID++
		name := fmt.Sprintf("User%d", id)
		client := &Client{id: id, name: name, conn: conn, ch: make(chan Message)}
		clients[id] = client
		clientsMutex.Unlock()

		joinMsg := Message{Type: "join", From: name}
		broadcastCh <- joinMsg

		go handleClient(client)
		go sendMessages(client)
	}
}

func broadcaster() {
	for msg := range broadcastCh {
		fmt.Println(stringify(msg))
		clientsMutex.Lock()
		for _, c := range clients {
			// لمنع self-echo: إذا رسالة من نفس العميل، تجاهلها
			if msg.Type == "message" && msg.From == c.name {
				continue
			}
			// لمنع العميل من رؤية join/leave نفسه
			if (msg.Type == "join" || msg.Type == "leave") && msg.From == c.name {
				continue
			}
			c.ch <- msg
		}
		clientsMutex.Unlock()
	}
}

func handleClient(c *Client) {
	reader := bufio.NewScanner(c.conn)
	for reader.Scan() {
		body := reader.Text()
		if body == "" {
			continue
		}
		msg := Message{Type: "message", From: c.name, Body: body}
		broadcastCh <- msg
	}
	clientsMutex.Lock()
	delete(clients, c.id)
	clientsMutex.Unlock()

	leaveMsg := Message{Type: "leave", From: c.name}
	broadcastCh <- leaveMsg
	c.conn.Close()
}

func sendMessages(c *Client) {
	for msg := range c.ch {
		data, _ := json.Marshal(msg)
		fmt.Fprintln(c.conn, string(data))
	}
}

func stringify(msg Message) string {
	data, _ := json.Marshal(msg)
	return string(data)
}
