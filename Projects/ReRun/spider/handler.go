package socket

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

func createServer(s *Spider) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlerPage)
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		s.handleWebSocket(w, r)
	})
	server := http.Server{
		Addr:    ":7359",
		Handler: mux,
	}
	return &server
}

// handleWebSocket handles WebSocket connections.
func (s *Spider) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to upgrade to WebSocket:", err)
		return
	}

	connection := Connection{
		conn: conn,
		ID:   r.RemoteAddr, // Use the remote address as a simple identifier
	}

	fmt.Printf("Adding %s\n", connection.ID)
	s.addConn <- connection

	go func() {
		defer func() {
			fmt.Printf("Removing %s\n", connection.ID)
			s.removeConn <- connection
		}()
		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					fmt.Printf("Connection closed normally: %s\n", connection.ID)
				} else {
					fmt.Printf("Error reading message from %s: %v\n", connection.ID, err)
				}
				return
			}

			fmt.Printf("%s Received: %s\n", connection.ID, message)

			s.BroadcastMessage(messageType, message, connection)
		}
	}()
}

func handlerPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, htmlContent)
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
		console.log("Connected to WebSocket spider.");
	};

	socket.onmessage = function(event) {
		if (event.data == "SPIDER RELOAD") {
			location.reload()
		} 
		console.log("-> " + event.data);
	};

	socket.onclose = function(event) {
		console.log("Disconnected from Spider.");
	};

	socket.onerror = function(event) {
		console.error('Spider error:', event);
	};

	function sendMessage() {
		if (socket.readyState === WebSocket.OPEN) {
			socket.send("Hello, spider!");
			console.log("Hello, spider!");
		} else {
			console.log("WebSocket is not open.");
		}
	}
</script>
</html>
`
