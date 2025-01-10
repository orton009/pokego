package main

import (
	"bufio"
	"fmt"
	"internal/pokecache"
	"os"
	"time"
)

const PokemonLocationEndpoint = "https://pokeapi.co/api/v2/location?offset=0&limit=20"

func commands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "display names of 20 location areas of pokemon",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "display names of previous 20 locations of pokemon",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "explore pokemone present in provided area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "catch a pokemon by name as first argument",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "get information of pokemon if already caught",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "view list of all pokemons caught by you",
			callback:    commandPokedex,
		},
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	config := config{Next: PokemonLocationEndpoint, Previous: "", cache: pokecache.NewCache(5 * time.Second), caughtPokemons: map[string]pokemon{}}
	for {
		fmt.Print("Pokedex > ")
		ok := scanner.Scan()
		if ok {
			words := cleanInput(scanner.Text())

			if command, ok := commands()[words[0]]; ok {
				err := command.callback(&config, words[1:])
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
			} else {
				fmt.Println("Unknown command ")
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}
