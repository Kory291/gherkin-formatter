package format

import (
	re "regexp"
	"slices"
	s "strings"

	"github.com/Kory291/gherkin-formatter/internal/configuration"
)

type Element int

const (
	ElementGiven Element = iota
	ElementWhen
	ElementThen
	ElementAnd
	ElementFeature
	ElementScenario
	ElementBackground
	ElementExamples
	ElementDescription
	ElementTag
	ElementEmpty
	ElementTable
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

	// currentELementMatcher := re.MustCompile(`^(given\s|when\s|then\s|and\s|feature:|scenario( outline)?:|background:|examples)`)
	// match := currentELementMatcher.FindString(line)
	// if match != "" {
	// 	match = s.TrimSuffix(match, ":")
	// 	return s.TrimSuffix(match, " outline")
	// }
	// tagMatcher := re.MustCompile(`^@[\d\w_.-]+`)
	// if tagMatcher.MatchString(line) {
	// 	return "tag"
	// }
	// tableMatcher := re.MustCompile(ElementRegex[ElementTable])
	// if tableMatcher.MatchString(line) {
	// 	return "table"
	// }
	// return "description"
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
	if configuration.IntendAnd && previousElement == ElementAnd && currentElement != ElementAnd {
		return true
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

	for lineNumber, line := range fileContent {
		cutLine := s.Trim(line, " ")

		if cutLine == "" {
			continue
		}

		tags := []string{}

		currentElement := getCurrentGherkinElement(cutLine)
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

		if currentElement == ElementTag && previousFoundElement == ElementTag {
			continue
		}

		// check if indentation has to be increased
		if increaseIntendation(currentElement, previousFoundElement, configuration) {
			// fmt.Println("..Increasing intendation")
			currentIntendation += 1
		}

		// check if intendation has to be decreased
		if decreaseIntendation(currentElement, previousFoundElement, configuration) && currentIntendation > 1 {
			// fmt.Println("..Decreasing intendation")
			currentIntendation -= 1
		}

		if addNewLine(currentElement, previousFoundElement) {
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

		if currentElement != ElementEmpty {
			previousFoundElement = currentElement
		}
	}
	return formattedFileContents, nil
}
