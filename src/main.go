package main

import (
	// "bufio"
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
	// "os"
	// "strings"
)

func main() {
	var start time.Time

	reader := bufio.NewReader(os.Stdin)
	recipes, baseElements, elementToTier, err := parseJSON("recipe.json")
	if err != nil {
		fmt.Println("Parser error!", err)
	}

	var multiple bool
	var amtOfMultiple int
	// semua resep
	fmt.Println("=== Parsed Recipes ===")
	for pair, result := range recipes {
		fmt.Printf("%s + %s = %s\n", pair[0], pair[1], result)
	}

	for el, tier := range elementToTier {
		if tier == 0 {
			baseElements[el] = true
		}
	}
	// semua elemen
	fmt.Println("=== All Elements ===")
	for el, tier := range elementToTier {
		if tier == 0 {
			baseElements[el] = true
		}
		fmt.Printf("%s (%d)\n", el, tier)
	}

	// starting elements
	fmt.Println("Base Elements:")
	for el := range baseElements {
		fmt.Println("-", el)
	}

	fmt.Println("Masukkan target elemen yang ingin dicari: ")
	// fmt.Scanln(&target)
	target, _ := reader.ReadString('\n')
	target = strings.TrimSpace(target)

	var method int
	fmt.Println("Pilih Metode (1. BFS, 2. DFS): ")
	fmt.Scanln(&method)
	fmt.Println("Berapa resep? (integer): ")
	fmt.Scanln(&amtOfMultiple)
	if amtOfMultiple == 1 {
		multiple = false
	} else if amtOfMultiple > 1 {
		multiple = true
	}

	if method == 1 {
		if !multiple {
			start = time.Now()
			path, err := bfs(target, recipes, baseElements)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Path to create", target)
				for _, step := range path {
					fmt.Println(" ", step)
				}
			}
		}
	} else if method == 2 {
		if !multiple {
			start = time.Now()
			fmt.Println("=== Single Recipe ===")
			path, err := dfs(target, recipes, baseElements, elementToTier)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Path to create", target)
				for _, step := range path {
					fmt.Println(" ", step)
				}
				tree := pathToTree(path) // convert the first path for visualization

				jsonBytes, err := json.MarshalIndent(tree, "", "  ")
				if err != nil {
					fmt.Println("JSON error:", err)
					return
				}
				fmt.Println(string(jsonBytes))
			}
		} else {
			// multiple paths for single target
			start = time.Now()
			fmt.Println("=== All Paths for Target ===")
			allPaths, err := dfsMultiplePaths(target, recipes, baseElements, amtOfMultiple, elementToTier)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Printf("Found %d path(s) to create %s:\n", len(allPaths), target)
				for i, path := range allPaths {
					fmt.Printf("Path %d:\n", i+1)
					for _, step := range path {
						fmt.Println("  ", step)
					}
					fmt.Println()
					tree := pathToTree(allPaths[i])

					jsonBytes, err := json.MarshalIndent(tree, "", "  ")
					if err != nil {
						fmt.Println("JSON error:", err)
						return
					}
					fmt.Println(string(jsonBytes))
				}
			}
		}
	}

	elapsed := time.Since(start)
	fmt.Printf("Execution time: %s\n", elapsed)
}
