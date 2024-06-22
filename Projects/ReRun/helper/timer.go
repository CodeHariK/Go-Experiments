package helper

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func TickerFunction(t time.Duration, fn func()) {
	done := make(chan os.Signal)
	signal.Notify(done, syscall.SIGINT, os.Kill, os.Interrupt)

	go func() {
		ticker := time.NewTicker(t)

		fn()
		for {
			select {
			case <-ticker.C:
				fn()
			case <-done:
				fmt.Println("Timer stopped")
				ticker.Stop()
				return
			}
		}
	}()
}
