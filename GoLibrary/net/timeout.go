package net

import (
	"context"
	"fmt"
	"net"
	"time"
)

const (
	serverAddr        = "localhost:8080"
	heartbeatInterval = 3 * time.Second
	connectionTimeout = 10 * time.Second
)

func Timeout() {
	// Start the server
	go timeoutStartServer()

	// Wait for a moment to ensure the server is running
	time.Sleep(1 * time.Second)

	// Start the client
	timeoutStartClient()
}

func timeoutStartServer() {
	listener, err := net.Listen("tcp", serverAddr)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on", serverAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Handle the connection in a new goroutine
		go timeoutHandleConnection(conn)
	}
}

func timeoutHandleConnection(conn net.Conn) {
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set a timeout for the connection
	go func() {
		select {
		case <-time.After(connectionTimeout):
			fmt.Println("Connection timed out")
		case <-ctx.Done():
			// Connection closed or heartbeat received, do nothing
		}
	}()

	// Read and respond to incoming messages (heartbeats)
	for {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading:", err)
			break // Break from the loop to ensure cancel() is called
		}

		message := string(buffer[:n])
		fmt.Printf("Received message: %s", message)

		// Respond to the heartbeat
		response := "Heartbeat response\n"
		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Error responding:", err)
			break // Break from the loop to ensure cancel() is called
		}

		// Reset the timeout on each received heartbeat
		cancel()
		ctx, cancel = context.WithCancel(context.Background())
	}

	// Ensure cancel() is called in case of errors
	cancel()
}

func timeoutStartClient() {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	// Start sending periodic heartbeat messages
	ticker := time.NewTicker(heartbeatInterval)
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
