package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Element struct {
	Name    string     `json:"name"`
	Recipes [][]string `json:"recipes"`
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	// Test Recipes
	// recipes := map[[2]string]string{
	// 	createKey("water", "earth"): "mud",
	// 	createKey("fire", "earth"):  "lava",
	// 	createKey("water", "fire"):  "steam",
	// 	createKey("air", "lava"):    "stone",
	// 	createKey("earth", "life"):  "human",
	// 	createKey("water", "stone"): "life",
	// }

	recipes, err := parseJSON("recipe.json")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var method int
	start := []string{"Water", "Fire", "Earth", "Air"}
	fmt.Println("Masukkan target elemen yang ingin dicari: ")
	// fmt.Scanln(&target)
	target, _ := reader.ReadString('\n')
	target = strings.TrimSpace(target)

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
