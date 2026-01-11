package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joaomiguelcurto/log-parser/internal/scanner"
)

func main() {
	searchFlag := flag.String("search", "ERROR", "Search term to search for a specific flag in the log file. (example: ERROR, CRITICAL, INFO, etc...)")
	pathFlag := flag.String("path", "", "File to the Log file to be analyzed.")
	flag.Parse()

	if *pathFlag == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Needs * so it doesnt try to split the address instead of the actual string
	searchTerms := strings.Split(*searchFlag, ",")
	path := *pathFlag

	fmt.Printf("Analyzing log file: %s\n", path)

	lineCount := 0
	flagCount := 0
	errorPercentage := 0.0
	linesPerSecond := 0.0

	analyzeStart := time.Now()

	// Callback
	err := scanner.ReadLog(path, func(line string) {
		// fmt.Printf(line, "\n")
		lineCount++

		for _, element := range searchTerms {
			if strings.Contains(strings.ToUpper(line), strings.ToUpper(element)) {
				flagCount++
			}
		}
	})

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	analyzeEnd := time.Now()

	analyzeDuration := analyzeEnd.Sub(analyzeStart)
	errorPercentage = (float64(flagCount) / float64(lineCount)) * 100.0
	linesPerSecond = float64(lineCount) / analyzeDuration.Seconds()

	PrintReport(*searchFlag, lineCount, flagCount, errorPercentage, analyzeDuration, linesPerSecond)

	fmt.Printf("Done Analyzing\n")
}

func PrintReport(searchTerms string, lineCount int, flagCount int, errorPercentage float64, analyzeDuration time.Duration, linesPerSecond float64) {
	fmt.Printf("----- Start Info ----- \n")
	fmt.Printf("Search terms: %s\n", searchTerms)
	fmt.Printf("----- End Info ----- \n")

	fmt.Printf("----- Start Report ----- \n")
	fmt.Printf("Lines: %d\n", lineCount)
	fmt.Printf("Lines with Flags: %d\n", flagCount)
	fmt.Printf("Flags Percentage: %.1f%%\n", errorPercentage)
	fmt.Printf("Analyze Duration: %s\n", analyzeDuration)
	fmt.Printf("Lines per Second: %.0f\n", linesPerSecond)
	fmt.Printf("----- End Report ----- \n")
}
