package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
)

const inputFile = "/Users/frankhmeidan/golang/advent_of_code/day_9/input.txt"

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

	maze, err := predictor.NewDiffCalculator()
	if err != nil {
		panic("could not build network")
	}
	steps := maze.StepsToFinish()
	bothSteps := maze.LcmStepsToFinish()
	fmt.Printf("steps to finish: %d, steps to finish in parallel: %d\n", steps, bothSteps)
}
