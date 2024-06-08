package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

func main() {
	baseUrl := "https://pokedex.org/"
	fmt.Println("Base url: " + baseUrl)

	resp, err := http.Get(baseUrl)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	defer resp.Body.Close()

	for i := 1; i <= 649; i++ {
		fmt.Println(newChromedp(baseUrl + "#/pokemon/" + fmt.Sprint(i)))
	}
}

func newChromedp(url string) string {
	ctx, cancel := chromedp.NewContext(
		context.Background(),
	)
	// to release the browser resources when
	// it is no longer needed
	defer cancel()
	var name, index, hp, attack, defense, speed, spAtk, spDef string
		var types []*cdp.Node
	// run the automation logic
	task := chromedp.Tasks{
		// visit the target page
		// detail-view-container
		chromedp.Navigate(url),
		chromedp.Reload(),
		// wait for the page to load
		chromedp.WaitReady(`#detail-view-container`, chromedp.ByID),
		chromedp.Sleep(1*time.Second),
		
		// get the HTML content
		chromedp.Text(`#detail-view-container h1`, &name, chromedp.ByID),
		// chromedp.Text(".monster-type:nth-of-type(1)", &types, chromedp.NodeVisible),
        // chromedp.Text(".monster-type:nth-of-type(2)", &temp, chromedp.NodeVisible, chromedp.ByQueryAll),
		// find all the .monster-type elements :nth-of-type(n)
		chromedp.Nodes(".monster-type", &types, chromedp.ByQueryAll),
		
        chromedp.Text(".detail-national-id span", &index),
        chromedp.Text(".detail-stats-row:nth-of-type(1) .stat-bar-fg", &hp),
        chromedp.Text(".detail-stats-row:nth-of-type(2) .stat-bar-fg", &attack),
        chromedp.Text(".detail-stats-row:nth-of-type(3) .stat-bar-fg", &defense),
        chromedp.Text(".detail-stats-row:nth-of-type(4) .stat-bar-fg", &speed),
        chromedp.Text(".detail-stats-row:nth-of-type(5) .stat-bar-fg", &spAtk),
        chromedp.Text(".detail-stats-row:nth-of-type(6) .stat-bar-fg", &spDef),
	}

	err := chromedp.Run(ctx,task);
	fmt.Println("Name:", name)
    var typez []string
for i, _ := range types {
    var text string
    chromedp.Text(i, &text)
    typez = append(typez, text)
	
}
	fmt.Println("Types:", typez)
    fmt.Println("Index:", index)
    fmt.Println("HP:", hp)
    fmt.Println("Attack:", attack)
    fmt.Println("Defense:", defense)
    fmt.Println("Speed:", speed)
    fmt.Println("Sp Atk:", spAtk)
    fmt.Println("Sp Def:", spDef)
	if err != nil {
		log.Fatal("Error while performing the automation logic:", err)
	}
	return name
}
