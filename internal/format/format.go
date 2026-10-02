package format

import (
	"log/slog"
	re "regexp"
	"slices"
	s "strings"

	"github.com/Kory291/gherkin-formatter/internal/configuration"
)

type Element string

const (
	ElementGiven       Element = "Given"
	ElementWhen        Element = "When"
	ElementThen        Element = "Then"
	ElementAnd         Element = "And"
	ElementFeature     Element = "Feature"
	ElementScenario    Element = "Scenario"
	ElementBackground  Element = "Background"
	ElementExamples    Element = "Examples"
	ElementDescription Element = "Description"
	ElementTag         Element = "Tag"
	ElementEmpty       Element = "Empty"
	ElementTable       Element = "Table"
	ElementComment     Element = "Comment"
)

var ElementRegex = map[Element]string{
	ElementGiven:       `^given\s`,
	ElementWhen:        `^when\s`,
	ElementThen:        `^then\s`,
	ElementAnd:         `^and\s`,
	ElementFeature:     `^feature:`,
	ElementScenario:    `^scenario( outline)?:`,
	ElementBackground:  `^background:`,
	ElementExamples:    `^examples`,
	ElementComment:     `^\s*#`,
	ElementDescription: ``,
	ElementTag:         `^@[\d\w_.-]`,
	ElementTable:       `^\|`,
	ElementEmpty:       ``,
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

func increaseIndentation(currentElement Element, previousElement Element, configuration configuration.Config) int {
	// find in which line we are
	// this is important if we have a change in the following cases:
	// Feature name -> Feature description
	// Feature -> Scenario
	// Scenario -> Given | When | Then

	// Special case for tags:
	// if a tag was before a scenario, we do not want to increase indentation for the scenario
	if currentElement == ElementEmpty || currentElement == ElementComment {
		return 0
	}
	if currentElement == previousElement {
		return 0
	}
	if currentElement == ElementScenario && previousElement == ElementTag {
		return 0
	}
	if (currentElement == ElementScenario || currentElement == ElementTag) && previousElement == ElementDescription {
		return 0
	}
	if currentElement == ElementScenario && previousElement == ElementTable {
		return 0
	}
	if currentElement == ElementTable && previousElement != ElementTable {
		return 1
	}
	if previousElement == ElementFeature || previousElement == ElementScenario || previousElement == ElementBackground || previousElement == ElementExamples {
		return 1
	}
	if !configuration.IndentAnd {
		return 0
	}
	// if currentElement == ElementAnd && previousElement != ElementAnd && previousElement != ElementTable {
	if currentElement == ElementAnd && previousElement != ElementAnd {
		return 1
	}
	return 0
}

func decreaseIndentation(currentElement Element, previousElement Element, configuration configuration.Config) int {
	if currentElement == ElementEmpty || currentElement == ElementComment {
		return 0
	}
	if configuration.IndentAnd && previousElement == ElementAnd {
		if currentElement == ElementTable {
			return 0
		} else if currentElement == ElementScenario || currentElement == ElementExamples || currentElement == ElementTag {
			return 2
		} else if currentElement != ElementAnd {
			return 1
		}
	}
	if previousElement == ElementTable && currentElement != ElementTable {
		return 1
	}
	if (currentElement == ElementScenario || currentElement == ElementTag) && (previousElement == ElementDescription || previousElement == ElementFeature || previousElement == ElementTag) {
		return 0
	}
	if currentElement == ElementScenario || currentElement == ElementExamples || currentElement == ElementTag {
		return 1
	}
	return 0
}

func addNewLine(currentElement Element, previousElement Element) bool {
	return (previousElement != currentElement) && (previousElement != ElementTag) && (currentElement == ElementScenario || currentElement == ElementBackground || currentElement == ElementExamples || currentElement == ElementTag)
}

func FormatFile(fileContent []string, configuration configuration.Config) ([]string, error) {
	currentIndentation := 0
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
		if indentationChange := increaseIndentation(currentElement, previousFoundElement, configuration); indentationChange > 0 {
			slog.Debug("Increasing indentation from ", "currentIndentation", currentIndentation, "indentationChange", indentationChange)
			if tableSource == "And" && currentElement == ElementAnd && configuration.IndentAnd {
				indentationChange -= 1
			}
			currentIndentation += indentationChange
		}

		// check if indentation has to be decreased
		if indentationChange := decreaseIndentation(currentElement, previousFoundElement, configuration); indentationChange > 0 {
			if (currentElement == ElementGiven || currentElement == ElementWhen || currentElement == ElementThen) && tableSource == "And" && configuration.IndentAnd {
				indentationChange += 1
			}
			if currentElement == ElementExamples || currentElement == ElementScenario || currentElement == ElementTag {
				if tableSource == "Step" {
					indentationChange += 1
				} else if tableSource == "And" && configuration.IndentAnd {
					indentationChange += 2
				} else if tableSource == "And" && !configuration.IndentAnd {
					indentationChange += 1
				}
			}
			currentIndentation -= indentationChange
			if currentIndentation < 0 {
				currentIndentation = 0
			}
			slog.Debug("Decreased indentation to ", "currentIntedantion", currentIndentation)
		}

		if addNewLine(currentElement, previousFoundElement) {
			formattedFileContents = append(formattedFileContents, "")
		}

		// set the new line with the required numbers of whitespaces
		slog.Debug("Write line with indentation", "line", cutLine, "currentIndentation", currentIndentation)
		newLine := s.Repeat(" ", currentIndentation*configuration.Indentation) + cutLine

		if len(tags) > 0 {
			for _, tag := range tags {
				newLine := s.Repeat(" ", currentIndentation*configuration.Indentation) + tag
				formattedFileContents = append(formattedFileContents, newLine)
			}
		} else {
			formattedFileContents = append(formattedFileContents, newLine)
		}
		if currentElement == ElementTable {
			if previousFoundElement == ElementExamples {
				tableSource = "Examples"
			} else if slices.Contains([]Element{ElementGiven, ElementWhen, ElementThen}, previousFoundElement) {
				tableSource = "Step"
			} else if previousFoundElement == ElementAnd {
				tableSource = "And"
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
