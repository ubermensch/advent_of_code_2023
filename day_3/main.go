package main

import (
	"day_3/schematic"
	"fmt"
	"log"
)

func main() {
	s, err := schematic.NewSchematic()
	if err != nil {
		log.Fatal(err)
	}
	valid, err := s.ValidPartNumbers()
	if err != nil {
		log.Fatal(err)
	}

	sum := 0
	for _, curr := range valid {
		sum += curr.Number
	}

	fmt.Println(fmt.Sprintf("Sum of valid part numbers is: %d", sum))
}
