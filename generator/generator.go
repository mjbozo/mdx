package generator

import (
	"fmt"
	"github.com/matt-bourke/mdx/ast"
	"log"
	"os"
)

type GeneratorConfig struct {
	Title          string
	Filename       string
	OutputFilename string
	CssFilename    string
	Favicon        string
	FontLink       string
}

func GenerateDocument(elements []ast.Component, config *GeneratorConfig) {
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

	if len(config.CssFilename) > 0 {
		file.WriteString(fmt.Sprintf(`
		<link rel="stylesheet" type="text/css" href="%s">`, config.CssFilename))
	}

	if len(config.Favicon) > 0 {
		file.WriteString(fmt.Sprintf(`
		<link rel="icon" href="static/favicon.ico">`, config.Favicon))
	}

	if len(config.FontLink) > 0 {
		file.WriteString(fmt.Sprintf(`
		<link rel="stylesheet" href="%s">`, config.FontLink))
	}

	file.WriteString(`
        <script>
            const handleClick = () => {
                console.log("Hello!");
            }
        </script>
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

	fmt.Printf("%d autogenerated bytes written to %s\n", n, config.OutputFilename)
}
