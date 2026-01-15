package report

import (
	"fmt"
	"os"
	"text/template"

	"github.com/joaomiguelcurto/log-parser/internal/models"
)

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Magenta = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

// Template format for the main report about the scan.
const reportTemplate = `
{{Cyan}}----- Start Report -----{{Reset}}
{{Red}}Path:{{Reset}}             {{.Path}}
{{Blue}}Total Lines:{{Reset}}      {{.LineCount}}
{{Green}}Lines per Second:{{Reset}} {{.FormatedLinesPerSecond}}
Duration:          {{.AnalyzeDuration}}

{{Yellow}}----- Process Breakdown -----{{Reset}}
{{- range .ProcessStats}}
Process: {{printf "%-30s" .Name}} | Count: {{.Count}}
{{- end}}
{{Cyan}}----- End Report -----{{Reset}}
`

// Prints the report and information about the scan.
func PrintReport(r models.Report) {
	funcs := template.FuncMap{
		"Red":     func() string { return Red },
		"Green":   func() string { return Green },
		"Yellow":  func() string { return Yellow },
		"Blue":    func() string { return Blue },
		"Cyan":    func() string { return Cyan },
		"Reset":   func() string { return Reset },
		"Magenta": func() string { return Magenta },
	}

	tmpl, err := template.New("report").Funcs(funcs).Parse(reportTemplate)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(os.Stdout, r)
	if err != nil {
		fmt.Printf("Error printing report: %v\n", err)
	}
}
