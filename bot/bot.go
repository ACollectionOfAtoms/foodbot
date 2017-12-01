package bot

import (
	"context"
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
		"Excellent statement yet, I've no idea what you mean!",
		"Alas, your words fall on deaf ears.",
		"Sorry what",
		"Eh...",
		"I've not a clue to what you mean!",
		"Truly I am sorry but, there is a disconnect between us.",
		"Ahem, excuse me?",
		"'Tis but words that tears us apart. Please rephase.",
	}
	rand.Seed(time.Now().Unix())
	randResponseIndex := rand.Intn(len(randomdResponses))
	return randomdResponses[randResponseIndex]
}

func parseEverythingAfter(s, prefix string) string {
	// parse everything after "s" in a string
	parsedString := ""
	reString := fmt.Sprintf("(%s).*", prefix)
	re := regexp.MustCompile(reString)
	match := re.FindString(s)
	if len(match) == 0 {
		return parsedString
	}
	parsedString = strings.TrimPrefix(match, prefix)
	return parsedString
}

func parseTheBest(s string) string {
	return parseEverythingAfter(s, "the best")
}

func parseNearest(s string) string {
	return parseEverythingAfter(s, "nearest")
}

func queryGoogleMaps(query string, b *Bot, ranking maps.RankBy) string {
	var name string
	var req *maps.NearbySearchRequest
	var minPrice maps.PriceLevel = "0"
	var maxPrice maps.PriceLevel = "4"
	if ranking == maps.RankByDistance {
		req = &maps.NearbySearchRequest{
			Location: &b.LatLong,
			Keyword:  query,
			MinPrice: minPrice,
			MaxPrice: maxPrice,
			RankBy:   ranking,
			OpenNow:  true,
			Type:     "restaurant",
		}
	} else if ranking == maps.RankByProminence {
		req = &maps.NearbySearchRequest{
			Location: &b.LatLong,
			Radius:   uint(200),
			Keyword:  query,
			MinPrice: minPrice,
			MaxPrice: maxPrice,
			RankBy:   ranking,
			OpenNow:  true,
			Type:     "restaurant",
		}
	} else {
		name = ""
	}
	res, e := b.GcClient.NearbySearch(context.Background(), req)
	if e != nil || len(res.Results) == 0 {
		fmt.Println(e)
		fmt.Println("No results!")
		name = ""
	} else {
		// TODO: provide address.
		fmt.Println(res.Results)
		name = res.Results[0].Name
	}
	return name
}

func randomFoodType() string {
	var foodTypes = []string{
		"burger",
		"sushi",
		"seafood",
		"comfort food",
		"chicken",
		"pizza",
		"italian",
		"healthy",
		"mediterranean",
		"african",
		"mexian",
		"chinese",
		"american",
		"thai",
		"greek",
		"indian",
		"latin",
		"cajun",
		"vietnamese",
	}
	rand.Seed(time.Now().Unix())
	randFoodIndex := rand.Intn(len(foodTypes))
	return foodTypes[randFoodIndex]
}

// Parse string for further processings
func (b *Bot) Parse(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	response := randomResponse()
	b.SetLocation(b.Location) // TODO: only do this when asked
	placeName := ""

	if strings.Contains(s, "where are you") {
		response = fmt.Sprintf("I am currently at %s", b.Location)
	}
	if strings.Contains(s, "the best") {
		query := parseTheBest(s)
		placeName = queryGoogleMaps(query, b, maps.RankByProminence)
		response = fmt.Sprintf("I'd say the best %s is at %s", query, placeName)
	}
	if strings.Contains(s, "nearest") {
		query := parseNearest(s)
		placeName = queryGoogleMaps(query, b, maps.RankByDistance)
		response = fmt.Sprintf("The nearest %s is probably at %s", query, placeName)
	}
	if strings.Contains(s, "where should i eat") {
		query := randomFoodType()
		placeName = queryGoogleMaps(query, b, maps.RankByDistance)
		attempts := 0
		for placeName == "" {
			if attempts == 3 {
				break
			}
			query = randomFoodType()
			placeName = queryGoogleMaps(query, b, maps.RankByDistance)
			attempts++
		}
		response = fmt.Sprintf("Why not try %s?", placeName)
	}
	if placeName == "" {
		response = "Sorry, I couldn't find anything!"
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
