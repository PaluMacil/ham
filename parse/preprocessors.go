package parse

import (
	"github.com/PaluMacil/ham/experiment"
	"github.com/reiver/go-porterstemmer"
	"regexp"
	"strings"
)

type PreprocessStemmer struct{}

func (p PreprocessStemmer) Process(ex *experiment.Experiment) {
	for i, m := range ex.Classes.Spam {
		ex.Classes.Spam[i] = p.processMessage(m)
	}
	for i, m := range ex.Classes.Ham {
		ex.Classes.Ham[i] = p.processMessage(m)
	}
	for i, t := range ex.Test.Cases {
		ex.Test.Cases[i].Text = p.processMessage(t.Text)
	}
}

func (p PreprocessStemmer) processMessage(original string) string {
	var words []string
	for _, word := range strings.Split(original, " ") {
		stem := porterstemmer.StemString(word)
		words = append(words, stem)
	}
	return strings.Join(words, " ")
}

type PreprocessRemovePunctuation struct{}

func (p PreprocessRemovePunctuation) Process(ex *experiment.Experiment) {
	for i, m := range ex.Classes.Spam {
		ex.Classes.Spam[i] = p.processMessage(m)
	}
	for i, m := range ex.Classes.Ham {
		ex.Classes.Ham[i] = p.processMessage(m)
	}
	for i, t := range ex.Test.Cases {
		ex.Test.Cases[i].Text = p.processMessage(t.Text)
	}
}

func (p PreprocessRemovePunctuation) processMessage(original string) string {
	reg, _ := regexp.Compile("[^a-zA-Z0-9 ]+")

	return reg.ReplaceAllString(original, "")
}

type PreprocessRemoveCommonWords struct{}

func (p PreprocessRemoveCommonWords) Process(ex *experiment.Experiment) {
	for i, m := range ex.Classes.Spam {
		ex.Classes.Spam[i] = p.processMessage(m)
	}
	for i, m := range ex.Classes.Ham {
		ex.Classes.Ham[i] = p.processMessage(m)
	}
	for i, t := range ex.Test.Cases {
		ex.Test.Cases[i].Text = p.processMessage(t.Text)
	}
}

func (p PreprocessRemoveCommonWords) processMessage(original string) string {
	var words []string
	for _, word := range strings.Split(original, " ") {
		if isCommon(word) {
			continue
		}
		words = append(words, word)
	}
	return strings.Join(words, " ")
}

func isCommon(word string) bool {
	commonWords := []string{"the", "of", "and", "a", "to", "in", "is", "you", "that", "it", "he", "was", "for", "on",
		"are", "as", "with", "his", "they", "i", "at", "be", "this", "have", "from", "or", "one", "had", "by", "word",
		"but", "not", "what", "all", "were", "we", "when", "your", "can", "said", "there", "use", "an", "each", "which",
		"she", "do", "how", "their", "if", "will", "up", "other", "about", "out", "many", "then", "them", "these", "so",
		"some", "her", "would", "make", "like", "him", "into", "time", "has", "look", "two", "more", "write", "go",
		"see", "number", "no", "way", "could", "people", "my", "than", "first", "water", "been", "call", "who", "oil",
		"its", "now", "find", "long", "down", "day", "did", "get", "come", "made", "may", "part"}

	for _, commonWord := range commonWords {
		if commonWord == strings.ToLower(word) {
			return true
		}
	}
	return false
}
