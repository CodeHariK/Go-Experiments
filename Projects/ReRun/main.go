package main

import (
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/codeharik/rerun/helper"
	"github.com/codeharik/rerun/spider"
	"github.com/codeharik/rerun/types"
	"github.com/codeharik/rerun/watcher"
)

func main() {
	flagKillPorts := flag.String("k", "", "Optional Kill Ports")
	flagReRunDelay := flag.Int("t", -1, "Optional Rerun Delay Time in Milliseconds[Min 100]")

	flag.Parse()

	nonFlagArgs := flag.Args()

	if len(nonFlagArgs) < 2 {
		fmt.Println("ReRun: Monitor a directory and automatically execute a command when files change, or rerun the command on a set interval.")
		flag.PrintDefaults()
		fmt.Println()
		fmt.Println("Usage: go run main.go [-k optional kill ports] [-t optional rerun delay time] <watch directory> <run command>")
		fmt.Println("Usage: go run main.go ../Hello \"go run ../Hello/main.go\"")
		fmt.Println("Usage: go run main.go -k=8080,3000 -t=4000 ../Hello \"go run ../Hello/main.go\"")
		fmt.Println()
		fmt.Println("Usage: rerun [-k optional kill ports] [-t optional rerun delay time] <watch directory> <run command>")
		fmt.Println("Usage: rerun -k=8080,3000 -t=4000 ../Hello \"go run ../Hello/main.go\"")
		fmt.Println("Usage: rerun ../Hello \"go run ../Hello/main.go\"")
		return
	}

	killPortsString := *flagKillPorts
	rerunTimer := time.Duration(*flagReRunDelay) * 1000000
	if rerunTimer > 0 && rerunTimer < time.Millisecond*100 {
		log.Fatal("Min 100 milliseconds delay required")
	}

	directory := flag.Arg(0)
	command := flag.Arg(1)
	killPorts := []int{}
	if killPortsString != "" {
		k, err := helper.ParseStringInts(killPortsString)
		if err == nil {
			killPorts = k
		}
	}

	rerun := types.ReRun{}

	var wg sync.WaitGroup
	defer wg.Wait()

	stdLogs := make(map[string][]string)

	spider := spider.NewSpider(directory, stdLogs, stdLogs)
	spider.StartSpider(&wg)

	w := watcher.NewWatcher(
		&rerun,
		command, rerunTimer, killPorts, directory,
		spider,
		stdLogs,
	)
	w.StartWatcher()
}
