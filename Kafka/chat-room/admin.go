package main

import (
	"context"
	"net"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
	topics "github.com/segmentio/kafka-go/topics"
)

type Admin struct {
	client *kafka.Client
	conn   *kafka.Conn
}

func NewAdmin(brokers []string) *Admin {
	client := &kafka.Client{
		Addr:    kafka.TCP(brokers...),
		Timeout: 10 * time.Second,
	}

	//--------------
	conn, err := kafka.Dial("tcp", brokers[0])
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		panic(err.Error())
	}
	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		panic(err.Error())
	}

	return &Admin{
		client: client,
		conn:   controllerConn,
	}
}

func (a *Admin) TopicExists(topic string) bool {
	ctx := context.Background()
	allTopics, err := topics.List(ctx, a.client)
	if err != nil {
		panic(err)
	}
	for _, metadata := range allTopics {
		if metadata.Name == topic {
			return true
		}
	}
	return false
}

func (a *Admin) CreateTopic(topic string) {
	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     3,
			ReplicationFactor: 1,
		},
	}

	err := a.conn.CreateTopics(topicConfigs...)
	if err != nil {
		panic(err.Error())
	}
}

func (a *Admin) Close() {
	a.conn.Close()
}
