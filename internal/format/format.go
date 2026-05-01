package format

import (
	// "fmt"
	re "regexp"
	"slices"
	s "strings"

	"github.com/Kory291/gherkin-formatter/internal/configuration"
)

func getCurrentGherkinElement(line string) string {
	line = s.ToLower(s.Trim(line, " "))
	if line == "" {
		return ""
	}
	currentELementMatcher := re.MustCompile(`^(given\s|when\s|then\s|and\s|feature:|scenario( outline)?:|background:|examples)`)
	match := currentELementMatcher.FindString(line)
	if match != "" {
		match = s.TrimSuffix(match, ":")
		return s.TrimSuffix(match, " outline")
	}
	tagMatcher := re.MustCompile(`^@[\d\w_.-]+`)
	if tagMatcher.MatchString(line) {
		return "tag"
	}
	tableMatcher := re.MustCompile(`^\|`)
	if tableMatcher.MatchString(line) {
		return "table"
	}
	return "description"
}

func increaseIntendation(line string, currentElement string, previousElement string, configuration configuration.Config) bool {
	// find in which line we are
	// this is important if we have a change in the following cases:
	// Feature name -> Feature description
	// Feature -> Scenario
	// Scenario -> Given | When | Then

	// Special case for tags:
	// if a tag was before a scenario, we do not want to increase intendation for the scenario
	line = s.Trim(line, " ")
	if line == "" {
		return false
	}
	if currentElement == previousElement {
		return false
	}
	if currentElement == "scenario" && previousElement == "tag" {
		return false
	}
	if (currentElement == "scenario" || currentElement == "tag") && previousElement == "description" {
		return false
	}
	if currentElement == "table" && previousElement != "table" {
		return true
	}
	if previousElement == "feature" || previousElement == "scenario" || previousElement == "background" || previousElement == "examples" {
		return true
	}
	if !configuration.IntendAnd {
		return false
	}
	return (currentElement == "and") && (previousElement != "and")
}

func decreaseIntendation(line string, currentElement string, previousElement string, configuration configuration.Config) bool {
	if !configuration.IntendAnd {
		return false
	}
	line = s.Trim(line, " ")
	if line == "" {
		return false
	}
	if previousElement == "table" && currentElement != "table" {
		return true
	}
	return currentElement == "scenario" || currentElement == "examples" || currentElement == "tag"
}

func addNewLine(currentElement string, previousElement string, configuration configuration.Config) bool {
	return (previousElement != currentElement) && (previousElement != "tag") && (currentElement == "scenario" || currentElement == "background" || currentElement == "examples" || currentElement == "tag")
}

func FormatFile(fileContent []string, configuration configuration.Config) ([]string, error) {
	currentIntendation := 0
	formattedFileContents := make([]string, 0)
	// tagMatcher := re.MustCompile(`@[\d\w_-.]+`)

	var previousFoundElement string

	for lineNumber, line := range fileContent {
		cutLine := s.Trim(line, " ")

		if cutLine == "" {
			continue
		}

		tags := []string{}

		currentElement := getCurrentGherkinElement(cutLine)
		// see if there are more tags in the following lines
		if currentElement == "tag" && previousFoundElement != "tag" {
			tagsMatches := re.MustCompile(`@[\d\w_.-]+`)

			// go to next lines
			for _, nextLine := range fileContent[lineNumber:] {
				lineTags := tagsMatches.FindAllString(nextLine, -1)

				tags = append(tags, lineTags...)

				nextElement := getCurrentGherkinElement(nextLine)
				// no tag following anymore can do other stuff
				if nextElement != "tag" {
					break
				}
			}
			if configuration.SortTags {
				slices.Sort(tags)
			}
		}

		if currentElement == "tag" && previousFoundElement == "tag" {
			continue
		}

		// check if indentation has to be increased
		if increaseIntendation(cutLine, currentElement, previousFoundElement, configuration) {
			// fmt.Println("..Increasing intendation")
			currentIntendation += 1
		}

		// check if intendation has to be decreased
		if decreaseIntendation(line, currentElement, previousFoundElement, configuration) && currentIntendation > 1 {
			// fmt.Println("..Decreasing intendation")
			currentIntendation -= 1
		}

		if addNewLine(currentElement, previousFoundElement, configuration) {
			formattedFileContents = append(formattedFileContents, "")
		}

		// set the new line with the required numbers of whitespaces
		newLine := s.Repeat(" ", currentIntendation*configuration.Intendation) + cutLine

		if len(tags) > 0 {
			for _, tag := range tags {
				newLine := s.Repeat(" ", currentIntendation*configuration.Intendation) + tag
				formattedFileContents = append(formattedFileContents, newLine)
			}
		} else {
			formattedFileContents = append(formattedFileContents, newLine)
		}

		if currentElement != "" {
			previousFoundElement = currentElement
		}
	}
	return formattedFileContents, nil
}
