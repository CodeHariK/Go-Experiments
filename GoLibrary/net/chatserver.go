package net

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

// User represents a user in the chat system.
type User struct {
	ID        int
	Username  string
	CreatedAt time.Time
}

// Message represents a message sent in the chat system.
type Message struct {
	ID        int
	User      User
	Content   string
	Timestamp time.Time
}

// ChatServer represents a basic chat server.
type ChatServer struct {
	clients      map[net.Conn]User
	addClient    chan net.Conn
	removeClient chan net.Conn
	broadcast    chan Message
	nextUserID   int
}

func NewChatServer() *ChatServer {
	return &ChatServer{
		clients:      make(map[net.Conn]User),
		addClient:    make(chan net.Conn),
		removeClient: make(chan net.Conn),
		broadcast:    make(chan Message),
		nextUserID:   1,
	}
}

func (cs *ChatServer) Start() {
	for {
		select {
		case conn := <-cs.addClient:
			username := cs.askForUsername(conn)
			user := User{ID: cs.nextUserID, Username: username, CreatedAt: time.Now()}
			cs.nextUserID++
			cs.clients[conn] = user
			go cs.handleClient(conn, user)
			cs.broadcastMessage(Message{User: user, Content: "joined the chat", Timestamp: time.Now()}, user)
		case conn := <-cs.removeClient:
			user := cs.clients[conn]
			delete(cs.clients, conn)
			closeConnection(conn)
			cs.broadcastMessage(Message{User: user, Content: "left the chat", Timestamp: time.Now()}, user)
		case message := <-cs.broadcast:
			cs.broadcastMessage(message, message.User)
		}
	}
}

func (cs *ChatServer) handleClient(conn net.Conn, user User) {
	defer func() {
		fmt.Println(conn.RemoteAddr(), "Connection closed")
		fmt.Fprintln(conn, "Connection closed")
		cs.removeClient <- conn
	}()

	cs.sendSystemMessage(conn, "Welcome to the chat, "+user.Username+"!")

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		messageContent := scanner.Text()

		if messageContent == "quit" || messageContent == "exit" {
			cs.sendSystemMessage(conn, "Goodbye, "+user.Username+"!")
			return
		}

		message := Message{User: user, Content: messageContent, Timestamp: time.Now()}

		// Set a write deadline of 10 seconds
		if err := conn.SetWriteDeadline(time.Now().Add(10 * time.Second)); err != nil {
			fmt.Println("Error setting write deadline:", err)
			return
		}

		// Write the message to the connection
		_, err := fmt.Fprintf(conn, "[%s] %s: %s\n", message.Timestamp.Format("15:04:05"), message.User.Username, message.Content)
		if err != nil {
			fmt.Println("Error writing to connection:", err)
			return
		}

		cs.broadcast <- message
	}
}

func (cs *ChatServer) broadcastMessage(message Message, sender User) {
	for conn, user := range cs.clients {
		if user == sender {
			fmt.Fprint(conn, "> ")
		} else {
			fmt.Fprint(conn, "* ")
		}
		fmt.Fprintf(conn, "[%s] %s: %s\n", message.Timestamp.Format("15:04:05"), message.User.Username, message.Content)
	}
}

func (cs *ChatServer) sendSystemMessage(conn net.Conn, message string) {
	fmt.Fprintf(conn, "[System] %s\n", message)
}

func (cs *ChatServer) askForUsername(conn net.Conn) string {
	cs.sendSystemMessage(conn, "Please enter your username:")
	username, _ := bufio.NewReader(conn).ReadString('\n')
	return strings.TrimSpace(username)
}

func closeConnection(conn net.Conn) {
	if err := conn.Close(); err != nil {
		fmt.Println("Error closing connection:", err)
	}
}

func ChatServerStart() {
	server := NewChatServer()

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Chat Server is listening on port 8080")

	go server.Start()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		server.addClient <- conn
	}
}
