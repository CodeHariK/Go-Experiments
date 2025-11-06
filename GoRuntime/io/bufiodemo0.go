package io

import (
	"bufio"
	"fmt"
)

func Bufiodemo0() {
	fmt.Println("Unbuffered I/O")
	w := new(Writer)
	w.Write([]byte{'a'})
	w.Write([]byte{'b'})
	w.Write([]byte{'c'})
	w.Write([]byte{'d'})
	fmt.Println("Buffered I/O")
	bw := bufio.NewWriterSize(w, 3)
	bw.Write([]byte{'a'})
	bw.Write([]byte{'b'})
	bw.Write([]byte{'c'})
	bw.Write([]byte{'d'})
	err := bw.Flush()
	if err != nil {
		panic(err)
	}
}

// Unbuffered I/O
// 1
// 1
// 1
// 1
// Buffered I/O
// 3
// 1
