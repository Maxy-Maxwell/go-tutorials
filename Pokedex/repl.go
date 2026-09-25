package main

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"os"
	"sort"
	"strings"

	"github.com/Maxy-Maxwell/go-tutorials/Pokedex/internal/callouts"
)

var pokedex map[string]catchResponse = map[string]catchResponse{}

func runRepl(conf *config) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Pokedex > ")

	for scanner.Scan() {
		input := scanner.Text()
		inputWords := cleanInput(input)

		if len(inputWords) == 0 {
			continue
		}

		cmdDetails, ok := conf.commands[inputWords[0]]

		if len(inputWords) < 2 {
			inputWords = append(inputWords, "")
		}

		if !ok {
			fmt.Println("Unknown command")
			fmt.Print("Pokedex > ")
			continue
		}

		err := cmdDetails.callback(conf, inputWords[1])

		if err != nil {
			fmt.Printf("Callback function err: %v\n", err.Error())
		}

		fmt.Print("Pokedex > ")
	}
}

func cleanInput(text string) []string {

	if text == "" {
		return []string{}
	}

	splitStrings := strings.Fields(text)
	results := []string{}

	for _, str := range splitStrings {
		str = strings.ToLower(str)

		results = append(results, str)
	}

	return results
}

func commandExit(c *config, _ string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *config, _ string) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")

	sortedKeys := []string{}

	for cmd := range c.commands {
		sortedKeys = append(sortedKeys, cmd)
	}

	sort.Strings(sortedKeys)

	for _, key := range sortedKeys {
		fmt.Printf("%s: %s\n", c.commands[key].name, c.commands[key].description)

	}
	return nil
}

func commandMap(c *config, _ string) error {

	if c.nextLocationURL == "" {
		fmt.Println("There's no more map values. Try 'mapb' to get the previous set of locations.")
	}

	parsedResponse := mapResponse{}
	err := callouts.DoGetCall[mapResponse](c.nextLocationURL, &parsedResponse)
	if err != nil {
		return err
	}

	c.previousLocationURL = c.nextLocationURL

	if parsedResponse.Next == nil {
		c.nextLocationURL = ""
	} else {
		c.nextLocationURL = *parsedResponse.Next
	}

	for _, result := range parsedResponse.Results {
		fmt.Println(result.Name)
	}

	return nil
}

func commandMapb(c *config, _ string) error {

	if c.previousLocationURL == "" {
		fmt.Println("There's no previous map values. Try 'map' to get the next set of locations.")
		return nil
	}

	parsedResponse := mapResponse{}
	err := callouts.DoGetCall[mapResponse](c.previousLocationURL, &parsedResponse)
	if err != nil {
		return err
	}

	if parsedResponse.Previous == nil {
		c.previousLocationURL = ""
	} else {
		c.previousLocationURL = *parsedResponse.Previous
	}

	if parsedResponse.Next == nil {
		c.nextLocationURL = ""
	} else {
		c.nextLocationURL = *parsedResponse.Next
	}

	for _, result := range parsedResponse.Results {
		fmt.Println(result.Name)
	}

	return nil
}

func commandExplore(c *config, arg string) error {
	if arg == "" {
		return fmt.Errorf("Missing 'map area'. Append a 'map area' to the explore command. e.g. 'explore pastoria-city-area'. Map areas can be found using the 'map' and 'mapb' commands.")
	}

	fmt.Printf("Exploring %v...\n", arg)

	parsedResponse := exploreResponse{}
	err := callouts.DoGetCall[exploreResponse](fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%v", arg), &parsedResponse)
	if err != nil {
		return err
	}

	if len(parsedResponse.Pokemon_encounters) == 0 {
		fmt.Println("Looks like we couldn't find any Pokemon in this area. Double check that the area is spelt correctly.")
	}

	for _, encounter := range parsedResponse.Pokemon_encounters {
		fmt.Printf(" - %v\n", encounter.Pokemon.Name)
	}

	return nil
}

func commandCatch(c *config, arg string) error {
	if arg == "" {
		return fmt.Errorf("Couldn't find a %v. Double check the spelling, or use 'explore' to see a list of available pokemon for a map area.", arg)
	}

	parsedResponse := catchResponse{}
	err := callouts.DoGetCall[catchResponse](fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%v", arg), &parsedResponse)
	if err != nil {
		return err
	}

	fmt.Printf("Throwing a Pokeball at %v...\n", arg)

	// Lowest base xp is 36 (blipbug/sunkern), so a lowest pokemon should be caught every ~2 attempts
	// Highest is 635 (blissey), so they should be caught once every ~36ish attempts
	expectedCatchRate := float64(parsedResponse.Base_experience) / 18.0

	rollForCatch := rand.IntN(int(expectedCatchRate-1+1)) + 1

	if rollForCatch == int(expectedCatchRate) {
		fmt.Printf("%v was caught!\n", arg)
		pokedex[arg] = parsedResponse
	} else {
		fmt.Printf("%v escaped!\n", arg)
	}

	return nil
}

func commandInspect(c *config, arg string) error {
	pokeData, ok := pokedex[arg]

	if arg == "" || !ok {
		return fmt.Errorf("You aint caught this fella (%v) yet big guy. Try catching one THEN I'll fill you in on it's stats. (checkout the 'catch' command)", arg)
	}

	fmt.Println("Name:", pokeData.Name)
	fmt.Println("Height:", pokeData.Height)
	fmt.Println("Weight:", pokeData.Weight)
	fmt.Println("Stats:")

	for _, stat := range pokeData.Stats {
		fmt.Printf("  - %v: %v\n", stat.Stat.Name, stat.Base_stat)
	}

	fmt.Println("Types:")

	for _, t := range pokeData.Types {
		fmt.Printf("  - %v\n", t.Type.Name)
	}

	return nil
}

func commandPokedex(c *config, arg string) error {
	sortedPokemon := []string{}

	for entry := range pokedex {
		sortedPokemon = append(sortedPokemon, entry)
	}

	sort.Strings(sortedPokemon)

	for _, key := range sortedPokemon {
		fmt.Printf("  - %s\n", key)
	}

	return nil
}
