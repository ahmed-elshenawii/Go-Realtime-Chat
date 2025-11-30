package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

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
		line := scanner.Text()

		fmt.Println(line)
	}
}
