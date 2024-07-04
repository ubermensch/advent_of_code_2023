package predictor

import "github.com/samber/lo"

type predictor interface {
	Next() int
}
type Series struct {
	predictor
	// History of measurements for a particular value
	History []int
	// telescoping list of deltas between one measurement and the next,
	// calculated until the deltas on a particular level converge to zero.
	deltas [][]int
}

func getDeltas(history []int) [][]int {
	var allDeltas [][]int
	allZeroes := false
	currDeltas := history

	lineDeltas := func(line []int) []int {
		var deltas []int
		for i := 0; i < len(line)-1; i++ {
			deltas = append(deltas, line[i+1]-line[i])
		}

		return deltas
	}

	for !allZeroes {
		nextDeltas := lineDeltas(currDeltas)
		allDeltas = append(allDeltas, nextDeltas)

		uniques := lo.Uniq(nextDeltas)
		allZeroes = len(uniques) == 1 && uniques[0] == 0
		currDeltas = nextDeltas
	}

	return allDeltas
}

func NewSeries(history []int) *Series {
	deltas := getDeltas(history)
	return &Series{
		History: history,
		deltas:  deltas,
	}
}

func (s *Series) Next() int {
	deltas := s.deltas
	nexts := []int{0}
	for i := len(deltas) - 1; i >= 0; i-- {
		nexts = append(nexts, deltas[i][len(deltas[i])-1]+nexts[len(nexts)-1])
	}
	return s.History[len(s.History)-1] + nexts[len(nexts)-1]
}
