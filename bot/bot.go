package bot

import (
	"os"
	"strings"
)

// A Bot to take care of responses
type Bot struct {
	location string
}

var googleMapsAPIKey = os.Getenv("GOOGLE_MAPS_API_KEY")

// Parse foodbot message
func (b Bot) Parse(s string) string {
	response := "uh sorry what m8?"
	if strings.Contains(s, `you are in`) {
		q := strings.Split(s, " ")
		loc := q[len(q)-1]
		response = loc
	}
	return response
}
