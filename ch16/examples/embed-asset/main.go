// embed: ship a file inside the binary, so the program is still one file to copy
package main

import (
	_ "embed"
	"fmt"
	"os"
	"text/template"
)

//go:embed templates/report.tmpl
var reportTemplate string

func main() {
	tmpl, err := template.New("report").Parse(reportTemplate)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	data := struct {
		Name  string
		Items []string
	}{"Monday", []string{"coffee", "stand-up", "code review"}}
	if err := tmpl.Execute(os.Stdout, data); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
