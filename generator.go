package mdx

import (
	"fmt"
	"log"
	"os"
)

type GeneratorConfig struct {
	Title          string
	InputFilename  string
	OutputFilename string
	Links          []map[string]string
}

func transformMDX(elements []component) string {
	// Note: Testing changing from body to div and putting more in the template.html
	// body := &Body{Children: elements}
	body := &div{Children: elements}
	htmlString := body.html()
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

	body := &body{Children: elements}
	n, writeErr := file.WriteString(body.Html())
	if writeErr != nil {
		log.Printf(writeErr.Error())
		return n, writeErr
	}

	file.WriteString(`
</html>
`)

	return n, nil
}
