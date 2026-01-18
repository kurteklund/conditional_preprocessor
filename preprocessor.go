package main

import (
	"log"
	"regexp"
	"slices"
)

var conditionalRegionRegExp = regexp.MustCompile(`\{\{#if\s*(!)?\s*(\w+)\s*\}\}((.|\n)*)\{\{#endif}}`)

func processSections(bookSections *MdBookTopItem, conditionalRegions []string) {
	for i := range bookSections.Sections {
		var section = &bookSections.Sections[i]
		processSection(section, conditionalRegions)
	}
}

func processSection(section *MdBookSection, conditionalRegions []string) {
	processChapter(&section.Chapter, conditionalRegions)
	for i := range section.Chapter.SubItems {
		subSection := &section.Chapter.SubItems[i]
		processSection(subSection, conditionalRegions)
	}
}

func processChapter(chapter *MdBookChapter, conditionalRegions []string) {
	chapter.Content = ProcessConditionalRegions(chapter.Content, conditionalRegions)
}

func ProcessConditionalRegions(text string, conditionalRegions []string) string {
	var regexpIndexes = conditionalRegionRegExp.FindStringSubmatchIndex(text)

	if regexpIndexes == nil {
		return text
	}

	if len(regexpIndexes) == 10 {
		var expression = text[regexpIndexes[0]:regexpIndexes[1]]
		var notIndication = regexpIndexes[2] != regexpIndexes[3]
		var regionName = text[regexpIndexes[4]:regexpIndexes[5]]
		var regionText = text[regexpIndexes[6]:regexpIndexes[7]]
		log.Println("Expression: " + expression)
		log.Printf("NotIndication: %v", notIndication)
		log.Println("Region Name: " + regionName)
		log.Println("Section Text: " + regionText)

		var showRegionText = false

		if notIndication {
			showRegionText = !slices.Contains(conditionalRegions, regionName)
		} else {
			showRegionText = slices.Contains(conditionalRegions, regionName)
		}

		// Build the result string!
		var result = text[0:regexpIndexes[0]] // Text before the conditional stuff
		if showRegionText {
			result += regionText
		}

		result += text[regexpIndexes[1]:]

		return result
	}

	return text
}
