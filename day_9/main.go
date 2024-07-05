package main

import (
	"bufio"
	"day_9/predictor"
	"fmt"
	"github.com/samber/lo"
	"os"
	"path"
	"strconv"
	"strings"
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

	var allSeries []*predictor.Series
	var histories [][]int
	for scanner.Scan() {
		line := scanner.Text()
		history := lo.Map(strings.Split(line, " "), func(s string, idx int) int {
			i, err := strconv.Atoi(s)
			if err != nil {
				panic("could not parse line")
			}
			return i
		})
		histories = append(histories, history)
	}

	for _, curr := range histories {
		allSeries = append(allSeries, predictor.NewSeries(curr))
	}

	totalNext := lo.Reduce(allSeries, func(total int, s *predictor.Series, idx int) int {
		next := s.Next()
		return total + next
	}, 0)

	totalPrev := lo.Reduce(allSeries, func(total int, s *predictor.Series, idx int) int {
		previous := s.Previous()
		return total + previous
	}, 0)

	fmt.Println("Total of all next predictions: " + strconv.Itoa(totalNext))
	fmt.Println("Total of all prev predictions: " + strconv.Itoa(totalPrev))
}
