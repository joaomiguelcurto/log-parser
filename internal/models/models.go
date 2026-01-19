package models

import (
	"time"
)

// Structure that will display the final report of the scan.
type Report struct {
	Path                   string
	CleanTerms             []string
	ProcessStats           []ProcessStat
	FormatedLineCount      string
	LineCount              int
	AnalyzeDuration        time.Duration
	LinesPerSecond         float64
	FormatedLinesPerSecond string
}

// Structure to hold the name of each process and the count of appearences.
type ProcessStat struct {
	Name  string
	Count int
}

// Each line of the log file is an entry.
// This structure contains various information about entry.
type LogEntry struct {
	Timestamp   time.Time
	Hostname    string
	ProcessName string
	PID         string
	Message     string
	Valid       bool
}
