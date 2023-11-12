package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	fmt.Println("Hello world")

	server, err := net.Listen("tcp", "localhost:7711")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer server.Close()

	fmt.Println("Server started")

	clients := []net.Conn{}
	addch := make(chan net.Conn, 1)

	go func() {
		for {
			client, err := server.Accept()
			if err != nil {
				fmt.Println(err)
			}
			addch <- client
		}
	}()

	for {
		select {
		case client := <-addch:
			clients = append(clients, client)
			_, _ = client.Write([]byte("Welcome\n"))
			go read(client)
		}
	}
}

func read(client net.Conn) {
	for {
		b := make([]byte, 1024)
		_, _ = client.Read(b)
		fmt.Println(string(b))
	}
}
