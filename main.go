package main

import (
	"bytes"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run . <filename>")
		return
	}

	fileName := os.Args[1]

	content, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println("ERROR: invalid data format")
		return
	}

	farm, err := ParseReader(bytes.NewReader(content))
	if err != nil {
		fmt.Println("ERROR: invalid data format")
		return
	}

	paths := FindAllPaths(farm)
	if len(paths) == 0 {
		fmt.Println("ERROR: invalid data format")
		return
	}

	fmt.Print(string(content))
	if len(content) > 0 && content[len(content)-1] != '\n' {
		fmt.Println()
	}
	fmt.Println()

	SimulateAndPrint(farm.AntsCount, paths, farm.End)
}
