package main

import (
	"fmt"
	"sync"
)

func dfs(target string, recipes map[[2]string]string, baseElements map[string]bool) ([]string, error) {
	type State struct {
		Elements map[string]bool
		Path     []string
	}

	var dfsHelper func(State, map[string]bool) []string

	dfsHelper = func(current State, visited map[string]bool) []string {
		if current.Elements[target] {
			return current.Path
		}

		stateKey := stateToString(current.Elements)
		if visited[stateKey] {
			return nil
		}
		visited[stateKey] = true

		elements := keys(current.Elements)
		for i := 0; i < len(elements); i++ {
			for j := i; j < len(elements); j++ {
				key := createKey(elements[i], elements[j])
				result, ok := recipes[key]
				if !ok || current.Elements[result] {
					continue
				}

				newElements := copySet(current.Elements)
				newElements[result] = true

				newPath := append([]string{}, current.Path...)
				newPath = append(newPath, fmt.Sprintf("%s + %s => %s", elements[i], elements[j], result))

				res := dfsHelper(State{newElements, newPath}, visited)
				if res != nil {
					return res
				}
			}
		}
		return nil
	}

	// Start DFS from initial state
	start := State{
		Elements: copySet(baseElements),
		Path:     []string{},
	}
	visited := make(map[string]bool)
	result := dfsHelper(start, visited)
	if result == nil {
		return nil, fmt.Errorf("No path found to create %s", target)
	}
	return result, nil
}

func dfsMultiplePaths(
	target string,
	recipes map[[2]string]string,
	baseElements map[string]bool,
	amtOfMultiple int,
) ([][]string, error) {
	type State struct {
		Elements map[string]bool
		Path     []string
	}

	var results [][]string
	var mu sync.Mutex
	var wg sync.WaitGroup
	var foundCount int // Counter for the number of paths found

	var dfsHelper func(State, map[string]bool)
	dfsHelper = func(current State, visited map[string]bool) {
		defer wg.Done()

		// Stop searching if the desired number of paths has been found
		mu.Lock()
		if foundCount >= amtOfMultiple {
			mu.Unlock()
			return
		}
		mu.Unlock()

		if current.Elements[target] {
			mu.Lock()
			if foundCount < amtOfMultiple {
				results = append(results, current.Path)
				foundCount++
			}
			mu.Unlock()
			return
		}

		stateKey := stateToString(current.Elements)
		if visited[stateKey] {
			return
		}
		visited[stateKey] = true

		elements := keys(current.Elements)
		for i := 0; i < len(elements); i++ {
			for j := i; j < len(elements); j++ {
				key := createKey(elements[i], elements[j])
				if result, ok := recipes[key]; ok && !current.Elements[result] {
					newElements := copySet(current.Elements)
					newElements[result] = true

					newPath := append([]string{}, current.Path...)
					newPath = append(newPath, fmt.Sprintf("%s + %s => %s", elements[i], elements[j], result))

					wg.Add(1)
					go dfsHelper(State{newElements, newPath}, copyMap(visited))
				}
			}
		}
	}

	initial := State{Elements: copySet(baseElements), Path: []string{}}
	wg.Add(1)
	go dfsHelper(initial, make(map[string]bool))
	wg.Wait()

	if len(results) == 0 {
		return nil, fmt.Errorf("No path found to create %s", target)
	}
	return results, nil
}
