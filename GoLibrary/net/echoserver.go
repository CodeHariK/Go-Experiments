package net

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

// Message represents the data structure to hold information for echoing
type EchoMessage struct {
	Buffer  []byte
	Length  int
	Address net.Addr
}

func EchoServer() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Echo server is listening on port 8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)

	// //=========
	// go func() {
	// 	fmt.Println("Copy")
	// 	io.Copy(conn, conn)
	// }()

	//=========
	scanner := bufio.NewScanner(conn)
	go func() {
		for scanner.Scan() {
			fmt.Println("Scanner")
			io.WriteString(conn, scanner.Text())
		}
	}()

	fmt.Println("==========")

	//=========
	for {
		fmt.Println("Buffer")
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading:", err)
			return
		}

		// Create a Message struct with the relevant information
		message := EchoMessage{
			Buffer:  buffer[:n],
			Length:  n,
			Address: conn.RemoteAddr(),
		}

		// Echo the entire message back to the client
		_, err = conn.Write(serializeMessage(message))
		if err != nil {
			fmt.Println("Error writing:", err)
			return
		}

		// Print information about the echoed message
		fmt.Printf("Echoed message to %s\n", messageToString(message))
	}
}

func serializeMessage(message EchoMessage) []byte {
	return []byte(fmt.Sprintf("%s (%d) -> %s", message.Address, message.Length, message.Buffer))
}

func messageToString(message EchoMessage) string {
	// Convert the message to a string for printing
	return fmt.Sprintf("{ Length: %d, Address: %s, Buffer: %s }", message.Length, message.Address, message.Buffer)
}
