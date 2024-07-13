package main

import (
	"bufio"
	"day_10/maze"
	"fmt"
	"os"
	"path"
)

const inputFile = "/Users/frankhmeidan/golang/advent_of_code/day_10/input.txt"

func fileScanner() (*bufio.Scanner, *os.File) {
	filePath := path.Join(inputFile)
	file, err := os.Open(filePath)
	if err != nil {
		panic("could not open file")
	}

	scanner := bufio.NewScanner(file)
	if err != nil {
		panic("could not scan file")
	}

	return scanner, file
}

func main() {
	scanner, file := fileScanner()
	defer file.Close()

	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	mz := maze.NewMaze(lines)
	fmt.Println(fmt.Sprintf("Furthest tile distance: %v", mz.FurthestTile().Distance()))
	fmt.Println(fmt.Sprintf("Area enclosed by loop: %v", mz.PathArea()))
}
