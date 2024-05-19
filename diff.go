package linediff

import (
	"fmt"
	"strings"
)

type Tag int

const (
	Same Tag = iota
	Added
	Removed
)

type DiffSet struct {
	segments []string
	tags     []Tag
}

func (s *DiffSet) String() string {
	var buf strings.Builder
	for i, segment := range s.segments {
		switch s.tags[i] {
		case Same:
			buf.WriteString(segment)
		case Added:
			buf.WriteString(fmt.Sprintf("(++%s++)", segment))
			if i > 0 && s.tags[i-1] == Removed {
				continue
			}
		case Removed:
			buf.WriteString(fmt.Sprintf("(--%s--)", segment))
		}
	}
	return buf.String()
}

func (s *DiffSet) Add(tag Tag, tokens ...string) {
	if len(tokens) == 0 {
		return
	}
	for _, token := range tokens {
		s.segments = append(s.segments, token)
	}
	for i := 0; i < len(tokens); i++ {
		s.tags = append(s.tags, tag)
	}
}

func (s *DiffSet) AddRemoval(tokens ...string) {
	s.Add(Removed, tokens...)
}

func (s *DiffSet) AddAddition(tokens ...string) {
	s.Add(Added, tokens...)
}

func (s *DiffSet) AddSimilarity(tokens ...string) {
	s.Add(Same, tokens...)
}

func Diff(a, b string) *DiffSet {
	return DiffSplit(a, b, SplitSpaces)
}

// DiffCrossConfidence determines how many tokens the cross comparison phase will look ahead before giving up.
// This can be tuned according to the length of the incoming data set to emit more or less verbose diffs.
var DiffCrossConfidence = 3

func DiffSplit(a, b string, split Splitter) *DiffSet {
	if split == nil {
		panic("nil splitter")
	}

	var (
		ds      = new(DiffSet)
		as      = split.Split(a)
		aOffset int
		bs      = split.Split(b)
		bOffset int
		maxi    = max(len(as), len(bs))
	)

loop:
	for i := 0; i < maxi; i++ {
		ai := i + aOffset
		bi := i + bOffset

		// Capacity diffs
		if ai >= len(as) {
			ds.AddAddition(bs[bi:]...)
			break loop
		}
		if bi >= len(bs) {
			ds.AddRemoval(as[ai:]...)
			break
		}

		// Comparisons
		if as[ai] == bs[bi] {
			ds.AddSimilarity(as[ai])
			continue
		}

		// Cross comparison
		for j := 1; j <= DiffCrossConfidence; j++ {
			if bi+j < len(bs) && as[ai] == bs[bi+j] {
				ds.AddAddition(bs[bi : bi+j]...)
				bOffset += j
				ds.AddSimilarity(as[ai])
				continue loop
			}
			if ai+j < len(as) && bs[bi] == as[ai+j] {
				ds.AddRemoval(as[ai : ai+j]...)
				aOffset += j
				ds.AddSimilarity(bs[bi])
				continue loop
			}
		}
		// Straight diff
		ds.AddRemoval(as[ai])
		ds.AddAddition(bs[bi])
	}
	return ds
}
