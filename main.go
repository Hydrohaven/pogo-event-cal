package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"google.golang.org/api/calendar/v3"
)

var legendaries = make(map[string]bool)
var eventCache = make(map[string]CalendarEvent)

func main() {
	//  ==== <PRE-STUFF> ====
	// Setup env variables
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "service-worker-auth.json")

	// Setup Legendary/Mythical Json
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
		legendaries[name] = true
	}

	// Load event cache
	if data, err = os.ReadFile("cache.json"); err == nil {
		_ = json.Unmarshal(data, &eventCache)
	}

	// ==== <1> HTTP Get ====
	// Fetching Initial Event Links
	// Fetch Events home page
	baseURL := "https://leekduck.com"
	res, err := http.Get(baseURL + "/events/")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	// ==== <2> Traverse DOM ====
	// Load HTMl document with Goquery
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Fetch event link and dates
	doc.Find(".event-header-item-wrapper").Slice(0, 5).Each(func(i int, s *goquery.Selection) {
		eventType := s.Find(".event-item-wrapper").Find("p").First().Text()
		slug, _ := s.Find("a").Attr("href")
		startDate, _ := s.Attr("data-event-start-date-check")
		endDate, _ := s.Attr("data-event-end-date")
		fmt.Printf("Item %-4d %-25s %s\n", i, eventType, baseURL+slug)

		// ==== <2.5> Skip events I don't like ====
		cleanType := strings.TrimSpace(eventType)
		if _, ok := stringToEventType[cleanType]; !ok {
			fmt.Printf("@ Skipped %v\n", cleanType)
			return
		}

		// ==== <3> Compare Curr to Cache ====
		if _, ok := eventCache[slug]; ok {
			fmt.Printf("Skipped %v\n", slug)
			return
		}

		parseEventData(baseURL, slug, startDate, endDate)

		// ==== <4> GCal POST Request ====
		postCalendarEvent(eventCache[slug])
	})

	// ==== <5> Encode new data to cache ====
	cacheData, _ := json.MarshalIndent(eventCache, "", "\t")
	_ = os.WriteFile("cache.json", cacheData, 0644) // 0644 gives rw perms
}

func parseEventData(baseURL string, slug string, startDate string, endDate string) {
	// Fetch event data
	res, err := http.Get(baseURL + slug)
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
		event.Title = strings.TrimSpace(s.Find(".page-title").First().Text())
		event.Description = parseDescription(s.Find(".event-description"))
		if t, err := time.Parse(time.RFC3339, startDate); err == nil {
			event.StartDate = t.Local()
		} else if t, err := time.ParseInLocation("2006-01-02T15:04:05", startDate, time.Local); err == nil {
			event.StartDate = t
		}
		if t, err := time.Parse(time.RFC3339, endDate); err == nil {
			event.EndDate = t.Local()
		} else if t, err := time.ParseInLocation("2006-01-02T15:04:05", endDate, time.Local); err == nil {
			event.EndDate = t
		}
		event.EventType = stringToEventType[s.Find(".page-tags .tag").First().Text()]
		event.Link = baseURL + slug

		parsePokemonData(s, &event)
	})

	eventCache[slug] = event
	fmt.Println(event)
	fmt.Print("\n\n")
}

func parsePokemonData(s *goquery.Selection, e *CalendarEvent) {
	if e.MegaList == nil {
		e.MegaList = make(map[string]struct{})
	}
	if e.LegendaryList == nil {
		e.LegendaryList = make(map[string]struct{})
	}
	if e.PokemonList == nil {
		e.PokemonList = make(map[string]struct{})
	}

	s.Find(".pkmn-list-item").Each(func(i int, node *goquery.Selection) {
		pkmn := node.Text()
		if ignoreGenesect(pkmn) {
			return
		}

		// Add to megaList if Mega
		if strings.HasPrefix(pkmn, "Mega") {
			e.MegaList[pkmn] = struct{}{}
			return
		}

		// If legendary, add to legendList, else, pokemonList
		if legendaries[pkmn] {
			e.LegendaryList[pkmn] = struct{}{}
		} else {
			e.PokemonList[pkmn] = struct{}{}
		}
	})
}

// Ignores stupid genesect drives cuz i DONT like genesect
func ignoreGenesect(s string) bool {
	if s == "Genesect" {
		return false
	}

	if strings.HasPrefix(s, "Genesect") {
		return true
	}

	return false
}

func postCalendarEvent(e CalendarEvent) {
	ctx := context.Background() // a service that manages the lifecycle, signals, and metadata of operations around api's and goroutines
	srv, err := calendar.NewService(ctx)
	if err != nil {
		log.Fatal(err)
	}

	var config map[string]string
	if data, err := os.ReadFile("config.json"); err == nil {
		_ = json.Unmarshal(data, &config)
	}
	calendarID := config["calendar_id"]

	gcalEvent := calendar.Event{
		Summary:     e.Title,
		Location:    e.Link,
		Description: e.Description,
		Start: &calendar.EventDateTime{
			DateTime: e.StartDate.Format(time.RFC3339)},
		End: &calendar.EventDateTime{
			DateTime: e.EndDate.Format(time.RFC3339),
		},
		ColorId: e.EventType.ColorId(),
	}

	event, err := srv.Events.Insert(calendarID, &gcalEvent).Do()
	if err != nil {
		log.Fatalf("Unable to create event. %v\n", err)
	}
	fmt.Printf("Event created: %s\n", event.HtmlLink)
}

func parseDescription(desc *goquery.Selection) string {
	var lines []string
	desc.Children().Each(func(i int, child *goquery.Selection) {
		if child.Is("ul") {
			child.Find("li").Each(func(j int, li *goquery.Selection) {
				lines = append(lines, "- "+strings.TrimSpace(li.Text()))
			})
		} else {
			text := strings.TrimSpace(child.Text())
			if text != "" {
				lines = append(lines, text)
			}
		}
	})
	return strings.Join(lines, "\n\n")
}
