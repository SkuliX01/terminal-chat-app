package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
)
  
var (
    clients = make(map[net.Conn]bool)
    broadcast = make(chan string)
    mutex = &sync.Mutex{}
)


func handleConnection(conn net.Conn) {
  
  defer conn.Close()
  clients[conn] = true
  
  reader := bufio.NewReader(conn)

  for {
    
    message, err := reader.ReadString('\n')

    if err != nil {
      delete(clients, conn)
      break
    }
    
    broadcast <- message
  }
}

func broadcaster() {
  for {
    message := <- broadcast
    mutex.Lock()


    for conn := range clients {
      fmt.Printf("%s says : %s", conn, message)
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

  fmt.Println("server started on port :8080")

  for {
    
    conn, err := listener.Accept()

    if err != nil {
      fmt.Println("Error Accepting connection: ", err)
      continue
    }

    go handleConnection(conn)
  }

}

