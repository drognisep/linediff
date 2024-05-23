package main

import (
	"fmt"
	"github.com/drognisep/linediff"
	"log"
	"os"
)

func main() {
	config := new(Config)
	flags := setupFlags(config)
	if err := flags.Parse(os.Args[1:]); err != nil {
		log.Println(err)
		flags.Usage()
		os.Exit(1)
	}
	if config.HelpRequested {
		flags.Usage()
		return
	}
	if config.BufferSize < 256 {
		log.Println("Buffer size", config.BufferSize, "is below the lower bound of 256")
		flags.Usage()
		os.Exit(1)
	}
	linediff.BufferSize = config.BufferSize
	if config.LookAheadMatching < 3 {
		log.Println("Lookahead matching threshold", config.LookAheadMatching, "is below the lower bound of 3")
		flags.Usage()
		os.Exit(1)
	}

	if len(config.InFile) > 0 {
		if err := runFileGeneration(config); err != nil {
			log.Println(err)
			flags.Usage()
			os.Exit(1)
		}
		return
	}

	if flags.NArg() < 2 {
		log.Println("Missing A and B sample arguments")
		flags.Usage()
		os.Exit(1)
	}
	r := diffRecord{A: flags.Arg(0), B: flags.Arg(1), splitter: getSplitter(*config)}
	fmt.Print(r.DiffHTML())
}
