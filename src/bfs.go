package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"sort"
	"strings"
	"time"
)

type Node struct {
	Element  string  `json:"element"`
	Recipe   string  `json:"recipe"`
	Children []*Node `json:"children"`
}

// single-thread BFS
func bfs(target string, recipes map[[2]string]string, baseElements map[string]bool) ([]string, int32, time.Duration, error) {
	// inisialisasi
	type Status struct {
		Elements map[string]bool // elemen yang dipunya sekarang
		Path     []string        // kombinasi resep yang sudah dicoba
	}

	startTime := time.Now()
	var visitedNodes atomic.Int32

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

		fmt.Println("Processing:", current.Path) // Debug

		if current.Elements[target] { // kalo udah ketemu
			elapsed := time.Since(startTime)
			return current.Path, visitedNodes.Load(), elapsed, nil
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
				visitedNodes.Add(1) // kalo ada elemen yang terbuat, tambah node
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

	return nil, visitedNodes.Load(), elapsed, fmt.Errorf("No path found to create %s", target)

}

// multi-thread BFS
func bfsMultipleRecipe(
	target string, 
	recipes map[[2]string]string, 
	baseElements map[string]bool, 
	elementToTier map[string]int, 
)	([][]string, int32, time.Duration, error) {
	type Status struct {
		Elements map[string]bool // elemen yang dipunya sekarang
		Path     []string        // kombinasi resep yang sudah dicoba
		Tree 	*Node          // tree untuk menyimpan langkah-langkah
	}

	startTime := time.Now()
	var visitedNodes atomic.Int32

	root := &Node{Element: "ROOT"}
	initial := Status{Elements: copySet(baseElements), Path: []string{}, Tree: root}

	visited := make(map[string]bool)
	var visitedMu sync.Mutex

	results := [][]string{}
	var resultsMu sync.Mutex

	queue := []Status{initial}
	var queueMu sync.Mutex

	// resultChan := make(chan []string, amtOfMultiple)
	// var foundCount atomic.Int32

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

			// fmt.Println("Current Path:", current.Path) // Debug
			visitedNodes.Add(1)

			if current.Elements[target] {
				// if foundCount.Add(1) <= int32(amtOfMultiple) {
				// 	resultChan <- current.Path
				// }
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
		return nil, visitedNodes.Load(), elapsed, fmt.Errorf("No path found to create %s", target)
	}

	return results, visitedNodes.Load(), elapsed, nil
}

func deduplicatePaths(target string, paths [][]string, amtOfMultiple int) [][]string {
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
			if (amtOfMultiple > 0) && (len(uniquePaths) >= amtOfMultiple) {
				break
			}
		}
	}
	return uniquePaths
}

func saveTreeToFile(tree *Node, filename string) {
	err := func() error {
		data, err := json.MarshalIndent(tree, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal tree: %w", err)
		}
		if err := os.WriteFile(filename, data, 0644); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}
		return nil
	}()
	if err != nil {
		fmt.Printf("Error saving tree to file: %v\n", err)
	} else {
		fmt.Printf("Tree successfully saved to %s\n", filename)
	}
}

func buildTreeFromPath(path []string) *Node {
	nodes := make(map[string]*Node)

	for _, step := range path {
		parts := strings.Split(step, " => ")
		if len(parts) != 2 {
			continue
		}
		ingredients := strings.Split(parts[0], " + ")
		if len(ingredients) != 2 {
			continue
		}
		result := parts[1]

		// pastikan node bahan dibuat dulu
		for _, ing := range ingredients {
			if _, ok := nodes[ing]; !ok {
				nodes[ing] = &Node{
					Element: ing,
					Recipe:  ing, // base element
				}
			}
		}

		// buat node hasil
		nodes[result] = &Node{
			Element:  result,
			Recipe:   step,
			Children: []*Node{nodes[ingredients[0]], nodes[ingredients[1]]},
		}
	}

	// asumsi hasil akhir path adalah target
	lastStep := path[len(path)-1]
	result := strings.Split(lastStep, " => ")[1]
	return nodes[result]
}
