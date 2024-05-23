package main

import (
	"fmt"
	"github.com/drognisep/linediff"
	"strings"
)

func getSplitter(config Config) linediff.Splitter {
	return linediff.SplitterFunc(func(tr *linediff.TokenReader) []string {
		var tokens []string
		sr := []rune(config.Delimiters)
		tokenFound := false
		for {
			tokenFound = false
			token, found := tr.Until(config.Delimiters)
			if found {
				tokenFound = true
				tokens = append(tokens, token)
			}
			for _, r := range sr {
				_, found = tr.AcceptToken(string(r))
				if found {
					tokenFound = true
					tokens = append(tokens, string(r))
					break
				}
			}
			if !tokenFound {
				break
			}
		}
		return tokens
	})
}

type diffRecord struct {
	A, B     string
	splitter linediff.Splitter
}

func (r *diffRecord) DiffHTML() string {
	diffs := linediff.DiffSplit(r.A, r.B, r.splitter).Iterator()
	var (
		seg  string
		tag  linediff.Tag
		next bool
		buf  strings.Builder
	)

	seg, tag, next = diffs.Next()
	for next {
		switch tag {
		case linediff.Same:
			buf.WriteString(seg)
		case linediff.Removed:
			buf.WriteString(fmt.Sprintf(`<span class="rem">%s</span>`, seg))
		case linediff.Added:
			buf.WriteString(fmt.Sprintf(`<span class="add">%s</span>`, seg))
		}
		seg, tag, next = diffs.Next()
	}
	return buf.String()
}
