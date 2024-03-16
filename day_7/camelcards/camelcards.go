package camelcards

import (
	"github.com/samber/lo"
	"golang.org/x/exp/maps"
)

type CardValue rune

const (
	A     CardValue = 'A'
	K     CardValue = 'K'
	Q     CardValue = 'Q'
	J     CardValue = 'J'
	T     CardValue = 'T'
	Nine  CardValue = '9'
	Eight CardValue = '8'
	Seven CardValue = '7'
	Six   CardValue = '6'
	Five  CardValue = '5'
	Four  CardValue = '4'
	Three CardValue = '3'
	Two   CardValue = '2'
)

type Card struct {
	value CardValue
}

type HandType string

const (
	fiveOfAKind  HandType = "five of a kind"
	fourOfAKind  HandType = "four of a kind"
	fullHouse    HandType = "full house"
	threeOfAKind HandType = "three of a kind"
	twoPair      HandType = "two pair"
	onePair      HandType = "one pair"
	highCard     HandType = "high card"
)

var HandTypes = []HandType{fiveOfAKind, fourOfAKind, fullHouse, threeOfAKind, twoPair, onePair, highCard}

type Hand struct {
	Cards    []*Card
	handType HandType
}

func handType(h *Hand) HandType {
	calcHandType := func(values []CardValue) HandType {
		valueHash := make(map[CardValue]int)
		for _, v := range values {
			_, ok := valueHash[v]
			if ok {
				valueHash[v] += 1
			} else {
				valueHash[v] = 1
			}
		}

		switch {
		case len(valueHash) == 1:
			return fiveOfAKind
		case lo.Max(maps.Values(valueHash)) == 4:
			return fourOfAKind
		case len(valueHash) == 2 && lo.Max(maps.Values(valueHash)) == 3:
			return fullHouse
		case len(valueHash) == 3 && lo.Max(maps.Values(valueHash)) == 3:
			return threeOfAKind
		case len(valueHash) == 3 && lo.CountBy(maps.Values(valueHash), func(v int) bool { return v == 2 }) == 2:
			return twoPair
		case len(valueHash) == 4:
			return onePair
		default:
			return highCard
		}
	}

	values := lo.Map(h.Cards, func(c *Card, _ int) CardValue {
		return c.value
	})
	handType := calcHandType(values)

	return handType
}

// Returns the hand type of the given hand `h`
func (h *Hand) Type() HandType {
	// If it's already been set, return that.
	if h.handType != "" {
		return h.handType
	}

	h.handType = handType(h)
	return h.handType
}
