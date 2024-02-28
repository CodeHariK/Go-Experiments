package main

import (
	"fmt"
	"time"

	"github.com/pkg/profile"
)

type Ball struct{ hits int }

func main() {
	defer profile.Start(profile.TraceProfile).Stop()
	fmt.Println(time.Now().Nanosecond())
	table := make(chan *Ball)

	// go player("ping", table)
	// go player("pong", table)
	for i := 1; i <= 26; i++ {
		go player(fmt.Sprintf("%c", i+64), table)
	}

	table <- new(Ball) // game on; toss the ball
	time.Sleep(1 * time.Second)
	<-table // game over; grab the ball
}

func player(name string, table chan *Ball) {
	for {
		ball := <-table
		ball.hits++

		fmt.Printf("%s  %03d  %x\n", name, ball.hits, time.Now().Nanosecond())

		time.Sleep(10 * time.Millisecond)

		table <- ball
	}
}
