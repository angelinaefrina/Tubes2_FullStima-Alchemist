package main

import (
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
