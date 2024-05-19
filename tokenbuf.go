package linediff

import (
	"github.com/drognisep/runebuffer"
	"io"
	"strings"
)

type Splitter interface {
	Split(s string) []string
}

type SplitterFunc func(s string) []string

func (f SplitterFunc) Split(s string) []string {
	return f(s)
}

var SplitSpaces = SplitterFunc(func(s string) []string {
	tr := NewStringTokenReader(s)
	var tokens []string
	for {
		token, found := tr.Until(" ")
		if !found {
			return tokens
		}
		tokens = append(tokens, token)
		for {
			space, found := tr.AcceptToken(" ")
			if !found {
				break
			}
			tokens = append(tokens, space)
		}
	}
})

type TokenReader struct {
	*runebuffer.RuneBuffer
}

func NewStringTokenReader(s string) *TokenReader {
	return NewStringTokenReaderWithSize(s, runebuffer.DefaultBufferSize)
}

func NewStringTokenReaderWithSize(s string, size int) *TokenReader {
	return NewTokenReaderWithSize(strings.NewReader(s), size)
}

func NewTokenReader(r io.Reader) *TokenReader {
	return NewTokenReaderWithSize(r, runebuffer.DefaultBufferSize)
}

func NewTokenReaderWithSize(r io.Reader, size int) *TokenReader {
	if r == nil {
		panic("nil reader")
	}
	return &TokenReader{
		RuneBuffer: runebuffer.NewRuneBufferWithSize(r, size),
	}
}

// AcceptToken returns exactly the token from the input if it exists.
// Otherwise, "" and false are returned, and all read runes are unread.
func (tr *TokenReader) AcceptToken(token string) (string, bool) {
	rs := []rune(token)

	for i := 0; i < len(rs); i++ {
		r, err := tr.ReadRune()
		if err != nil {
			tr.UnreadNumRunes(i + 1)
			return "", false
		}
		if rs[i] != r {
			tr.UnreadNumRunes(i + 1)
			return "", false
		}
	}
	return token, true
}

// Accept reads runes matching the accept list, returning the read token if any runes were read.
// The list is split to runes, and all unique runes are included in a match set.
func (tr *TokenReader) Accept(list string) (string, bool) {
	rs := []rune(list)
	matchSet := map[rune]bool{}
	var buf strings.Builder
	for _, r := range rs {
		matchSet[r] = true
	}
	for {
		r, err := tr.RuneBuffer.ReadRune()
		if err != nil {
			l := buf.Len()
			return buf.String(), l > 0
		}
		if matchSet[r] {
			buf.WriteRune(r)
			continue
		}
		tr.RuneBuffer.UnreadRune()
		l := buf.Len()
		return buf.String(), l > 0
	}
}

// TODO: Add UntilToken

// Until reads until a rune in the list string matches, returning the read token if any runes were read.
// The list is split to runes, and all unique runes are included in a match set.
func (tr *TokenReader) Until(list string) (string, bool) {
	rs := []rune(list)
	matchSet := map[rune]bool{}
	var buf strings.Builder
	for _, r := range rs {
		matchSet[r] = true
	}
	for {
		r, err := tr.RuneBuffer.ReadRune()
		if err != nil {
			l := buf.Len()
			return buf.String(), l > 0
		}
		if r != 0 && !matchSet[r] {
			buf.WriteRune(r)
			continue
		}
		tr.RuneBuffer.UnreadRune()
		l := buf.Len()
		return buf.String(), l > 0
	}
}
