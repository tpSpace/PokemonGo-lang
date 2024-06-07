package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/chromedp/chromedp"
	"golang.org/x/net/html"
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	// fmt.Println(string(body))
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return
	}
	fmt.Print(doc)
	
	fmt.Println(newChromedp(baseUrl+"#/pokemon/1"))
}

func newChromedp(url string) string {
	ctx, cancel := chromedp.NewContext(
		context.Background(),
	)
	// to release the browser resources when
	// it is no longer needed
	defer cancel()

	var html string
	err := chromedp.Run(ctx,
		// visit the target page
		// detail-view-container
		chromedp.Navigate(url),
		// wait for the page to load
		chromedp.Sleep(10*time.Millisecond),
		chromedp.WaitVisible(`#detail-view-container`, chromedp.ByID),
		// get the HTML content
		chromedp.InnerHTML(`#detail-view-container`, &html, chromedp.ByID),
	)

	if err != nil {
		log.Fatal("Error while performing the automation logic:", err)
	}
	return html
}
