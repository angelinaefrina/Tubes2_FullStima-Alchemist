package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func parseJSON(filename string) (map[[2]string]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	fmt.Println("Raw JSON Data:", string(byteValue)) // Debugging line

	var data map[string][]Element
	if err := json.Unmarshal(byteValue, &data); err != nil {
		return nil, err
	}
	fmt.Println("Unmarshalled Data:", data) // Debug

	recipes := make(map[[2]string]string)

	for tier, elements := range data {
		fmt.Println("Processing Tier:", tier) // Debug
		for _, elem := range elements {
			result := strings.TrimSpace(elem.Name)
			fmt.Println("Element Name:", result) // Debug
			for _, recipe := range elem.Recipes {
				if len(recipe) != 2 {
					fmt.Println("Skipping Malformed Recipe:", recipe) // Debug
					continue
				}
				key := createKey(recipe[0], recipe[1])
				fmt.Printf("Parsed Recipe: %v + %v = %v\n", recipe[0], recipe[1], result) // Debug
				recipes[key] = result
			}
		}
	}

	//fmt.Println("All Parsed Recipes:", recipes) // Debugging line
	return recipes, nil
}

// func main() {
// 	recipes, err := parseCSV("things.csv")
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		return
// 	}

// 	for k, v := range recipes {
// 		fmt.Printf("%v + %v = %v\n", k[0], k[1], v)
// 	}
// }
