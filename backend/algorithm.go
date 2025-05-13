package main

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
)

type Node struct {
	Element  string  `json:"element"`
	Recipe   string  `json:"recipe"`
	Children []*Node `json:"children"`
}

// ---------- Utilities ----------

func createKey(a, b string) [2]string {
	if a > b {
		a, b = b, a
	}
	return [2]string{a, b}
}

func keys(m map[string]bool) []string {
	k := make([]string, 0, len(m))
	for key := range m {
		k = append(k, key)
	}
	return k
}

func copySet(original map[string]bool) map[string]bool {
	dup := make(map[string]bool)
	for k, v := range original {
		dup[k] = v
	}
	return dup
}

func stateToString(set map[string]bool) string {
	keys := keys(set)
	sort.Strings(keys)
	return strings.Join(keys, ",")
}

func countNodes(node *Node) int32 {
	if node == nil {
		return 0
	}
	count := int32(1)
	for _, child := range node.Children {
		count += countNodes(child)
	}
	return count
}

func pathToTree(path []string) *Node {
	elementToNode := make(map[string]*Node)
	var root *Node

	for _, step := range path {
		parts := strings.Split(step, " => ")
		if len(parts) != 2 {
			continue
		}

		reactants := strings.Split(parts[0], " + ")
		if len(reactants) != 2 {
			continue
		}

		a, b, result := reactants[0], reactants[1], parts[1]

		left := elementToNode[a]
		if left == nil {
			left = &Node{Element: a, Recipe: a}
			elementToNode[a] = left
		}

		right := elementToNode[b]
		if right == nil {
			right = &Node{Element: b, Recipe: b}
			elementToNode[b] = right
		}

		node := &Node{
			Element:  result,
			Recipe:   fmt.Sprintf("%s + %s => %s", a, b, result),
			Children: []*Node{left, right},
		}

		elementToNode[result] = node
		root = node
	}

	return root
}

func buildRecipeMap(elements []ElementFromFandom) map[[2]string]string {
	recipeMap := make(map[[2]string]string)
	for _, el := range elements {
		for _, r := range el.Recipes {
			if len(r) == 2 {
				key := createKey(r[0], r[1])
				recipeMap[key] = el.Name
			}
		}
	}
	return recipeMap
}

func buildTierMap(elements []ElementFromFandom) map[string]int {
	tierMap := make(map[string]int)
	for _, el := range elements {
		tierMap[el.Name] = el.Tier
	}
	return tierMap
}

func printFoundPaths(pathMap map[string][][]string) {
	for target, paths := range pathMap {
		log.Printf("Found paths for: %s", target)
		for i, path := range paths {
			log.Printf("  Path %d:", i+1)
			for _, step := range path {
				log.Printf("    %s", step)
			}
		}
	}
}

func serializeTree(node *Node) string {
	if node == nil {
		return ""
	}
	if node.Children == nil || len(node.Children) == 0 {
		return node.Element
	}

	leftStr := serializeTree(node.Children[0])
	rightStr := serializeTree(node.Children[1])

	if leftStr > rightStr {
		leftStr, rightStr = rightStr, leftStr
	}

	return fmt.Sprintf("(%s + %s => %s)", leftStr, rightStr, node.Element)
}

// QUEUE FOR BFS
type queue []string

func (q *queue) isEmpty() bool {
	return len(*q) == 0
}
func (q *queue) enqueue(data string) {
	*q = append((*q), data)
}

func (q *queue) dequeue() (string, bool) {
	if q.isEmpty() {
		return "", false
	} else {
		// index := len(*q)
		// Print dequeued value
		dq := (*q)[0]
		// fmt.Printf("%s dequeued\n", dq)
		*q = (*q)[1:]
		return dq, true
	}
}

// ---------- BFS single path ----------

func bfs(target string, recipes map[[2]string]string, baseElements map[string]bool, elementToTier map[string]int) ([]string, int, error) {

	discovered := make(map[string]bool)
	queue := queue{}
	var searchPath []string
	targetTier := elementToTier[target]
	var nodeCount int32 = 0

	for elem := range baseElements {
		discovered[elem] = true
		queue.enqueue(elem)
	}

	for !queue.isEmpty() {
		elem, _ := queue.dequeue()
		atomic.AddInt32(&nodeCount, 1)
		if elem == target {
			return searchPath, int(nodeCount), nil
		}

		for current := range discovered {
			recipeKey := createKey(elem, current)
			result, ok := recipes[recipeKey]
			if !ok || discovered[result] {
				continue
			}

			resultTier, ok := elementToTier[result]
			if !ok || resultTier > targetTier {
				continue
			}

			if !discovered[result] {
				discovered[result] = true
				queue.enqueue(result)
				searchPath = append(searchPath, fmt.Sprintf("%s + %s => %s", recipeKey[0], recipeKey[1], result))

				if result == target {
					return searchPath, int(nodeCount), nil
				}
			}
		}
	}
	return nil, int(nodeCount), fmt.Errorf("no path found to create %s", target)
}

// ---------- DFS single path ----------

func dfs(target string, recipes map[[2]string]string, baseElements map[string]bool, elementToTier map[string]int) ([]string, int, error) {
	type State struct {
		Elements map[string]bool
		Path     []string
	}

	var nodeCount int32 = 0
	var dfsHelper func(State, map[string]bool) []string

	dfsHelper = func(current State, visited map[string]bool) []string {
		//fmt.Printf("DFS: %s", current.Path) // Debug
		atomic.AddInt32(&nodeCount, 1)
		if current.Elements[target] {
			return current.Path
		}

		stateKey := stateToString(current.Elements)
		if visited[stateKey] {
			return nil
		}
		visited[stateKey] = true

		elements := keys(current.Elements)
		//fmt.Printf("DFS: %s\n", elements)      // Debug
		//fmt.Printf("DFS: %d\n", len(elements)) // Debug
		sort.Strings(elements)
		for i := 0; i < len(elements); i++ {
			for j := i; j < len(elements); j++ {
				//fmt.Printf("DFS: %s\n", elements[j])         // Debug
				//fmt.Printf("%d", elementToTier[elements[i]]) // Debug
				//fmt.Printf("%d", elementToTier[elements[j]]) // Debug
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
				//fmt.Printf("DFS: %s + %s => %s\n", elements[i], elements[j], result) // Debug
				newPath = append(newPath, fmt.Sprintf("%s + %s => %s", elements[i], elements[j], result))

				//fmt.Printf(" LOOPDFS: %s\n", newPath) // Debug
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
		return nil, int(nodeCount), fmt.Errorf("no path found to create %s", target)
	}
	return result, int(nodeCount), nil
}

// ---------- BFS multi path (concurrent) ----------

func bfsMultiplePaths(
	target string,
	recipes map[[2]string]string,
	baseElements map[string]bool,
	amtOfMultiple int,
	elementToTier map[string]int,
) (map[string][][]string, int, error) {

	type Path struct {
		steps     []string
		last      string
		available map[string]bool
	}

	queue := make(chan Path, 1000)
	var mu sync.Mutex
	var wg sync.WaitGroup
	var visitedCount int32 = 0
	var stopSignal int32 = 0 // Atomic flag to signal workers to stop

	targetTier := elementToTier[target]
	allPaths := make(map[string][][]string)

	// Initialize available elements with base elements
	for elem := range baseElements {
		available := make(map[string]bool)
		for k := range baseElements {
			available[k] = true
		}
		queue <- Path{[]string{}, elem, available}
	}

	// Worker function to process the queue
	worker := func() {
		defer wg.Done()
		for path := range queue {
			if atomic.LoadInt32(&stopSignal) == 1 {
				return
			}

			current := path.last
			available := path.available

			for other := range available {
				key := createKey(current, other)
				result, ok := recipes[key]
				if !ok || elementToTier[result] > targetTier {
					continue
				}

				newSteps := append([]string{}, path.steps...)
				newSteps = append(newSteps, fmt.Sprintf("%s + %s => %s", current, other, result))

				newAvailable := make(map[string]bool)
				for k := range available {
					newAvailable[k] = true
				}
				newAvailable[result] = true

				mu.Lock()
				atomic.AddInt32(&visitedCount, 1)

				if result == target {
					alreadyExists := false
					for _, existing := range allPaths[target] {
						if equalStrings(existing, newSteps) {
							alreadyExists = true
							break
						}
					}

					// Cegah overfill karena race condition
					if !alreadyExists {
						if len(allPaths[target]) >= amtOfMultiple {
							mu.Unlock()
							return
						}

						allPaths[target] = append(allPaths[target], newSteps)
						log.Printf("Recipe found (%d/%d):\n", len(allPaths[target]), amtOfMultiple)
						for i, step := range newSteps {
							log.Printf("   %d. %s", i+1, step)
						}

						if len(allPaths[target]) >= amtOfMultiple {
							atomic.StoreInt32(&stopSignal, 1)
							mu.Unlock()
							return
						}
					}
				} else {
					queue <- Path{newSteps, result, newAvailable}
				}
				mu.Unlock()
			}
		}
	}

	// Start workers
	numWorkers := 4
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go worker()
	}

	go func() {
		wg.Wait()
		close(queue)
	}()

	wg.Wait()

	if len(allPaths[target]) == 0 {
		return nil, int(visitedCount), fmt.Errorf("no recipe found for %s", target)
	}

	paths := map[string][][]string{target: allPaths[target]}
	log.Printf("Found %d paths", len(paths[target]))
	printFoundPaths(paths)
	return paths, int(visitedCount), nil
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// ---------- DFS multi path (concurrent) ----------

func dfsMultiplePaths(
	target string,
	recipes map[[2]string]string,
	baseElements map[string]bool,
	amtOfMultiple int,
	elementToTier map[string]int,
) ([][]string, int, error) {
	seenTrees := make(map[string]bool)
	resultChan := make(chan []string, amtOfMultiple)
	doneChan := make(chan struct{})
	var results [][]string
	var nodeCount int32 = 0

	go func() {
		for path := range resultChan {
			results = append(results, path)
			// Check if the target is found
			log.Printf("Recipe found: %d steps (Total found: %d/%d)\n", len(path), len(results), amtOfMultiple)

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
			atomic.AddInt32(&nodeCount, 1)

			if current.Elements[target] {
				tree := pathToTree(current.Path)
				serialized := serializeTree(tree)

				if seenTrees[serialized] {
					continue
				}
				seenTrees[serialized] = true

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
		return nil, int(nodeCount), fmt.Errorf("no path found to create %s", target)
	}

	paths := map[string][][]string{target: results}
	log.Printf("Found %d paths", len(paths[target]))
	printFoundPaths(paths)
	return results, int(nodeCount), nil
}
