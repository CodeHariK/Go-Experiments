package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"sync"

	wa "github.com/codeharik/rerun/watcher"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
)

var (
	clients    = make(map[*websocket.Conn]bool)
	broadcast  = make(chan string)
	upgrader   = websocket.Upgrader{}
	cmdMutex   sync.Mutex
	currentCmd *exec.Cmd
	logCh      = make(chan string)
)

func main() {
	// Start file watcher
	go watchDirectory("../server")

	http.HandleFunc("/", handlePage)
	http.HandleFunc("/ws", handleConnections)

	// Start listening for incoming messages
	go handleMessages()

	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()
	go func() {
		defer wg.Done()
		for log := range logCh {
			broadcast <- log
		}
	}()

	// Start the server
	log.Println("HTTP server started on :7359")
	err := http.ListenAndServe(":7359", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func watchDirectory(directory string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)

	go executeCommand("go run ../server/main.go")

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)

					if currentCmd != nil && currentCmd.Process != nil {
						cmdMutex.Unlock()
						wa.PortKiller([]int{8080})
						err := currentCmd.Process.Kill()
						fmt.Println(err)
						err = currentCmd.Wait()
						fmt.Println(err)
					}

					go executeCommand("go run ../server/main.go")
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("fs error:", err)
			}
		}
	}()

	err = watcher.Add(directory)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func executeCommand(command string) {
	cmdMutex.Lock()

	fmt.Println(command)

	currentCmd = exec.Command("sh", "-c", command)

	stdoutPipe, err := currentCmd.StdoutPipe()
	if err != nil {
		logCh <- fmt.Sprintf("Error creating stdout pipe: %v", err)
		return
	}

	stderrPipe, err := currentCmd.StderrPipe()
	if err != nil {
		logCh <- fmt.Sprintf("Error creating stderr pipe: %v", err)
		return
	}

	// Start the command
	if err := currentCmd.Start(); err != nil {
		logCh <- fmt.Sprintf("Error starting command: %v", err)
		return
	}

	// Function to read from pipe and send to channel
	readPipe := func(pipe *bufio.Reader, prefix string) {
		for {
			line, err := pipe.ReadBytes('\n')
			if len(line) > 0 {
				logCh <- fmt.Sprintf("%s: %s", prefix, bytes.TrimRight(line, "\n"))
			}
			if err != nil {
				if err.Error() != "EOF" {
					logCh <- fmt.Sprintf("Error reading %s pipe: %v", prefix, err)
				}
				break
			}
		}
	}

	// Read stdout and stderr in separate goroutines
	stdoutReader := bufio.NewReader(stdoutPipe)
	stderrReader := bufio.NewReader(stderrPipe)

	go readPipe(stdoutReader, "STDOUT")
	go readPipe(stderrReader, "STDERR")

	fmt.Println("Wait")

	if err := currentCmd.Wait(); err != nil {
		logCh <- fmt.Sprintf("Error waiting for command: %v", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	clients[ws] = true

	for {
		messageType, message, err := ws.ReadMessage()
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}

		switch messageType {
		case websocket.TextMessage:
			log.Printf("recv: %s", message)
		case websocket.PingMessage:
			fmt.Println("Pong")
			err := ws.WriteMessage(websocket.PongMessage, nil)
			if err != nil {
				log.Printf("pong error: %v", err)
			}
		case websocket.PongMessage:
			log.Println("pong received")
		}
	}
}

func handleMessages() {
	for {
		msg := <-broadcast

		for client := range clients {
			m, _ := json.Marshal(msg)
			fmt.Println(string(m))
			err := client.WriteJSON(string(m))
			if err != nil {
				log.Printf("handleMessages error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func handlePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(200)

	fmt.Fprintln(w, `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>WebSocket Client</title>
</head>
<body style="background:black;color:white;text-align: center;align-content: center;font: 30px monospace;">
    <span>Hello</span>
    <script>
        let socket;
		let retryCount = 0;
		const maxRetries = 5;

		function connect() {
			socket = new WebSocket('ws://localhost:7359/ws');

			socket.onopen = function(event) {
				console.log('WebSocket connection established.');
				retryCount = 0;  // Reset the retry count upon successful connection
			};

			socket.onmessage = function(event) {
				console.log('Message from server:', event.data);
			};

			socket.onclose = function(event) {
				console.log('WebSocket closed:', event);
				if (retryCount < maxRetries) {
					const retryDelay = Math.pow(3, retryCount) * 400;
					console.log('Retrying in '+retryDelay+'ms...');
					setTimeout(connect, retryDelay);
					retryCount++;
				} else {
					console.log('Max retries reached. No further attempts will be made.');
				}
			};

			socket.onerror = function(error) {
				console.error('WebSocket error:', error);
				socket.close();
			};
		}

		connect();
    </script>
</body>
</html>`)
}
