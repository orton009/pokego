package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
)

func commandExit(_ *config, _ []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func FetchAndShowAreas(config *config, url string) error {
	bytes, ok := config.cache.Get(url)

	var areas []areaMetadata
	var locationResponse pokemonLocations

	if !ok {
		// fetch data from api if not found in cache
		data, err := FetchAllLocations(url)
		if err != nil {
			return err
		}
		locationResponse = data

		locations, err := FetchAreaMetadata(data)
		if err != nil {
			return err
		}
		areas = GetAreasFromLocations(locations)

		bytes, err := json.Marshal(LocationCache{
			Areas:            areas,
			PokemonLocations: locationResponse,
		})
		if err != nil {
			return err
		}
		config.cache.Add(url, bytes)

	} else {
		// decode data from cache
		locationCache := LocationCache{}

		err := json.Unmarshal(bytes, &locationCache)
		if err != nil {
			return err
		}

		locationResponse = locationCache.PokemonLocations
		areas = locationCache.Areas

	}

	config.Previous = locationResponse.Previous
	config.Next = locationResponse.Next

	for _, area := range areas {
		fmt.Println(area.Name)
	}

	return nil
}

func commandMap(config *config, _ []string) error {
	return FetchAndShowAreas(config, config.Next)
}

func commandMapb(config *config, _ []string) error {
	return FetchAndShowAreas(config, config.Previous)
}

func commandHelp(_ *config, _ []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	for _, c := range commands() {
		fmt.Println(c.name + ": " + c.description)
	}
	return nil
}

func commandExplore(c *config, params []string) error {
	if len(params) == 0 {
		return errors.New("expected area name as command argument")
	}
	area := params[0]

	pokemons := []string{}
	bytes, ok := c.cache.Get(area)

	if !ok {

		data, err := FetchPokemonsInArea(area)
		if err != nil {
			return err
		}
		for _, p := range data.PokemonEncounters {
			pokemons = append(pokemons, p.Pokemon.Name)
		}

		jsonData, err := json.Marshal(pokemons)
		if err != nil {
			return err
		}
		c.cache.Add(area, jsonData)

	} else {
		err := json.Unmarshal(bytes, &pokemons)
		if err != nil {
			return err
		}
	}

	fmt.Println("Exploring " + area + "...")
	fmt.Println("Found Pokemon:")
	for _, n := range pokemons {
		fmt.Println("\t- ", n)
	}

	return nil
}

func commandCatch(config *config, args []string) error {
	if len(args) == 0 {
		return errors.New("expected pokemon name as first argument")
	}
	name := args[0]

	pokemon, err := FetchPokemonByName(name)
	if err != nil {
		return err
	}

	_, ok := config.caughtPokemons[name]
	if !ok {
		// add logic for catching pokemon
		if rand.Intn(10) > 5 {
			fmt.Println("Throwing a Pokeball at " + name + "...")
			config.caughtPokemons[name] = pokemon
		} else {
			fmt.Println("Oops! you just missed..")
		}
	} else {
		fmt.Println("Pokemon already caught...")
	}

	return nil
}

func commandInspect(config *config, args []string) error {
	if len(args) == 0 {
		return errors.New("expected atleast one argument i.e. pokemon name")
	}
	name := args[0]

	p, ok := config.caughtPokemons[name]
	if !ok {
		fmt.Println("you have not caught that pokemon")
	} else {
		bytes, err := json.MarshalIndent(p, "", "\t")
		if err != nil {
			return err
		}
		fmt.Println(string(bytes))
	}

	return nil
}

func commandPokedex(config *config, _ []string) error {
	fmt.Println("Your Pokedex:")
	for k := range config.caughtPokemons {
		fmt.Println("- ", config.caughtPokemons[k].Name)
	}
	return nil
}
