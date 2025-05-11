package main

import (
	"fmt"
	"sync"
	"time"
)

type Node struct {
	Element  string  `json:"element"`
	Recipe   string  `json:"recipe"`
	Children []*Node `json:"children"`
}

// single-thread BFS
func bfs(target string, recipes map[[2]string]string, baseElements map[string]bool) ([]string, time.Duration, error) {
	// inisialisasi
	type Status struct {
		Elements map[string]bool // elemen yang dipunya sekarang
		Path     []string        // kombinasi resep yang sudah dicoba
	}

	startTime := time.Now()

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

		// fmt.Println("Processing:", current.Path) // Debug

		if current.Elements[target] { // kalo udah ketemu
			elapsed := time.Since(startTime)
			return current.Path, elapsed, nil
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
	elapsed := time.Since(startTime)

	return nil, elapsed, fmt.Errorf("No path found to create %s", target)

}

// multi-thread BFS
func bfsMultipleRecipe(
	target string, 
	recipes map[[2]string]string, 
	baseElements map[string]bool, 
	elementToTier map[string]int, 
)	([][]string, time.Duration, error) {
	type Status struct {
		Elements map[string]bool // elemen yang dipunya sekarang
		Path     []string        // kombinasi resep yang sudah dicoba
		Tree 	*Node          // tree untuk menyimpan langkah-langkah
	}

	startTime := time.Now()

	root := &Node{Element: "ROOT"}
	initial := Status{Elements: copySet(baseElements), Path: []string{}, Tree: root}

	visited := make(map[string]bool)
	var visitedMu sync.Mutex

	results := [][]string{}
	var resultsMu sync.Mutex

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
					if tier1 > tierHasil || tier2 > tierHasil || tierHasil > elementToTier[target] {
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

					newNode := &Node{
						Element:  hasil,
						Recipe:   fmt.Sprintf("%s + %s => %s", elems[i], elems[j], hasil),
						Children: []*Node{},
					}
					current.Tree.Children = append(current.Tree.Children, newNode)

					queueMu.Lock()
					queue = append(queue, Status{Elements: newElements, Path: newPath, Tree: newNode})
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

	elapsed := time.Since(startTime)

	if len(results) == 0 {
		return nil, elapsed, fmt.Errorf("No path found to create %s", target)
	}

	return results, elapsed, nil
}
