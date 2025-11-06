package io

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"net"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("Accepted connection from", conn.RemoteAddr())

	// Wrap the connection with a Rot13Reader
	r13 := &Rot13Reader{r: conn}

	if rand.Intn(4) < 2 {
		io.Copy(conn, r13)
	} else {
		// Create a bufio.Scanner to read lines
		scanner := bufio.NewScanner(r13)

		// Echo the Rot13 encrypted lines back to the client
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println("Received:", line)

			// Echo the Rot13 encrypted line back to the client
			conn.Write([]byte(line + "\n"))
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading input:", err)
		}
	}
}

func Rot13Echo() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Rot13 Echo Server is listening on port 8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn)
	}
}
