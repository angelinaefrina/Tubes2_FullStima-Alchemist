package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"sort"
	"strings"
)

type Status struct {
	Elements map[string]bool // elemen yang dipunya sekarang
	Path     []string        // kombinasi resep yang sudah dicoba
}

func bfsMultipleRecipe(
	target string, 
	recipes map[[2]string]string, 
	baseElements map[string]bool, 
	elementToTier map[string]int,
	amtOfMultiple int, 
)	([][]string, error) {

	initial := Status{Elements: copySet(baseElements), Path: []string{}}

	visited := make(map[string]bool)
	var visitedMu sync.Mutex

	results := [][]string{}
	var resultsMu sync.Mutex

	queue := []Status{initial}
	var queueMu sync.Mutex

	resultChan := make(chan []string, amtOfMultiple)
	var foundCount atomic.Int32

	worker := func() {
		for {
			queueMu.Lock()
			if len(queue) == 0 || foundCount.Load() >= int32(amtOfMultiple) {
				queueMu.Unlock()
				return
			}
			current := queue[0]
			queue = queue[1:]
			queueMu.Unlock()

			if current.Elements[target] {
				if foundCount.Add(1) <= int32(amtOfMultiple) {
					resultChan <- current.Path
				}
				continue
			}

			elems  := keys(current.Elements)
			for i := 0; i < len(elems); i++ {
				for j := i; j < len(elems); j++ {
					key := createKey(elems[i], elems[j])
					hasil, ok := recipes[key]
					if !ok || current.Elements[hasil] {
						continue
					}

					tier1, ok1 := elementToTier[elems[i]]
					tier2, ok2 := elementToTier[elems[j]]
					tierHasil, okH := elementToTier[hasil]
					if !ok1 || !ok2 || !okH {
						continue
					}
					if tier1 > tierHasil || tier2 > tierHasil {
						continue
					}

					newElements := copySet(current.Elements)
					newElements[hasil] = true
					newPath := append(append([]string{}, current.Path...), fmt.Sprintf("%s + %s => %s", elems[i], elems[j], hasil))
					newState := stateToString(newElements)

					visitedMu.Lock()
					if visited[newState] {
						visitedMu.Unlock()
						continue
					}
					visited[newState] = true
					visitedMu.Unlock()

					queueMu.Lock()
					queue = append(queue, Status{Elements: newElements, Path: newPath})
					queueMu.Unlock()
				}
			}
		}
	}

	// start workers
	numWorkers := 8
	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()
			worker()
		}()
	}

	// collect results
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for path := range resultChan {
		resultsMu.Lock()
		results = append(results, path)
		resultsMu.Unlock()
		if len(results) >= amtOfMultiple {
			break
		}
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("No path found to create %s", target)
	}

	return results, nil
}

func deduplicatePaths(target string, paths [][]string) [][]string {
	seen := make(map[string]bool)
	var uniquePaths [][]string

	for _, path := range paths {
		if len(path) == 0 {
			continue
		}
		lastStep := path[len(path)-1]

		parts := strings.Split(lastStep, " => ")
		if len(parts) != 2 {
			continue
		}
		ingredients := strings.Split(parts[0], " + ")
		if len(ingredients) != 2 {
			continue
		}

		sort.Strings(ingredients)
		key := strings.Join(ingredients, "+") + "=>" + parts[1]

		if !seen[key] {
			seen[key] = true
			uniquePaths = append(uniquePaths, path)
		}
	}
	return uniquePaths
}

// single-thread BFS
func bfs(target string, recipes map[[2]string]string, baseElements map[string]bool) ([]string, error) {
	// inisialisasi
	initialState := Status{
		Elements: copySet(baseElements), 
		Path: []string{},
	}
	queue := []Status{initialState}
	visited := map[string]bool{
		stateToString(initialState.Elements): true,
	}

	// loop queue
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current.Elements[target] { // kalo udah ketemu
			return current.Path, nil
		}

		// coba semua kombinasi dari elemen yang ada
		elems := keys(current.Elements)
		for i := 0; i < len(elems); i++ {
			for j := i + 1; j < len(elems); j++ {
				k := createKey(elems[i], elems[j])
				hasil, ok := recipes[k]
				if ok {
					if current.Elements[hasil] {
						continue
					}
				}
				newElements := copySet(current.Elements)
				newElements[hasil] = true
				newPath := append(append([]string{}, current.Path...), fmt.Sprintf("%s + %s => %s", elems[i], elems[j], hasil))
				
				stateKey := stateToString(newElements)
				if visited[stateKey] {
					continue
				}
				visited[stateKey] = true
				queue = append(queue, Status{
					Elements: newElements,
					Path:    newPath,
				})
			}
		}
	}

	return nil, fmt.Errorf("No path found to create %s", target)

}