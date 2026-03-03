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
	Rooms     map[string]*Room
}
