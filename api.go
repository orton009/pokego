package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type areaMetadata struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type stats struct {
	Stat struct {
		Name string `json:"name"`
	} `json:"stat"`

	BaseStat int `json:"base_stat"`
}

type pokemon struct {
	Name   string  `json:"name"`
	Height int     `json:"height"`
	Weight int     `json:"weight"`
	Stats  []stats `json:"stats"`
}

type location struct {
	Areas []areaMetadata `json:"areas"`
}

type pokemonResponse struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type locationMetadata struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type pokemonLocations struct {
	Count    int                `json:"count"`
	Next     string             `json:"next"`
	Previous string             `json:"previous"`
	Results  []locationMetadata `json:"results"`
}

type LocationCache struct {
	PokemonLocations pokemonLocations `json:"pokemonLocations"`
	Areas            []areaMetadata   `json:"areas"`
}

func FetchLocation(url string) (location, error) {
	res, err := http.Get(url)

	location := location{}
	if err != nil {
		return location, fmt.Errorf("error fetching pokemon location: %w", err)
	}

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&location)
	return location, err
}

func FetchPokemonByName(name string) (pokemon, error) {
	pokemon := pokemon{}
	res, err := http.Get("https://pokeapi.co/api/v2/pokemon/" + name)
	if err != nil {
		return pokemon, err
	}

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&pokemon)
	if err != nil {
		return pokemon, err
	}

	return pokemon, nil
}

func FetchPokemonsInArea(area string) (pokemonResponse, error) {
	url := "https://pokeapi.co/api/v2/location-area/" + area
	res, err := http.Get(url)
	result := pokemonResponse{}

	if err != nil {
		return result, err
	}
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func FetchAllLocations(url string) (pokemonLocations, error) {
	res, err := http.Get(url)
	pokemon := pokemonLocations{}

	if err != nil {
		return pokemon, fmt.Errorf("error fetching pokemon: %w", err)
	}

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&pokemon)
	if err != nil {
		return pokemon, fmt.Errorf("error decoding pokemon response: %w", err)
	}
	return pokemon, nil
}

func FetchAreaMetadata(p pokemonLocations) ([]location, error) {
	locations := []location{}

	for _, locationMetadata := range p.Results {
		// TODO: use routine to concurrently fetch all locations
		location, err := FetchLocation(locationMetadata.Url)
		if err != nil {
			return locations, err
		}
		locations = append(locations, location)
	}

	return locations, nil
}

func GetAreasFromLocations(locations []location) []areaMetadata {
	areas := []areaMetadata{}
	for _, l := range locations {
		areas = append(areas, l.Areas...)
	}
	return areas
}
