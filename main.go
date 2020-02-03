package main

import (
	"flag"
	"fmt"
	"github.com/PaluMacil/ham/analysis"
	"os"

	"github.com/PaluMacil/ham/parse"
)

func main() {
	var flagFilename string
	flag.StringVar(&flagFilename, "file", "textMsgs.data", "filename")
	var flagDelimiter string
	flag.StringVar(&flagDelimiter, "delimiter", "\t", "delimiter between class and words in data (default is tab)")
	var flagWriteToFile bool
	flag.BoolVar(&flagWriteToFile, "write", false, "write classes to files")
	flag.Parse()

	exp, err := parse.FromFile(flagFilename, flagDelimiter)
	if err != nil {
		fmt.Println("cannot parse file:", err)
		os.Exit(1)
	}

	analyses := analysis.Run(exp)

	if flagWriteToFile {
		if err := analyses.WriteToFile(); err != nil {
			fmt.Printf("could not write analysis to file: %s", err)
		}
	}
	fmt.Println("Done.")
}
