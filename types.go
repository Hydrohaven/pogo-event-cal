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
	Title         string              `json:"title"`
	StartDate     time.Time           `json:"startDate"`
	EndDate       time.Time           `json:"endDate"`
	EventType     EventType           `json:"eventType"`
	Link          string              `json:"link"`
	PokemonList   map[string]struct{} `json:"pokemonList"`
	LegendaryList map[string]struct{} `json:"legendaryList"`
	MegaList      map[string]struct{} `json:"megaList"`
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

	return fmt.Sprintf(msg, c.Title, c.StartDate, c.EndDate, c.EventType, c.Link, c.PokemonList, c.LegendaryList, c.MegaList)
}
