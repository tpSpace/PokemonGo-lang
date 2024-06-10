package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	Pokemon "tpSpace/PokemonGo-lang/entity"

	// "tpSpace/PokemonGo-lang/entity/pokemon"

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
		go fmt.Println(newChromedp(baseUrl + "#/pokemon/" + fmt.Sprint(i)))
	}
	// write a function to get the pokemon data with concurrencyq

}

func newChromedp(url string) Pokemon.Pokemon {
	ctx, cancel := chromedp.NewContext(
		context.Background(),
	)
	// to release the browser resources when
	// it is no longer needed
	defer cancel()

	var pokemon = Pokemon.GetPokemon()
	var temp string
	var temp2 []*cdp.Node
	// run the automation logic
	task := chromedp.Tasks{
		// visit the target page
		// detail-view-container
		chromedp.Navigate(url),
		chromedp.Reload(),
		// wait for the page to load
		chromedp.WaitReady(`#detail-view-container`, chromedp.ByID),
		chromedp.Sleep(200*time.Millisecond),
		
		// get the HTML content
		chromedp.Text(`#detail-view-container h1`, &pokemon.General.Name, chromedp.ByID),
		chromedp.Evaluate(`(() => {
            const monsterTypes = document.querySelectorAll('.detail-types > span.monster-type');
            const visibleMonsterTypes = [];
            monsterTypes.forEach(type => {
                const style = window.getComputedStyle(type);
                if (style.display !== 'none' && style.visibility !== 'hidden' && style.opacity !== '0') {
                    visibleMonsterTypes.push(type.textContent);
                }
            });
            return visibleMonsterTypes;
        })()`, &pokemon.General.Types),
		// get the profile data
		
		
		// get the text content of the element
		chromedp.ActionFunc(func (ctx context.Context) error {
			
						
			// get the general data
			chromedp.Text(".detail-national-id span", &temp, chromedp.ByQuery).Do(ctx)
			// the index is in the format #001 therefore we have to slice the first character
			pokemon.General.Index, _ = strconv.Atoi(temp[1:])
			// convert the text to an integer using js
			chromedp.Evaluate(`parseInt(document.querySelector(".detail-stats-row:nth-of-type(2) .stat-bar-fg").textContent)`, &pokemon.General.Attack, chromedp.EvalAsValue).Do(ctx)
        	chromedp.Text(".detail-stats-row:nth-of-type(1) .stat-bar-fg", &temp).Do(ctx)
			pokemon.General.HP, _ = strconv.Atoi(temp)
			chromedp.Text(".detail-stats-row:nth-of-type(3) .stat-bar-fg", &temp).Do(ctx)
			pokemon.General.Defense, _ = strconv.Atoi(temp)
			chromedp.Text(".detail-stats-row:nth-of-type(4) .stat-bar-fg", &temp).Do(ctx)
			pokemon.General.Speed, _ = strconv.Atoi(temp)
			chromedp.Text(".detail-stats-row:nth-of-type(5) .stat-bar-fg", &temp).Do(ctx)
			pokemon.General.Sp_Atk, _ = strconv.Atoi(temp)
			chromedp.Text(".detail-stats-row:nth-of-type(6) .stat-bar-fg", &temp).Do(ctx)
			pokemon.General.Sp_Def, _ = strconv.Atoi(temp)
			chromedp.Text(".detail-below-header .monster-species", &pokemon.General.Monster_Species).Do(ctx)
			chromedp.Text(".detail-below-header .monster-description", &pokemon.General.Monster_Description).Do(ctx)
			fmt.Println("hii")
			
			// get the profile data
			// chromedp.Text(".detail-below-header .monster-minutia", &temp, chromedp.ByQueryAll).Do(ctx)
			
			chromedp.Nodes(`#detail-view .detail-view-fg .mui-panel .detail-panel-content .detail-below-header .monster-minutia span`, &temp2, chromedp.ByQueryAll).Do(ctx)
			// go through temp2 and get the profile data
			// for _, node := range temp2 {
				
			// 	// get the children of the node
			// 	for i := int64(0); i < node.ChildNodeCount; i++ {
					
			// 		// class, ok := node.Children[0].Attribute("class")
			// 		// if ok && class == "monster-minutia" {
						
			// 		// }
			// 		// text := node.Children[i].Children[0].NodeValue
			// 		// 	fmt.Println(text)
			// 		// check if class is monster-minutia
			// 		class, ok := node.Children[i].Attribute("class")
			// 		if class == "monster-minutia" && ok {
			// 			for _, child := range node.Children[i].Children {
			// 				if child.NodeType == 3 { // Node.TEXT_NODE
			// 					text := child.NodeValue
			// 					fmt.Println(text)
			// 				}
			// 			}
			// 		}
			// 	}
			// }
			for _, node := range temp2 {
				text := node.Children[0].NodeValue
				fmt.Println(text)
			}
			return nil
		}),
		
	}
	
	err := chromedp.Run(ctx,task);
	if err != nil {
		log.Fatal("Error while performing the automation logic:", err)
	}
	return pokemon
}
