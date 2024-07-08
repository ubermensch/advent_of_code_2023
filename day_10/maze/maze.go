package maze

import (
	"github.com/samber/lo"
	"strings"
)

type Tile struct {
	// what kind of tile this is? One of [., F, -, 7, |, J, L, S]
	label rune

	// co-ordinates of this Tile in the maze.
	x int
	y int

	// references to the neighbour tiles
	// going left and right along the path (might be nil)
	left  *Tile
	right *Tile

	// number of steps from the start of the maze to this tile's location
	distance int
}

type Maze struct {
	// references to the tile in this maze (0-indexed).
	// e.g. tiles[3][5] should have Tile with x == 3, y == 5
	tiles [][]*Tile
	// starting tile
	start *Tile
	// the longest continuous loop,
	// starting from the start tile and ending back at it
	loop []*Tile
}

func tilesFromLine(line string) []*Tile {
	return lo.Map(strings.Split(line, ""), func(s string, i int) *Tile {
		return &Tile{
			label: rune(s[0]),
			x:     i,
		}
	})
}

func findStart(tiles [][]*Tile) *Tile {
	for _, line := range tiles {
		for _, tile := range line {
			if tile.label == rune('S') {
				return tile
			}
		}
	}
	return nil
}

func (m *Maze) findLoop() []*Tile {
	return []*Tile{}
}

func (m *Maze) FurthestTile() *Tile {
	return nil
}

func NewMaze(lines []string) *Maze {
	tiles := lo.Map(lines, func(s string, i int) []*Tile {
		lineTiles := tilesFromLine(s)
		lo.ForEach(lineTiles, func(t *Tile, _ int) {
			t.y = i
		})
		return lineTiles
	})

	start := findStart(tiles)
	maze := &Maze{
		tiles: tiles,
		start: start,
		loop:  []*Tile{},
	}
	maze.loop = maze.findLoop()

	return maze
}
