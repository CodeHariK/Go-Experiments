package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Event is the Input Payload representation
type Event struct {
	Name string `json:"name"`
}

// Our Lambda function now Accepts a APIGatewayProxyRequest and outputs a Response instead
func lambdaHandler(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// This is our orignial event
	var input Event
	// We can find the users Payload inside event.Body as a String, Marshal it into our wanted Event format.
	if err := json.Unmarshal([]byte(event.Body), &input); err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	output := fmt.Sprintf("Hello From Lambda %s, the variable is %s", input.Name, os.Getenv("my-cool-variable"))

	return events.APIGatewayProxyResponse{
		Body:       output,
		StatusCode: 200,
	}, nil
}

func main() {
	// Start the lambda handler
	lambda.Start(lambdaHandler)
}
