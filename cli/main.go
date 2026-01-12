package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joaomiguelcurto/log-parser/internal/parser"
	"github.com/joaomiguelcurto/log-parser/internal/scanner"
)

// Structure that will display the final report of the scan.
type report struct {
	path            string
	cleanTerms      []string
	flagMap         map[string]int
	lineCount       int
	flagCount       int
	errorPercentage float64
	analyzeDuration time.Duration
	linesPerSecond  float64
}

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

	// TODO: Find use for the terms (OBSOLETE RIGHT NOW)
	var cleanTerms []string

	for _, element := range searchTerms {
		cleanTerms = append(cleanTerms, strings.TrimSpace(strings.ToUpper(element)))
	}

	fmt.Printf("Analyzing log file: %s\n", path)

	lineCount := 0
	linesPerSecond := 0.0

	// make instead of var so it initializes instead of just declaration
	flagMap := make(map[string]int)

	/*for _, element := range cleanTerms {
		flagMap[element] = 0
	}*/

	analyzeStart := time.Now()

	parsed := parser.LogEntry{}

	p := parser.LinuxParser{}

	cleanProcess := ""
	processNameIndex := 0

	// Callback
	err := scanner.ReadLog(path, func(line string) {
		// fmt.Printf(line, "\n")
		lineCount++

		// upperLine := strings.ToUpper(line)

		parsed = p.Parse(line)

		processNameIndex = strings.Index(parsed.Process[:], "[")

		if processNameIndex != -1 {
			// [ found so it slices it (example: sshd[1234] -> sshd)
			cleanProcess = parsed.Process[:processNameIndex]
			flagMap[cleanProcess]++
		} else {
			// No [ found so the process is already clean
			flagMap[parsed.Process]++
		}

		/*
			for _, element := range cleanTerms {
				// if it finds the flag then it skips the rest of the line (flags usually found at the start)
				if strings.Contains(upperLine, element) {
					// flagMap[element]++
					break
				}
			}*/
	})

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	analyzeEnd := time.Now()

	analyzeDuration := analyzeEnd.Sub(analyzeStart)
	linesPerSecond = float64(lineCount) / analyzeDuration.Seconds()

	fmt.Println(parsed)
	fmt.Println(parsed.Timestamp)
	fmt.Println(parsed.Hostname)
	fmt.Println(parsed.Process)
	fmt.Println(parsed.Message)

	results := report{
		flagMap:         flagMap,
		path:            path,
		cleanTerms:      cleanTerms,
		lineCount:       lineCount,
		analyzeDuration: analyzeDuration,
		linesPerSecond:  linesPerSecond,
	}

	PrintReport(results)

	fmt.Printf("Done Analyzing\n")
}

// Prints the report and information about the scan.
func PrintReport(r report) {
	fmt.Printf("----- Start Info ----- \n")
	fmt.Printf("Search terms: %s\n", r.cleanTerms)
	fmt.Printf("----- End Info ----- \n")

	fmt.Printf("----- Start Terms Info ----- \n")
	for name, count := range r.flagMap {
		fmt.Printf("Search term and count - %s: %d\n", name, count)
	}
	fmt.Printf("----- End Terms Info ----- \n")

	fmt.Printf("----- Start Report ----- \n")
	fmt.Printf("Lines: %d\n", r.lineCount)
	fmt.Printf("Analyze Duration: %s\n", r.analyzeDuration)
	fmt.Printf("Lines per Second: %.0f\n", r.linesPerSecond)
	fmt.Printf("----- End Report ----- \n")
}
