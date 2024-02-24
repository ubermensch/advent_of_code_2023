package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/thoas/go-funk"
	"os"
	"path"
	"strconv"
	"strings"
)

const inputFile = "/Users/frankhmeidan/golang/advent_of_code/day_2/input.txt"

const (
	maxGreen = 13
	maxRed   = 12
	maxBlue  = 14
)

var gameChannel = make(chan *Game)
var games []*Game

type Game struct {
	id         int
	blueDrawn  int
	greenDrawn int
	redDrawn   int
	power      int
	isPossible bool
}

type Draw struct {
	blue  int
	green int
	red   int
}

func (draw *Draw) setColors(drawStr string) error {
	var finalErr error = nil

	pieces := strings.Split(drawStr, ",")
	funk.ForEach(pieces, func(s string) {
		pieces := strings.Split(strings.Trim(s, " "), " ")
		color := pieces[1]
		count, err := strconv.Atoi(pieces[0])
		if err != nil {
			finalErr = errors.New("could not parse draw for string: " + drawStr)
			return
		}
		switch color {
		case "green":
			draw.green = count
		case "blue":
			draw.blue = count
		case "red":
			draw.red = count
		default:
			finalErr = errors.New("could not parse draw for string: " + drawStr)
		}
	})

	return finalErr
}

func parseDraw(drawStr string) (*Draw, error) {
	draw := &Draw{
		green: 0,
		red:   0,
		blue:  0,
	}

	err := draw.setColors(drawStr)
	return draw, err
}

func setIsPossible(game *Game) {
	if game.blueDrawn <= maxBlue &&
		game.redDrawn <= maxRed &&
		game.greenDrawn <= maxGreen {
		game.isPossible = true
	} else {
		game.isPossible = false
	}
}

// Turns the game text line from the input file into a Game struct and returns
// the pointer to it
func newGame(gameText string) (*Game, error) {
	idAndGame := strings.Split(gameText, ":")
	idStr, gameStr := idAndGame[0], idAndGame[1]

	gameId, err := strconv.Atoi(strings.Split(idStr, " ")[1])
	if err != nil {
		return nil, errors.New("could not determine game ID")
	}

	drawStrings := strings.Split(gameStr, ";")
	var draws []*Draw
	for _, currDraw := range drawStrings {
		draw, err := parseDraw(currDraw)
		if err != nil {
			return nil, errors.New("could not parse game string: " + currDraw)
		}
		draws = append(draws, draw)
	}

	// Find the maximum number drawn of each color in this game
	blueMax, greenMax, redMax := 0, 0, 0
	funk.ForEach(draws, func(draw *Draw) {
		if draw.blue > blueMax {
			blueMax = draw.blue
		}
		if draw.green > greenMax {
			greenMax = draw.green
		}
		if draw.red > redMax {
			redMax = draw.red
		}
	})

	game := &Game{
		id:         gameId,
		blueDrawn:  blueMax,
		greenDrawn: greenMax,
		redDrawn:   redMax,
		power:      blueMax * greenMax * redMax,
		isPossible: false,
	}

	setIsPossible(game)
	return game, nil
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

	var lines []string

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	for _, line := range lines {
		go func(l string) {
			currGame, err := newGame(l)
			if err != nil {
				panic("could not create game")
			}
			gameChannel <- currGame
		}(line)
	}
	for i := 0; i < len(lines); i++ {
		games = append(games, <-gameChannel)
	}

	possible := funk.Filter(games, func(game *Game) bool {
		return game.isPossible
	})

	sumOfPossibleIds := funk.Reduce(
		possible,
		func(acc int, g *Game) int {
			return acc + g.id
		},
		0,
	)

	sumOfPowers := funk.Reduce(
		games,
		func(acc int, g *Game) int {
			return acc + g.power
		},
		0,
	)

	fmt.Println(
		fmt.Sprintf(
			"Possible game IDs: %v\nSum of IDs: %d\nSum of Powers: %d\n",
			funk.Map(possible, func(g *Game) int { return g.id }),
			sumOfPossibleIds,
			sumOfPowers,
		),
	)
}
