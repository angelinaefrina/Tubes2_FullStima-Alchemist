package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	var start time.Time

	// reader := bufio.NewReader(os.Stdin)
	// recipes, baseElements, elementToTier, err := parseJSON("recipe.json")
	// if err != nil {
	// 	fmt.Println("Parser error!", err)
	// }

	var multiple bool
	var amtOfMultiple int
	// // semua resep
	// fmt.Println("=== Parsed Recipes ===")
	// for pair, result := range recipes {
	// 	fmt.Printf("%s + %s = %s\n", pair[0], pair[1], result)
	// }

	// for el, tier := range elementToTier {
	// 	if tier == 0 {
	// 		baseElements[el] = true
	// 	}
	// }
	// // semua elemen
	// fmt.Println("=== All Elements ===")
	// for el, tier := range elementToTier {
	// 	if tier == 0 {
	// 		baseElements[el] = true
	// 	}
	// 	fmt.Printf("%s (%d)\n", el, tier)
	// }

	// // starting elements
	// fmt.Println("Base Elements:")
	// for el := range baseElements {
	// 	fmt.Println("-", el)
	// }
	// Step 1: Scrape and generate JSON and SVGs
	log.Println("Starting scraper...")
	scraper()
	log.Println("Scraping completed.")

	// Step 2: Parse the freshly created JSON
	log.Println("Parsing recipe.json...")
	recipes, baseElements, elementToTier, err := parseJSON("recipe.json")
	if err != nil {
		log.Fatalf("failed to parse JSON: %v", err)
	}
	log.Printf("Parsed %d recipes, %d base elements, %d elements with tier info.",
		len(recipes), len(baseElements), len(elementToTier))


	fmt.Println("Masukkan target elemen yang ingin dicari: ")
	reader := bufio.NewReader(os.Stdin)
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
			fmt.Println("=== Single Recipe ===")
			path, recipe, err := bfs(target, recipes, baseElements, elementToTier)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Recipe to create", target)
				for _, step := range recipe {
					fmt.Println(" ", step)
				}
				tree := pathToTreeBFS(recipe) 
				fmt.Println("Total nodes visited:", len(path))
				saveTreeToFile(tree, fmt.Sprintf("%s.json", target))
			}
		} else {
			// multiple paths for single target
			start = time.Now()
			fmt.Println("=== All Paths for Target ===")
			allRecipe, nodes, err := bfsMultiplePaths(target, recipes, baseElements, elementToTier)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				recipes, ok := allRecipe[target]
				if !ok || len(recipes) == 0 {
					fmt.Println("No recipe found for", target)
					return
				}
				fmt.Printf("Found %d recipe(s) to create %s:\n", len(recipes), target)

				for i, recipe := range recipes {
					if i >= amtOfMultiple {
						break
					}
					fmt.Printf("Recipe %d:\n", i+1)
					for _, step := range recipe {
						fmt.Println(" ", step.From1, "+", step.From2, "=>", step.Result)
					}

					tree := pathToTreeBFS(recipe)
					filename := fmt.Sprintf("%s_%d.json", target, i+1)
					saveTreeToFile(tree, filename)
				}
				fmt.Printf("Total nodes visited: %d\n", nodes)
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
				// nodesVisited += int(countNodes(tree))
				saveTreeToFile(tree, fmt.Sprintf("%s.json", target))
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
					// tree := pathToTree(allPaths[i])

					// jsonBytes, err := json.MarshalIndent(tree, "", "  ")
					// if err != nil {
					// 	fmt.Println("JSON error:", err)
					// 	return
					// }
					// fmt.Println(string(jsonBytes))
				}
				for i := range allPaths {
					tree := pathToTree(allPaths[i])
					// nodesVisited += int(countNodes(tree))
					saveTreeToFile(tree, fmt.Sprintf("%s_%d.json", target, i+1))
				}
			}
		}
	}

	// fmt.Printf("Total nodes yang dikunjungi: %d\n", nodesVisited)
	elapsed := time.Since(start)
	fmt.Printf("Waktu pencarian: %s\n", elapsed)
}
