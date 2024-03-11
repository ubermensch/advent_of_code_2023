package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/samber/lo"
	"os"
	"path"
	"slices"
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

// Traverses the set of maps given, finding the location number
// matching the given seed number
func seedLocation(seed int64, maps []*Map) (int64, error) {
	var location, currSource, currTarget int64
	var sourceType, targetType string

	currSource = seed
	sourceType = "seed"

	// traverse the maps, where target of current source
	// becomes next source. Stop when we have found the location value.
	for {
		currMap, ok := lo.Find(maps, func(m *Map) bool {
			return m.source == sourceType
		})
		if !ok {
			return 0, errors.New("map for source not found")
		}
		currTarget = currMap.Find(currSource)
		targetType = currMap.target

		if targetType == "location" {
			location = currTarget
			break
		} else {
			currSource = currTarget
			sourceType = targetType
		}
	}

	return location, nil
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

	var seedLocations []int64
	for _, currSeed := range seeds {
		currLocation, err := seedLocation(currSeed, maps)
		if err != nil {
			panic("could not find location for seed")
		}
		seedLocations = append(seedLocations, currLocation)
	}

	lowest := slices.Min(seedLocations)
	fmt.Println("Lowest location is: ", lowest)
}
