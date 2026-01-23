package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"path"
)

func check(e error) {
	if e != nil {
		log.Printf("Error: (%v)\n", e)
		panic(e)
	}
}

func CreateTextFile(filetPath string, content string) {
	f, err := os.Create(filetPath)
	check(err)

	n, err := f.WriteString(content + "\n")
	check(err)
	if n < 0 {
		panic("tRAMS")
	}
	f.Sync()
}

func ReadTextFile(path string) string {
	result := ""
	f, _ := os.Open(path)
	// Create a new Scanner for the file.
	scanner := bufio.NewScanner(f)
	// Loop over all lines in the file and print them.
	for scanner.Scan() {
		result += scanner.Text()
	}

	return result
}

func readFileLinesFile(filePath string) []string {
	readFile, err := os.Open(filePath)

	check(err)

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string

	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}

	readFile.Close()

	return fileLines
}

func readTextFile(filePath string) string {
	lines := readFileLinesFile(filePath)
	text := ""

	for _, line := range lines {
		text += line
	}

	return text
}

func createTmpSubFolderAndReturnPath() string {
	//  Get the system temp directory (cross-platform)
	const subdirectory = "conditional_preprocessor"
	var subDir = path.Join(os.TempDir(), subdirectory)
	err := os.MkdirAll(subDir, os.ModePerm)
	if err != nil {
		log.Fatalf("failed to create temp subdirectory " + subdirectory)
	}

	return subDir
}

func GetConditionalRegions(topItem MdBookTopItem) []string {
	var preprocessor = topItem.Config.Preprocessor
	if preprocessor.Test != nil && preprocessor.Test.ConditionalRegions != nil {
		return preprocessor.Test.ConditionalRegions
	}

	return []string{}
}

func GetVariableDeclarations(topItem MdBookTopItem) []VarNameAndValue {
	var preprocessor = topItem.Config.Preprocessor
	if preprocessor.Test != nil && preprocessor.Test.Variables != nil {
		return preprocessor.Test.Variables
	}

	return []VarNameAndValue{}
}

func main() {
	debugSaveStdInToJsonFile := false
	tempSubDir := createTmpSubFolderAndReturnPath()
	debugInputJsonFileName := path.Join(tempSubDir, "input.json")

	logFileName := path.Join(tempSubDir, "log.txt")
	f, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	if len(os.Args) > 1 {
		if os.Args[1] == "supports" {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	os.Stderr.WriteString(" INFO Kurt was here\n")

	// For debug purposes, read the json from the file instead
	// jsonText := ReadTextFile(debugInputJsonFileName)
	// When not debugging, read the json from stdin
	jsonText := readJsonFromStdIn()
	if debugSaveStdInToJsonFile {
		log.Println("Save stdin to file: " + debugInputJsonFileName)
		CreateTextFile(debugInputJsonFileName, jsonText)
	}

	var book []MdBookTopItem
	errJson := json.Unmarshal([]byte(jsonText), &book)
	check(errJson)
	if len(book) == 2 {
		// The input json is an slice with 2 items.
		// The first item is configuration, parameters to the preprocessor and other stuff
		// The second item is the "content" of the book, the part that should be exported
		conditionalRegions := GetConditionalRegions(book[0])
		variableDeclarations := GetVariableDeclarations(book[0])
		bookSections := book[1]
		processSections(&bookSections, conditionalRegions, variableDeclarations)
		// writeBookSectionsToFile(bookSections, "/tmp/mdbook_out.json")
		writeBookSectionsStdOut(bookSections)
	}
}
