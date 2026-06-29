package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	// Fetch Events home page
	baseURL := "https://leekduck.com"
	res, err := http.Get(baseURL + "/events/")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Fetched HTML")
	}
	defer res.Body.Close()

	// Load HTMl document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Loaded HTML")

	// Fetch event link and dates
	doc.Find(".event-header-item-wrapper").Slice(0, 10).Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		eventType := s.Find(".event-item-wrapper").Find("p").First().Text()
		link, _ := s.Find("a").Attr("href")
		fmt.Printf("Item %-4d %-25s %s\n", i, eventType, baseURL+link)

		parseEventData(baseURL + link)
	})
}

func parseEventData(url string) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find CalenderEvent
	event := CalendarEvent{}
	doc.Find(".page-content").Each(func(i int, s *goquery.Selection) {
		event.title = s.Find(".page-title").Find("p").First().Text()
		event.link = url

		parsePokemonData(s, event)
	})

	fmt.Println(event)
}

func parsePokemonData(s *goquery.Selection, e CalendarEvent) {
	s.Find(".pkmn-list-item").Each(func(i int, node *goquery.Selection) {
		pkmn := node.Text()
		fmt.Println(pkmn)
	})
	fmt.Println("\n")
}
