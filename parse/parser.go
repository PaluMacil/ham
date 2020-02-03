package parse

import (
	"bufio"
	"fmt"
	"github.com/PaluMacil/ham/experiment"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

const ratioToTrain = .75

func FromFile(filename, delimiter string) (experiment.Experiment, error) {
	file, err := os.Open(filename)
	if err != nil {
		return experiment.Experiment{}, fmt.Errorf("reading %s: %w", filename, err)
	}
	defer file.Close()
	return Parse(file, delimiter)
}

func Parse(reader io.Reader, delimiter string) (experiment.Experiment, error) {

	scanner := bufio.NewScanner(reader)
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		// skip empty lines
		if line == "" {
			continue
		}
		lines = append(lines, line)
	}
	// randomize slice in-place
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(lines), func(i, j int) { lines[i], lines[j] = lines[j], lines[i] })

	numberToTrain := int(float64(len(lines)) * ratioToTrain)
	var ex experiment.Experiment
	for i, line := range lines {
		parts := strings.Split(line, delimiter)
		// eliminate lines without a class
		if len(parts) < 2 {
			fmt.Println("skipped line because it had not class:", parts[0])
			continue
		}
		thisClass, err := experiment.ClassType(parts[0])
		if err != nil {
			return ex, fmt.Errorf("checking class type on shuffled row %d: %w", i, err)
		}
		if len(parts[1]) == 0 {
			return ex, fmt.Errorf("empty text on shuffled row %d", i)
		}
		if i < numberToTrain {
			// training sets
			if thisClass == experiment.SpamClass {
				ex.Classes.Spam = append(ex.Classes.Spam, parts[1])
			}
			if thisClass == experiment.HamClass {
				ex.Classes.Ham = append(ex.Classes.Ham, parts[1])
			}
		} else {
			// test sets
			testCase := experiment.TestCase{
				Class: thisClass,
				Text:  parts[1],
			}
			ex.Test.Cases = append(ex.Test.Cases, testCase)
		}

	}

	return ex, nil
}
