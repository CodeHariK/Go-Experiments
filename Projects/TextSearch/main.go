package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"
)

func main() {
	var dumpPath string
	flag.StringVar(&dumpPath, "p", "enwiki.xml.gz", "wiki abstract dump path")
	flag.Parse()

	log.Println("Running Full Text Search")

	start := time.Now()
	docs, err := LoadDocuments(dumpPath)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Loaded %d documents in %v", len(docs), time.Since(start))

	start = time.Now()
	idx := make(Index)
	idx.Add(docs)

	log.Printf("Indexed %d documents in %v", len(docs), time.Since(start))

	start = time.Now()

	scanner(idx, docs, start)
}

func scanner(idx Index, docs []document, start time.Time) {
	fmt.Println("\n> Type something and press Enter (Ctrl+D to exit):")

	// Create a scanner to read input from the user
	scanner := bufio.NewScanner(os.Stdin)

	// Keep scanning for input until an error or Ctrl+D is pressed
	for scanner.Scan() {
		text := scanner.Text()

		if text == "clear" {
			clearConsole()
			fmt.Println("> Type something and press Enter (Ctrl+D to exit):")
			continue
		}

		fmt.Println("You entered:", text)

		matchedIDs := idx.Search(text)
		log.Printf("\nSearch found %d documents in %v\n", len(matchedIDs), time.Since(start))

		for _, id := range matchedIDs {
			doc := docs[id]
			log.Printf("%d\t%s\n%s\n%s\n\n", id, doc.Title, doc.Text, doc.URL)
		}

		fmt.Println("\n> Type something and press Enter (Ctrl+D to exit):")
	}

	// Check for any errors that may have occurred during scanning
	if err := scanner.Err(); err != nil {
		fmt.Println("Error:", err)
	}
}

func clearConsole() {
	switch runtime.GOOS {
	case "linux", "darwin":
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		fmt.Println("Clearing console is not supported on this platform.")
	}
}
