package linediff

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDiffSet_AddRemoval(t *testing.T) {
	ds := new(DiffSet)
	ds.AddRemoval()
	assert.Len(t, ds.segments, 0)
	assert.Len(t, ds.tags, 0)

	ds.AddRemoval("a")
	assert.Len(t, ds.segments, 1)
	assert.Len(t, ds.tags, 1)
	assert.Equal(t, "a", ds.segments[0])
	assert.Equal(t, Removed, ds.tags[0])

	ds.AddRemoval("string", "here")
	assert.Len(t, ds.segments, 3)
	assert.Len(t, ds.tags, 3)
	assert.Equal(t, "a", ds.segments[0])
	assert.Equal(t, Removed, ds.tags[0])
	assert.Equal(t, "string", ds.segments[1])
	assert.Equal(t, Removed, ds.tags[1])
	assert.Equal(t, "here", ds.segments[2])
	assert.Equal(t, Removed, ds.tags[2])
}

func TestDiffSet_AddAddition(t *testing.T) {
	ds := new(DiffSet)
	ds.AddAddition()
	assert.Len(t, ds.segments, 0)
	assert.Len(t, ds.tags, 0)

	ds.AddAddition("a")
	assert.Len(t, ds.segments, 1)
	assert.Len(t, ds.tags, 1)
	assert.Equal(t, "a", ds.segments[0])
	assert.Equal(t, Added, ds.tags[0])

	ds.AddAddition("string", "here")
	assert.Len(t, ds.segments, 3)
	assert.Len(t, ds.tags, 3)
	assert.Equal(t, "a", ds.segments[0])
	assert.Equal(t, Added, ds.tags[0])
	assert.Equal(t, "string", ds.segments[1])
	assert.Equal(t, Added, ds.tags[1])
	assert.Equal(t, "here", ds.segments[2])
	assert.Equal(t, Added, ds.tags[2])
}

func TestDiffSet_AddSame(t *testing.T) {
	ds := new(DiffSet)
	ds.AddSimilarity()
	assert.Len(t, ds.segments, 0)
	assert.Len(t, ds.tags, 0)

	ds.AddSimilarity("a")
	assert.Len(t, ds.segments, 1)
	assert.Len(t, ds.tags, 1)
	assert.Equal(t, "a", ds.segments[0])
	assert.Equal(t, Same, ds.tags[0])

	ds.AddSimilarity("string", "here")
	assert.Len(t, ds.segments, 3)
	assert.Len(t, ds.tags, 3)
	assert.Equal(t, "a", ds.segments[0])
	assert.Equal(t, Same, ds.tags[0])
	assert.Equal(t, "string", ds.segments[1])
	assert.Equal(t, Same, ds.tags[1])
	assert.Equal(t, "here", ds.segments[2])
	assert.Equal(t, Same, ds.tags[2])
}

func TestDiff(t *testing.T) {
	assert.Panics(t, func() {
		DiffSplit("a", "b", nil)
	})
	tests := map[string]struct {
		A      string
		B      string
		Result string
	}{
		"No difference": {
			A:      "a string here",
			B:      "a string here",
			Result: "a string here",
		},
		"Here added": {
			A:      "a string",
			B:      "a string here",
			Result: "a string(++ ++)(++here++)",
		},
		"Here removed": {
			A:      "a string here",
			B:      "a string",
			Result: "a string(-- --)(--here--)",
		},
		"A added": {
			A:      "string here",
			B:      "a string here",
			Result: "(++a++)(++ ++)string here",
		},
		"A removed": {
			A:      "a string here",
			B:      "string here",
			Result: "(--a--)(-- --)string here",
		},
		"Replacement": {
			A:      "a string here",
			B:      "some string here",
			Result: "(--a--)(++some++) string here",
		},
		"Transposition": {
			A:      "a string here",
			B:      "string here a",
			Result: "(--a--)(-- --)string here(++ ++)(++a++)",
		},
		"Space difference": {
			A:      "a string  here",
			B:      "a  string here",
			Result: "a (++ ++)string (-- --)here",
		},
		"Complete replacement": {
			A:      "a string here",
			B:      "the quick brown fox jumps over the lazy goose",
			Result: "(--a--)(++the++) (--string--)(++quick++) (--here--)(++brown++)(++ ++)(++fox++)(++ ++)(++jumps++)(++ ++)(++over++)(++ ++)(++the++)(++ ++)(++lazy++)(++ ++)(++goose++)",
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			ds := Diff(tc.A, tc.B)
			assert.Equal(t, tc.Result, ds.String())
		})
	}
}

func TestDiff_LongStrings(t *testing.T) {
	tests := map[string]struct {
		A      string
		B      string
		Result string
	}{
		"Same": {
			A:      "a really long string that goes around here",
			B:      "a really long string that goes around here",
			Result: "a really long string that goes around here",
		},
		"Transposition": {
			A:      "a really long string that goes around here",
			B:      "really long string that goes around here a",
			Result: "(--a--)(-- --)really long string that goes around here(++ ++)(++a++)",
		},
		"Inner transposition": {
			A:      "a really long string that goes around here",
			B:      "a really long that string goes around here",
			Result: "a really long (++that++)(++ ++)string (--that--)(-- --)goes around here",
		},
		"Complete Replacement": {
			A:      "a really long string that goes around here",
			B:      "the quick brown fox jumps over the lazy goose",
			Result: "(--a--)(++the++) (--really--)(++quick++) (--long--)(++brown++) (--string--)(++fox++) (--that--)(++jumps++) (--goes--)(++over++) (--around--)(++the++) (--here--)(++lazy++)(++ ++)(++goose++)",
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			ds := Diff(tc.A, tc.B)
			assert.Equal(t, tc.Result, ds.String())
		})
	}
}
