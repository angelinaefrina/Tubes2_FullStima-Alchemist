package main

func dfs(current State, target string, recipes map[[2]string]string, visited map[string]bool) []string {
	if current.Elements[target] {
		return current.Path
	}

	stateKey := stateToString(current.Elements)
	if visited[stateKey] {
		return nil
	}
	visited[stateKey] = true

	elements := keys(current.Elements)
	for i := 0; i < len(elements); i++ {
		for j := i + 1; j < len(elements); j++ {
			key := createKey(elements[i], elements[j])
			if result, exists := recipes[key]; exists && !current.Elements[result] {
				newElements := copySet(current.Elements)
				newElements[result] = true

				newPath := append([]string{}, current.Path...)
				newPath = append(newPath, elements[i]+" + "+elements[j]+" => "+result)

				res := dfs(State{newElements, newPath}, target, recipes, visited)
				if res != nil {
					return res
				}
			}
		}
	}
	return nil
}
