package spider

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
	conn *websocket.Conn
	ID   string
}

type Spider struct {
	mu         sync.Mutex
	conns      map[string]Connection
	addConn    chan Connection
	removeConn chan Connection
}

func NewSpider() *Spider {
	return &Spider{
		conns:      make(map[string]Connection),
		addConn:    make(chan Connection),
		removeConn: make(chan Connection),
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Check the origin of the request and decide whether to accept or reject it
		return true
	},
}

func (s *Spider) StartSpider(wg *sync.WaitGroup) {
	server := createServer(s)

	go s.manageConnections()

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Spider started on :7359")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("ListenAndServe error: %v\n", err)
		}
		fmt.Println("Spider stopped")
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-done

		// Create a context with a timeout to allow the server to shut down gracefully
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			fmt.Printf("Spider forced to shutdown: %v\n", err)
		}
	}()
}

func (s *Spider) manageConnections() {
	for {
		select {
		case newConn := <-s.addConn:
			s.mu.Lock()
			s.conns[newConn.ID] = newConn // Add to map
			s.mu.Unlock()
			fmt.Printf("Add new connection %s\n", newConn.ID)
		case removeConn := <-s.removeConn:
			s.mu.Lock()
			if conn, ok := s.conns[removeConn.ID]; ok {
				conn.conn.Close()              // Close WebSocket connection
				delete(s.conns, removeConn.ID) // Remove from map
			}
			s.mu.Unlock()
			fmt.Printf("Remove connection %s\n", removeConn.ID)
		}
	}
}

// BroadcastMessage sends a message to all active connections.
func (s *Spider) BroadcastMessage(messageType int, message []byte, sender Connection) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, conn := range s.conns {
		if conn == sender {
			continue
		}
		err := conn.conn.WriteMessage(messageType, []byte(fmt.Sprintf("%s %s", sender.ID, string(message))))
		fmt.Printf("-> %s %s\n", sender.ID, string(message))

		if err != nil {
			fmt.Printf("Error broadcasting message to %s: %v\n", conn.ID, err)
			s.removeConn <- conn
		}
	}
}
