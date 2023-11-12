package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

type Client struct {
	net.Conn
	id int
}

type Msg struct {
	senderID int
	text     string
}

func main() {
	fmt.Println("Hello world")

	server, err := net.Listen("tcp", "localhost:7711")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer server.Close()

	fmt.Println("Server started")

	clients := make(map[int]Client)
	addch := make(chan Client, 1)
	msgch := make(chan Msg, 1)

	go func() {
		for cnt := 0; ; cnt++ {
			conn, err := server.Accept()
			if err != nil {
				fmt.Printf("error accepting client: %v\n", err)
			}
			addch <- Client{Conn: conn, id: cnt}
		}
	}()

	for {
		select {
		case client := <-addch:
			clients[client.id] = client
			_, _ = client.Write([]byte("Welcome\n"))
			go read(client, msgch)

		case msg := <-msgch:
			fmt.Println(msg)
			for _, c := range clients {
				if msg.senderID == c.id {
					continue
				}
				c.Write([]byte(msg.text))
			}
		}
	}
}

func read(client Client, msgch chan Msg) {
	for {
		b := make([]byte, 1024)
		_, _ = client.Read(b)
		txt := strings.TrimRight(string(b), "\n\r")
		msgch <- Msg{senderID: client.id, text: txt}
	}
}
