package analysis

import (
	"github.com/PaluMacil/ham/experiment"
	"github.com/PaluMacil/ham/parse"
	"math"
	"strings"
)

type Preprocessor interface {
	Process(ex *experiment.Experiment)
}

type Analysis struct {
	Name        string
	TrainingSet TrainingSet
	TestSet     TestSet
}

func (a Analysis) Test(set experiment.TestSet) TestSet {
	results := TestSet{
		MessageTotal:          len(set.Cases),
		CorrectHam:            0,
		CorrectSpam:           0,
		IncorrectHam:          0,
		IncorrectSpam:         0,
		PercentageCorrectHam:  0,
		PercentageCorrectSpam: 0,
	}

	var hamScore float64
	var spamScore float64
	for _, sms := range set.Cases {
		hamScore, spamScore = 0, 0
		words := strings.Split(sms.Text, " ")
		for _, word := range words {
			// skip word if it isn't in the vocabulary
			if !a.TrainingSet.Vocabulary.Contains(word) {
				continue
			}

			// ham (add logs to prevent underflow of float with lots of multiplying)
			hamScore = hamScore + math.Log10(a.TrainingSet.Ham.WordProbabilities[word])

			// spam (add logs to prevent underflow of float with lots of multiplying)
			spamScore = spamScore + math.Log10(a.TrainingSet.Spam.WordProbabilities[word])
		}
		// unlog
		hamScore = math.Pow(10, hamScore)
		spamScore = math.Pow(10, spamScore)

		// if the algorithm says this is *ham*
		if hamScore > spamScore {
			if sms.Class == experiment.HamClass {
				results.CorrectHam = results.CorrectHam + 1
			} else {
				results.IncorrectHam = results.IncorrectHam + 1
			}
		} else { // algorithm says this is *spam*
			if sms.Class == experiment.SpamClass {
				results.CorrectSpam = results.CorrectSpam + 1
			} else {
				results.IncorrectSpam = results.IncorrectSpam + 1
			}
		}
	}
	results.PercentageCorrectHam = float64(results.CorrectHam) / float64(results.CorrectHam+results.IncorrectSpam)
	results.PercentageCorrectSpam = float64(results.CorrectSpam) / float64(results.CorrectSpam+results.IncorrectHam)

	return results
}

type TestSet struct {
	MessageTotal          int
	CorrectHam            int
	CorrectSpam           int
	IncorrectHam          int
	IncorrectSpam         int
	PercentageCorrectHam  float64
	PercentageCorrectSpam float64
}

type TrainingSet struct {
	MessageTotal int
	Spam         Class
	Ham          Class
	Vocabulary   Vocabulary
}

type Vocabulary []string

func (v Vocabulary) Contains(word string) bool {
	for _, w := range v {
		if w == word {
			return true
		}
	}
	return false
}

type WordFrequency map[string]int

func (wf WordFrequency) Probability(v Vocabulary) Probability {
	p := make(map[string]float64)
	for _, vocabWord := range v {
		p[vocabWord] = float64(wf[vocabWord]+1) / float64(len(wf)+len(v))
	}

	return p
}

type Probability map[string]float64

type Class struct {
	// MessageTotal is the message total
	MessageTotal int
	// PofC is the probability of this class out of the total in the training set
	PofC float64
	// WordFrequency represents words and how many times they occur in this class
	WordFrequency WordFrequency
	// WordProbabilities is a lookup of words and their probability of occurring in this class
	WordProbabilities Probability
}

type Analyses []Analysis

func Run(ex experiment.Experiment) Analyses {
	var analyses Analyses

	// first analyze with no preprocessors
	defaultAnalysis := analysisFrom(ex)
	defaultAnalysis.Name = "Default Analysis (no preprocessing)"
	analyses = append(analyses, defaultAnalysis)

	// copy experiment for multiple types of preprocessing
	ex2, ex3, ex4, ex5 := ex, ex, ex, ex

	// analyze with punctuation removed
	parse.PreprocessRemovePunctuation{}.Process(&ex2)
	a2 := analysisFrom(ex2)
	a2.Name = "No Punctuation Analysis"
	analyses = append(analyses, a2)

	// analyze with stemmer
	parse.PreprocessStemmer{}.Process(&ex3)
	a3 := analysisFrom(ex3)
	a3.Name = "Stemmer Analysis"
	analyses = append(analyses, a3)

	// analyze with punctuation removed and then stemmer
	parse.PreprocessStemmer{}.Process(&ex4)
	parse.PreprocessRemovePunctuation{}.Process(&ex4)
	a4 := analysisFrom(ex4)
	a4.Name = "Stemmer and No Punctuation Analysis"
	analyses = append(analyses, a4)

	// analyze with common words removed
	parse.PreprocessRemoveCommonWords{}.Process(&ex5)
	a5 := analysisFrom(ex5)
	a5.Name = "Remove 100 Most Common English Words"
	analyses = append(analyses, a5)

	return analyses
}

func analysisFrom(ex experiment.Experiment) Analysis {
	totalTrainingMessages := len(ex.Classes.Ham) + len(ex.Classes.Spam)
	vocabulary := vocabularyFrom(ex.Classes.Ham, ex.Classes.Spam)
	hamFrequency := wordFrequencyFrom(ex.Classes.Ham)
	hamProbabilities := hamFrequency.Probability(vocabulary)
	spamFrequency := wordFrequencyFrom(ex.Classes.Spam)
	spamProbabilities := spamFrequency.Probability(vocabulary)

	analysis := Analysis{
		TrainingSet: TrainingSet{
			MessageTotal: totalTrainingMessages,
			Ham: Class{
				MessageTotal:      len(ex.Classes.Ham),
				PofC:              float64(len(ex.Classes.Ham)) / float64(totalTrainingMessages),
				WordFrequency:     hamFrequency,
				WordProbabilities: hamProbabilities,
			},
			Spam: Class{
				MessageTotal:      len(ex.Classes.Spam),
				PofC:              float64(len(ex.Classes.Spam)) / float64(totalTrainingMessages),
				WordFrequency:     spamFrequency,
				WordProbabilities: spamProbabilities,
			},
			Vocabulary: vocabulary,
		},
	}
	analysis.TestSet = analysis.Test(ex.Test)

	return analysis
}

func vocabularyFrom(messageLists ...[]string) Vocabulary {
	keys := make(map[string]bool)
	var vocabulary Vocabulary
	for _, messageList := range messageLists {
		for _, msg := range messageList {
			for _, word := range strings.Split(msg, " ") {
				if word == "" {
					continue
				}
				if _, exists := keys[word]; !exists {
					keys[word] = true
					vocabulary = append(vocabulary, word)
				}
			}
		}
	}

	return vocabulary
}

func wordFrequencyFrom(messageList []string) WordFrequency {
	frequency := make(map[string]int)
	for _, msg := range messageList {
		for _, word := range strings.Split(msg, " ") {
			if word == "" {
				continue
			}
			occurrences, exists := frequency[word]
			if exists {
				frequency[word] = occurrences + 1
			} else {
				frequency[word] = 1
			}
		}
	}

	return frequency
}
