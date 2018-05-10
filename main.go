package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
)

var (
	// ErrNameNotProvided is thrown when a name is not provided
	ErrNameNotProvided = errors.New("no name was provided in the HTTP body")
)

// ChallengeHandler Handle Slack challenge request.
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

	// TODO: check if type = url_verification
	if strings.Contains(request.Body, "challenge") {
		return ChallengeHandler(request)
	}

	if len(request.Body) < 1 {
		return events.APIGatewayProxyResponse{}, ErrNameNotProvided
	}

	// Use Slack events API to process response
	slackToken := os.Getenv("SLACK_API_KEY")
	api := slack.New(slackToken)
	params := slack.PostMessageParameters{}
	eventsAPIEvent, e := slackevents.ParseEvent(json.RawMessage(request.Body))
	if e != nil {
		log.Printf("Error!, %s", e)
	}
	innerEvent := eventsAPIEvent.InnerEvent
	switch ev := innerEvent.Data.(type) {
	// This is an EventsAPI specific event
	case *slackevents.AppMentionEvent:
		api.PostMessage(ev.Channel, "hi", params)
		break
	// There are Events API specific MessageEvents
	// https://api.slack.com/events/message.channels
	case *slackevents.MessageEvent:
		if strings.Contains(ev.Text, "pizza") {
			api.PostMessage(ev.Channel, "I like pizza too.", params)
		}
		break
	// This is an Event shared between RTM and the EventsAPI
	case *slack.ChannelCreatedEvent:
		api.PostMessage(ev.Channel.ID, "Oh hey, a new channel!", params)
		break
	default:
		fmt.Println(innerEvent.Type)
		fmt.Println("no event to match")
	}
	return events.APIGatewayProxyResponse{
		Body:       "Event:" + string(innerEvent.Type) + ", recieved!",
		StatusCode: 200,
	}, nil
}

func main() {
	// blocking call
	lambda.Start(Handler)
}
