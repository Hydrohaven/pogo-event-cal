package main

import (
	"fmt"
)

var startFlag bool = false
var option string

func startCLI() {
	if !startFlag {
		welcomeMsg := fmt.Sprintf(`
	                                            _                  _ 
 _ __   ___   __ _  ___         _____   _____ _ __ | |_       ___ __ _| |
| '_ \ / _ \ / _ ` + "`" + `|/ _ \ _____ / _ \ \ / / _ \ '_ \| __|____ / __/ _` + "`" + ` | |
| |_) | (_) | (_| | (_) |_____|  __/\ V /  __/ | | | ||_____| (_| (_| | |
| .__/ \___/ \__, |\___/       \___| \_/ \___|_| |_|\__|     \___\__,_|_|
|_|          |___/                                                       
	`)
		fmt.Println(welcomeMsg)
		fmt.Println("<> Welcome! What would you like to do? <>")
		startFlag = true
	}

	fmt.Println("(1) Run calendar sync")
	fmt.Println("(2) See all calendar events")
	fmt.Println("(3) Delete specific calendar event")
	fmt.Println("(4) Delete all calendar events")
	fmt.Println("(5) Exit")

	fmt.Print("<> ")
	_, err := fmt.Scanln(&option)
	if err != nil {
		fmt.Println("<> Error reading input: <>", err)
		return
	}

	options()
}

func options() {
	switch option {
	case "1":
		startSync()
	case "2":
		return
	case "3":
		return
	case "4":
		deleteAllEvents()
	case "5":
		// fmt.Println("<> Shutting off... <>")
		return
	}

	options()
}
