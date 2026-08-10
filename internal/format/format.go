package format

import (
	re "regexp"
	"slices"
	s "strings"
	"log/slog"

	"github.com/Kory291/gherkin-formatter/internal/configuration"
)

type Element string

const (
	ElementGiven Element = "Given"
	ElementWhen Element = "When"
	ElementThen Element = "Then"
	ElementAnd Element = "And"
	ElementFeature Element = "Feature"
	ElementScenario Element = "Scenario"
	ElementBackground Element = "Background"
	ElementExamples Element = "Examples"
	ElementDescription Element = "Description"
	ElementTag Element = "Tag"
	ElementEmpty Element = "Empty"
	ElementTable Element = "Table"
)

var ElementRegex = map[Element]string{
	ElementGiven: `^given\s`,
	ElementWhen: `^when\s`,
	ElementThen: `^then\s`,
	ElementAnd: `^and\s`,
	ElementFeature: `^feature:`,
	ElementScenario: `^scenario( outline)?:`,
	ElementBackground: `^background:`,
	ElementExamples: `^examples`,
	ElementDescription: ``,
	ElementTag: `^@[\d\w_.-]`,
	ElementTable: `^\|`,
	ElementEmpty: ``,
}

func getCurrentGherkinElement(line string) Element {
	line = s.ToLower(s.Trim(line, " "))
	if line == "" {
		return ElementEmpty
	}

	for element, regex := range ElementRegex {
		elementMatcher := re.MustCompile(regex)
		match := elementMatcher.FindString(line)
		if match == "" {
			continue
		}
		return element
	}  

	return ElementDescription
}

func increaseIntendation(currentElement Element, previousElement Element, configuration configuration.Config) bool {
	// find in which line we are
	// this is important if we have a change in the following cases:
	// Feature name -> Feature description
	// Feature -> Scenario
	// Scenario -> Given | When | Then

	// Special case for tags:
	// if a tag was before a scenario, we do not want to increase intendation for the scenario
	if currentElement == ElementEmpty {
		return false
	}
	if currentElement == previousElement {
		return false
	}
	if currentElement == ElementScenario && previousElement == ElementTag {
		return false
	}
	if (currentElement == ElementScenario || currentElement == ElementTag) && previousElement == ElementDescription {
		return false
	}
	if currentElement == ElementScenario && previousElement == ElementTable {
		return false
	}
	if currentElement == ElementTable && previousElement != ElementTable {
		return true
	}
	if previousElement == ElementFeature || previousElement == ElementScenario || previousElement == ElementBackground || previousElement == ElementExamples {
		return true
	}
	if !configuration.IntendAnd {
		return false
	}
	return (currentElement == ElementAnd) && (previousElement != ElementAnd)
}

func decreaseIntendation(currentElement Element, previousElement Element, configuration configuration.Config) bool {
	if currentElement == ElementEmpty {
		return false
	}
	if configuration.IntendAnd && previousElement == ElementAnd {
		if currentElement == ElementTable {
			return false
		} else if currentElement != ElementAnd {
			return true
		}	
	}
	if previousElement == ElementTable && currentElement != ElementTable {
		return true
	}
	return currentElement == ElementScenario || currentElement == ElementExamples || currentElement == ElementTag
}

func addNewLine(currentElement Element, previousElement Element) bool {
	return (previousElement != currentElement) && (previousElement != ElementTag) && (currentElement == ElementScenario || currentElement == ElementBackground || currentElement == ElementExamples || currentElement == ElementTag)
}

func FormatFile(fileContent []string, configuration configuration.Config) ([]string, error) {
	currentIntendation := 0
	formattedFileContents := make([]string, 0)

	var previousFoundElement Element
	var tableSource string

	for lineNumber, line := range fileContent {
		cutLine := s.Trim(line, " ")

		if cutLine == "" {
			continue
		}

		tags := []string{}
		slog.Debug("Working on line:", "cutLine", cutLine)
		currentElement := getCurrentGherkinElement(cutLine)
		// set source for table - either step or example
				// see if there are more tags in the following lines
		if currentElement == ElementTag && previousFoundElement != ElementTag {
			tagsMatches := re.MustCompile(`@[\d\w_.-]+`)

			// go to next lines
			for _, nextLine := range fileContent[lineNumber:] {
				lineTags := tagsMatches.FindAllString(nextLine, -1)

				tags = append(tags, lineTags...)

				nextElement := getCurrentGherkinElement(nextLine)
				// no tag following anymore can do other stuff
				if nextElement != ElementTag {
					break
				}
			}
			if configuration.SortTags {
				slices.Sort(tags)
			}
		}

		slog.Debug("Elements: ", "currentElement", currentElement, "previousFoundElement", previousFoundElement)
		if currentElement == ElementTag && previousFoundElement == ElementTag {
			continue
		}

		// check if indentation has to be increased
		if increaseIntendation(currentElement, previousFoundElement, configuration) {
			currentIntendation += 1
			slog.Debug("Increase Intendation to", "currentIntendation", currentIntendation)
		}

		// check if intendation has to be decreased
		if decreaseIntendation(currentElement, previousFoundElement, configuration) && currentIntendation > 1 {
			slog.Debug("current tableSource", "tableSource", tableSource)
			if tableSource == "Step" {
				currentIntendation -= 2
			} else {
				currentIntendation -= 1
			}
			slog.Debug("Decrease Intendation to", "currentIntendation", currentIntendation)
		}

		if addNewLine(currentElement, previousFoundElement) {
			formattedFileContents = append(formattedFileContents, "")
		}

		// set the new line with the required numbers of whitespaces
		slog.Debug("Write line with intendation", "line", cutLine, "currentIntendation", currentIntendation)
		newLine := s.Repeat(" ", currentIntendation*configuration.Intendation) + cutLine

		if len(tags) > 0 {
			for _, tag := range tags {
				newLine := s.Repeat(" ", currentIntendation*configuration.Intendation) + tag
				formattedFileContents = append(formattedFileContents, newLine)
			}
		} else {
			formattedFileContents = append(formattedFileContents, newLine)
		}
		if currentElement == ElementTable {
			if previousFoundElement == ElementExamples {
				tableSource = "Examples"
			} else if slices.Contains([]Element{ElementGiven, ElementWhen, ElementThen, ElementAnd}, previousFoundElement) {
				tableSource = "Step"
			}
		} else if currentElement != ElementEmpty {
			if previousFoundElement == ElementTable {
				slog.Debug("Unsetting tableSource for", "previousFoundElement", previousFoundElement, "currentElement", currentElement)
				tableSource = ""
			}
		}
		if currentElement == ElementEmpty {
			slog.Debug("Skipping something becuase of empty line")
			continue
		}
		previousFoundElement = currentElement
	}
	return formattedFileContents, nil
}
