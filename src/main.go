package main

import (
	// "bufio"
	"bufio"
	"fmt"
	"os"
	"strings"
	// "os"
	// "strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	recipes, baseElements, err := parseJSON("recipe.json")
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
			fmt.Println("=== Single Recipe ===")
			path, err := dfs(target, recipes, baseElements)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Path to create", target)
				for _, step := range path {
					fmt.Println(" ", step)
				}
			}
		} else {
			// multiple paths for single target
			fmt.Println("=== All Paths for Target ===")
			allPaths, err := dfsMultiplePaths(target, recipes, baseElements, amtOfMultiple)
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
				}
			}
		}
	}
}
