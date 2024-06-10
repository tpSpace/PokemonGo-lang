package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
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
			
			// get the profile data
			// chromedp.Text(".detail-below-header .monster-minutia", &temp, chromedp.ByQueryAll).Do(ctx)
			
			chromedp.Nodes(`#detail-view .detail-view-fg .mui-panel .detail-panel-content .detail-below-header .monster-minutia span`, &temp2, chromedp.ByQueryAll).Do(ctx)
			var count = 0

			for _, node := range temp2 {
				switch count {
				case 0:
					// trim the text to get the value
					re := regexp.MustCompile(`[0-9.]+`)
					valueStr := re.FindString(node.Children[0].NodeValue)
					value, _ := strconv.ParseFloat(valueStr, 32)
					pokemon.Profile.Height = float32(value)
				case 1:
					re := regexp.MustCompile(`[0-9.]+`)
					valueStr := re.FindString(node.Children[0].NodeValue)
					value, _ := strconv.ParseFloat(valueStr, 32)
					pokemon.Profile.Weight = float32(value)
				case 2:
					re := regexp.MustCompile(`[0-9.]+`)
					valueStr := re.FindString(node.Children[0].NodeValue)
					value, _ := strconv.ParseFloat(valueStr, 32)
					
					pokemon.Profile.Catch_Rate = float32(value)
				case 3:
					re := regexp.MustCompile(`(\d+(\.\d+)?)%`)
					matches := re.FindAllString(node.Children[0].NodeValue, -1)

					for _, match := range matches {
						valueStr := strings.TrimSuffix(match, "%")
						value, _ := strconv.ParseFloat(valueStr, 32)
						pokemon.Profile.Gender_Ratio = append(pokemon.Profile.Gender_Ratio, float32(value))
					}
				case 4:
					re := regexp.MustCompile(`[a-zA-Z]+`)
					matches := re.FindAllString(node.Children[0].NodeValue, -1)
					pokemon.Profile.Egg_Groups = append(pokemon.Profile.Egg_Groups, matches...)
				case 5:
					pokemon.Profile.Hatch_Steps, _ = strconv.Atoi(node.Children[0].NodeValue)
				case 6:
					re := regexp.MustCompile(`[a-zA-Z]+`)
					matches := re.FindAllString(node.Children[0].NodeValue, -1)
					pokemon.Profile.Abilities = append(pokemon.Profile.Abilities, matches...)
				case 7:
					// split the string by "," and then trim the string to get the value the substring may have space between the characters
					for _, value := range strings.Split(node.Children[0].NodeValue, ",") {
						pokemon.Profile.EVs = append(pokemon.Profile.EVs, strings.TrimSpace(value))
					}
				default: 
					fmt.Println("Error: ", node.Children[0].NodeValue)
				}
				count++				
			}
			// get the damge when attacked data

			const getDamageMultipliersScript = `(function() {
				var multipliers = [];
				var rows = document.querySelectorAll(".when-attacked-row");
				rows.forEach(function(row) {
					var typeElements = row.querySelectorAll(".monster-type");
					var multiplierElements = row.querySelectorAll(".monster-multiplier");
					for (var i = 0; i < typeElements.length; i++) {
						var type = typeElements[i].innerText;
						var multiplier = parseFloat(multiplierElements[i].innerText.replace("x", ""));
						multipliers.push({Type: type, Multipler: multiplier});
					}
				});
				return multipliers;
			})();`
			// Wait for the element to be visible
		chromedp.WaitVisible(".when-attacked")

		// Query for the damage multipliers and populate the struct
		chromedp.ActionFunc(func(ctx context.Context) error {
			var multipliers []struct {
				Type      string
				Multipler float32
			}
			err := chromedp.Evaluate(getDamageMultipliersScript, &multipliers).Do(ctx)
			if err != nil {
				return err
			}
			// fmt.Println(multipliers)
			for _, multiplier := range multipliers {
				switch multiplier.Type {
				case "ice":
					// remove "x"
					pokemon.DamgeWhenAttacked.Ice = multiplier.Multipler
				case "electric":
					pokemon.DamgeWhenAttacked.Electric = multiplier.Multipler
				case "fire":
					pokemon.DamgeWhenAttacked.Fire = multiplier.Multipler
				case "water":
					pokemon.DamgeWhenAttacked.Water = multiplier.Multipler
				case "flying":
					pokemon.DamgeWhenAttacked.Flying = multiplier.Multipler
				case "fairy":
					pokemon.DamgeWhenAttacked.Fairy = multiplier.Multipler
				case "psychic":
					pokemon.DamgeWhenAttacked.Psychic = multiplier.Multipler
				case "fighting":
					pokemon.DamgeWhenAttacked.Fighting = multiplier.Multipler
				case "ground":
					pokemon.DamgeWhenAttacked.Ground = multiplier.Multipler
				case "grass":
					pokemon.DamgeWhenAttacked.Grass = multiplier.Multipler
				}
			}
			return nil
		}).Do(ctx)


			return nil
		}),
		
	}
	
	err := chromedp.Run(ctx,task);
	if err != nil {
		log.Fatal("Error while performing the automation logic:", err)
	}
	return pokemon
}
