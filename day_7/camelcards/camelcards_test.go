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

var expectedTypes = []HandType{twoPair, highCard, onePair, threeOfAKind, fullHouse}

func TestHand_Type(t *testing.T) {
	for i, hand := range testHands {
		if hand.Type() != expectedTypes[i] {
			t.Fatalf("[TestHand_Type] expected: %s, actual: %s", expectedTypes[i], hand.Type())
		}
	}
}
