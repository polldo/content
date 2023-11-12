package main

import (
	"fmt"
	"net"
	"os"
	"strings"
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
	msgch := make(chan string, 1)

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
			go read(client, msgch)

		case msg := <-msgch:
			fmt.Println(msg)
			for _, c := range clients {
				c.Write([]byte(msg))
			}
		}
	}
}

func read(client net.Conn, msgch chan string) {
	for {
		b := make([]byte, 1024)
		_, _ = client.Read(b)
		txt := strings.TrimRight(string(b), "\n\r")
		msgch <- txt
	}
}
