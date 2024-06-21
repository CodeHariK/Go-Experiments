package socket

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

// Connection represents a WebSocket connection.
type Connection struct {
	Conn *websocket.Conn
	ID   string
}

var (
	mu sync.Mutex

	conns       []Connection
	connections = make(chan []Connection)
	addConn     = make(chan Connection)
	removeConn  = make(chan Connection)

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// Check the origin of the request and decide whether to accept or reject it
			return true
		},
	}
)

func manageConnections() {
	for {
		select {
		case newConn := <-addConn:
			fmt.Printf("Add new connection %s\n", newConn.ID)
			mu.Lock()
			conns = append(conns, newConn)
			mu.Unlock()
			connections <- conns
		case removeConn := <-removeConn:
			fmt.Printf("Remove connection %s\n", removeConn.ID)
			mu.Lock()
			for i, conn := range conns {
				if conn.ID == removeConn.ID {
					conn.Conn.Close()
					conns = append(conns[:i], conns[i+1:]...)
					break
				}
			}
			mu.Unlock()
			connections <- conns
		}
	}
}

// broadcastMessage sends a message to all active connections.
func broadcastMessage(messageType int, message []byte, sender Connection) {
	for _, conn := range conns {
		if conn == sender {
			continue
		}
		err := conn.Conn.WriteMessage(messageType, []byte(fmt.Sprintf("%s %s", conn.ID, string(message))))
		fmt.Printf("Sent to : %s %s\n", conn.ID, string(message))

		if err != nil {
			fmt.Printf("Error broadcasting message to %s: %v\n", conn.ID, err)
			removeConn <- conn
		}
	}
}

// Handler function
func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, htmlContent)
}

func StartServer() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer close(stop)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(w, r)
	})

	go manageConnections()

	go func() {
		time.Sleep(time.Second * 20)
		fmt.Println("Hello")
		broadcastMessage(1, []byte("Hello"), Connection{})
	}()

	fmt.Println("Server started on :7359")

	server := http.Server{
		Addr:    ":7359",
		Handler: mux,
	}

	// Run server in a goroutine to be able to shut it down gracefully
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("ListenAndServe error: %v\n", err)
		}
	}()

	// Block until we receive a signal
	<-stop

	fmt.Println("Shutting down server...")

	// Create a context with a timeout to allow the server to shut down gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the server gracefully
	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Server forced to shutdown: %v\n", err)
	}

	fmt.Println("Server stopped")
}

// handleWebSocket handles WebSocket connections.
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to upgrade to WebSocket:", err)
		return
	}

	connection := Connection{
		Conn: conn,
		ID:   r.RemoteAddr, // Use the remote address as a simple identifier
	}

	fmt.Printf("Adding %s\n", connection.ID)
	addConn <- connection

	go func() {
		defer func() {
			fmt.Printf("Removing %s\n", connection.ID)
			removeConn <- connection
		}()
		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("Error reading message:", err)
				return
			}

			fmt.Printf("%s Received: %s\n", connection.ID, message)

			broadcastMessage(messageType, message, connection)
		}
	}()
}

const htmlContent = `
<!DOCTYPE html>
<html lang="en">
<head>
    <title>My Web Page</title>
</head>
<style>
	body {
		background: black;
	}
</style>
<body>
	<button onclick="sendMessage()">Send Message</button>
</body>
<script>
	let socket = new WebSocket("ws://localhost:7359/ws");

	socket.onopen = function(event) {
		console.log("Connected to WebSocket server.");
	};

	socket.onmessage = function(event) {
		if (event.data == "reload") {
			location.reload()
		} 
		console.log("Received from server: " + event.data);
	};

	socket.onclose = function(event) {
		console.log("Disconnected from WebSocket server.");
	};

	socket.onerror = function(event) {
		console.error('WebSocket error:', event);
	};

	function sendMessage() {
		if (socket.readyState === WebSocket.OPEN) {
			socket.send("Hello, server!");
			console.log("Message sent.");
		} else {
			console.log("WebSocket is not open.");
		}
	}
</script>
</html>
`
