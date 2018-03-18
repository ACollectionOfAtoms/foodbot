package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestChallengeHandler(t *testing.T) {
	req := events.APIGatewayProxyRequest{
		Body: `{
				"challenge": "lol"
			}`,
	}
	res, err := ChallengeHandler(req)
	if err != nil {
		t.Fail()
	}
	challenge := res.Body
	if challenge != "lol" {
		t.Fail()
	}
}
