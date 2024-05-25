package main

import (
	"fmt"
	"sync"

	"github.com/gocolly/colly"
)

func main() {
	// Initialize a new collector
	c := colly.NewCollector()

	// Use a wait group to handle multiple requests concurrently
	var wg sync.WaitGroup

	// Set up the callback for when the HTML element is found
	c.OnHTML("h1.detail-panel-header", func(e *colly.HTMLElement) {
		// Extract the Pokémon name from the <h1> tag
		pokemonName := e.Text
		fmt.Println("Pokemon Name:", pokemonName)
	})

	// Set up the callback for when a request is made
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Loop over the first 10 Pokémon IDs
	for i := 1; i <= 10; i++ {
		wg.Add(1) // Add to the wait group
		go func(id int) {
			defer wg.Done() // Mark this goroutine as done in the wait group
			// Construct the URL
			url := fmt.Sprintf("https://pokedex.org/#/pokemon/%d", id)
			// Visit the URL
			c.Visit(url)
		}(i)
	}

	// Wait for all goroutines to finish
	wg.Wait()
}
