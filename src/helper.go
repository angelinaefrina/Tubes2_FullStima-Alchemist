package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

func createKey(a, b string) [2]string {
	if a > b {
		a, b = b, a
	}
	return [2]string{a, b}
}

func makeSet(arr []string) map[string]bool {
	set := map[string]bool{}
	for _, v := range arr {
		set[v] = true
	}
	return set
}

func copySet(original map[string]bool) map[string]bool {
	dup := make(map[string]bool)
	for k, v := range original {
		dup[k] = v
	}
	return dup
}

func keys(m map[string]bool) []string {
	k := make([]string, 0, len(m))
	for key := range m {
		k = append(k, key)
	}
	return k
}

func stateToString(set map[string]bool) string {
	keys := keys(set)
	sort.Strings(keys)
	return strings.Join(keys, ",")
}

func copyMap(original map[string]bool) map[string]bool {
	copied := make(map[string]bool)
	for k, v := range original {
		copied[k] = v
	}
	return copied
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
