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
		select {
		case <-interrupt:
			return
		default:
			err := connectAndListen()
			if err != nil {
				log.Printf("Connection error: %v", err)
				time.Sleep(time.Second * 2)
			} else {
				break
			}
		}
	}
}

func connectAndListen() error {
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
		case <-ticker.C:
			err := c.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				log.Println("write:", err)
				return nil
			}
		case <-done:
			log.Println("interrupt")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Printf("write close: %v", err)
				return err
			}
			return nil
		}
	}
}
