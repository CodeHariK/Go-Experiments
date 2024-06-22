package watcher

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/codeharik/rerun/helper"
	socket "github.com/codeharik/rerun/spider"
	"github.com/fsnotify/fsnotify"
)

type watcher struct {
	shellProcess *os.Process
	childProcess *os.Process
	counter      int32

	command       string
	reRunDuration time.Duration
	killPorts     []int
	directory     string

	spider *socket.Spider
}

func NewWatcher(
	command string,
	reRunDuration time.Duration,
	killPorts []int,
	directory string,

	spider *socket.Spider,
) *watcher {
	return &watcher{
		command:       command,
		reRunDuration: reRunDuration,
		killPorts:     killPorts,
		directory:     directory,

		spider: spider,
	}
}

func (w *watcher) StartWatcher() {
	done := make(chan os.Signal)
	signal.Notify(done, syscall.SIGINT, os.Kill, os.Interrupt)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	if w.reRunDuration >= time.Millisecond*100 {
		helper.TickerFunction(
			w.reRunDuration,
			func() {
				w.runCommand()
			},
		)
	}

	var wg sync.WaitGroup
	defer wg.Wait()
	wg.Add(1)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					fmt.Println("Watcher.Events closed:", err)
					return
				}

				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
					fmt.Println("modified file:", event.Name)
					w.runCommand()
				}

				if event.Op&fsnotify.Create == fsnotify.Create {
					info, err := os.Stat(event.Name)
					if err == nil && info.IsDir() {
						err = AddRecursive(watcher, event.Name)
						if err != nil {
							fmt.Println("error adding directory:", err)
						}
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					fmt.Println("Watcher.Errors closed")
					return
				}
				if err != nil {
					fmt.Println("Watcher.Errors : ", err)
				}
			case <-done:
				defer wg.Done()
				w.childProcess.Kill()
				w.shellProcess.Kill()
				time.Sleep(300 * time.Millisecond)
				fmt.Println("Watcher stopped")
				if err := watcher.Close(); err != nil {
					fmt.Println("Failed to stop watcher")
				}
				return
			}
		}
	}()

	err = AddRecursive(watcher, w.directory)
	if err != nil {
		log.Fatal(err)
	}
}

func (w *watcher) runCommand() {
	// helper.ClearScreen()

	atomic.AddInt32(&w.counter, 1)

	a := atomic.LoadInt32(&w.counter)

	fmt.Printf("\n%d %s [Rerun:%s]\n\n", a, w.command, w.reRunDuration)

	w.spider.BroadcastMessage(fmt.Sprintf("ReRun %d", a), socket.Connection{ID: "SPIDER"})

	KillProcess(w.shellProcess)
	KillProcess(w.childProcess)
	PortKiller(w.killPorts)

	stdo := stdOutSave{fn: func(s string) {
		fmt.Println("\n\nout")

		w.spider.BroadcastMessage(fmt.Sprintf("Output %s", s), socket.Connection{ID: "SPIDER"})
	}}
	stde := stdErrSave{fn: func(s string) {
		fmt.Println("\n\nerr")
		w.spider.BroadcastMessage(fmt.Sprintf("Error %s", s), socket.Connection{ID: "SPIDER"})
	}}

	cmd := ExecCommand(w.command, stdo, stde)

	helper.Spinner(time.Millisecond * 400)

	CopyProcess(cmd, &w.shellProcess, &w.childProcess)

	fmt.Printf("...\n\n")
}
