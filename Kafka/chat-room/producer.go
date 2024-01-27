package main

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	client *kafka.Writer
	topic  string
}

func NewProducer(brokers []string, topic string) *Producer {
	client := &kafka.Writer{
		Addr:        kafka.TCP(brokers...),
		Topic:       topic,
		Balancer:    &kafka.RoundRobin{},
		Logger:      kafka.LoggerFunc(logf),
		ErrorLogger: kafka.LoggerFunc(logf),
	}

	return &Producer{client: client, topic: topic}
}

func (p *Producer) SendMessage(user, message string) {
	msg := kafka.Message{
		Key:   []byte(user),
		Value: []byte(message),
		// Partition: 2,
	}

	err := p.client.WriteMessages(context.Background(),
		msg,
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}
}

func (p *Producer) Close() {
	p.client.Close()
}

func logf(msg string, a ...interface{}) {
	// fmt.Printf(msg, a...)
	fmt.Println()
}
