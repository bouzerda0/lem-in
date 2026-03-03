package main

import (
	"fmt"
	"sort"
	"strings"
)

func SimulateAndPrint(antsCount int, paths [][]string, endRoom string) {
	if len(paths) == 0 {
		return
	}

	// Sort paths by length (shortest first) for optimal distribution.
	sort.Slice(paths, func(i, j int) bool {
		return len(paths[i]) < len(paths[j])
	})

	// Distribute ants greedily: assign each ant to the path that currently
	// yields the lowest completion time (pathLen - 1 + antsOnPath).
	assigned := make([]int, len(paths)) // ants assigned per path
	for i := 0; i < antsCount; i++ {
		best := 0
		bestCost := (len(paths[0]) - 1) + assigned[0]
		for j := 1; j < len(paths); j++ {
			cost := (len(paths[j]) - 1) + assigned[j]
			if cost < bestCost {
				best = j
				bestCost = cost
			}
		}
		assigned[best]++
	}

	// Simulate turn by turn.
	type ant struct {
		id   int
		path []string
		step int // current index in path (0 = start)
	}

	var moving []*ant
	nextID := 1

	for {
		var moves []string

		// Advance all moving ants one step.
		for _, a := range moving {
			a.step++
			moves = append(moves, fmt.Sprintf("L%d-%s", a.id, a.path[a.step]))
		}

		// Launch new ants from each path that still has capacity.
		for i := range paths {
			if assigned[i] > 0 {
				a := &ant{id: nextID, path: paths[i], step: 1}
				moving = append(moving, a)
				moves = append(moves, fmt.Sprintf("L%d-%s", a.id, a.path[1]))
				nextID++
				assigned[i]--
			}
		}

		// Remove ants that reached the end.
		var stillMoving []*ant
		for _, a := range moving {
			if a.step < len(a.path)-1 {
				stillMoving = append(stillMoving, a)
			}
		}
		moving = stillMoving

		if len(moves) == 0 {
			break
		}
		fmt.Println(strings.Join(moves, " "))
	}
}
