package main

import (
	"fmt"
)

type Node struct {
	Element  string  `json:"element"`
	Recipe   string  `json:"recipe"`
	Children []*Node `json:"children"`
}

func dfs(target string, recipes map[[2]string]string, baseElements map[string]bool, elementToTier map[string]int) ([]string, error) {
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
				if elementToTier[elements[i]] >= elementToTier[target] || elementToTier[elements[j]] >= elementToTier[target] {
					continue
				}
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
	elementToTier map[string]int,
) ([][]string, error) {
	resultChan := make(chan []string, amtOfMultiple)
	doneChan := make(chan struct{})
	var results [][]string

	go func() {
		for path := range resultChan {
			results = append(results, path)
			fmt.Printf("Recipe found: %d steps (Total found: %d/%d)\n", // Debug
				len(path), len(results), amtOfMultiple)

			if len(results) >= amtOfMultiple {
				close(doneChan)
				return
			}
		}
	}()

	go func() {
		defer close(resultChan)

		type State struct {
			Elements map[string]bool
			Path     []string
		}

		stack := []State{{
			Elements: copySet(baseElements),
			Path:     []string{},
		}}

		visited := make(map[string]bool)

		for len(stack) > 0 {
			select {
			case <-doneChan:
				return
			default:
			}

			current := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			if current.Elements[target] {
				resultChan <- current.Path
				continue
			}

			stateKey := stateToString(current.Elements)
			if visited[stateKey] {
				continue
			}

			visited[stateKey] = true

			elements := keys(current.Elements)

			for i := len(elements) - 1; i >= 0; i-- {
				for j := len(elements) - 1; j >= i; j-- {
					if elementToTier[elements[i]] >= elementToTier[target] || elementToTier[elements[j]] >= elementToTier[target] {
						continue
					}
					key := createKey(elements[i], elements[j])
					result, ok := recipes[key]

					if !ok || current.Elements[result] {
						continue
					}

					newElements := copySet(current.Elements)
					newElements[result] = true

					newPath := append([]string{}, current.Path...)
					newPath = append(newPath, fmt.Sprintf("%s + %s => %s", elements[i], elements[j], result))

					// Push to stack
					stack = append(stack, State{
						Elements: newElements,
						Path:     newPath,
					})
				}
			}
		}
	}()

	<-doneChan

	if len(results) == 0 {
		return nil, fmt.Errorf("No path found to create %s", target)
	}

	return results, nil
}
