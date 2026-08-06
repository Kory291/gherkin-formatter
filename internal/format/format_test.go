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

func TestIncreaseIntendation(t *testing.T) {
	// Define test cases
	testCases := []struct {
		currentElement  Element
		previousElement Element
		expected        bool
	}{
		{ElementFeature, ElementEmpty, false},
		{ElementScenario, ElementFeature, true},
		{ElementGiven, ElementScenario, true},
		{ElementWhen, ElementGiven, false},
		{ElementThen, ElementWhen, false},
		{ElementAnd, ElementThen, false},
		{ElementBackground, ElementFeature, true},
		{ElementExamples, ElementScenario, true},
	}

	for _, tc := range testCases {
		actual := increaseIntendation(tc.currentElement, tc.previousElement, configuration.Config{})
		if actual != tc.expected {
			t.Errorf("For currentElement %v and previousElement %v, expected %v but got %v", tc.currentElement, tc.previousElement, tc.expected, actual)
		}
	}
}

func TestDecreaseIntendation(t *testing.T) {
	// Define test cases
	testCases := []struct {
		currentElement  Element
		previousElement Element
		expected        bool
	}{
		{ElementGiven, ElementAnd, true},
		{ElementWhen, ElementAnd, true},
		{ElementThen, ElementAnd, true},
		{ElementAnd, ElementGiven, false},
		{ElementAnd, ElementWhen, false},
		{ElementEmpty, ElementAnd, false},
		{ElementScenario, ElementThen, true},
		{ElementGiven, ElementGiven, false},
		{ElementWhen, ElementWhen, false},
		{ElementThen, ElementThen, false},
	}

	for _, tc := range testCases {
		actual := decreaseIntendation(tc.currentElement, tc.previousElement, configuration.Config{IntendAnd: true})
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