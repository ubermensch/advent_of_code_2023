package main

import (
	"bufio"
	"day_8/network"
	"fmt"
	"os"
	"path"
)

const inputFile = "/Users/frankhmeidan/golang/advent_of_code/day_8/input.txt"

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

	maze, err := network.NewNetwork(scanner)
	if err != nil {
		panic("could not build network")
	}
	steps := maze.StepsToFinish()
	fmt.Printf("steps to finish: %d\n", steps)
}
