package camelcards

import (
	"testing"
)

var testHands = []Hand{
	{
		Cards: []*Card{
			&Card{'4'}, &Card{'9'}, &Card{'A'}, &Card{'4'}, &Card{'9'},
		},
	},
	{
		Cards: []*Card{
			&Card{'6'}, &Card{'7'}, &Card{'5'}, &Card{'9'}, &Card{'4'},
		},
	},
	{
		Cards: []*Card{
			&Card{'Q'}, &Card{'2'}, &Card{'4'}, &Card{'2'}, &Card{'9'},
		},
	},
	{
		Cards: []*Card{
			&Card{'A'}, &Card{'7'}, &Card{'7'}, &Card{'7'}, &Card{'Q'},
		},
	},
	{
		Cards: []*Card{
			&Card{'A'}, &Card{'A'}, &Card{'A'}, &Card{'Q'}, &Card{'Q'},
		},
	},
}

var testBids = []*Bid{
	{Hand: &testHands[0], BidAmount: 20},
	{Hand: &testHands[1], BidAmount: 50},
	{Hand: &testHands[2], BidAmount: 40},
	{Hand: &testHands[3], BidAmount: 10},
	{Hand: &testHands[4], BidAmount: 100},
}

var expectedTypes = []HandType{twoPair, highCard, onePair, threeOfAKind, fullHouse}

func TestHand_Type(t *testing.T) {
	for i, hand := range testHands {
		if hand.Type() != expectedTypes[i] {
			t.Fatalf("[TestHand_Type] expected: %s, actual: %s", expectedTypes[i], hand.Type())
		}
	}
}

func TestHand_IsStrongerThan(t *testing.T) {
	strongerMap := map[[2]int]bool{
		{0, 1}: true,
		{0, 3}: false,
		{4, 3}: true,
	}

	for hands, expectStronger := range strongerMap {
		isStronger, err := testHands[hands[0]].IsStrongerThan(&testHands[hands[1]])
		if err != nil {
			t.Fatalf("[TestHand_IsStrongerThan] unexpected error '%s'", err.Error())
		}

		if isStronger != expectStronger {
			t.Fatalf("[TestHand_IsStrongerThan] failed IsStrongerThan expectation")
		}
	}
}

func Test_SortByStrength(t *testing.T) {
	toTest := SortByStrength(testBids)

	for i := 1; i < len(toTest); i++ {
		fail, err := toTest[i-1].Hand.IsStrongerThan(toTest[i].Hand)
		if err != nil {
			t.Fatalf("[Test_SortByStrength] error during strength check")
		}
		if fail {
			t.Fatalf("[Test_SortByStrength] failed: hand at index %d stronger than at index %d", i-1, i)
		}
	}
}
