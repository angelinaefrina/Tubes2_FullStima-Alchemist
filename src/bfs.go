package main

import (
	"fmt"
	"sync"
)

// single-thread BFS
func bfs(
	target string, 
	recipes map[[2]string]string, 
	baseElements map[string]bool, 
	elementToTier map[string]int,
	) ([]string, []string, error) {
	// inisialisasi
	type Status struct {
		Elements map[string]bool // elemen yang dipunya sekarang
		Path     []string        // kombinasi resep yang sudah dicoba
	}

	initialState := Status{
		Elements: copySet(baseElements), 
		Path: []string{},
	}
	queue := []Status{initialState}
	visited := map[string]bool{
		stateToString(initialState.Elements): true,
	}

	discoveredElements := copySet(baseElements)
	triedCombos := make(map[[2]string]bool)
	var searchRoute []string

	// loop queue
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		// fmt.Println("Processing:", current.Path) // Debug

		if current.Elements[target] { // kalo udah ketemu
			return current.Path, searchRoute, nil
		}

		// coba semua kombinasi dari elemen yang ada
		elems := keys(current.Elements)

		for i := 0; i < len(elems); i++ {
			for j := i; j < len(elems); j++ {
				k := createKey(elems[i], elems[j])

				comboKey := [2]string{elems[i], elems[j]}
				if triedCombos[comboKey] {
					continue // already tried this combo, skip
				}
				triedCombos[comboKey] = true

				hasil, ok := recipes[k]
				if !ok {
					continue
				}
				
				if discoveredElements[hasil] {
					continue // already discovered this element, skip
				}
				tier1, ok1 := elementToTier[elems[i]]
				tier2, ok2 := elementToTier[elems[j]]
				tierHasil, okH := elementToTier[hasil]
				if !ok1 || !ok2 || !okH {
					continue
				}
				if tier1 > tierHasil || tier2 > tierHasil || tierHasil > elementToTier[target] {
					continue
				}
				
				discoveredElements[hasil] = true // mark as discovered
				combination := fmt.Sprintf("%s + %s = %s", elems[i], elems[j], hasil)
				searchRoute = append(searchRoute, combination)

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

	return nil, searchRoute, fmt.Errorf("No path found to create %s", target)
}


// multi-thread BFS
func bfsMultiplePaths(
	target string, 
	recipes map[[2]string]string, 
	baseElements map[string]bool, 
	elementToTier map[string]int, 
)	([][]string, []string, error) {
	// inisialisasi}
	type Status struct {
		Elements map[string]bool // elemen yang dipunya sekarang
		Path     []string        // kombinasi resep yang sudah dicobah
	}

	// startTime := time.Now()

	initial := Status{Elements: copySet(baseElements), Path: []string{}}

	visited := make(map[string]bool)
	var visitedMu sync.Mutex

	results := [][]string{}
	var resultsMu sync.Mutex

	searchRoute := []string{}
	var searchRouteMu sync.Mutex

	queue := []Status{initial}
	var queueMu sync.Mutex

	worker := func() {
		for {
			queueMu.Lock()
			if len(queue) == 0 {
				queueMu.Unlock()
				return
			}
			current := queue[0]
			queue = queue[1:]
			queueMu.Unlock()

			// fmt.Println("Processing:", current.Path) // Debug
			// visitedNodes.Add(1)

			if current.Elements[target] {
				resultsMu.Lock()
				results = append(results, current.Path)
				resultsMu.Unlock()
				continue
			}

			elems  := keys(current.Elements)
			triedCombos := make(map[[2]string]bool)

			for i := 0; i < len(elems); i++ {
				for j := i; j < len(elems); j++ {
					comboKey := [2]string{elems[i], elems[j]}
					if triedCombos[comboKey] {
						continue
					}
					triedCombos[comboKey] = true

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
					if tier1 > tierHasil || tier2 > tierHasil || tierHasil > elementToTier[target] {
						continue
					}

					recipeStr := fmt.Sprintf("%s + %s = %s", elems[i], elems[j], hasil)
					searchRouteMu.Lock()
					searchRoute = append(searchRoute, recipeStr)
					searchRouteMu.Unlock()

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

	wg.Wait()

	// elapsed := time.Since(startTime)

	if len(results) == 0 {
		return nil, searchRoute, fmt.Errorf("No path found to create %s", target)
	}

	return results, searchRoute, nil
}
