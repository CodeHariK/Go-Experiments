package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	client *kafka.Reader
	topic  string
}

func NewConsumer(brokers []string, topic string) *Consumer {
	groupID := uuid.New().String()

	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: groupID,
		Topic:   topic,
		Dialer:  dialer,
		// Partition: 2,
	})

	return &Consumer{client: r, topic: topic}
}

func (c *Consumer) PrintMessages() {
	for {
		m, err := c.client.ReadMessage(context.Background())
		if err != nil {
			fmt.Printf(err.Error())
			break
		}

		fmt.Printf("%d(%d): %s = %s\n", m.Offset, m.Partition, string(m.Key), string(m.Value))
	}
}

func (c *Consumer) Close() {
	c.client.Close()
}
