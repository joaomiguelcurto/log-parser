package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/joaomiguelcurto/log-parser/internal/models"
	"github.com/joaomiguelcurto/log-parser/internal/parser"
	"github.com/joaomiguelcurto/log-parser/internal/report"
	"github.com/joaomiguelcurto/log-parser/internal/scanner"
	"github.com/joaomiguelcurto/log-parser/internal/utils"
)

func main() {
	searchFlag := flag.String("search", "", "Search term to search for a specific flag in the log file. (example: ERROR, CRITICAL, INFO, etc...)")
	pathFlag := flag.String("path", "", "File to the Log file to be analyzed.")
	flag.Parse()

	if *pathFlag == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Needs * so it doesnt try to split the address instead of the actual string
	searchTerms := strings.Split(*searchFlag, ",")
	path := *pathFlag

	var cleanTerms []string

	for _, element := range searchTerms {
		trimmed := strings.TrimSpace(element)
		if trimmed != "" {
			cleanTerms = append(cleanTerms, strings.ToUpper(trimmed))
		}
	}

	hasTerms := false

	if len(cleanTerms) != 0 {
		hasTerms = true
	}

	fmt.Printf("Analyzing log file: %s\n", path)

	lineCount := 0
	linesPerSecond := 0.0

	// make instead of var so it initializes instead of just declaration
	processMap := make(map[string]int)

	analyzeStart := time.Now()

	parsed := models.LogEntry{}
	p := parser.LinuxParser{}

	// Callback
	err := scanner.ReadLog(path, func(line string) {
		lineCount++

		parsed = p.Parse(line)

		if hasTerms == true {
			for _, element := range cleanTerms {
				if strings.Contains(strings.ToUpper(parsed.Message), element) {
					processMap[element+" "+parsed.ProcessName]++
					break
				}
			}
		} else {
			processMap[parsed.ProcessName]++
		}
	})

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	analyzeEnd := time.Now()

	analyzeDuration := analyzeEnd.Sub(analyzeStart)
	linesPerSecond = float64(lineCount) / analyzeDuration.Seconds()

	// Sort the processes by count
	processStats := []models.ProcessStat{}

	for name, count := range processMap {
		processStats = append(processStats, models.ProcessStat{
			Name:  name,
			Count: count,
		})
	}

	sort.Slice(processStats, func(i, j int) bool {
		return processStats[i].Count > processStats[j].Count
	})

	results := models.Report{
		ProcessStats:           processStats,
		Path:                   path,
		CleanTerms:             cleanTerms,
		LineCount:              lineCount,
		AnalyzeDuration:        analyzeDuration,
		LinesPerSecond:         linesPerSecond,
		FormatedLinesPerSecond: utils.FormatNumberSimple(linesPerSecond),
	}

	report.PrintReport(results)

	fmt.Printf("Done Analyzing\n")
}
