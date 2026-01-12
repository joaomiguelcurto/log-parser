package parser

import (
	"strings"
	"time"
)

// Each line of the log file is an entry.
// This structure contains various information about entry.
type LogEntry struct {
	Timestamp time.Time
	Hostname  string
	Process   string
	Message   string
}

// Interface for the any parser that will be created.
type LogParser interface {
	Parse(line string) LogEntry
	GetName() string
}

type LinuxParser struct{}

func (p LinuxParser) GetName() string {
	return "Linux/Syslog"
}

func (p LinuxParser) Parse(line string) LogEntry {
	r := LogEntry{}

	const timeLayout = "Jan _2 15:04:05"

	r.Timestamp, _ = time.Parse(timeLayout, line[0:15])

	// Finds the first Colon after 15th character of the string
	colonIndex := strings.Index(line[15:], ":")

	// Finds the first space after 16th character of the string
	hostnameIndex := strings.Index(line[16:], " ")

	if colonIndex != -1 {
		// Offset the colonIndex by 15 so it gives the actual message.
		actualColonPos := colonIndex + 15

		// Everything between the first space after the timestamp until the colon.
		r.Process = strings.TrimSpace(line[hostnameIndex+16 : actualColonPos])

		// Everything between the 15th position (end of the timestamp) and the position of the start of message is the Hostname.
		r.Hostname = strings.TrimSpace(line[15 : hostnameIndex+16])

		// Everything after the colon is the message.
		r.Message = strings.TrimSpace(line[actualColonPos+1:])
	} else {
		r.Hostname = "Unknown"
		r.Process = "Unknown"
		r.Message = line[15:]
	}

	return r
}
