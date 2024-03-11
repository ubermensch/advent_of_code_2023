package main

import (
	"bufio"
	"fmt"
	"github.com/samber/lo"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

const inputFile = "/Users/frankhmeidan/golang/advent_of_code/day_5/input.txt"

type lookupRange struct {
	destRange   int64
	sourceRange int64
	rangeLength int64
}

type Map struct {
	// source and target are 'seed', 'soil', 'water', 'location' etc.
	source string
	target string
	// the partial lookup map
	lookup []*lookupRange
}

func (m *Map) Find(k int64) int64 {
	for _, currRange := range m.lookup {
		// if k is in the range...
		if k >= currRange.sourceRange && k <= currRange.sourceRange+currRange.rangeLength {
			// ...return the value in the destRange in the same position offset
			// from beginning of source range.
			return k - currRange.sourceRange + currRange.destRange
		}
	}

	return k
}

// Returns a pointer to a new `Map` struct
func NewMap(source string, target string, lookupRanges [][]int64) *Map {
	lookup := lo.Map(
		lookupRanges,
		func(curr []int64, i int) *lookupRange {
			return &lookupRange{
				destRange:   curr[0],
				sourceRange: curr[1],
				rangeLength: curr[2],
			}
		},
	)

	return &Map{
		source: source,
		target: target,
		lookup: lookup,
	}
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

func getSeeds(s string) []int64 {
	seedNums := strings.Split(strings.Split(s, ":")[1], " ")
	return lo.Map(seedNums[1:], func(sn string, i int) int64 {
		i, err := strconv.Atoi(sn)
		if err != nil {
			panic("cannot parse seeds string")
		}
		return int64(i)
	})
}

func parseMapScanner(scanner *bufio.Scanner, mapLine string) (*Map, error) {
	pieces := strings.Split(mapLine, " ")
	sourceAndTarget := strings.Split(pieces[0], "-")
	source, target := sourceAndTarget[0], sourceAndTarget[2]
	fmt.Println("source: ", source, " target: ", target)

	var lookupRanges [][]int64

	var line string
	for scanner.Scan() {
		line = scanner.Text()
		if line == "" {
			break
		}
		destTargetPieces := lo.Map(
			strings.Split(line, " "),
			func(s string, i int) int64 {
				currInt, err := strconv.Atoi(s)
				if err != nil {
					panic("could not convert dest or target string to int")
				}
				return int64(currInt)
			},
		)
		lookupRanges = append(lookupRanges, destTargetPieces)
	}

	return NewMap(source, target, lookupRanges), nil
}

func main() {
	var seeds []int64
	var maps []*Map

	scanner, file := fileScanner()
	defer file.Close()

	scanner.Scan()
	line := scanner.Text()
	seeds = getSeeds(line)
	for scanner.Scan() {
		line = scanner.Text()

		// gap between map definitions
		if line == "" {
			continue
		}

		// start of a new map definition
		if strings.Contains(line, "map:") {
			newMap, err := parseMapScanner(scanner, line)
			if err != nil {
				panic("could not create new map from scanner")
			}
			maps = append(maps, newMap)
		}
	}

	log.Println(seeds)
	log.Println(maps)
}
