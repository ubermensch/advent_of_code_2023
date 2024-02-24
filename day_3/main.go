package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
)

const inputFile = "/Users/frankhmeidan/golang/advent_of_code/day_3/input.txt"

// A point in the schematic grid. The value for this point
// could be either:
// * A digit (part of a number)
// * A symbol (non-`.`, not letter, not digit).
type Point struct {
	x     int  // x-position in the Schematic line (0-indexed)
	y     int  // y-position in the Schematic, the line number (0-indexed)
	value rune // the value of this point
}

// The full schematic document
type Schematic struct {
	points [][]Point
}

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

func newSchematic() (*Schematic, error) {
	scanner, file := fileScanner()
	defer file.Close()

	var points [][]Point
	y := 0

	for scanner.Scan() {
		var linePoints []Point
		line := scanner.Text()
		for char, x := range line {
			linePoints = append(linePoints, Point{x: int(x), y: y, value: rune(char)})
		}

		points = append(points, linePoints)
		y += 1
	}

	return &Schematic{
		points: points,
	}, nil
}

func main() {
	schematic, _ := newSchematic()
	fmt.Printf("Schematic : %v", schematic.points)
}
