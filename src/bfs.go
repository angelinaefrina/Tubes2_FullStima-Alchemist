package main

import (
	"fmt"
	"sync"
	"time"
)

type Status struct {
	Elements map[string]bool // elemen yang dipunya sekarang
	Path     []string        // kombinasi resep yang sudah dicoba
}

// multi-thread BFS
func bfsMultiRecipe(target string, recipes map[[2]string]string, baseElements map[string]bool) [][]string {
	var results [][]string
	var resultsMutex sync.Mutex

	visited := make(map[string]bool)
	var visitedMutex sync.Mutex

	queueChannel := make(chan Status, 1000)
	resultChannel := make(chan []string, 100)
	var waitGroup sync.WaitGroup

	// worker routine
	worker := func() {
		defer waitGroup.Done()
		for current := range queueChannel {

			if current.Elements[target] { // kalo udah ketemu
				resultsMutex.Lock()
				results = append(results, current.Path)
				resultsMutex.Unlock()
				// resultChannel <- current.Path
				continue
			}

			elems := keys(current.Elements)
			for i := 0; i < len(elems); i++ {
				for j := i + 1; j < len(elems); j++ {
					k := createKey(elems[i], elems[j])
					hasil, ok := recipes[k]
					if !ok || current.Elements[hasil] {
						continue
					}

					newElements := copySet(current.Elements)
					newElements[hasil] = true
					newPath := append(append([]string{}, current.Path...), fmt.Sprintf("%s + %s => %s", elems[i], elems[j], hasil))
					stateKey := stateToString(newElements)

					visitedMutex.Lock()
					if visited[stateKey] {
						visitedMutex.Unlock()
						continue
					}
					visited[stateKey] = true
					visitedMutex.Unlock()

					queueChannel <- Status{
						Elements: newElements,
						Path:	newPath,
					}
				}
			}
		}
	}

	// start workers
	numWorkers := 8
	waitGroup.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go worker()
	}

	// inisialisasi
	initialState := Status{
		Elements: copySet(baseElements), 
		Path: []string{},
	}
	queueChannel <- initialState

	// close result collector
	go func() {
		waitGroup.Wait()
		close(resultChannel)
	}()

	// close queue channel when input done
	go func() {
		time.Sleep(2 * time.Second)
		close(queueChannel)
	}()
	
	for path := range resultChannel {
		resultsMutex.Lock()
		results = append(results, path)
		resultsMutex.Unlock()
	}

	return results

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