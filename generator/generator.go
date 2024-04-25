package generator

import (
	"fmt"
	"log"
	"mdx/ast"
	"os"
)

func GenerateDocument(filename string, elements []ast.Component) {
	file, fileErr := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if fileErr != nil {
		log.Fatal(fileErr.Error())
	}

	defer file.Close()

	file.WriteString(`<html>
    <head>
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

	fmt.Printf("%d autogenerated bytes written to %s\n", n, filename)
}
