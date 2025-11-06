package net

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
)

var Buffered = true

type SFileServer struct{}

func (fs *SFileServer) start() {
	ln, _ := net.Listen("tcp", ":3000")

	for {
		conn, _ := ln.Accept()

		fmt.Printf("Accept %s\n", conn.LocalAddr().String())

		if Buffered {
			go fs.readLoopBuffered(conn)
		} else {
			go fs.readLoopUnbuffered(conn)
		}
	}
}

func (fn *SFileServer) readLoopUnbuffered(conn net.Conn) {
	buf := make([]byte, 2048)
	for {
		n, _ := conn.Read(buf)
		file := buf[:n]
		fmt.Println(file)
		fmt.Printf("\nrec %d bytes over the net\n", n)
	}
}

func (fn *SFileServer) readLoopBuffered(conn net.Conn) {
	buf := new(bytes.Buffer)
	for {
		var size int64
		binary.Read(conn, binary.LittleEndian, &size)
		n, _ := io.CopyN(buf, conn, size)
		fmt.Println(buf.Bytes())
		fmt.Printf("\nrec %d bytes over the net\n", n)
	}
}

func sendFileUnbuffered(size int) {
	file := make([]byte, size)

	io.ReadFull(rand.Reader, file)

	conn, _ := net.Dial("tcp", ":3000")

	n, _ := conn.Write(file)
	fmt.Printf("\nwritten %d bytes over the net\n", n)
}

func sendFileBuffered(size int) {
	file := make([]byte, size)

	io.ReadFull(rand.Reader, file)

	conn, _ := net.Dial("tcp", ":3000")

	// n, _ := io.Copy(conn, bytes.NewReader(file))
	binary.Write(conn, binary.LittleEndian, int64(size))
	n, _ := io.CopyN(conn, bytes.NewReader(file), int64(size))
	fmt.Printf("written %d bytes over the net\n", n)
}

func LargeFileStreaming() {
	go func() {
		time.Sleep(4 * time.Second)

		if Buffered {
			sendFileBuffered(400000)
		} else {
			sendFileUnbuffered(4000)
		}
	}()

	server := &SFileServer{}
	server.start()
}
