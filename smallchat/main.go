package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type User struct {
	Client
	nick string
}

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

	users := make(map[int]User)
	addch := make(chan User, 1)
	quitch := make(chan Client, 1)
	msgch := make(chan Msg, 1)

	go func() {
		for cnt := 0; ; cnt++ {
			conn, err := server.Accept()
			if err != nil {
				fmt.Printf("error accepting client: %v\n", err)
			}
			addch <- User{Client: Client{Conn: conn, id: cnt}, nick: fmt.Sprintf("user-%d", cnt)}
		}
	}()

	for {
		select {
		case user := <-addch:
			users[user.id] = user
			_, _ = user.Write([]byte("Welcome\n"))
			go read(user.Client, msgch, quitch)

		case client := <-quitch:
			delete(users, client.id)
			_ = client.Close()
			fmt.Printf("Client[%d] disconnected\n", client.id)

		case msg := <-msgch:
			if err := process(users, msg); err != nil {
				err := fmt.Sprintf("Error: %s\n", err.Error())
				_, _ = users[msg.senderID].Write([]byte(err))
			}
		}
	}
}

func process(users map[int]User, msg Msg) error {
	sender, ok := users[msg.senderID]
	if !ok {
		fmt.Printf("sender of message[%v] has gone\n", msg)
		return nil
	}

	if strings.HasPrefix(msg.text, "/") {
		args := strings.Split(msg.text, " ")
		switch args[0] {
		case "/nick":
		default:
			return fmt.Errorf("cmd[%v] is unknwon", args[0])
		}
	}

	txt := fmt.Sprintf("%s: %s\n", sender.nick, msg.text)
	for _, c := range users {
		if msg.senderID == c.id {
			continue
		}
		c.Write([]byte(txt))
	}

	return nil
}

func read(client Client, msgch chan Msg, quitch chan Client) {
	reader := bufio.NewReader(client.Conn)
	for {
		txt, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading client[%d]:%v\n", client.id, err)
			quitch <- client
			return
		}

		txt = strings.TrimRight(txt, "\n\r")
		msgch <- Msg{senderID: client.id, text: txt}
	}
}
