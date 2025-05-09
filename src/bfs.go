package main

import (
	"fmt"
)

type Status struct {
	Elements map[string]bool // elemen yang dipunya sekarang
	Path     []string        // kombinasi resep yang sudah dicoba
}

func bfs(target string, recipes map[[2]string]string, baseElements(map[string]bool)) ([]string, error) {
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