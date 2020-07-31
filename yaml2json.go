package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {

	sourceFile := openSourceFile("./test_dir/test_file.yml")
	destFile := openDestFile("./target_test_file.json")

	writeIndentToJSON(destFile)

	pastIndentLevel, indentLevel := -2, 0

	scanner := bufio.NewScanner(sourceFile)
	for scanner.Scan() {

		//scanning source file line
		line := scanner.Text()

		//getting info from file line
		indentLevel = getIndentLevel(line)
		firstMember := getFirstMember(line)
		secondMember := getSecondMember(line)

		//writing the file

		if pastIndentLevel > indentLevel {
			writeUnindentToJSON(pastIndentLevel, indentLevel, destFile)
		}

		if pastIndentLevel >= indentLevel {
			writeNewLineToJSON(destFile)
		}

		writePadToJSON(indentLevel, destFile)

		writeFirstMemberToJSON(firstMember, destFile)
		writeSeparatorToJSON(destFile)

		if secondMember == "" {
			writeIndentToJSON(destFile)
		} else {
			writeSecondMemberToJSON(secondMember, destFile)
		}

		//Preparing for possible unindent
		pastIndentLevel = indentLevel

		//info output
		fmt.Println("Indent Level: ", indentLevel, "\tFirst Member: ", firstMember, "\tSecond Member: ", secondMember)
	}

	writeUnindentToJSON(0, -2, destFile)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	closeFile(sourceFile)
	closeFile(destFile)
}

var indentRegExp = regexp.MustCompile(`(?P<indent> *)`)
var firstMemberRegExp = regexp.MustCompile(`^ *(?P<first_member>\w+):`)
var secondMemberRegExp = regexp.MustCompile(`: (?P<second_member>[-.\w\d]*)$`)
var alpha = "abcdefghijklmnopqrstuvwxyz"

func openSourceFile(filename string) (file *os.File) {
	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}

	return file
}

func openDestFile(filename string) (destFile *os.File) {
	destFile, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	return
}

func closeFile(fp *os.File) {
	fp.Close()
}

func getIndentLevel(line string) (indentLevel int) {
	match := indentRegExp.FindStringSubmatch(line)
	result := make(map[string]string)

	for i, name := range indentRegExp.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}

	indentLevel = len(result["indent"])

	return

}

func getFirstMember(line string) (firstMember string) {
	match := firstMemberRegExp.FindStringSubmatch(line)
	result := make(map[string]string)
	for i, name := range firstMemberRegExp.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}

	firstMember = result["first_member"]

	return
}

func getSecondMember(line string) (secondMember string) {
	match := secondMemberRegExp.FindStringSubmatch(line)
	result := make(map[string]string)

	if len(match) > 0 {
		for i, name := range secondMemberRegExp.SubexpNames() {
			if i != 0 && name != "" {
				result[name] = match[i]
			}
		}
		secondMember = result["second_member"]
	} else {
		secondMember = ""
	}
	return
}

func writeIndentToJSON(file *os.File)  { writeToJSON("{\n", file) }
func writeNewLineToJSON(file *os.File) { writeToJSON(",\n", file) }

func writeUnindentToJSON(pastIndentLevel int, indentLevel int, file *os.File) {
	for i := 1; i <= (pastIndentLevel-indentLevel)/2; i++ {
		writeToJSON("\n"+strings.Repeat(" ", convertPaddingToJSON(pastIndentLevel-2*i))+"}", file)
	}
}

func writeSeparatorToJSON(file *os.File) { writeToJSON(": ", file) }

func writePadToJSON(padding int, file *os.File) {
	writeToJSON(strings.Repeat(" ", convertPaddingToJSON(padding)), file)
}

func writeFirstMemberToJSON(firstMember string, file *os.File) {
	writeToJSON("\""+firstMember+"\"", file)
}

func writeSecondMemberToJSON(secondMember string, file *os.File) {
	if secondMember == "true" || secondMember == "false" {
		writeToJSON(secondMember, file)
		return
	}

	if strings.Contains(alpha, strings.ToLower(string(secondMember[0]))) {
		writeToJSON("\""+secondMember+"\"", file)
		return
	}

	writeToJSON(secondMember, file)
	return
}

func writeToJSON(text string, file *os.File) {
	_, err := file.Write([]byte(text)) //TODO - could be a list
	if err != nil {
		log.Fatal(err)
	}
	return
}

func convertPaddingToJSON(yamlPadding int) (jsonPadding int) { jsonPadding = 4 + yamlPadding*2; return }
