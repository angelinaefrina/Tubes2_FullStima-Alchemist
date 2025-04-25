package main

import "fmt"

func main() {
	// Test Recipes
	recipes := map[[2]string]string{
		createKey("water", "earth"): "mud",
		createKey("fire", "earth"):  "lava",
		createKey("water", "fire"):  "steam",
		createKey("air", "lava"):    "stone",
		createKey("earth", "life"):  "human",
		createKey("water", "stone"): "life",
	}

	// recipes, err := parseCSV("things.csv")
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }

	fmt.Println("Recipes:", recipes)
	var method int
	var target string
	start := []string{"water", "fire", "earth", "air"}
	fmt.Println("Masukkan target elemen yang ingin dicari: ")
	fmt.Scanln(&target)

	fmt.Println("Pilih Metode (1. BFS, 2. DFS): ")
	fmt.Scanln(&method)
	if method == 1 {
		result := bfs(start, target, recipes)
		if result == nil {
			println("No recipe found for", target)
		} else {
			for _, step := range result {
				println(" -", step)
			}
		}
	} else if method == 2 {
		result := dfs(State{makeSet(start), []string{}}, target, recipes, make(map[string]bool))
		if result == nil {
			println("No recipe found for", target)
		} else {
			for _, step := range result {
				println(" -", step)
			}
		}
	} else {
		return
	}
}
