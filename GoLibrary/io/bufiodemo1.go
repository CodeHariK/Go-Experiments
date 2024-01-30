package io

import (
	"bufio"
	"fmt"
)

// 1. The buffer is full
// As soon as the buffer is full, the write operation takes place.

// 2. The buffer has space after write
// If the buffer still has space after the last write, it will not attempt to complete that write until specifically urged to do so by the Flush() method.

// 3. A write larger than buffer capacity is made
// If a write is larger than buffer capacity,​ the buffer is skipped because there is no need to buffer.

// Writer type used to initialize buffer writer
type Writer int

func (*Writer) Write(p []byte) (n int, err error) {
	fmt.Printf("Writing: %s\n", p)
	return len(p), nil
}

func Bufiodemo1() {
	// declare a buffered writer
	// with buffer size 4
	w := new(Writer)
	bw := bufio.NewWriterSize(w, 4)

	// Case 1: Writing to buffer until full
	fmt.Printf("%d\n", bw.Available())
	bw.Write([]byte{'1'})
	fmt.Printf("%d\n", bw.Available())
	bw.Write([]byte{'2'})
	fmt.Printf("%d\n", bw.Available())
	bw.Write([]byte{'3'})
	fmt.Printf("%d\n", bw.Available())
	bw.Write([]byte{'4'}) // write - buffer is full
	fmt.Printf("%d\n", bw.Available())

	// Case 2: Buffer has space
	bw.Write([]byte{'5'})
	fmt.Printf("%d\n", bw.Available())
	err := bw.Flush() // forcefully write remaining
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d\n", bw.Available())

	// Case 3: (too) large write for buffer
	// Will skip buffer and write directly
	bw.Write([]byte("12345"))
	fmt.Printf("%d\n", bw.Available())
}
