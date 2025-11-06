package net

import (
	"bufio"
	"fmt"
	"net"
)

func simpleHandleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("Accepted connection from", conn.RemoteAddr())

	// Create a new scanner to read data from the connection
	scanner := bufio.NewScanner(conn)

	// Read data from the client in a loop
	for scanner.Scan() {
		data := scanner.Text()
		fmt.Printf("[%s]: %s\n", conn.RemoteAddr(), data)
	}

	// Check for any errors encountered during scanning
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading from connection:", err)
	}
}

func SimpleHttp() {
	// Listen for incoming connections on port 8080
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Scan Server is listening on port 8080")

	// Accept connections in a loop
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Handle each connection in a new goroutine
		go handleConnection(conn)
	}
}
