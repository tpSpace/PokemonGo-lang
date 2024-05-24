package main

import (
	fmt
	"github.com/gocolly/colly"
)

func main() {
	fmt.Println("Hello, World!")
	// crawl data using 
	// id is number
	// link https://pokedex.org/#/pokemon/:id 
	// example https://pokedex.org/#/pokemon/1
	// get the data from the website and store it in the database
	// and save it in json file 

	c := colly.NewCollector()
	
}