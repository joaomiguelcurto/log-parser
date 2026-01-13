package report

import (
	"fmt"
	"os"
	"text/template"

	"github.com/joaomiguelcurto/log-parser/internal/models"
)

// Template format for the main report about the scan.
const reportTemplate = `
----- Start Report -----
Path:             {{.Path}}
Total Lines:      {{.LineCount}}
Lines per Second: {{printf "%s" .FormatedLinesPerSecond}}
Duration:         {{.AnalyzeDuration}}

----- Process Breakdown -----
{{- range .ProcessStats}}
Process: {{printf "%-30s" .Name}} | Count: {{.Count}}
{{- end}}
----- End Report -----
`

// Prints the report and information about the scan.
func PrintReport(r models.Report) {
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
