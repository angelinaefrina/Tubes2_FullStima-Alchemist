package main

import (
	"encoding/json"
	// "fmt"
	"io"
	"os"
	"strings"
)

type Element struct {
	Name    string     `json:"name"`
	Recipes [][]string `json:"recipes"`
}

func parseJSON(filename string) (map[[2]string]string, map[string]bool, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, nil, err
	}

	var data map[string][]Element
	if err := json.Unmarshal(byteValue, &data); err != nil {
		return nil, nil, err
	}

	recipes := make(map[[2]string]string)
	ingredients := make(map[string]bool)
	results := make(map[string]bool)
	baseElements := make(map[string]bool)

	for tier, elements := range data {
		for _, elem := range elements {
			result := strings.TrimSpace(elem.Name)
			if result == "Time" {
				continue
			}
			results[result] = true

			if tier == "Starting elements" {
				baseElements[result] = true
			}
			for _, recipe := range elem.Recipes {
				if len(recipe) != 2 {
					continue
				}
				if recipe[0] == "Time" || recipe[1] == "Time" { // skip yg ada Time nya
					continue
				}
				key := createKey(recipe[0], recipe[1])
				recipes[key] = result
				ingredients[recipe[0]] = true
				ingredients[recipe[1]] = true
			}
		}
	}

	return recipes, baseElements, nil
}

// parser nya udah aku modify biar ga ada element Time termasuk resep yang ada Time nya