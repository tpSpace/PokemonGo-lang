package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

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

	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	fmt.Println("Error: ", err)
	// 	return
	// }

	// // fmt.Println(string(body))
	// doc, err := html.Parse(bytes.NewReader(body))
	// if err != nil {
	// 	fmt.Println("Error parsing HTML:", err)
	// 	return
	// }
	// fmt.Println(doc)
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
	var html string;
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
		chromedp.OuterHTML(`#detail-view-container h1`, &html, chromedp.ByID),
	}

	err := chromedp.Run(ctx,task);

	if err != nil {
		log.Fatal("Error while performing the automation logic:", err)
	}
	return html
}
