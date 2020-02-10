package experiment

import (
	"fmt"
)

type Class int

func (c Class) String() string {
	switch c {
	case HamClass:
		return "ham"
	case SpamClass:
		return "spam"
	default:
		return ""
	}
}

func ClassType(str string) (Class, error) {
	switch str {
	case HamClass.String():
		return HamClass, nil
	case SpamClass.String():
		return SpamClass, nil
	default:
		return HamClass, fmt.Errorf("invalid class: %s", str)
	}
}

const (
	HamClass Class = iota
	SpamClass
)

type Experiment struct {
	Classes Classes
	Test    TestSet
}

type Classes struct {
	Ham  []string
	Spam []string
}

type TestSet struct {
	Cases []TestCase
}

type TestCase struct {
	Class Class
	Text  string
}
