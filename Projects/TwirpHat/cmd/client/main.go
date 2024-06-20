package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"twirphat/rpc/haberdasher"

	"github.com/twitchtv/twirp"
)

var (
	hat *haberdasher.Hat
	err error
)

func main() {
	jsonclient := haberdasher.NewHaberdasherJSONClient(
		"http://localhost:8080",
		&http.Client{},
		twirp.WithClientPathPrefix("/dash"),
	)

	protoclient := haberdasher.NewHaberdasherProtobufClient(
		"http://localhost:8080",
		http.DefaultClient,
		twirp.WithClientPathPrefix("/dash"),
	)

	// data := map[string]int{"inches": 12}
	// jsonData, _ := json.Marshal(data)
	jsonData, _ := json.Marshal(haberdasher.Size{Inches: 7})
	url := "http://localhost:8080/dash/twirphat.haberdasher.Haberdasher/MakeHat"
	resp, _ := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &hat)
	fmt.Println(string(body))
	fmt.Println(hat)
	resp.Body.Close()

	clients := []haberdasher.Haberdasher{jsonclient, protoclient}

	for i := 0; i < 5; i++ {
		hat, err = clients[i%len(clients)].MakeHat(context.Background(), &haberdasher.Size{Inches: int32(i*3 + 1)})
		if err != nil {
			if twerr, ok := err.(twirp.Error); ok {
				if twerr.Meta("retryable") != "" {
					// Log the error and go again.
					log.Printf("got error %q, retrying", twerr)
					continue
				}
			}
		}
		fmt.Printf("%+v\n", hat)
	}
}
