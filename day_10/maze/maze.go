package maze

import (
	"github.com/samber/lo"
	"math"
	"slices"
	"strings"
)

const (
	Wall              = '.'
	TopLeftCorner     = 'F'
	Horizontal        = '-'
	TopRightCorner    = '7'
	Vertical          = '|'
	BottomRightCorner = 'J'
	BottomLeftCorner  = 'L'
	Start             = 'S'
)

type Tile struct {
	// what kind of tile is this? One of [., F, -, 7, |, J, L, S]
	label rune

	// co-ordinates of this Tile in the maze.
	x int
	y int

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

func (t *Tile) canMoveFrom(source *Tile) bool {
	if source.x == t.x-1 {
		// source is left of target
		if slices.Contains([]rune{Horizontal, TopRightCorner, BottomRightCorner}, t.label) &&
			slices.Contains([]rune{Horizontal, TopLeftCorner, BottomLeftCorner, Start}, source.label) {
			return true
		}
	}

	if source.x == t.x+1 {
		// source is right of target
		if slices.Contains([]rune{Horizontal, TopLeftCorner, BottomLeftCorner}, t.label) &&
			slices.Contains([]rune{Horizontal, TopRightCorner, BottomRightCorner, Start}, source.label) {
			return true
		}
	}

	if source.y == t.y-1 {
		// source is above target
		if slices.Contains([]rune{Vertical, BottomLeftCorner, BottomRightCorner}, t.label) &&
			slices.Contains([]rune{Vertical, TopLeftCorner, TopRightCorner, Start}, source.label) {
			return true
		}
	}

	if source.y == t.y+1 {
		// source is below target
		if slices.Contains([]rune{Vertical, TopLeftCorner, TopRightCorner}, t.label) &&
			slices.Contains([]rune{Vertical, BottomRightCorner, BottomLeftCorner, Start}, source.label) {
			return true
		}
	}

	return false
}

func (t *Tile) canMoveTo(target *Tile) bool {
	return target.canMoveFrom(t)
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

func (m *Maze) neighboursOf(t *Tile) []*Tile {
	var ns []*Tile
	lo.ForEach([][]int{{t.x - 1, t.y}, {t.x, t.y - 1}, {t.x, t.y + 1}, {t.x + 1, t.y}}, func(nCoords []int, i int) {
		x, y := nCoords[0], nCoords[1]
		if x < 0 || y < 0 {
			return
		}
		if y >= len(m.tiles) || x >= len(m.tiles[0]) {
			return
		}

		possN := m.tiles[y][x]
		if m.tiles[nCoords[0]][nCoords[1]] != nil {
			ns = append(ns, possN)
		}
	})
	return ns
}

// algorithm:
// 1) start from S, set it's distance (0)
// 2) find adjacent available to move to, that doesn't have a distance set (doesn't matter which).
// 3) label it's distance (1).
// 4) move neighbour without distance (should only be 1 left)
// 5) label it's distance (prev + 1)
// 6) repeat steps 4 and 5 until reaching `S` again.
func (m *Maze) findLoop() []*Tile {
	start := m.start
	currDistance := 0

	curr := start
	loopTiles := []*Tile{}

	for curr != nil {
		curr.distance = currDistance
		loopTiles = append(loopTiles, curr)

		next, _ := lo.Find(m.neighboursOf(curr), func(n *Tile) bool {
			return curr.canMoveTo(n) && n.distance == 0
		})

		currDistance += 1
		curr = next
	}
	loopTiles = append(loopTiles, start)
	return loopTiles
}

func (t *Tile) Distance() int {
	return t.distance
}

func (m *Maze) FurthestTile() *Tile {
	if len(m.loop) == 0 || m.loop == nil {
		return nil
	}

	maxDist := 0
	var furthest *Tile

	for i, t := range m.loop {
		dist := int(math.Min(float64(i), float64(len(m.loop)-1-i)))
		t.distance = dist

		if dist > maxDist {
			maxDist = dist
			furthest = t
		}
	}

	return furthest
}

// returns the area enclosed by the loop,
// implements the shoelace algorithm:
// https://en.wikipedia.org/wiki/Shoelace_formula
func (m *Maze) PathArea() float64 {
	// sum the cross difference of each pair of path co-ordinates
	doubleArea := lo.Reduce(
		m.loop,
		func(a int, t *Tile, i int) int {
			// if we've reached the start at the end of the loop,
			// we're done.
			if i == len(m.loop)-1 {
				return a
			}
			// otherwise, we find the cross difference of t(i) and t(i+1). e.g.
			// if coords are t(i) = (1, 3) and t(i+1) = (6, 1),
			// then cross product is (1 * 1) - (6 * 3)
			nextT := m.loop[i+1]
			return a + (t.x * nextT.y) - (t.y * nextT.x)
		},
		0,
	)
	area := math.Abs(float64(doubleArea / 2))

	// tiles on the loop don't count as enclosed,
	// so deduct those (don't double-count the start point)
	area = area - float64(len(m.loop)-1)
	return area
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
