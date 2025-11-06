package main

import (
	"context"
	"log"
	"os"
	"os/signal"

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
		nats.UserCredentials("../../NGS-Default-CLI.creds"),
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

	stream, err := js.Stream(ctx, "orders")
	Fatal(err)
	consumer, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Name:    "orders_processor",
		Durable: "orders_processor",
	})
	Fatal(err)

	cctx, err := consumer.Consume(func(msg jetstream.Msg) {
		log.Printf("Recieved " + msg.Subject())
		msg.Ack()
	})
	Fatal(err)
	defer cctx.Stop()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
