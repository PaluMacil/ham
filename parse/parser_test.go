package parse_test

import (
	"path"
	"strings"
	"testing"

	"github.com/PaluMacil/ham/parse"
)

func TestParse(t *testing.T) {
	shortText := strings.NewReader(`
spam	request a ham for me
ham	buy buy buy me a ring
ham	I love carrots!
spam	how now, brown cow?
`)
	experiment, err := parse.Parse(shortText, "\t")
	if err != nil {
		t.Errorf("parsing shortText: %s", err)
	}
	hamPlusSpamLength := len(experiment.Classes.Ham) + len(experiment.Classes.Spam)
	if hamPlusSpamLength != 3 {
		t.Errorf("counting ham plus spam: expected 3, got %d", hamPlusSpamLength)
	}
	if len(experiment.Test.Cases) != 1 {
		t.Errorf("counting test cases: expected 1, got %d", len(experiment.Test.Cases))
	}
}

const filename = "textMsgs.data"

func TestFromFile(t *testing.T) {
	experiment, err := parse.FromFile(path.Join("..", filename), "\t")
	if err != nil {
		t.Errorf("parsing %s: %s", filename, err)
	}
	totalCasesAllTypes := len(experiment.Classes.Ham) + len(experiment.Classes.Spam) + len(experiment.Test.Cases)
	if totalCasesAllTypes != 5574 {
		t.Errorf("counting all types of cases: expected 5574, got %d", totalCasesAllTypes)
	}
}
