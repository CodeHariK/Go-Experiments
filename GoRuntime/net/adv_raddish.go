package net

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type Command struct {
	Fields []string
	Result chan string
}

func raddishServer(commands chan Command) {
	advdata := make(map[string]string)
	for cmd := range commands {
		if len(cmd.Fields) < 2 {
			go func() {
				cmd.Result <- "Expected at least 2 arguments"
			}()
			continue
		}

		fmt.Println("Got Command ", cmd)

		fs := cmd.Fields

		switch fs[0] {

		case "GET":
			key := fs[1]
			value := advdata[key]
			cmd.Result <- fmt.Sprintf("%s\n", value)

		case "SET":
			if len(fs) < 3 {
				cmd.Result <- "Expected Value\n"
				continue
			}
			key := fs[1]
			value := fs[2]
			advdata[key] = value
			cmd.Result <- fmt.Sprintf("%s\n", value)

		case "DEL":
			delete(advdata, fs[1])
			cmd.Result <- ""

		default:
			cmd.Result <- ""
		}
	}
}

func advhandle(commands chan Command, conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		ln := scanner.Text()
		fs := strings.Fields(ln)

		result := make(chan string)
		commands <- Command{
			Fields: fs,
			Result: result,
		}

		fmt.Fprintln(conn, <-result)
	}
}

func AdvRaddish() {
	li, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalln(err)
	}
	defer li.Close()

	commands := make(chan Command)
	go raddishServer(commands)

	for {
		conn, err := li.Accept()
		if err != nil {
			log.Fatalln(err)
		}

		go advhandle(commands, conn)
	}
}
