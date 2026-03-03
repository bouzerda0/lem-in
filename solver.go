package main
type edge struct{ from, to string }


type bfsState struct {
	node string
	side byte 
}

type parentInfo struct {
	node string
	side byte
}

func bfsResidual(farm *Farm, flow map[edge]bool) []string {
	// Build set of interior nodes that carry flow (capacity consumed).
	usedInterior := make(map[string]bool)
	for e, active := range flow {
		if active {
			if e.from != farm.Start && e.from != farm.End {
				usedInterior[e.from] = true
			}
			if e.to != farm.Start && e.to != farm.End {
				usedInterior[e.to] = true
			}
		}
	}

	// BFS with virtual states.
	start := bfsState{farm.Start, 0}
	visited := map[bfsState]bool{start: true}
	parent := map[bfsState]parentInfo{}
	queue := []bfsState{start}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		u := cur.node

		for _, v := range farm.Rooms[u].Links {
	
			canLeaveForward := cur.side == 0 || cur.side == 'o'
			canLeaveBackward := cur.side == 0 || cur.side == 'i'

			// Try forward edge u→v (no flow on it).
			if canLeaveForward && !flow[edge{u, v}] {
				var next bfsState
				if v == farm.End || v == farm.Start || !usedInterior[v] {
					next = bfsState{v, 0} // free node
				} else {
					next = bfsState{v, 'i'} // enter used-interior forward → in-side
				}
				if !visited[next] {
					visited[next] = true
					parent[next] = parentInfo{u, cur.side}
					if v == farm.End {
						return buildPathFromBFS(parent, start, next)
					}
					queue = append(queue, next)
				}
			}

			// Try backward edge: flow exists on v→u, cancel it.
			if canLeaveBackward && flow[edge{v, u}] {
				var next bfsState
				if v == farm.End || v == farm.Start || !usedInterior[v] {
					next = bfsState{v, 0}
				} else {
					next = bfsState{v, 'o'} // enter used-interior backward → out-side
				}
				if !visited[next] {
					visited[next] = true
					parent[next] = parentInfo{u, cur.side}
					if v == farm.End {
						return buildPathFromBFS(parent, start, next)
					}
					queue = append(queue, next)
				}
			}
		}
	}
	return nil
}

func buildPathFromBFS(parent map[bfsState]parentInfo, start, end bfsState) []string {
	var path []string
	for s := end; s != start; {
		path = append(path, s.node)
		p := parent[s]
		s = bfsState{p.node, p.side}
	}
	path = append(path, start.node)
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}


func augment(flow map[edge]bool, path []string) {
	for i := 0; i < len(path)-1; i++ {
		u, v := path[i], path[i+1]
		rev := edge{v, u}
		if flow[rev] {
			delete(flow, rev) // edge cancellation
		} else {
			flow[edge{u, v}] = true
		}
	}
}

func extractPaths(farm *Farm, flow map[edge]bool) [][]string {
	next := make(map[string][]string)
	for e, active := range flow {
		if active {
			next[e.from] = append(next[e.from], e.to)
		}
	}

	var paths [][]string
	for _, first := range next[farm.Start] {
		path := []string{farm.Start}
		cur := first
		for cur != farm.End {
			path = append(path, cur)
			nbs := next[cur]
			if len(nbs) == 0 {
				break
			}
			cur = nbs[0]
			next[path[len(path)-1]] = nbs[1:]
		}
		path = append(path, farm.End)
		paths = append(paths, path)
	}
	return paths
}

func turnsNeeded(ants int, paths [][]string) int {
	k := len(paths)
	if k == 0 {
		return int(^uint(0) >> 1)
	}
	lengths := make([]int, k)
	for i, p := range paths {
		lengths[i] = len(p) - 1
	}
	// Sort ascending (insertion sort — k is tiny).
	for i := 1; i < k; i++ {
		for j := i; j > 0 && lengths[j] < lengths[j-1]; j-- {
			lengths[j], lengths[j-1] = lengths[j-1], lengths[j]
		}
	}
	for turns := lengths[k-1]; ; turns++ {
		cap := 0
		for _, l := range lengths {
			if c := turns - l; c > 0 {
				cap += c
			}
		}
		if cap >= ants {
			return turns
		}
	}
}

func FindAllPaths(farm *Farm) [][]string {
	flow := make(map[edge]bool)
	var bestPaths [][]string
	bestTurns := int(^uint(0) >> 1)

	for {
		aug := bfsResidual(farm, flow)
		if aug == nil {
			break
		}
		augment(flow, aug)
		paths := extractPaths(farm, flow)
		turns := turnsNeeded(farm.AntsCount, paths)

		if turns < bestTurns {
			bestTurns = turns
			bestPaths = paths
		} else {
			break
		}
	}
	return bestPaths
}
