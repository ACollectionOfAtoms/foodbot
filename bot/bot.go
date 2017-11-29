package bot

import (
	"fmt"
	"strings"

	"googlemaps.github.io/maps"
)

// A Bot to take care of responses
type Bot struct {
	GcClient *maps.Client
	Location string
}

// Parse string for further processings
func (b *Bot) Parse(s string) string {
	response := "uh sorry what m8?"
	if strings.Contains(s, `you are in`) {
		q := strings.Split(s, " ")
		loc := q[len(q)-1]
		err := b.SetLocation(loc)
		if err != nil {
			response = "Sorry, I don't know where that is!"
		} else {
			response = fmt.Sprintf("Ok, I'm in %s", b.Location)
		}
	}
	if strings.Contains(s, `where are you`) {
		response = fmt.Sprintf("I am currently in %s", b.Location)
	}
	return response
}

// Set Location for Bot
func (b *Bot) SetLocation(q string) error {
	var e error
	b.Location = q
	return e
}
