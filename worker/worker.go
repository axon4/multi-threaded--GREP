package worker

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Result struct {
	Path       string
	Line       string
	LineNumber int
}

type Results struct {
	Results []Result
}

func NewResult(path string, line string, lineNumber int) Result {
	return Result{path, line, lineNumber}
}

func FindInFile(path string, query string) *Results {
	file, err := os.Open(path)

	if err != nil {
		fmt.Println("error opening file:", err)

		return nil
	}

	results := Results{make([]Result, 0)}
	scanner := bufio.NewScanner(file)
	lineNumber := 1

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), query) {
			result := NewResult(path, scanner.Text(), lineNumber)
			results.Results = append(results.Results, result)
		}

		lineNumber++
	}

	if len(results.Results) == 0 {
		return nil
	} else {
		return &results
	}
}