package main

import (
	"sync"
)

type State struct {
	Elements map[string]bool // The current set of elements
	Path     []string        // The path taken to reach this state
}

func bfs(start []string, target string, recipes map[[2]string]string) []string {
	queue := []State{
		{Elements: makeSet(start), Path: []string{}},
	}
	visited := map[string]bool{}
	var mu sync.Mutex // Mutex to protect shared resources
	var wg sync.WaitGroup

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		// Debug
		// fmt.Println("Current State:", current)
		// fmt.Println("Visited States:", visited)

		if current.Elements[target] {
			return current.Path
		}

		elements := keys(current.Elements)
		for i := 0; i < len(elements); i++ {
			for j := i; j < len(elements); j++ { // Allow i == j for self-combination
				wg.Add(1) // Increment WaitGroup counter
				go func(i, j int) {
					defer wg.Done() // Decrement counter when goroutine finishes
					key := createKey(elements[i], elements[j])
					if result, exists := recipes[key]; exists && !current.Elements[result] {
						newElements := copySet(current.Elements)
						newElements[result] = true

						stateKey := stateToString(newElements)
						mu.Lock()
						if visited[stateKey] {
							mu.Unlock()
							return
						}
						visited[stateKey] = true
						mu.Unlock()

						newPath := append([]string{}, current.Path...)
						newPath = append(newPath, elements[i]+" + "+elements[j]+" => "+result)

						mu.Lock()
						queue = append(queue, State{Elements: newElements, Path: newPath})
						mu.Unlock()
					}
				}(i, j)
			}
		}
		wg.Wait() // Wait for all goroutines to finish
	}
	return nil
}
