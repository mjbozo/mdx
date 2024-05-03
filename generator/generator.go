package generator

import (
	"fmt"
	"log"
	"os"

	"github.com/matt-bourke/mdx/ast"
)

type GeneratorConfig struct {
	Title          string
	InputFilename  string
	OutputFilename string
	Links          []map[string]string
}

func TransformMDX(elements []ast.Component) string {
	body := &ast.Body{Children: elements}
	htmlString := body.Html()
	return htmlString
}

func GenerateHtml(elements []ast.Component, config *GeneratorConfig) (int, error) {
	file, fileErr := os.OpenFile(config.OutputFilename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if fileErr != nil {
		log.Fatal(fileErr.Error())
	}

	defer file.Close()

	file.WriteString(`<html>
    <head>
		<meta charset="UTF-8" />
		<meta name="viewport" content="width=device-width,initial-scale=1" />
		<meta name="description" content="" />`)

	if len(config.Title) > 0 {
		file.WriteString(fmt.Sprintf(`
		<title>%s</title>`, config.Title))
	}

	for _, link := range config.Links {
		linkString := "<link "
		for name, value := range link {
			if len(value) > 0 {
				linkString += fmt.Sprintf("%s=\"%s\" ", name, value)
			}
		}
		linkString += ">"
		file.WriteString(linkString)
	}

	file.WriteString(`
    </head>
`)

	body := &ast.Body{Children: elements}
	n, writeErr := file.WriteString(body.Html())
	if writeErr != nil {
		log.Fatal(writeErr.Error())
	}

	file.WriteString(`
</html>
`)

	return n, nil
}
