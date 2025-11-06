package main

import (
	"fmt"
	"os"

	"github.com/pkg/profile"
)

func main() {
	defer profile.Start(profile.TraceProfile).Stop()

	const n = 500

	leftmost := make(chan int)
	right := leftmost
	left := leftmost

	for i := 0; i < n; i++ {
		right = make(chan int)
		go pass(left, right)
		left = right
	}

	go sendFirst(right)
	fmt.Fprintln(os.Stderr, <-leftmost)
}

func pass(left, right chan int) {
	v := 1 + <-right
	left <- v
}

func sendFirst(ch chan int) { ch <- 0 }
