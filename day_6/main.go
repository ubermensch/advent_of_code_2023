// https://adventofcode.com/2023/day/6
package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/samber/lo"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
)

const inputFile = "/Users/frankhmeidan/golang/advent_of_code/day_6/input.txt"

var wgGetStrategies sync.WaitGroup

type Strategy struct {
	// the number of millis to hold the button
	pressTime int
	// distance the boat will travel after releasing the button
	distance int
}

type Race struct {
	// the number of millis the race lasts
	time int
	// the distance to beat
	distance int
	// the possible strategies to finish the race in `time` millis
	strategies []*Strategy
}

// returns the possible strategies of Race `r` that beat the benchmark distance
// `r.distance`
func (r *Race) winningStrategies() []*Strategy {
	return lo.Filter(r.strategies, func(s *Strategy, _ int) bool {
		return s.distance > r.distance
	})
}

// calculates the possible strategies and writes them to `r.strategies`
func (r *Race) setStrategies() {
	var strategies []*Strategy

	for time := 0; time <= r.time; time++ {
		// when button is let go, boat will move at
		// 1 millimetre per millisecond of pressTime
		speed := time

		// will travel distance of (speed * time left)
		distance := (r.time - time) * speed

		strategies = append(strategies, &Strategy{
			pressTime: time,
			distance:  distance,
		})
	}

	r.strategies = strategies
}

func getRaces(times []int, distances []int) ([]*Race, error) {
	var races []*Race

	if len(times) != len(distances) {
		return nil, errors.New("mismatched times and distances")
	}

	// build the races, then concurrently calculate the possible strategies
	for i := 0; i < len(times); i++ {
		races = append(races, &Race{
			time:       times[i],
			distance:   distances[i],
			strategies: []*Strategy{},
		})
	}
	wgGetStrategies.Add(len(races))
	for _, race := range races {
		go func(r *Race) {
			r.setStrategies()
			wgGetStrategies.Done()
		}(race)
	}
	wgGetStrategies.Wait()

	return races, nil
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

func main() {
	scanner, file := fileScanner()
	defer file.Close()

	// get the int values from the input lines for time and distance
	getInts := func(line string) []int {
		strs := lo.Reject(strings.Split(line, " "), func(s string, i int) bool {
			return s == ""
		})
		return lo.Map(strs, func(curr string, i int) int {
			currInt, err := strconv.Atoi(curr)
			if err != nil {
				panic("could not turn string into int from input file")
			}
			return currInt
		})
	}

	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	times := getInts(strings.Split(lines[0], ":")[1])
	distances := getInts(strings.Split(lines[1], ":")[1])

	races, _ := getRaces(times, distances)
	solution := 1
	for i, race := range races {
		solution *= len(race.winningStrategies())
		log.Println(fmt.Sprintf("race %d has %d winning strategies", i+1, len(race.winningStrategies())))
	}

	log.Println("Multiplication of winning strategies count is ", solution)

}
