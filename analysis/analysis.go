package analysis

import (
	"github.com/PaluMacil/ham/experiment"
	"strings"
)

type Preprocessor interface {
	Process(original string) string
}

type Analysis struct {
	TrainingSet TrainingSet
}

type TrainingSet struct {
	MessageTotal int
	Spam         Class
	Ham          Class
	Vocabulary   map[string]int
}

type Vocabulary map[string]int

func (v Vocabulary) Probabilities(wordCount int) map[string]float64 {
	p := make(map[string]float64)
	for key, frequency := range v {
		p[key] = float64(frequency) / float64(wordCount)
	}

	return p
}

type Class struct {
	// MessageTotal is the message total
	MessageTotal int
	// PofC is the probability of this class out of the total in the training set
	PofC float64
	// WordCount is the count of words in this class
	WordCount int
	// Vocabulary represents words and how many times they occur in this class
	Vocabulary Vocabulary
	// WordProbabilities is a lookup of words and their probability of occurring in this class
	WordProbabilities map[string]float64
}

type Analyses []Analysis

func (a Analyses) WriteToFile() error {
	return nil
}

func Run(ex experiment.Experiment, preprocessors ...Preprocessor) Analyses {
	analyses := make([]Analysis, len(preprocessors))

	// always analyze with no preprocessors
	defaultAnalysis := analysisFrom(ex)

	analyses = append(analyses, defaultAnalysis)

	return analyses
}

func analysisFrom(ex experiment.Experiment) Analysis {
	totalTrainingMessages := len(ex.Classes.Ham) + len(ex.Classes.Spam)
	classWordCounts := ex.Classes.WordCounts()
	hamVocabulary := VocabularyFrom(ex.Classes.Ham)
	hamProbabilities := hamVocabulary.Probabilities(classWordCounts.Ham)
	spamVocabulary := VocabularyFrom(ex.Classes.Spam)
	spamProbabilities := spamVocabulary.Probabilities(classWordCounts.Spam)

	return Analysis{
		TrainingSet: TrainingSet{
			MessageTotal: totalTrainingMessages,
			Ham: Class{
				MessageTotal:      len(ex.Classes.Ham),
				PofC:              float64(len(ex.Classes.Ham)) / float64(totalTrainingMessages),
				WordCount:         classWordCounts.Ham,
				Vocabulary:        hamVocabulary,
				WordProbabilities: hamProbabilities,
			},
			Spam: Class{
				MessageTotal:      len(ex.Classes.Spam),
				PofC:              float64(len(ex.Classes.Spam)) / float64(totalTrainingMessages),
				WordCount:         classWordCounts.Spam,
				Vocabulary:        spamVocabulary,
				WordProbabilities: spamProbabilities,
			},
			Vocabulary: VocabularyFrom(ex.Classes.Ham, ex.Classes.Spam),
		},
	}
}

func VocabularyFrom(messageLists ...[]string) Vocabulary {
	vocabulary := make(map[string]int)
	for _, messageList := range messageLists {
		for _, msg := range messageList {
			for _, word := range strings.Split(msg, " ") {
				if word == "" {
					continue
				}
				occurrences, exists := vocabulary[word]
				if exists {
					vocabulary[word] = occurrences + 1
				} else {
					vocabulary[word] = 1
				}
			}
		}
	}

	return vocabulary
}
