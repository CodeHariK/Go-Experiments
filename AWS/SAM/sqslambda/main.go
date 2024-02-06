package main

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// our lambda handler that listens on SQS events
func lambdaHandler(event events.SQSEvent) error {
	log.Println(event)

	return nil
}

func main() {
	// Start the lambda handler
	lambda.Start(lambdaHandler)
}
