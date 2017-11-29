package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ACollectionOfAtoms/foodbot/bot"
	"github.com/nlopes/slack"
)

var slackAPIKey = os.Getenv("SLACK_API_KEY")

func main() {
	var BotID string
	// the string here is a channel ID :D
	bots := make(map[string]bot.Bot)

	api := slack.New(slackAPIKey)
	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)
	api.SetDebug(true)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		fmt.Print("Event Received: ")
		switch ev := msg.Data.(type) {
		case *slack.ConnectedEvent:
			BotID = ev.Info.User.ID

		case *slack.ChannelJoinedEvent:
			channelID := ev.Channel.ID
			channelName := ev.Channel.Name
			message := fmt.Sprintf("Yes hello, #%s. Eh, where am I?", channelName)
			bots[channelID] = bot.Bot{}
			rtm.SendMessage(rtm.NewOutgoingMessage(message, channelID))

		case *slack.MessageEvent:
			fmt.Println(bots)
			channel := ev.Msg.Channel
			if _, in := bots[channel]; !in {
				bots[channel] = bot.Bot{}
			}
			b := bots[channel]
			incomingMessage := ev.Msg.Text
			wasAddressed := strings.HasPrefix(incomingMessage, "<@"+BotID+">")
			if wasAddressed {
				message := b.Parse(incomingMessage)
				go rtm.SendMessage(rtm.NewOutgoingMessage(message, channel))
			}
		case *slack.PresenceChangeEvent:
			fmt.Printf("Presence Change: %v\n", ev)

		case *slack.LatencyReport:
			fmt.Printf("Current latency: %v\n", ev.Value)

		case *slack.RTMError:
			fmt.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			return

		default:
			// Ignore other events..
			// fmt.Printf("Unexpected: %v\n", msg.Data)
		}
	}
}
