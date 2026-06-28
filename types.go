package main

import (
	"fmt"
	"time"
)

type EventType int

const (
	Event EventType = iota
	RaidBattle
	RaidDay
	RaidHour
	CommunityDay
	MaxMonday
	SpotlightHour
	Special // GO Fest, other big events like this
)

var stringToEventType = map[string]EventType{
	"Event":          Event,
	"Raid Battles":   RaidBattle,
	"Raid Day":       RaidDay,
	"Raid Hour":      RaidHour,
	"Community Day":  CommunityDay,
	"Max Mondays":    MaxMonday,
	"Spotlight Hour": SpotlightHour,
	"Special":        Special,
}

func (e EventType) String() string {
	names := [...]string{
		"Event",
		"Raid Battles",
		"Raid Day",
		"Raid Hour",
		"Community Day",
		"Max Mondays",
		"Spotlight Hour",
		"Special",
	}
	if e < 0 || int(e) >= len(names) {
		return fmt.Sprintf("EventType(%d)", e)
	}
	return names[e]
}

type CalendarEvent struct {
	title         string
	startDate     time.Time // includes date and time
	endDate       time.Time
	eventType     EventType //
	link          string
	pokemonList   []string // Populated by all non-legendary pokemon on the page
	legendaryList []string // Populated by all legendary/mythic pokemon on the page
	megaList      []string // Populated by all mega pokemon on page (including legendaries)
}

func (c CalendarEvent) String() string {
	msg := "Title\t\t%v\n" +
		"Start Date\t%v\n" +
		"End Date\t%v\n" +
		"Event Type\t%v\n" +
		"Link\t\t%v\n" +
		"Pokemon\t\t%v\n" +
		"Legendaries\t%v\n" +
		"Megas\t\t%v"

	return fmt.Sprintf(msg, c.title, c.startDate, c.endDate, c.eventType, c.link, c.pokemonList, c.legendaryList, c.megaList)
}
