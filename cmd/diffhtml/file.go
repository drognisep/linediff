package main

import (
	_ "embed"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"text/template"
)

var (
	ErrInvalidColIndex = errors.New("invalid column index")
	//go:embed output.gohtml
	templText string
	templ     = template.Must(template.New("html").Funcs(map[string]any{
		"plusOne": func(i int) int {
			return i + 1
		},
	}).Parse(templText))
)

func runFileGeneration(config *Config) error {
	if config.ACol < 0 {
		return fmt.Errorf("%w: column A index '%d' is invalid", ErrInvalidColIndex, config.ACol)
	}
	if config.BCol < 0 {
		return fmt.Errorf("%w: column B index '%d' is invalid", ErrInvalidColIndex, config.BCol)
	}
	if config.ACol == config.BCol {
		return fmt.Errorf("%w: column indexes cannot be the same", ErrInvalidColIndex)
	}
	in, err := os.Open(config.InFile)
	if err != nil {
		return fmt.Errorf("%w: failed to open input file '%s'", err, config.InFile)
	}
	defer func() {
		_ = in.Close()
	}()
	csvr := csv.NewReader(in)

	out, err := os.Create(config.OutFile)
	if err != nil {
		return fmt.Errorf("failed to create output file '%s': %w", config.OutFile, err)
	}
	defer func() {
		_ = out.Close()
	}()

	log.Println("Reading input file...")
	var (
		splitter    = getSplitter(*config)
		diffRecords = make([]diffRecord, 0, 1024)
		first       = config.SkipFirstRow
		maxCol      = max(config.ACol, config.BCol)
		i           int
	)

	for {
		record, err := csvr.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to read CSV record: %w", err)
		}
		if first {
			first = false
			continue
		}
		if maxCol >= len(record) {
			return fmt.Errorf("one or more column index is out of bounds for row %d", i)
		}
		diffRecords = append(diffRecords, diffRecord{A: record[config.ACol], B: record[config.BCol], splitter: splitter})
	}

	log.Println("Generating HTML...")
	err = templ.Execute(out, map[string]any{
		"FileName": config.InFile,
		"Records":  diffRecords,
		"HeaderA":  config.ALabel,
		"HeaderB":  config.BLabel,
	})
	if err != nil {
		return fmt.Errorf("failed to generate HTML file: %w", err)
	}
	log.Println("Done")
	return nil
}
