package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
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
		var a, b, result string
		fmt.Sscanf(step, "%s + %s => %s", &a, &b, &result)

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

// ---------- BFS single path ----------

func bfs(target string, recipes map[[2]string]string, baseElements map[string]bool, _ map[string]int) ([]string, error) {
	type State struct {
		Elements map[string]bool
		Path     []string
	}

	initialState := State{Elements: copySet(baseElements), Path: []string{}}
	queue := []State{initialState}
	visited := map[string]bool{stateToString(initialState.Elements): true}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current.Elements[target] {
			return current.Path, nil
		}

		elems := keys(current.Elements)
		for i := 0; i < len(elems); i++ {
			for j := i; j < len(elems); j++ {
				k := createKey(elems[i], elems[j])
				result, ok := recipes[k]
				if !ok {
					continue
				}
				if current.Elements[result] {
					continue
				}

				newElements := copySet(current.Elements)
				newElements[result] = true
				newPath := append([]string{}, current.Path...)
				newPath = append(newPath, fmt.Sprintf("%s + %s => %s", elems[i], elems[j], result))

				stateKey := stateToString(newElements)
				if visited[stateKey] {
					continue
				}
				visited[stateKey] = true
				queue = append(queue, State{Elements: newElements, Path: newPath})
			}
		}
	}

	return nil, fmt.Errorf("no path found to create %s", target)
}

// ---------- DFS single path ----------

func dfs(target string, recipes map[[2]string]string, baseElements map[string]bool, _ map[string]int) ([]string, error) {
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

		elems := keys(current.Elements)
		for i := 0; i < len(elems); i++ {
			for j := i; j < len(elems); j++ {
				k := createKey(elems[i], elems[j])
				result, ok := recipes[k]
				if !ok || current.Elements[result] {
					continue
				}

				newElements := copySet(current.Elements)
				newElements[result] = true
				newPath := append([]string{}, current.Path...)
				newPath = append(newPath, fmt.Sprintf("%s + %s => %s", elems[i], elems[j], result))

				res := dfsHelper(State{newElements, newPath}, visited)
				if res != nil {
					return res
				}
			}
		}
		return nil
	}

	start := State{Elements: copySet(baseElements), Path: []string{}}
	visited := make(map[string]bool)
	result := dfsHelper(start, visited)
	if result == nil {
		return nil, fmt.Errorf("no path found to create %s", target)
	}
	return result, nil
}

// ---------- BFS multi path (concurrent) ----------

func bfsMultiplePaths(target string, recipes map[[2]string]string, baseElements map[string]bool, _ map[string]int) ([][]string, error) {
	type State struct {
		Elements map[string]bool
		Path     []string
	}

	var results [][]string
	var resultsMu sync.Mutex
	var wg sync.WaitGroup
	queue := make(chan State, 100)
	visited := sync.Map{}

	worker := func() {
		defer wg.Done()
		for current := range queue {
			if current.Elements[target] {
				resultsMu.Lock()
				results = append(results, current.Path)
				resultsMu.Unlock()
				continue
			}

			stateKey := stateToString(current.Elements)
			if _, loaded := visited.LoadOrStore(stateKey, true); loaded {
				continue
			}

			elems := keys(current.Elements)
			for i := 0; i < len(elems); i++ {
				for j := i; j < len(elems); j++ {
					k := createKey(elems[i], elems[j])
					result, ok := recipes[k]
					if !ok || current.Elements[result] {
						continue
					}

					newElements := copySet(current.Elements)
					newElements[result] = true
					newPath := append([]string{}, current.Path...)
					newPath = append(newPath, fmt.Sprintf("%s + %s => %s", elems[i], elems[j], result))

					queue <- State{Elements: newElements, Path: newPath}
				}
			}
		}
	}

	initial := State{Elements: copySet(baseElements), Path: []string{}}
	queue <- initial

	numWorkers := 8
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go worker()
	}

	go func() {
		wg.Wait()
		close(queue)
	}()

	wg.Wait()

	if len(results) == 0 {
		return nil, fmt.Errorf("no paths found to create %s", target)
	}
	return results, nil
}

// ---------- DFS multi path (concurrent) ----------

func dfsMultiplePaths(target string, recipes map[[2]string]string, baseElements map[string]bool, amt int, _ map[string]int) ([][]string, error) {
	type State struct {
		Elements map[string]bool
		Path     []string
	}

	var results [][]string
	var resultsMu sync.Mutex
	var wg sync.WaitGroup
	workQueue := make(chan State, 100)
	visited := sync.Map{}

	worker := func() {
		defer wg.Done()
		for current := range workQueue {
			if current.Elements[target] {
				resultsMu.Lock()
				if len(results) < amt {
					results = append(results, current.Path)
				}
				resultsMu.Unlock()
				continue
			}

			stateKey := stateToString(current.Elements)
			if _, loaded := visited.LoadOrStore(stateKey, true); loaded {
				continue
			}

			elems := keys(current.Elements)
			for i := 0; i < len(elems); i++ {
				for j := i; j < len(elems); j++ {
					k := createKey(elems[i], elems[j])
					result, ok := recipes[k]
					if !ok || current.Elements[result] {
						continue
					}

					newElements := copySet(current.Elements)
					newElements[result] = true
					newPath := append([]string{}, current.Path...)
					newPath = append(newPath, fmt.Sprintf("%s + %s => %s", elems[i], elems[j], result))

					resultsMu.Lock()
					if len(results) >= amt {
						resultsMu.Unlock()
						return
					}
					resultsMu.Unlock()

					workQueue <- State{Elements: newElements, Path: newPath}
				}
			}
		}
	}

	initial := State{Elements: copySet(baseElements), Path: []string{}}
	workQueue <- initial

	numWorkers := 8
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go worker()
	}

	go func() {
		wg.Wait()
		close(workQueue)
	}()

	wg.Wait()

	if len(results) == 0 {
		return nil, fmt.Errorf("no paths found to create %s", target)
	}
	return results, nil
}
