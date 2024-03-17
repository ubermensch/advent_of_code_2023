package main

import (
	"bufio"
	"os"
	"path"
)

const inputFile = "/Users/frankhmeidan/golang/advent_of_code/day_7/input.txt"

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
}
