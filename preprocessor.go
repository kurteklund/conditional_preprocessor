package main

import (
	"regexp"
	"slices"
	"strings"
)

var conditionalRegionRegExp = regexp.MustCompile(`{{if\s*(!)?\s*(\w+)\s*\}}((.|\n)+?){{endif}}`)

func processSections(bookSections *MdBookTopItem, conditionalRegions []string, varNameAndValues []VarNameAndValue) {
	for i := range bookSections.Sections {
		var section = &bookSections.Sections[i]
		processSection(section, conditionalRegions, varNameAndValues)
	}
}

func processSection(section *MdBookSection, conditionalRegions []string, varNameAndValues []VarNameAndValue) {
	processChapter(&section.Chapter, conditionalRegions, varNameAndValues)
	for i := range section.Chapter.SubItems {
		subSection := &section.Chapter.SubItems[i]
		processSection(subSection, conditionalRegions, varNameAndValues)
	}
}

func processChapter(chapter *MdBookChapter, conditionalRegions []string, varNameAndValues []VarNameAndValue) {
	chapter.Content = processVariables(chapter.Content, varNameAndValues)
	chapter.Content = processConditionalRegions(chapter.Content, conditionalRegions)
}

func processVariables(text string, varNameAndValues []VarNameAndValue) string {
	for _, varNameAndValue := range varNameAndValues {
		var textToReplace = "{{" + varNameAndValue.Name + "}}"
		text = strings.ReplaceAll(text, textToReplace, varNameAndValue.Value)
	}

	return text
}

func processConditionalRegions(text string, conditionalRegions []string) string {
	text, replaced := replaceFirstRegion(text, conditionalRegions)

	for replaced == true {
		text, replaced = replaceFirstRegion(text, conditionalRegions)
	}

	return text
}

func replaceFirstRegion(text string, conditionalRegions []string) (string, bool) {
	var regexpIndexes = conditionalRegionRegExp.FindStringSubmatchIndex(text)

	if regexpIndexes == nil {
		return text, false
	}

	if len(regexpIndexes) == 10 {
		// var expression = text[regexpIndexes[0]:regexpIndexes[1]]
		var notIndication = regexpIndexes[2] != regexpIndexes[3]
		var regionName = text[regexpIndexes[4]:regexpIndexes[5]]
		var regionText = text[regexpIndexes[6]:regexpIndexes[7]]
		// log.Println("Expression: " + expression)
		// log.Printf("NotIndication: %v", notIndication)
		// log.Println("Region Name: " + regionName)
		// log.Println("Section Text: " + regionText)

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

		return result, true
	}

	return text, false
}
