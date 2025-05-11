package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)


// single-thread BFS
func bfs(
	target string,
	recipes map[[2]string]string,
	baseElements map[string]bool,
	elementsToTier map[string]int,
) ([]string, []RecipeStep, error) {

	// inisialisasi
	discovered := make(map[string]bool)
	queue := queue{}
	var searchPath []string
	recipeSteps := make(map[string]RecipeStep)

	targetTier := elementsToTier[target]

	for elem := range baseElements {
		discovered[elem] = true
		queue.enqueue(elem)
	}

	for !queue.isEmpty(){
		// fmt.Println("Queue:", queue) // Debug
		// fmt.Println("Discovered:", discovered) // Debug
		elem, _ := queue.dequeue() // dequeue elemen pertama dan diproses

		// fmt.Println("Processing:", elem) // Debug

		if elem == target {
			break
		}

		for current := range discovered {
			// fmt.Println("Discovered Elements:", current) // Debug

			recipeKey := createKey(elem, current)
			recipeResult, ok := recipes[recipeKey]
			if !ok || discovered[recipeResult] {
				continue
			}

			resultTier, ok := elementsToTier[recipeResult]
			if !ok || resultTier > targetTier {
				continue
			}

			if !discovered[recipeResult] {
				discovered[recipeResult] = true
				// fmt.Println("Discovered:", recipeResult) // Debug
				queue.enqueue(recipeResult) // enqueue elemen yang baru dibuat
				step := RecipeStep{elem, current, recipeResult}
				searchPath = append(searchPath, 
					fmt.Sprintf("%s + %s => %s", recipeKey[0], recipeKey[1], recipeResult)) // tambahin ke jalur penelusuran
				recipeSteps[recipeResult] = step

				if recipeResult == target {
					break
				}
			}
		}
	}
	if _, found := recipeSteps[target]; !found {
		return nil, nil, fmt.Errorf("Path to %s not found", target)
	}

	visited := make(map[string]bool)
	var steps []RecipeStep
	var collectSteps func(string)
	collectSteps = func(elem string) {
		if visited[elem] {
			return
		}
		step, ok := recipeSteps[elem]
		if !ok {
			return
		}
		visited[elem] = true
		collectSteps(step.From1)
		collectSteps(step.From2)
		steps = append(steps, step)
	}
	collectSteps(target)

	return searchPath, steps, nil
}

// multi-thread BFS
func bfsMultiplePaths(
	target string,
	recipes map[[2]string]string,
	baseElements map[string]bool,
	elementToTier map[string]int,
) (map[string][][]RecipeStep, int32, error) {

	type Path struct {
		steps []RecipeStep 
		last  string
		available map[string]bool
	}

	queue := make(chan Path, 1000)
	var mu sync.Mutex
	var wg sync.WaitGroup
	var visitedCount int32 = 0 

	targetTier := elementToTier[target]
	allPaths := make(map[string][][]RecipeStep)

	// inisialisasi available elements dengan base elements
	for elem := range baseElements {
		available := make(map[string]bool)
		for k := range baseElements {
			available[k] = true
		}
		queue <- Path{[]RecipeStep{}, elem, available}
	}

	// worker function untuk memproses queue
	worker := func() {
		defer wg.Done()
		for path := range queue {
			current := path.last
			available := path.available

			for other := range available {
				key := createKey(current, other)
				result, ok := recipes[key]
				if !ok || elementToTier[result] > targetTier {
					continue
				}

				newSteps := append([]RecipeStep{}, path.steps...)
				newSteps = append(newSteps, RecipeStep{current, other, result})

				newAvailable := make(map[string]bool)
				for k := range available {
					newAvailable[k] = true
				}
				newAvailable[result] = true

				mu.Lock()
				atomic.AddInt32(&visitedCount, 1) // menghitung node

				if result == target {
					alreadyExists := false
					for _, existing := range allPaths[target] {
						if equalSteps(existing, newSteps) { // cek apakah langkah sudah ada
							alreadyExists = true
							break
						}
					}
					if !alreadyExists {
						allPaths[target] = append(allPaths[target], newSteps)
					}
				} else {
					// lanjut cari meskipun sudah ditemukan
					queue <- Path{newSteps, result, newAvailable}
				}
				mu.Unlock()
			}
		}
	}
	
	// goroutine
	numWorkers := 4
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go worker()
	}

	go func() {
		wg.Wait()
		close(queue)
	}()

	// tunggu semua worker selesai
	time.Sleep(time.Second) // tunggu agar semua item diproses

	if len(allPaths[target]) == 0 {
		return nil, visitedCount, fmt.Errorf("no recipe found for %s", target)
	}
	return map[string][][]RecipeStep{target: allPaths[target]}, visitedCount, nil
}
