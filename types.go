package main

import (
	"fmt"
	"time"
)

type EventType int

const (
	Default EventType = iota
	Event
	RaidBattle
	RaidDay
	RaidHour
	CommunityDay
	MaxMonday
	SpotlightHour
	PokemonGoFest
)

var stringToEventType = map[string]EventType{
	"Default":         Default,
	"Event":           Event,
	"Raid Battles":    RaidBattle,
	"Raid Day":        RaidDay,
	"Raid Hour":       RaidHour,
	"Community Day":   CommunityDay,
	"Max Mondays":     MaxMonday,
	"Spotlight Hour":  SpotlightHour,
	"Pokémon GO Fest": PokemonGoFest,
}

func (e EventType) String() string {
	names := [...]string{
		"Default",
		"Event",
		"Raid Battles",
		"Raid Day",
		"Raid Hour",
		"Community Day",
		"Max Mondays",
		"Spotlight Hour",
		"PokemonGoFest",
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
	pokemonList   map[string]struct{} // Populated by all non-legendary pokemon on the page
	legendaryList map[string]struct{} // Populated by all legendary/mythic pokemon on the page
	megaList      map[string]struct{} // Populated by all mega pokemon on page (including legendaries)
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
