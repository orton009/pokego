package main

import "strings"

func cleanInput(text string) []string {
	return strings.Split(
		strings.ToLower(
			strings.TrimSpace(text)),
		" ")
}
