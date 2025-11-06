package io

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func funcToWithIO() {
	defer timeTrack(time.Now(), "funcToWithIO")
	file, err := os.Open("file.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()

	data := make([]byte, 100)
	for {
		_, err := file.Read(data)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func funcToWithBufio() {
	defer timeTrack(time.Now(), "funcToWithBufio")
	file, err := os.Open("file.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	data := make([]byte, 100)
	for {
		_, err := reader.Read(data)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func createFile() {
	defer timeTrack(time.Now(), "createFile")
	file, err := os.Create("file.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	for i := 0; i < 1000000; i++ {
		file.Write([]byte("Hello World!"))
	}
}

func Bufiodemo3() {
	createFile()
	funcToWithIO()
	funcToWithBufio()
	os.Remove("file.txt")
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
