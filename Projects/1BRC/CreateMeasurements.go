package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"

	data "1brc/data"
)

func create_measurements(size int) {
	file, err := os.Create("measurements.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		os.Exit(1)
	}
	defer file.Close()

	// writer := bufio.NewWriter(file)
	writer := bufio.NewWriterSize(file, 32*1024)

	for i := 0; i < size; i++ {
		station := data.Stations[rand.Intn(len(data.Stations))]
		fmt.Fprintf(writer, "%s;%.1f\n", station.ID, station.Measurement())
	}

	writer.Flush()
}
