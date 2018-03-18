package main

import (
	"encoding/json"
	"errors"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	// ErrNameNotProvided is thrown when a name is not provided
	ErrNameNotProvided = errors.New("no name was provided in the HTTP body")
)

// Handle Slack challenge request.
func ChallengeHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	body := request.Body
	type ChallengeResponse struct {
		Challenge string
	}
	var r ChallengeResponse
	err := json.Unmarshal([]byte(body), &r)
	if err != nil {
		s := "Unable to parse challenge JSON!"
		return events.APIGatewayProxyResponse{
			Body:       s,
			StatusCode: 500,
		}, err
	}
	return events.APIGatewayProxyResponse{
		Body:       r.Challenge,
		StatusCode: 200,
	}, nil
}

// Handler handles Lambda functions calls
// We use aws API gateway request/responses from the
// events package.
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// stdout and stderr go to Cloudwatch!
	log.Printf("Processing Lambda request %s\n", request.RequestContext.RequestID)

	if strings.Contains(request.Body, "challenge") {
		return ChallengeHandler(request)
	}

	if len(request.Body) < 1 {
		return events.APIGatewayProxyResponse{}, ErrNameNotProvided
	}

	return events.APIGatewayProxyResponse{
		Body:       "Hello " + request.Body,
		StatusCode: 200,
	}, nil

}

func main() {
	// blocking call
	lambda.Start(Handler)
}
