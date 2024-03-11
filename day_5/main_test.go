package main

import "testing"

func TestMap_Find(t *testing.T) {
	testLookups := []*lookupRange{
		{
			destRange:   int64(10),
			sourceRange: int64(0),
			rangeLength: int64(10),
		},
	}
	testMap := Map{
		source: "source",
		target: "target",
		lookup: testLookups,
	}

	inputs := []int64{-1, 0, 1, 9, 15, 20}
	wants := []int64{-1, 10, 11, 19, 15, 20}

	for i, input := range inputs {
		actual := testMap.Find(input)
		if wants[i] != actual {
			t.Fatalf("[TestMap_Find] for input %d, wanted %d, got %d", input, wants[i], actual)
		}
	}
}
