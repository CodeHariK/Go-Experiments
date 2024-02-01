package io

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// Rot13Reader is a custom reader that performs Rot13 encryption on the input stream.
type Rot13Reader struct {
	r io.Reader
}

func (rr *Rot13Reader) Read(p []byte) (n int, err error) {
	n, err = rr.r.Read(p)
	for i := 0; i < n; i++ {
		c := p[i]

		switch {
		case 'A' <= c && c <= 'Z':
			p[i] = (c-'A'+13)%26 + 'A'
		case 'a' <= c && c <= 'z':
			p[i] = (c-'a'+13)%26 + 'a'
		}
	}

	return n, err
}

func Rot13Demo() {
	// Create a Rot13Reader wrapping stdin
	r13 := &Rot13Reader{r: os.Stdin}

	// Wrap the Rot13Reader with a bufio.Scanner
	scanner := bufio.NewScanner(r13)

	// Scan lines and print the Rot13 encrypted lines
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
	}
}
