package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"strconv"
	"sync"
)

type Row string

const inputFile = "/Users/frankhmeidan/golang/advent_of_code/day_1/input.txt"

var wg sync.WaitGroup
var rowChannel = make(chan int)

// Given a Row, find the int produced by concatenating the
// first and last digit (could be the same digit).
func calcRow(row Row) {
	defer wg.Done()

	re := regexp.MustCompile("[0-9]")
	var digits []int
	for _, i := range re.FindAllString(string(row), -1) {
		conv, err := strconv.Atoi(i)
		if err != nil {
			log.Fatal("could not parse line " + row)
			return
		}

		digits = append(digits, conv)
	}
	d1, d2 := digits[0], digits[0]
	if len(digits) > 1 {
		d2 = digits[len(digits)-1]
	}

	final, err := strconv.Atoi(fmt.Sprintf("%d%d", d1, d2))
	if err != nil {
		log.Fatal("could not parse line " + row)
		return
	}

	fmt.Println(" > line scanned: " + string(row))

	// TODO Work out the first and last digits and push into channel
	rowChannel <- final
}

func main() {
	filePath := path.Join(inputFile)
	file, err := os.Open(filePath)
	if err != nil {
		panic("could not open file")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if err != nil {
		panic("could not scan file")
	}

	var lines []Row
	for scanner.Scan() {
		line := scanner.Text()
		row := Row(line)
		lines = append(lines, row)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	wg.Add(len(lines))

	sum := 0
	for _, line := range lines {
		go calcRow(line)

	}

	for i := 0; i < len(lines); i++ {
		currVal := <-rowChannel
		fmt.Printf(" > got value: %d\n", currVal)
		sum += currVal

	}

	wg.Wait()

	fmt.Printf("Sum is: %d\n", sum)
}
