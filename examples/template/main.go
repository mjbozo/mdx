package main

import (
	"fmt"
	"os"
	"text/template"

	"github.com/mjbozo/mdx"
)

func main() {
	htmlString, _ := mdx.Transform("template.mdx")

	// insert htmlString into template.html {{ .Content }} section
	htmlTemplate, _ := os.ReadFile("template.html")

	t, _ := template.New("template.mdx").Parse(string(htmlTemplate))

	data := struct {
		Content string
	}{
		Content: htmlString,
	}

	outputFile := "output.html"
	file, _ := os.OpenFile(outputFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	_ = t.Execute(file, data)

	fmt.Printf("Template populated\n")
}
