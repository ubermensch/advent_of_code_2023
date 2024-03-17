package main

import (
	"bufio"
	"day_7/camelcards"
	"fmt"
	"github.com/samber/lo"
	"os"
	"path"
	"strings"
)

const inputFile = "/Users/frankhmeidan/golang/advent_of_code/day_7/input.txt"

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

	var bids []*camelcards.Bid
	for scanner.Scan() {
		line := scanner.Text()
		pieces := strings.Split(line, " ")
		cards, bidAmount := pieces[0], pieces[1]
		newBid, err := camelcards.NewBid(cards, bidAmount)
		if err != nil {
			panic("could not create new bid for: " + cards + " " + bidAmount)
		}
		bids = append(bids, newBid)
	}

	bids = camelcards.SortByStrength(bids)
	totalWinnings := lo.Reduce(bids, func(total int, bid *camelcards.Bid, i int) int {
		return total + (bid.BidAmount * (i + 1))
	}, 0)

	fmt.Println("Total winnings: ", totalWinnings)
}
