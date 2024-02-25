package schematic

import (
	"errors"
	"fmt"
	"github.com/samber/lo"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// PartNumber is the contiguous series of digits representing
// an integer that we're looking for in the schema.
type PartNumber struct {
	Number  int
	points  []*Point
	isValid bool
}
type Schematic struct {
	// The 2 dimensional grid of points represented by this schematic
	points [][]*Point

	// The part numbers hidden in this schematic, i.e. the contiguous
	// sequences of digits representing numbers.
	partNumbers []*PartNumber
}

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

func (pn *PartNumber) adjacentPoints() ([][]int, error) {
	if len(pn.points) == 0 {
		return nil, errors.New("empty part number (points not set)")
	}
	origin := pn.points[0]
	points := [][]int{
		{origin.x - 1, origin.y - 1},
		{origin.x - 1, origin.y},
		{origin.x - 1, origin.y + 1},
		{origin.x, origin.y - 1},
		{origin.x, origin.y + 1},
		{origin.x + 1, origin.y - 1},
		{origin.x + 1, origin.y},
		{origin.x + 1, origin.y + 1},
	}

	return points, nil
}

func (s *Schematic) getPoint(x int, y int) (*Point, error) {
	return s.points[x][y], nil
}

func (s *Schematic) setPartNumberValidity() error {
	if len(s.partNumbers) == 0 {
		return errors.New("no part numbers set on schematic")
	}
	for _, pn := range s.partNumbers {
		valid, err := s.isPartNumberValid(pn)
		if err != nil {
			return err
		}
		pn.isValid = valid
	}
	return nil
}

// Given part number, is it valid in this schematic? i.e. is the first
// point of this part number adjacent to a symbol point?
func (s *Schematic) isPartNumberValid(pn *PartNumber) (bool, error) {
	adjacent, err := pn.adjacentPoints()
	if err != nil {
		return false, err
	}
	hasAdjacentSymbol := lo.SomeBy(
		adjacent,
		func(point []int) bool {
			row, col := point[0], point[1]

			// Not valid if the adjacent point is outside of bounds
			if row >= len(s.points) || row < 0 {
				return false
			}
			if col >= len(s.points[row]) || col < 0 {
				return false
			}

			// Otherwise, part number is valid if the
			// adjacent point is a symbol
			if s.points[row][col].isSymbol() {
				return true
			}

			return false
		},
	)

	// Part number is valid if we have at least 1 adjacent symbol
	return hasAdjacentSymbol, nil
}

// Finds the part numbers hidden in the Schematic
func (s *Schematic) setPartNumbers() error {
	var parts []*PartNumber
	for x := 0; x < len(s.points); x++ {
		row := s.points[x]
		rowPartNumbers, err := partNumbersFromRow(row)
		if err != nil {
			return errors.New(
				fmt.Sprintf("could not get part numbers from row %d", x),
			)
		}
		parts = append(parts, rowPartNumbers...)
	}
	s.partNumbers = parts
	err := s.setPartNumberValidity()
	if err != nil {
		return err
	}

	return nil
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
			Number:  nInt,
			points:  row[index : index+len(num)],
			isValid: false,
		}

		partNumbers = append(partNumbers, pn)
	}

	return partNumbers, nil
}

func (s *Schematic) ValidPartNumbers() ([]*PartNumber, error) {
	if len(s.partNumbers) == 0 {
		return nil, errors.New("part numbers not yet determined for schematic")
	}

	return lo.Filter(
		s.partNumbers,
		func(pn *PartNumber, i int) bool {
			return pn.isValid
		},
	), nil
}

// Build a new Schematic from the input file
func NewSchematic() (*Schematic, error) {
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

	schematic := &Schematic{
		points:      points,
		partNumbers: []*PartNumber{},
	}

	err := schematic.setPartNumbers()
	if err != nil {
		return nil, err
	}

	return schematic, nil
}
