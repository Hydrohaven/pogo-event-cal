package main

import (
	"fmt"
	"net/http"
)

func main() {
	// Fetch HTML Content
	resp, _ := http.Get("https://leekduck.com/events/")
	defer resp.Body.Close()

	fmt.Print(resp)
}
