package schematic

import (
	"bufio"
	"os"
	"path"
)

const inputFile = "/Users/frankhmeidan/golang/advent_of_code/day_3/input.txt"

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
