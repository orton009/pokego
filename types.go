package main

import (
	"internal/pokecache"
)

type config struct {
	Next           string
	Previous       string
	cache          *pokecache.Cache
	caughtPokemons map[string]pokemon
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}
