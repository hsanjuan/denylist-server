package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Usage: program <source_file> <source_URL>")
		os.Exit(1)
	}

	localFilePath := os.Args[1]
	remoteURL := os.Args[2]

	localLines, err := readLines(localFilePath)
	if err != nil {
		log.Fatalf("Error reading local file: %v", err)
	}

	remoteLines, err := readRemoteFile(remoteURL)
	if err != nil {
		log.Fatalf("Error reading remote file: %v", err)
	}

	newLines := findMissingLinesInA(localLines, remoteLines, false)
	appendLinesWithPrefixToFile(localFilePath, newLines, "")

	// Find lines in the local file that are not in the remote file and append them with `!`
	missingLines := findMissingLinesInA(remoteLines, localLines, true)
	appendLinesWithPrefixToFile(localFilePath, missingLines, "!")

	// Print new lines appended to the file to stdout
	for _, line := range newLines {
		fmt.Println(line)
	}
	for _, line := range missingLines {
		fmt.Println("!" + line)
	}
}

func readLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	reader := bufio.NewReader(file)
	skipHeaders := true // Start by skipping headers

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		line = strings.TrimSpace(line)

		if line == "---" {
			skipHeaders = false // Stop skipping headers when `---` is found
			continue
		}

		if !skipHeaders {
			lines = append(lines, line)
		}
	}

	return lines, nil
}

func readRemoteFile(url string) ([]string, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var lines []string
	reader := bufio.NewReader(response.Body)
	skipHeaders := true // Start by skipping headers

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		line = strings.TrimSpace(line)

		if line == "---" {
			skipHeaders = false // Stop skipping headers when `---` is found
			continue
		}

		if !skipHeaders {
			lines = append(lines, line)
		}
	}

	return lines, nil
}

func findMissingLinesInA(a, b []string, ignoreNeg bool) []string {
	var newLines []string
	aSet := make(map[string]struct{})
	bSet := make(map[string]struct{})

	// file a should have all the lines in b
	for _, line := range a {
		aSet[line] = struct{}{}
	}

	// all the lines in b. if a line is negated
	// then that line is not in b anymore.
	for _, line := range b {
		if ignoreNeg && line[0] == '!' {
			delete(bSet, line[1:])
		} else {
			bSet[line] = struct{}{}
		}
	}

	// for every line in b, spit those that are not in a.
	for line := range bSet {
		if _, ok := aSet[line]; !ok {
			newLines = append(newLines, line)
		}
	}

	return newLines
}

func appendLinesWithPrefixToFile(filePath string, lines []string, prefix string) {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("Error opening file for appending: %v", err)
	}
	defer file.Close()

	for _, line := range lines {
		_, err := file.WriteString(prefix + line + "\n")
		if err != nil {
			log.Fatalf("Error appending lines to local file: %v", err)
		}
	}
}
