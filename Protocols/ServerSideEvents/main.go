package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir(".")))

	http.HandleFunc("/random", randomHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func randomHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Type", "text/event-stream")

	// send a random number every 2 seconds
	for {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		fmt.Fprintf(w, "data: %d \n\n", r.Intn(100))
		w.(http.Flusher).Flush()
		time.Sleep(2 * time.Second)
	}
}
