package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/joaomiguelcurto/log-parser/internal/models"
	"github.com/joaomiguelcurto/log-parser/internal/parser"
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

	// TODO: Find use for the terms (OBSOLETE RIGHT NOW)
	var cleanTerms []string

	for _, element := range searchTerms {
		cleanTerms = append(cleanTerms, strings.TrimSpace(strings.ToUpper(element)))
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
		// fmt.Printf(line, "\n")
		lineCount++

		// upperLine := strings.ToUpper(line)

		parsed = p.Parse(line)

		processMap[parsed.ProcessName]++

		/*
			for _, element := range cleanTerms {
				// if it finds the flag then it skips the rest of the line (flags usually found at the start)
				if strings.Contains(upperLine, element) {
					// processMap[element]++
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
		ProcessStats:    processStats,
		Path:            path,
		CleanTerms:      cleanTerms,
		LineCount:       lineCount,
		AnalyzeDuration: analyzeDuration,
		LinesPerSecond:  linesPerSecond,
	}

	PrintReport(results)

	fmt.Printf("Done Analyzing\n")
}

// Prints the report and information about the scan.
func PrintReport(r models.Report) {
	const reportTemplate = `
----- Start Report -----
Path:             {{.Path}}
Total Lines:      {{.LineCount}}
Lines per Second: {{printf "%.0f" .LinesPerSecond}}
Duration:         {{.AnalyzeDuration}}

----- Process Breakdown -----
{{- range .ProcessStats}}
Process: {{printf "%-30s" .Name}} | Count: {{.Count}}
{{- end}}
----- End Report -----
`

	// Create the template
	tmpl, err := template.New("report").Parse(reportTemplate)
	if err != nil {
		panic(err)
	}

	// Execute the template
	err = tmpl.Execute(os.Stdout, r)
	if err != nil {
		fmt.Printf("Error printing report: %v\n", err)
	}
}
