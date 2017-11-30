package bot

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"googlemaps.github.io/maps"
)

// A Bot to take care of responses
type Bot struct {
	GcClient maps.Client
	LatLong  maps.LatLng
	Location string
}

func randomResponse() string {
	randomdResponses := []string{
		"Sorry what",
		"Eh...",
		"I've not a clue to what you mean!",
		"Truly I am sorry but, there is a disconnect between us.",
		"Ahem, excuse me?",
	}
	rand.Seed(time.Now().Unix())
	randResponseIndex := rand.Intn(len(randomdResponses))
	return randomdResponses[randResponseIndex]
}

func parseTheBest(s string) string {
	// parse everythign after "the best" in a string
	parsedString := ""
	re := regexp.MustCompile("(the best).*")
	match := re.FindString(s)
	if len(match) == 0 {
		return parsedString
	}
	parsedString = strings.TrimPrefix(match, "the best")
	return parsedString
}

// Parse string for further processings
func (b *Bot) Parse(s string) string {
	response := randomResponse()
	b.SetLocation(b.Location) // TODO: only do this when asked
	if strings.Contains(s, `where are you`) {
		response = fmt.Sprintf("I am currently in %s", b.Location)
	}
	if strings.Contains(s, "the best") {

	}
	return response
}

// SetLocation for Bot
func (b *Bot) SetLocation(ltlng string) error {
	// TODO: do this by accepting strings like "Austin" or "Dumbo, Brooklyn"
	latlong, e := maps.ParseLatLng(ltlng)
	if e != nil {
		fmt.Println(e)
	}
	b.LatLong = latlong
	return e
}
