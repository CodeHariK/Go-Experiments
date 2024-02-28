package net

import (
	"fmt"
	"net"
	"time"
)

func Heartbeat() {
	// Start the server
	go startServer()

	// Wait for a moment to ensure the server is running
	time.Sleep(1 * time.Second)

	// Start the client
	startClient()
}

func startServer() {
	// Server listens on port 8080
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on port 8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Handle the connection in a new goroutine
		go hearthandleConnection(conn)
	}
}

func hearthandleConnection(conn net.Conn) {
	defer conn.Close()

	// Read and respond to incoming messages
	for {
		buffer := make([]byte, 1024)
		_, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading:", err)
			return
		}

		message := string(buffer)
		fmt.Printf("Received message: %s", message)

		// Respond to the heartbeat
		response := "Heartbeat response\n"
		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Error responding:", err)
			return
		}
	}
}

func startClient() {
	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	// Start sending periodic heartbeat messages
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// Send a heartbeat message to the server
		message := "Heartbeat\n"
		_, err := conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error sending heartbeat:", err)
			return
		}

		// Read the response from the server
		buffer := make([]byte, 1024)
		_, err = conn.Read(buffer)
		if err != nil {
			fmt.Println("Error receiving heartbeat response:", err)
			return
		}

		response := string(buffer)
		fmt.Printf("Received heartbeat response: %s", response)
	}
}
