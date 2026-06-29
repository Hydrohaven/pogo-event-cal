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

func (e EventType) ColorId() string {
	switch e {
	case Event:
		return "7"
	case RaidBattle:
		return "2"
	case RaidDay:
		return "10"
	case RaidHour:
		return "10"
	case CommunityDay:
		return "1"
	case MaxMonday:
		return "4"
	case SpotlightHour:
		return "1"
	case PokemonGoFest:
		return "9"
	default:
		return "6"
	}
}

type CalendarEvent struct { // struct tags customized how they're formatted json, camelCase is proper json, they would be Pascal without the tags!
	Title         string              `json:"title"`
	Description   string              `json:"description"`
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

type Config struct {
	CalendarID string `json:"calendar_id"`
}
