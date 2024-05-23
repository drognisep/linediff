package main

import (
	"fmt"
	"github.com/drognisep/linediff"
	flag "github.com/spf13/pflag"
)

type Config struct {
	HelpRequested     bool
	BufferSize        int
	LookAheadMatching int
	InFile            string
	ACol              int
	ALabel            string
	BCol              int
	BLabel            string
	Delimiters        string
	OutFile           string
	SkipFirstRow      bool
}

func setupFlags(config *Config) *flag.FlagSet {
	flags := flag.NewFlagSet("diffhtml", flag.ExitOnError)
	flags.BoolVarP(&config.HelpRequested, "help", "h", false, "Prints this usage information.")
	flags.IntVar(&config.BufferSize, "buffer", linediff.BufferSize, "Sets the read buffer size for diff samples in runes. This should be greater than or equal to the maximum sample size.")
	flags.IntVar(&config.LookAheadMatching, "matchahead", linediff.DiffCrossConfidence, "Sets the matching lookahead threshold for diffing. A larger threshold reduces performance, but tends to reduce diff size for inputs with less variance.")
	flags.StringVar(&config.InFile, "csv", "", "Specifies a CSV file should be read instead of arguments. Must be used with 'col-a' and 'col-b'.")
	flags.StringVarP(&config.OutFile, "out", "o", "index.html", "Specifies an output file for generation. Only used when the 'csv' option is specified.")
	flags.IntVarP(&config.ACol, "col-a", "a", -1, "Specifies the (0-indexed) A column for comparison. Only useful with the 'csv' option.")
	flags.StringVar(&config.ALabel, "header-a", "A", "Sets a label for sample A header.")
	flags.IntVarP(&config.BCol, "col-b", "b", -1, "Specifies the (0-indexed) B column for comparison. Only useful with the 'csv' option.")
	flags.StringVar(&config.BLabel, "header-b", "B", "Sets a label for sample B header.")
	flags.StringVar(&config.Delimiters, "delim", " ", "Specifies a custom input token delimiter. Each rune in this string is used to separate input terms for comparison. Defaults to space delimiting terms.")
	flags.BoolVar(&config.SkipFirstRow, "skip-header", true, "Skips the first row of a CSV file as the header.")

	flags.Usage = func() {
		fmt.Printf(`diffhtml generates an HTML page representing a table of diffs and their inputs.

USAGE
	diffhtml A B
	diffhtml --csv -a 0 -b 1

Be sure to escape long strings appropriately for your terminal when passing strings with spaces to the program.

SINGLE PAIR DIFF
If the 'csv' option is not used, then the first two arguments are expected to be A and B strings, respectively.
Argument A is diffed against argument B, and the diff is output to STDOUT with the default representation.

Example:
A: a simple string
B: a less simple string

This output will be printed to STDOUT with the inputs above.
a <span class="add">less</span><span class="add"> </span>simple string

CSV FILE DIFF
If the 'csv' option is used, then each row (excluding the first, assumed to be a header) will be part of a generated HTML table.
Both inputs will be shown as received, as well as the final diff view.

Example:
A: a simple string
B: a less simple string

This output (within the context of a surrounding table) will be generated to file.
<tr>
<td>a simple string</td>
<td>a less simple string</td>
<td>a <span class="add">less</span><span class="add"> </span>simple string</td>
</tr>

FLAGS
%s`, flags.FlagUsages())
	}

	return flags
}
