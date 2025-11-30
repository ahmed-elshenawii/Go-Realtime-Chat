package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type Message struct {
	Type string `json:"type"`
	From string `json:"from,omitempty"`
	Body string `json:"body,omitempty"`
}

func main() {
	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	fmt.Println("Connected to server!")
	go receive(conn)

	reader := bufio.NewScanner(os.Stdin)
	for reader.Scan() {
		text := reader.Text()
		if text == "" {
			continue
		}
		fmt.Fprintln(conn, text)
	}

	if err := reader.Err(); err != nil {
		fmt.Println("Error reading input:", err)
	}
}

func receive(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		var msg Message
		line := scanner.Text()
		err := json.Unmarshal([]byte(line), &msg)
		if err != nil {
			fmt.Println("Invalid JSON:", line)
			continue
		}

		switch msg.Type {
		case "join":
			fmt.Printf("%s joined the chat\n", msg.From)
		case "leave":
			fmt.Printf("%s left the chat\n", msg.From)
		case "message":
			fmt.Printf("%s: %s\n", msg.From, msg.Body)
		}
	}
}
