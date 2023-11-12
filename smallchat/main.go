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

	conn, err := server.Accept()
	if err != nil {
		fmt.Println(err)
	}

	_, _ = conn.Write([]byte("Welcome\n"))

	b := make([]byte, 1024)
	_, _ = conn.Read(b)
	fmt.Println(string(b))
}
