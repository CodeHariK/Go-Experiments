// package main

// import (
// 	"log"

// 	"github.com/gorilla/websocket"
// )

// func main() {
// 	// Connect to the WebSocket server
// 	u := "ws://localhost:7359/ws"
// 	c, _, err := websocket.DefaultDialer.Dial(u, nil)
// 	if err != nil {
// 		log.Fatal("dial:", err)
// 	}
// 	defer c.Close()

// 	done := make(chan struct{})

// 	// Read messages from the server
// 	go func() {
// 		defer close(done)
// 		for {
// 			_, message, err := c.ReadMessage()
// 			if err != nil {
// 				log.Println("read:", err)
// 				return
// 			}
// 			log.Printf("recv: %s", message)
// 		}
// 	}()

// 	<-done
// }

package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	for {
		err := connectAndListen(interrupt)
		if err != nil {
			log.Printf("Connection error: %v", err)
			time.Sleep(time.Second * 2) // Wait before trying to reconnect
		} else {
			break
		}
	}
}

func connectAndListen(interrupt chan os.Signal) error {
	u := "ws://localhost:7359/ws"
	c, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		return err
	}
	defer c.Close()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, os.Kill, syscall.SIGINT)

	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Printf("read: %v", err)
				return
			}
			fmt.Printf(string(message))
		}
	}()

	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return nil
		case <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(`{"type":"ping"}`))
			if err != nil {
				log.Println("write:", err)
				return nil
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Printf("write close: %v", err)
				return err
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return nil
		}
	}
}
