package mdx

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type GeneratorConfig struct {
	Title          string
	InputFilename  string
	OutputFilename string
	Links          []map[string]string
}

func transformMDX(elements []component) string {
	body := &div{Children: elements}
	htmlString := body.Html()
	return htmlString
}

func generateHtml(elements []component, config *GeneratorConfig) (int, error) {
	file, fileErr := os.OpenFile(config.OutputFilename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if fileErr != nil {
		log.Println(fileErr.Error())
		return 0, fileErr
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
		linkString := "\n        <link "
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

	body := &body{Children: elements}
	n, writeErr := file.WriteString(strings.ReplaceAll(body.Format(1), "\n\n", "\n"))
	if writeErr != nil {
		log.Printf(writeErr.Error())
		return n, writeErr
	}

	file.WriteString(`
</html>
`)

	return n, nil
}
