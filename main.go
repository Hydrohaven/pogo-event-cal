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

var legendaries = make(map[string]bool)         // Map of all legendaries, mythicals, and ultra beasts
var tempCache = make(map[string]bool)           // Temporary cache deleted at program close
var eventCache = make(map[string]CalendarEvent) // Long term cache to keep track of all logged events over 1 month
var ctx context.Context
var srv *calendar.Service
var calendarID string
var isCalSynced bool = false

func main() {
	//  ==== <PRE-STUFF> ====
	// Setup env variables
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "service-worker-auth.json")

	// Start CLI
	startCLI()
}

func syncCalendar() {
	var err any
	ctx := context.Background() // a service that manages the lifecycle, signals, and metadata of operations around api's and goroutines
	srv, err = calendar.NewService(ctx)
	if err != nil {
		log.Fatal(err)
	}

	var config map[string]string
	if data, err := os.ReadFile("config.json"); err == nil {
		_ = json.Unmarshal(data, &config)
	}
	calendarID = config["calendar_id"]
	fmt.Println("<> Calendar is synced <>")
}

func startSync() {
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
	doc.Find(".event-header-item-wrapper").Slice(0, goquery.ToEnd).Each(func(i int, s *goquery.Selection) {
		eventType := s.Find(".event-item-wrapper").Find("p").First().Text()
		slug, _ := s.Find("a").Attr("href")
		startDate, _ := s.Attr("data-event-start-date-check")
		endDate, _ := s.Attr("data-event-end-date")
		fmt.Printf("Item %-4d %-25s %s\n", i, eventType, baseURL+slug)

		// ==== <2.5> Skip events that meet certain conditions ====
		// Condition 1: Events I don't like, aka not in the stringToEventType array
		cleanType := strings.TrimSpace(eventType)
		if _, ok := stringToEventType[cleanType]; !ok {
			fmt.Printf("<SKIP> I don't like this event type, %v <SKIP>\n", cleanType)
			return
		}

		// Condition 2: Events more than 3 weeks in the future
		futureTime := time.Now().AddDate(0, 0, 21)

		if localize(startDate).Compare(futureTime) != -1 {
			fmt.Print("<SKIP> Event too far into the future <SKIP>\n")
			return
		}

		// Condition 3: Eveents more than 2 months old
		//  p.s. actually doing this for events on 01/01/01 for some reason, LEEKDUCK GLITCH
		pastTime := time.Now().AddDate(0, 0, -21)

		if localize(startDate).Compare(pastTime) != 1 {
			fmt.Print("<SKIP> Event too far into the past <SKIP>\n")
			return
		}

		// ==== <3> Skip events already logged in cache by slug ====
		// Note: This will have to be refactored in the future in the case we
		//       log and event that we want to update an incomplete event that
		//       we logged. For now, I'm just only adding events no more than
		//       3 weeeks in the future though!

		// Short-term Cache: Made because sometimes I ran into weird duplicate ghost events
		if _, ok := tempCache[slug]; ok {
			fmt.Printf("<SKIP> %v is already logged <SKIP>\n", slug)
			return
		}

		// Long-term Cache
		if _, ok := eventCache[slug]; ok {
			fmt.Printf("<SKIP> %v is already logged <SKIP>\n", slug)
			return
		}

		parseEventData(baseURL, slug, startDate, endDate)

		// ==== <4> GCal POST Request ====
		postEvent(eventCache[slug])
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
		event.Title = strings.ReplaceAll(event.Title, "\u00a0", " ")
		event.Description = parseDescription(s.Find(".event-description"))
		event.StartDate = localize(startDate)
		event.EndDate = localize(endDate)
		event.EventType = stringToEventType[s.Find(".page-tags .tag").First().Text()]
		event.Link = baseURL + slug

		parsePokemonData(s, &event)
	})

	eventCache[slug] = event
	tempCache[slug] = true

	// Print Parsed Event
	// fmt.Println(event)
	// fmt.Print("\n\n")
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

func postEvent(e CalendarEvent) {
	// If first time using calendar, Load calendar data and env variable
	if !isCalSynced {
		syncCalendar()
		isCalSynced = true
	}

	// Make Mega Raids dark green
	var colorId string
	if strings.Contains(e.Title, "Mega Raids") {
		colorId = Basil
	} else {
		colorId = e.EventType.ColorId()
	}

	// Make shadow raids show only on weekends
	var repeat []string
	if strings.Contains(e.Title, "Shadow Raids") {
		colorId = Grape
		endDate := e.StartDate.AddDate(0, 0, 7)

		for day := e.StartDate; day.Compare(endDate) != 1; day = day.Add(time.Hour * 24) {
			// Lowk i think its never gonna hit sunday before saturday
			if day.Weekday() == time.Saturday || day.Weekday() == time.Sunday {
				e.StartDate = localize(day)
				e.EndDate = localize(day.Add(time.Hour * 24))
				repeat = []string{"RRULE:FREQ=WEEKLY;COUNT=4"}
				break
			}
		}
	}

	// Create and post new pogo event
	gcalEvent := calendar.Event{
		Summary:     e.Title,
		Location:    e.Link,
		Description: e.Description,
		Start: &calendar.EventDateTime{
			DateTime: e.StartDate.UTC().Format(time.RFC3339),
			TimeZone: "UTC",
		},
		End: &calendar.EventDateTime{
			DateTime: e.EndDate.UTC().Format(time.RFC3339),
			TimeZone: "UTC",
		},
		ColorId:    colorId,
		Recurrence: repeat,
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

func deleteAllEvents() {
	if !isCalSynced {
		syncCalendar()
		isCalSynced = true
	}

	var check string
	fmt.Println("<> Are you sure you want to clear the calendar? <>")
	fmt.Println("Y/N?")
	_, err := fmt.Scanln(&check)
	if err != nil {
		log.Fatal(err)
	}

	switch check {
	case "Y", "y":
		// Page Token shennanigans cuz if I'm doing back-to-back large requests to GCal,
		//  Google creates a new page to put modified data on and sets the initial one as "tombstone"
		//  as to let other systems know that that page no longer has any data on it
		// If I don't manually specify the next page token with req.PageToken(), my code
		//  would only look at the tombstone page, thus ending the function pre-maturely.
		var pageToken string
		for {
			// 1. Build the request
			req := srv.Events.List(calendarID)
			if pageToken != "" {
				req.PageToken(pageToken)
			}

			// 2. Fetch the current page
			events, listErr := req.Do() // Populates events.NextPageToken
			if listErr != nil {
				log.Fatalf("Unable to retrieve events: %v", listErr)
			}

			// 3. Delete everything found on THIS page
			for _, item := range events.Items {
				fmt.Printf("Deleted: %v\n", item.Summary)
				err := srv.Events.Delete(calendarID, item.Id).Do()
				if err != nil {
					log.Printf("Could not delete event %s: %v\n", item.Summary, err)
				}
			}

			// 4. Move to the next page, or break if we are completely done
			pageToken = events.NextPageToken
			if pageToken == "" {
				break
			}
		}
		fmt.Println("<> Cleared all calendar events <>")
	}
}

// Localizes any string or time.Time object to your local time.
func localize(val any) time.Time {
	switch v := val.(type) {
	case string:
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			return t.Local()
		}
		if t, err := time.ParseInLocation("2006-01-02T15:04:05", v, time.Local); err == nil {
			return t
		}
	case time.Time:
		return v.Local()
	}
	return time.Time{}
}
