package main

import (
	"flag"
	"fmt"
	"github.com/PaluMacil/ham/analysis"
	"github.com/fatih/color"
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

	for _, a := range analyses {
		c := color.New(color.FgCyan).Add(color.Underline)
		c.Printf("Analysis: %s\n", a.Name)
		fmt.Println("Vocabulary has", len(a.TrainingSet.Vocabulary), "words")
		fmt.Println("\nTraining Set:")
		fmt.Printf("\t%d of %d messages were spam (%.2f%%)\n\n",
			a.TrainingSet.Spam.MessageTotal,
			a.TrainingSet.MessageTotal,
			a.TrainingSet.Spam.PofC*100)
		fmt.Println("Test Set:")
		fmt.Println("\tCorrect Ham:", a.TestSet.CorrectHam)
		fmt.Println("\tCorrect Spam:", a.TestSet.CorrectSpam)
		fmt.Println("\tIncorrect Ham (actually was spam):", a.TestSet.IncorrectHam)
		fmt.Println("\tIncorrect Spam (actually was ham):", a.TestSet.IncorrectSpam)
		fmt.Printf("\tPercentage Correct Ham: %.2f%%\n", a.TestSet.PercentageCorrectHam*100)
		fmt.Printf("\tPercentage Correct Spam: %.2f%%\n", a.TestSet.PercentageCorrectSpam*100)
		bold := color.New(color.FgGreen, color.Bold)
		bold.Printf("\tOverall Accuracy: %.2f%%\n",
			100*(float64(a.TestSet.CorrectSpam)+float64(a.TestSet.CorrectHam))/
				float64(a.TestSet.MessageTotal))
		fmt.Println()
	}
	fmt.Println("\nDone.")
}
