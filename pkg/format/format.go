package format

import (
	s "strings"

	"github.com/Kory291/gherkin-formatter/pkg/configuration"
)

func FormatFile(fileContent []string, configuration configuration.Config) ([]string, error) {
	currentIntendation := 0
	formattedFileContents := make([]string, 0)
	var previousFoundElement string

	for _, line := range fileContent {
		cutLine := s.Trim(line, " ")
		given := s.HasPrefix(s.ToLower(cutLine), "given")
		when := s.HasPrefix(s.ToLower(cutLine), "when")
		then := s.HasPrefix(s.ToLower(cutLine), "then")
		and := s.HasPrefix(s.ToLower(cutLine), "and")
		if ((previousFoundElement == "feature") || 
			(previousFoundElement == "scenario") || 
			(previousFoundElement == "background") ||
			(given && configuration.IntendAnd && previousFoundElement != "given")) {
			currentIntendation += 1
		}
		scenario := s.HasPrefix(s.ToLower(cutLine), "scenario:")
		feature := s.HasPrefix(s.ToLower(cutLine), "feature:")
		if (previousFoundElement == "and" && scenario && configuration.IntendAnd) && currentIntendation > 0{
			currentIntendation -= 1
		}
		if scenario && currentIntendation > 0 {
			currentIntendation -= 1
		}
		newLine := s.Repeat(" ", currentIntendation * configuration.Intendation) + cutLine
		formattedFileContents = append(formattedFileContents, newLine)
		if given {
			previousFoundElement = "given"
		} else if when {
			previousFoundElement = "when"
		} else if then {
			previousFoundElement = "then"
		} else if and {
			previousFoundElement = "and"
		} else if scenario {
			previousFoundElement = "scenario"
		} else if feature {
			previousFoundElement = "feature"
		} else {
			previousFoundElement = ""
		}
	}
	return formattedFileContents, nil
}