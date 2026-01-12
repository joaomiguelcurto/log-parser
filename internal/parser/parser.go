package parser

import (
	"strings"
	"time"

	"github.com/joaomiguelcurto/log-parser/internal/models"
)

// Interface for the any parser that will be created.
type LogParser interface {
	Parse(line string) models.LogEntry
	GetName() string
}

type LinuxParser struct{}

func (p LinuxParser) GetName() string {
	return "Linux/Syslog"
}

func (p LinuxParser) Parse(line string) models.LogEntry {
	r := models.LogEntry{}

	const timeLayout = "Jan _2 15:04:05"

	r.Timestamp, _ = time.Parse(timeLayout, line[0:15])

	// Finds the first Colon after 15th character of the string
	colonIndex := strings.Index(line[15:], ":")

	// Offset the colonIndex by 15 so it gives the actual message.
	actualColonPos := colonIndex + 15

	// Finds the first space after 16th character of the string
	hostnameIndex := strings.Index(line[16:], " ")

	// Everything between the first space after the timestamp until the colon.
	process := strings.TrimSpace(line[hostnameIndex+16 : actualColonPos])

	PIDIndex := strings.Index(process[:], "[")

	if colonIndex != -1 {
		if PIDIndex != -1 {
			// [ found so it slices it (example: sshd[1234] -> sshd)
			r.PID = process[PIDIndex:]

			r.ProcessName = process[:PIDIndex]
		} else {
			// No [ found so the process is already clean
			r.PID = "EMPTY"
			r.ProcessName = process
		}

		// Everything between the 15th position (end of the timestamp) and the position of the start of message is the Hostname.
		r.Hostname = strings.TrimSpace(line[15 : hostnameIndex+16])

		// Everything after the colon is the message.
		r.Message = strings.TrimSpace(line[actualColonPos+1:])

	} else {
		r.Hostname = "Unknown"
		r.PID = "Unknown"
		r.ProcessName = "Unknown"
		r.Message = line[15:]
	}

	return r
}
