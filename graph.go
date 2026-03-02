package main

type Room struct {
	Name  string
	x, y  int
	Links []string
}

type Farm struct {
	AntsCount int
	Start     string
	End       string
	Rooms     map[string]*Room // Find Room with name
}
type Path struct {
	Rooms  []string // exp : []string{"##start", "3", "4", "##end"}
	Length int      // 3 -> 4 -> end = 3
}

type PathData struct {
	Rooms     []string
	AntsCount int
}

type Ant struct {
	ID   int
	Path []string
	Step int
}
