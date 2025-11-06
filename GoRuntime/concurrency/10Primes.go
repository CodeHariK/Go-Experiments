package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func repeatFunc[T any, K any](done <-chan K, fn func() T) <-chan T {
	stream := make(chan T)

	println("repeatFunc stream")

	go func() {
		defer close(stream)

		for {
			select {
			case <-done:
				return
			case stream <- fn():
				println("-r")
			}
		}
	}()

	return stream
}

func take[T any, K any](done <-chan K, stream <-chan T, n int) <-chan T {
	taken := make(chan T)

	println("take stream")

	go func() {
		defer close(taken)

		for i := 0; i < n; i++ {
			select {
			case <-done:
				return
			case taken <- <-stream:
				println("---t")
			}
		}
	}()

	return taken
}

func primeFinder(done <-chan int, randIntStream <-chan int) <-chan int {
	isPrime := func(randomInt int) bool {
		for i := randomInt - 1; i > 1; i-- {
			if randomInt%i == 0 {
				return false
			}
		}
		return true
	}

	primes := make(chan int)
	go func() {
		defer close(primes)
		for {
			select {
			case <-done:
				return
			case randomInt := <-randIntStream:
				println("--p")
				if isPrime(randomInt) {
					primes <- randomInt
				}
			}
		}
	}()

	return primes
}

func fanIn[T any](done <-chan int, channels ...<-chan T) <-chan T {
	var wg sync.WaitGroup
	fannedInStream := make(chan T)

	transfer := func(c <-chan T) {
		defer wg.Done()
		for i := range c {
			select {
			case <-done:
				return
			case fannedInStream <- i:
			}
		}
	}

	for _, c := range channels {
		wg.Add(1)
		go transfer(c)
	}

	go func() {
		wg.Wait()
		close(fannedInStream)
	}()

	return fannedInStream
}

func Primes10() {
	start := time.Now()
	done := make(chan int)

	defer close(done)

	randNumFetcher := func() int { return rand.Intn(500000000) }
	randIntStream := repeatFunc(done, randNumFetcher)

	// for r := range take(done, randIntStream, 10) {
	// 	fmt.Println(r)
	// }

	// for r := range take(done, primeFinder(done, randIntStream), 10) {
	// 	fmt.Println(r)
	// }

	// fan out
	CPUCount := runtime.NumCPU()
	primeFinderChannels := make([]<-chan int, CPUCount)
	for i := 0; i < CPUCount; i++ {
		primeFinderChannels[i] = primeFinder(done, randIntStream)
	}
	// fan in
	fannedInStream := fanIn(done, primeFinderChannels...)
	for ra := range take(done, fannedInStream, 10) {
		println(ra)
	}

	fmt.Println(time.Since(start))
}
