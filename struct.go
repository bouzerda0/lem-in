package main

// Room represents a room in the ant farm.
type Room struct {
	Name  string
	X, Y  int
	Links []string // connected room names
}

// Farm holds the parsed ant farm data.
type Farm struct {
	AntsCount int
	Start     string
	End       string
	Rooms     map[string]*Room
}

// edge is a directed edge between two rooms.
type edge struct{ from, to string }

// bfsState tracks a node + virtual side for node-splitting BFS.
// side: 0=free, 'i'=in-half, 'o'=out-half of a used node.
type bfsState struct {
	node string
	side byte
}

// parentInfo stores the BFS predecessor for path reconstruction.
type parentInfo struct {
	node string
	side byte
}
