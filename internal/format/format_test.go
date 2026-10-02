package format

import (
	"testing"

	"github.com/Kory291/gherkin-formatter/internal/configuration"
)

func TestGetCurrentGherkinElement(t *testing.T) {
	testValues := map[Element]string{
		ElementGiven:       "Given I have a precondition",
		ElementWhen:        "When I perform an action",
		ElementThen:        "Then I expect a result",
		ElementAnd:         "And I have another step",
		ElementFeature:     "Feature: My feature",
		ElementScenario:    "Scenario: My scenario",
		ElementBackground:  "Background: My background",
		ElementExamples:    "Examples:",
		ElementDescription: "This is a description.",
		ElementTag:         "@mytag",
		ElementTable:       "| Column1 | Column2 |",
		ElementEmpty:       "",
	}
	for expectedElement, line := range testValues {
		actualElement := getCurrentGherkinElement(line)
		if actualElement != expectedElement {
			t.Errorf("Expected element %v for line '%s', but got %v", expectedElement, line, actualElement)
		}
	}
}

func TestIncreaseIndentation(t *testing.T) {
	// Define test cases
	testCases := []struct {
		currentElement  Element
		previousElement Element
		expected        int
	}{
		{ElementFeature, ElementEmpty, 0},
		{ElementScenario, ElementFeature, 1},
		{ElementGiven, ElementScenario, 1},
		{ElementWhen, ElementGiven, 0},
		{ElementThen, ElementWhen, 0},
		{ElementAnd, ElementThen, 0},
		{ElementBackground, ElementFeature, 1},
		{ElementExamples, ElementScenario, 1},
	}

	for _, tc := range testCases {
		actual := increaseIndentation(tc.currentElement, tc.previousElement, configuration.Config{})
		if actual != tc.expected {
			t.Errorf("For currentElement %v and previousElement %v, expected %v but got %v", tc.currentElement, tc.previousElement, tc.expected, actual)
		}
	}
}

func TestDecreaseIndentation(t *testing.T) {
	// Define test cases
	testCases := []struct {
		currentElement  Element
		previousElement Element
		expected        int
	}{
		{ElementGiven, ElementAnd, 1},
		{ElementWhen, ElementAnd, 1},
		{ElementThen, ElementAnd, 1},
		{ElementAnd, ElementGiven, 0},
		{ElementAnd, ElementWhen, 0},
		{ElementEmpty, ElementAnd, 0},
		{ElementScenario, ElementThen, 1},
		{ElementGiven, ElementGiven, 0},
		{ElementWhen, ElementWhen, 0},
		{ElementThen, ElementThen, 0},
	}

	for _, tc := range testCases {
		actual := decreaseIndentation(tc.currentElement, tc.previousElement, configuration.Config{IndentAnd: true})
		if actual != tc.expected {
			t.Errorf("For currentElement %v and previousElement %v, expected %v but got %v", tc.currentElement, tc.previousElement, tc.expected, actual)
		}
	}
}

func TestAddNewLine(t *testing.T) {
	// Define test cases
	testCases := []struct {
		currentElement  Element
		previousElement Element
		expected        bool
	}{
		{ElementScenario, ElementFeature, true},
		{ElementBackground, ElementFeature, true},
		{ElementExamples, ElementScenario, true},
		{ElementTag, ElementFeature, true},
		{ElementGiven, ElementScenario, false},
		{ElementWhen, ElementGiven, false},
		{ElementThen, ElementWhen, false},
	}

	for _, tc := range testCases {
		actual := addNewLine(tc.currentElement, tc.previousElement)
		if actual != tc.expected {
			t.Errorf("For currentElement %v and previousElement %v, expected %v but got %v", tc.currentElement, tc.previousElement, tc.expected, actual)
		}
	}
}
