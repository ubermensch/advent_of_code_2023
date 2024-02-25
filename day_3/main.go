package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"unicode"
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

func (p *Point) isDot() bool {
	return p.value == '.'
}
func (p *Point) isDigit() bool {
	return unicode.IsDigit(p.value)
}
func (p *Point) isSymbol() bool {
	return unicode.IsSymbol(p.value)
}

// The full schematic document
type Schematic struct {
	points [][]*Point
}

func (s *Schematic) getPoint(x int, y int) (*Point, error) {
	return s.points[x][y], nil
}

type PartNumber struct {
	number  int
	points  []*Point
	isValid bool
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

func partNumbersFromRow(row []*Point) ([]*PartNumber, error) {
	rowStr := ""
	for _, curr := range row {
		rowStr += string(curr.value)
	}
	re := regexp.MustCompile("\\D+")
	splits := re.Split(rowStr, -1)
	var numbers []string
	for _, curr := range splits {
		if len(curr) > 0 {
			numbers = append(numbers, curr)
		}
	}

	var partNumbers []*PartNumber
	for _, num := range numbers {
		index := strings.Index(rowStr, num)
		nInt, err := strconv.Atoi(num)
		if err != nil {
			return nil, errors.New("could not read number")
		}

		pn := &PartNumber{
			number:  nInt,
			points:  row[index : index+len(num)],
			isValid: false,
		}

		partNumbers = append(partNumbers, pn)
	}

	return partNumbers, nil
}

// Finds the part numbers hidden in the Schematic
func (s *Schematic) getPartNumbers() ([]*PartNumber, error) {
	var parts []*PartNumber
	for x := 0; x < len(s.points); x++ {
		row := s.points[x]
		rowPartNumbers, err := partNumbersFromRow(row)
		if err != nil {
			return nil, errors.New(
				fmt.Sprintf("could not get part numbers from row %d", x),
			)
		}
		parts = append(parts, rowPartNumbers...)
	}
	return parts, nil
}

// Build a new Schematic from the input file
func newSchematic() (*Schematic, error) {
	scanner, file := fileScanner()
	defer file.Close()

	var points [][]*Point
	y := 0

	for scanner.Scan() {
		var linePoints []*Point
		line := scanner.Text()
		for x, char := range line {
			linePoints = append(linePoints, &Point{x: x, y: y, value: char})
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
	_, err := schematic.getPartNumbers()
	if err != nil {
		log.Fatal(err)
	}
}
