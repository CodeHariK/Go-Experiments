package io

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func HelloWorldIO() {
	file, err := os.Create("example.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	bufferedWriter := bufio.NewWriterSize(file, 10)

	data := []byte("Hello, World!\n")
	fmt.Println(len(data))
	for i := 0; i < 5; i++ {
		bufferedWriter.Write(data)

		fileInfo, _ := file.Stat()
		content, _ := os.ReadFile("example.txt")
		fmt.Println(string(content))
		fmt.Printf("%x %d\n", time.Now().Second(), fileInfo.Size())

		time.Sleep(time.Second * 1)
	}

	fmt.Printf("Data written to file.")

	time.Sleep(time.Second * 10)
	os.Remove("example.txt")
}
