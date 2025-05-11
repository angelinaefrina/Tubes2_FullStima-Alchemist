package main

import (
	"encoding/json"
	"regexp"
	"strconv"

	// "fmt"
	"io"
	"os"
	"strings"
)

type Element struct {
	Name    string     `json:"name"`
	Recipes [][]string `json:"recipes"`
}

func parseJSON(filename string) (map[[2]string]string, map[string]bool, map[string]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, nil, err
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, nil, nil, err
	}

	var data map[string][]Element
	if err := json.Unmarshal(byteValue, &data); err != nil {
		return nil, nil, nil, err
	}

	recipes := make(map[[2]string]string)
	ingredients := make(map[string]bool)
	baseElements := make(map[string]bool)
	elementToTier := make(map[string]int)

	var tierRegex = regexp.MustCompile(`Tier (\d+)`)
	for tier, elements := range data {
		var tierNum int
		if tier == "Starting elements" {
			tierNum = 0
		} else {
			match := tierRegex.FindStringSubmatch(tier)
			if len(match) == 2 {
				tierNum, _ = strconv.Atoi(match[1])
			}
		}
		for _, elem := range elements {
			result := strings.TrimSpace(elem.Name)
			if result == "Time" {
				continue
			}
			elementToTier[result] = tierNum

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

	return recipes, baseElements, elementToTier, nil
}

// parser nya udah aku modify biar ga ada element Time termasuk resep yang ada Time nya
