package format

import (
	// "fmt"
	re "regexp"
	s "strings"

	"github.com/Kory291/gherkin-formatter/pkg/configuration"
)

func getCurrentGherkinElement(line string) string {
	line = s.ToLower(s.Trim(line, " "))
	currentELementMatcher := re.MustCompile(`^(given\s|when\s|then\s|and\s|feature:|scenario( outline)?:|background:|examples)`)
	match := currentELementMatcher.FindString(line)
	match = s.TrimSuffix(match, ":")
	return s.TrimSuffix(match, " outline")
}

func increaseIntendation(line string, currentElement string, previousElement string, configuration configuration.Config) bool {
	// find in which line we are
	// this is important if we have a change in the following cases:
	// Feature name -> Feature description
	// Feature -> Scenario
	// Scenario -> Given | When | Then
	line = s.Trim(line, " ")
	if line == "" {
		return false
	}
	if previousElement == "feature" || previousElement == "scenario" || previousElement == "background" || previousElement == "examples" {
		return true
	}
	if !configuration.IntendAnd {
		return false
	}
	return (currentElement == "and") && (previousElement != "and")
}

func decreaseIntendation(line string, currentElement string, configuration configuration.Config) bool {
	if !configuration.IntendAnd {
		return false
	}
	line = s.Trim(line, " ")
	if line == "" {
		return false
	}
	return currentElement == "scenario" || currentElement == "examples"
}

func FormatFile(fileContent []string, configuration configuration.Config) ([]string, error) {
	currentIntendation := 0
	formattedFileContents := make([]string, 0)
	// tagMatcher := re.MustCompile(`@[\d\w_-]+`)
	var previousFoundElement string


	for _, line := range fileContent {
		cutLine := s.Trim(line, " ")

		currentElement := getCurrentGherkinElement(cutLine)
		if increaseIntendation(cutLine, currentElement, previousFoundElement, configuration) {
			// fmt.Println("..Increasing intendation")
			currentIntendation += 1
		}
		if increaseIntendation(cutLine, currentElement, previousFoundElement, configuration) && currentElement == "scenario" {
			// fmt.Println("..Increasing intendation")
			currentIntendation += 1
		}
		if decreaseIntendation(line, currentElement, configuration) && currentIntendation > 1 {
			// fmt.Println("..Decreasing intendation")
			currentIntendation -= 1
		}
		if currentElement == "and" && decreaseIntendation(line, currentElement, configuration) && currentIntendation > 1 {
			// fmt.Println("..Decreasing intendation")
			currentIntendation -= 1
		}
		newLine := s.Repeat(" ", currentIntendation * configuration.Intendation) + cutLine
		formattedFileContents = append(formattedFileContents, newLine)
		previousFoundElement = currentElement
	}
	return formattedFileContents, nil
}