package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	Pokemon "tpSpace/PokemonGo-lang/entity"

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
	var pokemon []Pokemon.Pokemon
	for i := 1; i <= 3; i++ { //649
		pokedex := newChromedp(baseUrl + "#/pokemon/" + fmt.Sprint(i))
		// save the data to the pokedex.json as a json file
		fmt.Println(pokedex)
		pokemon = append(pokemon, pokedex)
		// save pokemon data to the pokedex.json file as a json file
		// cod
		
	}
	pokedexJson, err := json.MarshalIndent(pokemon, "", " ")
    if err != nil {
        fmt.Println("Error: ", err)
        return
    }

    // Write the JSON data to the pokedex.json file in correct format and with proper indentation
	
    err = os.WriteFile("pokedex.json", pokedexJson, 0644)
    if err != nil {
        fmt.Println("Error: ", err)
        return
    }
	// write a function to get the pokemon data with concurrency
	
}

func newChromedp(url string) Pokemon.Pokemon {
	// open the browser
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		// chromedp.Flag("headless", false),
		chromedp.Flag("disable-gpu", false),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	// create a new context

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()
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
		chromedp.Sleep(100*time.Millisecond),
		
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
			chromedp.WaitVisible(".detail-view-fg .mui-panel .detail-panel-content .detail-below-header .monster-minutia span")
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
				case "rock":
					pokemon.DamgeWhenAttacked.Rock = multiplier.Multipler
				case "steel":
					pokemon.DamgeWhenAttacked.Steel = multiplier.Multipler
				case "poison":
					pokemon.DamgeWhenAttacked.Poison = multiplier.Multipler
				case "ghost":
					pokemon.DamgeWhenAttacked.Ghost = multiplier.Multipler
				case "dark":
					pokemon.DamgeWhenAttacked.Dark = multiplier.Multipler
				case "dragon":
					pokemon.DamgeWhenAttacked.Dragon = multiplier.Multipler
				case "bug":
					pokemon.DamgeWhenAttacked.Bug = multiplier.Multipler
				case "normal":
					pokemon.DamgeWhenAttacked.Normal = multiplier.Multipler
				default:
					fmt.Println("Error: ", multiplier.Type)
				}
			
			}
			return nil
		}).Do(ctx)
		// get the evolutions data
		const getEvolutionsScript = `(function() {
			const data = [];
			document.querySelectorAll('.evolution-label').forEach(el => {
				const text = el.innerText;
				const from = text.split(' evolves into ')[0].trim();
				const to = text.split(' evolves into ')[1].split(' at level ')[0].trim();
				const level = parseInt(text.split(' at level ')[1].trim());
				data.push({ from, to, level });
			});
			return data;
		})();`
		// Wait for the element to be visible
		chromedp.WaitVisible(".evolution-row")
		var evolutions []struct {
			Level int
			From  string
			To    string
		}
		// Query for the evolutions and populate the struct
		chromedp.Evaluate(getEvolutionsScript, &evolutions).Do(ctx)
		for _, evolution := range evolutions {
			pokemon.Evolutions = append(pokemon.Evolutions, Pokemon.Evolutions{
				Level: evolution.Level,
				From:  evolution.From,
				To:    evolution.To,
			})
		}
		fmt.Println("Hello world")
			return nil
		}),
		// chromedp.Sleep(5 * time.Second),
		chromedp.WaitVisible(`div.monster-moves .moves-row`),    // Wait for the moves container to be visible
		// get the natural moves data
		chromedp.ActionFunc(func(ctx context.Context) error {
			var movesData []struct {
				Level       string
				Name        string
				Type        string
				Power       string
				Accuracy    string
				PP          string
				Description string
			}
			
			chromedp.WaitVisible(`div.monster-moves .moves-row`) // Wait for the moves container to be visible
			chromedp.Evaluate(`
			(function() {
				const rows = document.querySelectorAll('div.monster-moves .moves-row');
				let data = [];
				rows.forEach(row => {
					const move = {};
					move.level = row.querySelector('.moves-inner-row > span:nth-child(1)').innerText;
					move.name = row.querySelector('.moves-inner-row > span:nth-child(2)').innerText;
					move.type = row.querySelector('.moves-inner-row > span.monster-type').innerText;
					const stats = row.querySelector('.moves-row-detail .moves-row-stats');
					if (stats) {
						move.power = stats.querySelector('span:nth-child(1)').innerText.split(': ')[1];
						move.accuracy = stats.querySelector('span:nth-child(2)').innerText.split(': ')[1];
						move.pp = stats.querySelector('span:nth-child(3)').innerText.split(': ')[1];
					}
					const description = row.querySelector('.moves-row-detail .move-description');
					if (description) {
						move.description = description.innerText;
					}
					data.push({
						Level: move.level,
						Name: move.name,
						Type: move.type,
						Power: move.power,
						Accuracy: move.accuracy,
						PP: move.pp,
						Description: move.description
					});
				});
				return data;
			})();
		`, &movesData, chromedp.EvalAsValue).Do(ctx)
		// fmt.Println(movesData)
			// now map it to Pokemon.Natural_Moves
			// convert data to 
			for _, move := range movesData {
				// iterate through the move data and fix if the value is not a number then set it to -1 
				var moveData Pokemon.Natural_Moves
				if move.Power == "" || move.Power == "N/A" {
					moveData.Type.Power = -1
				} else {
					moveData.Type.Power, _ = strconv.Atoi(move.Power)
				}
				if move.Accuracy == "" || move.Accuracy == "N/A" {
					moveData.Type.Acc = -1
				} else {
					// remove % from the string
					moveData.Type.Acc, _ = strconv.Atoi(strings.TrimSuffix(move.Accuracy, "%"))
				}
				if move.PP == "" || move.PP == "N/A" {
					moveData.Type.PP = -1
				} else {
					moveData.Type.PP, _ = strconv.Atoi(move.PP)
				}
				moveData.Index, _ = strconv.Atoi(move.Level)
				moveData.Move = move.Name
				moveData.Type.Description = move.Description
				pokemon.Natural_Moves = append(pokemon.Natural_Moves, moveData)
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