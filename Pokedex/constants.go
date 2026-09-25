package main

type cliCommand struct {
	name        string
	description string
	callback    func(*config, string) error
}

type config struct {
	commands            map[string]cliCommand
	nextLocationURL     string
	previousLocationURL string
}

func initConstants(conf *config) error {
	// Set constant values
	validCommands := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Retrieve the next batch of map areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Retrieve the previous batch of map areas",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Retrieve a list of Pokemon for a given map area. E.g. 'explore pastoria-city-area'",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Try your luck at catching a Pokemon. Not sure what to catch? Try the 'explore' command. E.g. 'catch tentacool'",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Learn about a Pokemon you've caught before. e.g. 'inspect tentacool'",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "List all the pokemon you've caught!",
			callback:    commandPokedex,
		},
	}

	*conf = config{
		commands:            validCommands,
		nextLocationURL:     "https://pokeapi.co/api/v2/location-area?limit=20",
		previousLocationURL: "",
	}

	return nil
}
