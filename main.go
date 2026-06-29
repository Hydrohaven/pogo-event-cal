package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var lgnds = make(map[string]bool)

func main() {
	// ==== PRE-STUFF ==== //
	// <1> Setup Legendary/Mythical Json
	data, err := os.ReadFile("legendaries.json")
	if err != nil {
		log.Fatal(err)
	}

	var list []string
	err = json.Unmarshal(data, &list)
	if err != nil {
		log.Fatal(err)
	}

	for _, name := range list {
		lgnds[name] = true
	}

	// <2> Fetching Initial Event Links
	// Fetch Events home page
	baseURL := "https://leekduck.com"
	res, err := http.Get(baseURL + "/events/")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	// Load HTMl document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Fetch event link and dates
	doc.Find(".event-header-item-wrapper").Slice(0, 10).Each(func(i int, s *goquery.Selection) {
		eventType := s.Find(".event-item-wrapper").Find("p").First().Text()
		link, _ := s.Find("a").Attr("href")
		fmt.Printf("Item %-4d %-25s %s\n", i, eventType, baseURL+link)

		parseEventData(baseURL + link)
	})
}

func parseEventData(url string) {
	// Fetch event data
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Popular CalenderEvent
	event := CalendarEvent{}
	doc.Find(".page-content").Each(func(i int, s *goquery.Selection) {
		event.title = strings.TrimSpace(s.Find(".page-title").First().Text())
		event.eventType = stringToEventType[s.Find(".tag .raid-day").First().Text()]
		event.link = url

		parsePokemonData(s, &event)
	})

	fmt.Println(event)
	fmt.Println("\n")
}

func parsePokemonData(s *goquery.Selection, e *CalendarEvent) {
	if e.megaList == nil {
		e.megaList = make(map[string]struct{})
	}
	if e.legendaryList == nil {
		e.legendaryList = make(map[string]struct{})
	}
	if e.pokemonList == nil {
		e.pokemonList = make(map[string]struct{})
	}

	s.Find(".pkmn-list-item").Each(func(i int, node *goquery.Selection) {
		pkmn := node.Text()

		// Add to megaList if Mega
		if strings.HasPrefix(pkmn, "Mega") {
			e.megaList[pkmn] = struct{}{}
		}

		// If legendary, add to legendList, else, pokemonList
		if lgnds[pkmn] {
			e.legendaryList[pkmn] = struct{}{}
		} else {
			e.pokemonList[pkmn] = struct{}{}
		}
	})
}
