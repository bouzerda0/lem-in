package main

func BFS(farm *Farm, forbidden map[string]bool) []string {
	queue := []string{farm.Start}

	visited := make(map[string]bool)
	parent := make(map[string]string)

	for room := range forbidden {
		visited[room] = true
	}
	visited[farm.Start] = true

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if curr == farm.End {
			var path []string
			step := farm.End
			
			for step != farm.Start {
				path = append(path, step)
				step = parent[step]
			}
			path = append(path, farm.Start)

			for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
				path[i], path[j] = path[j], path[i]
			}
			return path
		}

		for _, neighbor := range farm.Rooms[curr].Links {
			if !visited[neighbor] {
				visited[neighbor] = true
				parent[neighbor] = curr 
				queue = append(queue, neighbor)
			}
		}
	}
	return nil
}

func FindAllPaths(farm *Farm) [][]string {
	var paths [][]string
	forbidden := make(map[string]bool)

	for {
		path := BFS(farm, forbidden)
		if path == nil {
			break
		}
		paths = append(paths, path)

		for i := 1; i < len(path)-1; i++ {
			forbidden[path[i]] = true
		}
	}
	return paths
}
