package main

import (
	"context"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func Fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Info(err error) {
	if err != nil {
		log.Print(err)
	}
}

func main() {
	// Connect to a server
	nc, err := nats.Connect(
		"connect.ngs.global",
		nats.Name("Order Publisher"),
		nats.UserCredentials("../NGS-Default-CLI.creds"),
	)
	Fatal(err)
	defer nc.Close()

	js, err := jetstream.New(nc)
	Fatal(err)

	ctx := context.Background()
	_, err = js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:        "orders",
		Description: "Messages from orders",
		Subjects: []string{
			"orders.>",
		},
		MaxBytes: 1024 * 1024 * 1024,
	})
	Fatal(err)

	for i := 0; ; i++ {
		_, err := js.Publish(ctx,
			fmt.Sprintf("orders.%d", i),
			[]byte("Hello World"))
		Info(err)
		log.Printf("Published %d", i)
	}
}
