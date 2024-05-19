package linediff

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestSplitSpaces(t *testing.T) {
	tests := map[string]struct {
		input  string
		tokens []string
	}{
		"No tokens": {
			input:  "",
			tokens: nil,
		},
		"Hello world": {
			input:  "hello world",
			tokens: []string{"hello", " ", "world"},
		},
		"Hello space world": {
			input:  "hello  world",
			tokens: []string{"hello", " ", " ", "world"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokens := SplitSpaces.Split(tc.input)
			assert.Equal(t, tc.tokens, tokens)
		})
	}
}

func TestNewTokenReader(t *testing.T) {
	assert.Panics(t, func() {
		NewTokenReader(nil)
	})
}

func TestTokenReader_Accept(t *testing.T) {
	tests := map[string]struct {
		text       string
		acceptList string
		output     string
		found      bool
	}{
		"Not found": {
			text:       "abc",
			acceptList: "123",
			output:     "",
			found:      false,
		},
		"Found prefix": {
			text:       "cabdef",
			acceptList: "abc",
			output:     "cab",
			found:      true,
		},
		"Full string": {
			text:       "cab",
			acceptList: "abc",
			output:     "cab",
			found:      true,
		},
		"Not found later": {
			text:       "cabdef",
			acceptList: "def",
			output:     "",
			found:      false,
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			tr := NewTokenReader(strings.NewReader(tc.text))
			output, found := tr.Accept(tc.acceptList)
			assert.Equal(t, tc.output, output)
			assert.Equal(t, tc.found, found)
		})
	}
}

func TestTokenReader_Until(t *testing.T) {
	tests := map[string]struct {
		text      string
		untilList string
		output    string
		found     bool
	}{
		"Full string": {
			text:      "abc",
			untilList: "123",
			output:    "abc",
			found:     true,
		},
		"Hit excluded prefix": {
			text:      "cabdef",
			untilList: "abc",
			output:    "",
			found:     false,
		},
		"Capture prefix": {
			text:      "cabdef",
			untilList: "def",
			output:    "cab",
			found:     true,
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			tr := NewTokenReader(strings.NewReader(tc.text))
			output, found := tr.Until(tc.untilList)
			assert.Equal(t, tc.output, output)
			assert.Equal(t, tc.found, found)
		})
	}
}

func TestTokenReader_AcceptToken(t *testing.T) {
	tests := map[string]struct {
		input       string
		searchToken string
		token       string
		found       bool
	}{
		"Empty string": {
			input:       "",
			searchToken: "blah",
			token:       "",
			found:       false,
		},
		"Not next token": {
			input:       "some blah",
			searchToken: "blah",
			token:       "",
			found:       false,
		},
		"Found first": {
			input:       "blah some",
			searchToken: "blah",
			token:       "blah",
			found:       true,
		},
		"Case sensitive": {
			input:       "BLAH some",
			searchToken: "blah",
			token:       "",
			found:       false,
		},
		"Found without residual": {
			input:       "blahblahblah",
			searchToken: "blah",
			token:       "blah",
			found:       true,
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			tr := NewStringTokenReader(tc.input)
			token, found := tr.AcceptToken(tc.searchToken)
			assert.Equal(t, tc.token, token)
			assert.Equal(t, tc.found, found)
		})
	}
}
