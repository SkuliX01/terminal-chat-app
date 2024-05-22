package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

var (
	clients   = make(map[net.Conn]string)
	broadcast = make(chan string)
	mutex     = &sync.Mutex{}
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// Read the username
	username, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Failed reading username:", err)
		return
	}
	username = strings.TrimSpace(username)
	clients[conn] = username

	fmt.Printf("%s connected\n", username)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			delete(clients, conn)
			fmt.Printf("%s disconnected\n", username)
			break
		}

		message = strings.TrimSpace(message)

		// Broadcast the message only if it's not empty
		if message != "" {
			broadcast <- fmt.Sprintf("%s says: %s", username, message)
		}
	}
}

func broadcaster() {
	for {
		message := <-broadcast
		mutex.Lock()
		for conn := range clients {
			// Send the message to each client
			fmt.Fprintln(conn, message)
		}
		mutex.Unlock()
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
	defer listener.Close()

	go broadcaster()

	fmt.Println("Server started on :8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

