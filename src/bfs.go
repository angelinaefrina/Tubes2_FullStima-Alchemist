package main

type State struct {
	Elements map[string]bool
	Path     []string
}

func bfs(start []string, target string, recipes map[[2]string]string) []string {
	queue := []State{
		{Elements: makeSet(start), Path: []string{}},
	}
	visited := map[string]bool{}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current.Elements[target] {
			return current.Path
		}

		elements := keys(current.Elements)
		for i := 0; i < len(elements); i++ {
			for j := i + 1; j < len(elements); j++ {
				key := createKey(elements[i], elements[j])
				if result, exists := recipes[key]; exists && !current.Elements[result] {
					newElements := copySet(current.Elements)
					newElements[result] = true

					stateKey := stateToString(newElements)
					if visited[stateKey] {
						continue
					}
					visited[stateKey] = true

					newPath := append([]string{}, current.Path...)
					newPath = append(newPath, elements[i]+" + "+elements[j]+" => "+result)

					queue = append(queue, State{Elements: newElements, Path: newPath})
				}
			}
		}
	}
	return nil
}
