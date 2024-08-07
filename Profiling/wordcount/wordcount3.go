package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"unicode"

	"github.com/pkg/profile"
)

var buf [1]byte

func readbyte3(r io.Reader) (rune, error) {
	_, err := r.Read(buf[:])
	return rune(buf[0]), err
}

func wordcount3() {
	// defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	defer profile.Start(profile.MemProfile, profile.MemProfileRate(1), profile.ProfilePath(".")).Stop()

	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("could not open file %q: %v", os.Args[1], err)
	}

	words := 0
	inword := false
	b := bufio.NewReader(f)
	for {
		r, err := readbyte3(b)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("could not read file %q: %v", os.Args[1], err)
		}
		if unicode.IsSpace(r) && inword {
			words++
			inword = false
		}
		inword = unicode.IsLetter(r)
	}
	fmt.Printf("%q: %d words\n", os.Args[1], words)
}
