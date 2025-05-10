package main

import (
	// "bufio"
	"fmt"
	// "os"
	// "strings"
)

func main() {
	recipes, baseElements, err := parseJSON("recipe.json")
	if err != nil {
		fmt.Println("Parser error!", err)
	}

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

	target := "Stone" // example target

	// path, err := bfs(target, recipes, baseElements)
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// } else {
	// 	fmt.Println("Path to create", target)
	// 	for _, step := range path {
	// 		fmt.Println(" ", step)
	// 	}
	// }
	
	paths := bfsMultiRecipe(target, recipes, baseElements)
	for i, p := range paths {
		fmt.Printf("Path %d to create %s:\n", i+1, target)
		for _, step := range p {
			fmt.Println(" ", step)
		}
	}
}
