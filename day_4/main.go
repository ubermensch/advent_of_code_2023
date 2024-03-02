package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

var cardChannel = make(chan int)

const inputFile = "/Users/frankhmeidan/golang/advent_of_code/day_4/input.txt"

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

// Returns the number of elements from `target` that are present in
// `src` slice.
func containsCount(src []int, target []int) int {
	count := 0
	for _, curr := range target {
		if slices.Contains(src, curr) {
			count += 1
		}
	}
	return count
}

func numsFromString(numStr string) []int {
	re := regexp.MustCompile("[0-9]+")
	var nums []int
	for _, num := range re.FindAllString(numStr, -1) {
		conv, err := strconv.Atoi(num)
		if err != nil {
			log.Fatal("could not parse numbers string: " + numStr)
			return nil
		}

		nums = append(nums, conv)
	}
	return nums
}

func calcCard(line string) {
	pieces := strings.Split(line, ":")
	results := pieces[1]
	drawnAndBet := strings.Split(results, "|")
	drawnNums, betNums := numsFromString(drawnAndBet[0]), numsFromString(drawnAndBet[1])
	matches := containsCount(drawnNums, betNums)

	// First match is worth 1 point.
	// Every subsequent match doubles the score.
	score := 0
	for i := 0; i < matches; i++ {
		if score == 0 {
			score = 1
		} else {
			score *= 2
		}
	}
	fmt.Print(fmt.Sprintf("%d --> | ", score))
	cardChannel <- score
}

// https://adventofcode.com/2023/day/4
func main() {
	scanner, file := fileScanner()
	defer file.Close()

	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	for _, line := range lines {
		go calcCard(line)
	}

	sum := 0
	for i := 0; i < len(lines); i++ {
		currVal := <-cardChannel
		fmt.Print(fmt.Sprintf("%d <-- | ", currVal))
		sum += currVal
	}

	fmt.Printf("\n\nSum is: %d\n", sum)
}
