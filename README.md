# Ham (SMS spam classifier)

## Summary

The purpose of this project is to demonstrate a simple probabilistic SMS spam
classifier in Go. This supervised learning activity is accomplished using a
Naive Bayes classifier. This simple algorithm assumes conditional independence,
does a lot of counts, calculates some ratios, and multiplies them together.

### Project and Course Name

Project 1: Exploratory Data Analysis

CIS 678, Winter 2020

### Implementation

#### Organization

This project is organized into an experiment, analysis, and the main package and 
tests. Theses areas of the application have the following purposes:

- Experiment: holding the shuffled SMS messages, divided by set (train or test) and 
class (spam or ham).
- Analysis: holding calculations and statistics as well as a trained model for 
training classes and finally test results.
- main: connecting options to experiment code, making copies of the original 
experiment for preprocessor comparisons, and printing of the final output in a user-
friendly format.
- tests: parser tests ensure that ingested data is understood correctly.

#### Underflow Prevention

Underflow can occur when many small numbers are multiplied against each other. 
In most languages the rounding remains accurate for floating point underflow, 
but the loss in precision can round a small number to zero, causing multiplication 
to wipe out the other values around it. This is especially common in Naive Bayes 
analysis of text because you are multiplying a large number of small probabilities.

The code below implements the relationship `log(a * b) = log(a) + log(b)` in 
order to prevent underflow.

```
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
```

#### Experiments

##### Preprocessing: Remove Punctuation

Removing punctuation has almost no affect on the results. This might be due to 
the two classes tending to have similar punctuation. Also, many text
messages don't use any punctuation.

##### Preprocessing: Porter's Stemmer

The stemmer removes variations on words, combining them into one stem word. 
The idea behind this is that the ideas conveyed probably belong together in the 
analysis. However, there is a slight dip in accuracy (about a percent) when using 
a stemmer. This was one of the most surprising results to me. My theory is that 
the large quantity of slang and misspelling in sms cause some unexpected behavior 
from the stemmer. I would test with this preprocessor again on other types of 
text.

##### Preprocessing: Combination

I tested the combination of the stemmer and punctuation removal, assuming that 
they would both help and that combined they would be even better. The lack of 
accuracy uplift from the stemmer means that the combination was also not helpful.

##### Preprocessing: Remove 100 Most Common

This preprocessor removes the 100 most common English words from all messages. 
While the accuracy did get a small list from this adjustment, the difference 
was very small, and the improvement was mostly to ham identification whereas 
more ham was incorrectly identifies as spam. For this reason, I would not use 
this in a production environment.

##### Train vs Test Set Ratio

Changing the size of the training vs test sets (stored in a parser constant 
`const ratioToTrain = .75`) didn't have a large effect on the accuracy of the 
classifier. Even a drastic change from .75 to .25 only resulted in about one 
percent less accuracy. I assumed that this is due to even the smaller training 
set being sufficiently representative to train a valid model. Since there did 
not seem to be a good reason to make this setting configurable, I left it in a
constant instead of making it a commandline flag.  

#### Third Party Libraries

##### github.com/fatih/color

This color library adds ANSI-standard color codes to terminal output, allowing easier 
visual separation of information in the command output.

##### github.com/reiver/go-porterstemmer

This project implements an algorithm for suffix stripping in Go. The complexity of this 
stemmer project exceeds the complexity of this project, so it was a great candidate for
using a well tested, specialized external library.

### Analysis

#### Performance

Performance in the training phase is not a primary focus of this type of project. Ideally,
one can take as much time as they need to train a model. The application of the model is 
more important, but once you have a model, you don't need to recalculate it. This application 
takes about 0.86 seconds to run through all 4180 messages with five different preprocessors 
on a computer built in 2013. 

#### Output

```
Analysis: Default Analysis (no preprocessing)
Vocabulary has 13125 words

Training Set:
        564 of 4180 messages were spam (13.49%)

Test Set:
        Correct Ham: 1204
        Correct Spam: 136
        Incorrect Ham (actually was spam): 47
        Incorrect Spam (actually was ham): 7
        Percentage Correct Ham: 99.42%
        Percentage Correct Spam: 74.32%
        Overall Accuracy: 96.13%

Analysis: No Punctuation Analysis
Vocabulary has 9877 words

Training Set:
        564 of 4180 messages were spam (13.49%)

Test Set:
        Correct Ham: 1207
        Correct Spam: 137
        Incorrect Ham (actually was spam): 46
        Incorrect Spam (actually was ham): 4
        Percentage Correct Ham: 99.67%
        Percentage Correct Spam: 74.86%
        Overall Accuracy: 96.41%

Analysis: Stemmer Analysis
Vocabulary has 6972 words

Training Set:
        564 of 4180 messages were spam (13.49%)

Test Set:
        Correct Ham: 1208
        Correct Spam: 120
        Incorrect Ham (actually was spam): 63
        Incorrect Spam (actually was ham): 3
        Percentage Correct Ham: 99.75%
        Percentage Correct Spam: 65.57%
        Overall Accuracy: 95.27%

Analysis: Stemmer and No Punctuation Analysis
Vocabulary has 6953 words

Training Set:
        564 of 4180 messages were spam (13.49%)

Test Set:
        Correct Ham: 1208
        Correct Spam: 120
        Incorrect Ham (actually was spam): 63
        Incorrect Spam (actually was ham): 3
        Percentage Correct Ham: 99.75%
        Percentage Correct Spam: 65.57%
        Overall Accuracy: 95.27%

Analysis: Remove 100 Most Common English Words
Vocabulary has 6867 words

Training Set:
        564 of 4180 messages were spam (13.49%)

Test Set:
        Correct Ham: 1198
        Correct Spam: 153
        Incorrect Ham (actually was spam): 30
        Incorrect Spam (actually was ham): 13
        Percentage Correct Ham: 98.93%
        Percentage Correct Spam: 83.61%
        Overall Accuracy: 96.92%


Done.

```

#### Effectiveness and Accuracy

This application is fairly accurate at predicting ham or spam classes for an 
SMS message. If one guessed ham for every message, the accuracy of this static 
approach would be just over 86% due to the total frequency of ham. However, 
this application is able to make a smart prediction with about 95% accuracy. 
Additionally, when mistakes are made, they are only rarely false positives for 
spam, so the user will not be losing large quantities of valid email. The user 
will still need to manually go through some spam, but the majority is caught.  

### License, Limitations, and Usage

The purpose of this project was academic in nature. I don't recommend using 
this code in production. I have licensed this code with an MIT license, so 
reuse is permissible. If you are in an academic institution, you might have 
additional guidelines to follow.

### Author Details

Daniel Wolf