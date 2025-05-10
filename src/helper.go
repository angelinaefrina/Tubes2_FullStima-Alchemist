package main

import (
	"fmt"
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
