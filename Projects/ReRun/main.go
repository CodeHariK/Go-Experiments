// package main

// import (
// 	"flag"
// 	"fmt"
// 	"log"
// 	"sync"
// 	"time"

// 	"github.com/codeharik/rerun/helper"
// 	"github.com/codeharik/rerun/spider"
// 	"github.com/codeharik/rerun/watcher"
// )

// func main() {
// 	flagKillPorts := flag.String("k", "", "Optional Kill Ports")
// 	flagReRunDelay := flag.Int("t", -1, "Optional Rerun Delay Time in Milliseconds[Min 100]")

// 	flag.Parse()

// 	nonFlagArgs := flag.Args()

// 	if len(nonFlagArgs) < 2 {
// 		fmt.Println("ReRun: Monitor a directory and automatically execute a command when files change, or rerun the command on a set interval.")
// 		flag.PrintDefaults()
// 		fmt.Println()
// 		fmt.Println("Usage: go run main.go [-k optional kill ports] [-t optional rerun delay time] <watch directory> <run command>")
// 		fmt.Println("Usage: go run main.go ../Hello \"go run ../Hello/main.go\"")
// 		fmt.Println("Usage: go run main.go -k=8080,3000 -t=4000 ../Hello \"go run ../Hello/main.go\"")
// 		fmt.Println()
// 		fmt.Println("Usage: rerun [-k optional kill ports] [-t optional rerun delay time] <watch directory> <run command>")
// 		fmt.Println("Usage: rerun -k=8080,3000 -t=4000 ../Hello \"go run ../Hello/main.go\"")
// 		fmt.Println("Usage: rerun ../Hello \"go run ../Hello/main.go\"")
// 		return
// 	}

// 	killPortsString := *flagKillPorts
// 	rerunTimer := time.Duration(*flagReRunDelay) * 1000000
// 	if rerunTimer > 0 && rerunTimer < time.Millisecond*100 {
// 		log.Fatal("Min 100 milliseconds delay required")
// 	}

// 	directory := flag.Arg(0)
// 	command := flag.Arg(1)
// 	killPorts := []int{}
// 	if killPortsString != "" {
// 		k, err := helper.ParseStringInts(killPortsString)
// 		if err == nil {
// 			killPorts = k
// 		}
// 	}

// 	var wg sync.WaitGroup
// 	defer wg.Wait()

// 	stdLogs := make(map[string][]string)

// 	spider := spider.NewSpider(directory, stdLogs, stdLogs)
// 	spider.StartSpider(&wg)

// 	w := watcher.NewWatcher(
// 		command, rerunTimer, killPorts, directory,
// 		spider,
// 		stdLogs,
// 	)
// 	w.StartWatcher()
// }

///////
///////
///////
///////

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
	go watchDirectory("./TwirpHat")

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

	go executeCommand("go run TwirpHat/cmd/server/main.go")

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

					go executeCommand("go run TwirpHat/cmd/server/main.go")
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
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
		var msg string
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
	}
}

func handleMessages() {
	for {
		msg := <-broadcast

		for client := range clients {
			fmt.Println(msg)
			m, _ := json.Marshal(msg)
			fmt.Println(string(m))
			err := client.WriteJSON(`Hello I am fine`)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func handlePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(200)
	// fmt.Fprint(w, htmlContent)

	http.ServeFile(w, r, "/Users/Shared/Go/Go-Experiments/Projects/ReRun/spider/spider.html")
}
