package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

// ProgressReader is a reader that tracks the progress of the read operation
type ProgressReader struct {
	reader io.Reader
	total  int64
	read   int64
}

func (pr *ProgressReader) Read(p []byte) (n int, err error) {
	n, err = pr.reader.Read(p)
	pr.read += int64(n)
	pr.printProgress()
	return n, err
}

func (pr *ProgressReader) printProgress() {
	if pr.total > 0 {
		percentage := float64(pr.read) / float64(pr.total) * 100
		fmt.Printf("\r%.2f%% complete", percentage)
	}
}

func downloadFile(url, filePath string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP status code: %d", response.StatusCode)
	}

	// Create the output file
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the total file size from the Content-Length header
	fileSize, err := strconv.Atoi(response.Header.Get("Content-Length"))
	if err != nil {
		// Content-Length header not present or invalid, unable to show progress
		fmt.Println("Unable to determine file size. Downloading without progress.")
		_, err = io.Copy(out, response.Body)
		return err
	}

	fmt.Printf("\nDownloading %s...\n", filePath)

	// Create a proxy reader to track the progress
	progressReader := &ProgressReader{reader: response.Body, total: int64(fileSize)}

	// Copy the contents while tracking progress
	_, err = io.Copy(out, progressReader)
	if err != nil {
		return err
	}

	fmt.Printf("\nFile downloaded successfully: %s\n", filePath)
	return nil
}

func DownloadFile(path string) error {
	url := "https://dumps.wikimedia.org/enwiki/latest/enwiki-latest-abstract10.xml.gz"

	err := downloadFile(url, path)
	if err != nil {
		fmt.Println("Error:", err)
	}
	return err
}
