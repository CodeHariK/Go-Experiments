package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

// Event is the Input Payload representation
type Event struct {
	Name string `json:"name"`
}

// Our Lambda function
// Output Whatever you want, error or str, or struct and error etc.
func lambdaHandler(event Event) (string, error) {
	// Worlds simplest Lambda, Concatinate Event Name with Hello
	output := fmt.Sprintf("Hello From Lambda %s", event.Name)

	return output, nil
}

func main() {
	// Start the lambda handler
	lambda.Start(lambdaHandler)
}
