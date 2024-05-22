
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	var address string

	fmt.Println("Enter IP address with port of chat you want to connect to:")
	fmt.Scanln(&address)

	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Failed to connect to the server:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter your username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	// Send username to server
	fmt.Fprintln(conn, username)

	fmt.Printf("Connected to: %s as %s\n", conn.RemoteAddr(), username)

	go func() {
		serverReader := bufio.NewReader(conn)
		for {
			message, err := serverReader.ReadString('\n')
			if err != nil {
				fmt.Println("Disconnected from server")
				os.Exit(0)
			}
			fmt.Print(message)
		}
	}()

	for {
		fmt.Print("> ")
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message)
		if message != "" {
			// Send message to server only if it's not empty
			fmt.Fprintln(conn, message)
		}
	}
}

