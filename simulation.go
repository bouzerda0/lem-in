package main

import (
	"fmt"
	"strings"
)

func DistributeAnts(paths [][]string, totalAnts int) []PathData {
	var activePaths []PathData

	for _, p := range paths {
		activePaths = append(activePaths, PathData{Rooms: p, AntsCount: 0})
	}

	for i := 0; i < totalAnts; i++ {
		bestPathIdx := 0
		bestCost := len(activePaths[0].Rooms) + activePaths[0].AntsCount

		for j := 1; j < len(activePaths); j++ {
			cost := len(activePaths[j].Rooms) + activePaths[j].AntsCount
			if cost < bestCost {
				bestPathIdx = j
				bestCost = cost
			}
		}
		activePaths[bestPathIdx].AntsCount++
	}

	return activePaths
}

func SimulateAndPrint(pathsData []PathData, totalAnts int) {
	var movingAnts []*Ant
	nextAntID := 1

	for {
		var turnMoves []string

		for _, ant := range movingAnts {
			ant.Step++
			roomName := ant.Path[ant.Step]
			turnMoves = append(turnMoves, fmt.Sprintf("L%d-%s", ant.ID, roomName))
		}

		for i := range pathsData {
			if pathsData[i].AntsCount > 0 {
				newAnt := &Ant{
					ID:   nextAntID,
					Path: pathsData[i].Rooms,
					Step: 1,
				}
				movingAnts = append(movingAnts, newAnt)
				turnMoves = append(turnMoves, fmt.Sprintf("L%d-%s", newAnt.ID, newAnt.Path[newAnt.Step]))

				nextAntID++
				pathsData[i].AntsCount--
			}
		}

		var stillMoving []*Ant
		for _, ant := range movingAnts {
			if ant.Step < len(ant.Path)-1 {
				stillMoving = append(stillMoving, ant)
			}
		}
		movingAnts = stillMoving

		if len(turnMoves) == 0 {
			break
		}

		fmt.Println(strings.Join(turnMoves, " "))
	}
}
