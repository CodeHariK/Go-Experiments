package net

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

var data = make(map[string]string)

func handle(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		ln := scanner.Text()

		fs := strings.Fields(ln)
		if len(fs) < 2 {
			continue
		}

		switch fs[0] {
		case "GET":
			key := fs[1]
			value := data[key]
			fmt.Fprintf(conn, "%s\n", value)
		case "SET":
			if len(fs) < 3 {
				conn.Write([]byte("Expected Value\n"))
				continue
			}
			key := fs[1]
			value := fs[2]
			data[key] = value
		case "DEL":
			delete(data, fs[1])
		default:
			io.WriteString(conn, "Invalid Command "+fs[0]+"\n")
		}

	}
}

func SimpleRaddish() {
	li, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalln(err)
	}
	defer li.Close()

	for {
		conn, err := li.Accept()
		if err != nil {
			log.Fatalln(err)
		}

		//------------------
		// handle(conn) // Handles only one connection
		go handle(conn) // Handles many connection but with race condition
		//------------------
	}
}
