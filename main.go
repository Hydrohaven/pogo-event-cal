package main

import (
	"fmt"
)

func main() {
	event := CalendarEvent{
		title:     "Sample Event",
		eventType: RaidBattle,
	}

	fmt.Print(event)
}
