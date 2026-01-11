package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/joaomiguelcurto/log-parser/internal/scanner"
)

func main() {
	// os.Args[0] is the program name, os.Args[1] is the first argument
	if len(os.Args) < 2 {
		fmt.Printf("Usage: lp <file-path> \n")

		os.Exit(1)
	}

	path := os.Args[1]

	fmt.Printf("Analizing log file: %s\n", path)

	lineCount := 0
	errorCount := 0

	// Callback pattern
	err := scanner.ReadLog(path, func(line string) {
		// fmt.Printf(line, "\n")
		lineCount++

		if strings.Contains(line, "ERROR") {
			errorCount++
		}
	})

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Report: \n")
	fmt.Printf("Lines: %d\n", lineCount)
	fmt.Printf("Lines with Errors: %d\n", errorCount)

	fmt.Printf("Done\n")
}
