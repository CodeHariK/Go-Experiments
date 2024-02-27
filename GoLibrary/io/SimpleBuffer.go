package io

import (
	"fmt"
)

// SimpleBuffer is a basic buffer implementation.
type SimpleBuffer struct {
	data     []byte
	readPos  int
	writePos int
}

// NewSimpleBuffer creates a new SimpleBuffer with a given capacity.
func NewSimpleBuffer(capacity int) *SimpleBuffer {
	return &SimpleBuffer{
		data: make([]byte, capacity),
	}
}

// Write appends data to the buffer.
func (b *SimpleBuffer) Write(data []byte) (int, error) {
	if len(data) > len(b.data)-b.writePos {
		// Not enough space in the buffer, you might want to resize it here.
		return 0, fmt.Errorf("buffer overflow")
	}

	copy(b.data[b.writePos:], data)
	b.writePos += len(data)
	return len(data), nil
}

// Read reads data from the buffer.
func (b *SimpleBuffer) Read(dst []byte) (int, error) {
	if len(dst) > len(b.data)-b.readPos {
		// Not enough data in the buffer.
		return 0, fmt.Errorf("buffer underflow")
	}

	copy(dst, b.data[b.readPos:])
	b.readPos += len(dst)
	return len(dst), nil
}

func SimpleBufferMain() {
	buffer := NewSimpleBuffer(1024)

	// Writing to the buffer
	dataToWrite := []byte("Hello, buffer!")
	written, err := buffer.Write(dataToWrite)
	if err != nil {
		fmt.Println("Error writing to buffer:", err)
		return
	}
	fmt.Printf("Written %d bytes to buffer\n", written)

	// Reading from the buffer
	readSize := 8
	readData := make([]byte, readSize)
	read, err := buffer.Read(readData)
	if err != nil {
		fmt.Println("Error reading from buffer:", err)
		return
	}
	fmt.Printf("Read %d bytes from buffer: %s\n", read, readData[:read])
}
