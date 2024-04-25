package main

import (
	"fmt"
	"mdx/generator"
	"mdx/lexer"
	"mdx/parser"
	"mdx/token"
	"os"
	"strings"
)

type InvalidFileError struct {
	error
}

func (e *InvalidFileError) Error() string {
	return "Invalid file type. File must have .md or .mdx extension"
}

// TODO: Ordered list only parses when list starts with single digit number
// TODO: Format HTML correctly when generating file
// TODO: Blockquote are not able to be nested at the moment
// TODO: Ordered list elements only render <p> and only parse continuous lists (doesn't continue over empty lines)
// TODO: Unordered list can't be immediately followed by horizontal rule
// TODO: Button child is hard coded as fragment, doesn't read child components yet
// TODO: Need a way to put data in <head> section
// TODO: Div doesn't return fragment when no closing bracket exists

func main() {

	err := Generate("sample.mdx")
	if err != nil {
		fmt.Println("Error occurred: %s\n", err.Error())
	}

	// data, err := os.ReadFile("./sample.mdx")
	// if err != nil {
	// 	fmt.Printf("File error: %s\n", err.Error())
	// }

	// l := lexer.New(string(data))
	// p := parser.New(l)
	// elements := p.Parse()

	// tok := l.NextToken()
	// for tok.Type != token.EOF {
	// 	fmt.Println(tok.String())
	// 	tok = l.NextToken()
	// }
	// header := &ast.Header{Text: "My First Auto Generated File!"}
	// paragraph := &ast.Paragraph{Text: "How cool is this?!"}
	// button := &ast.Button{OnClick: "handleClick", Child: &ast.Fragment{String: "Click Me"}}
	// code := &ast.CodeBlock{Content: content}

	// root := &ast.Div{Children: []ast.Component{header, paragraph, button, code}}
	// generator.GenerateDocument("MyFirstAutoGenFile.html", elements)
}

func Generate(filename string) error {
	if !(strings.HasSuffix(filename, ".md") || strings.HasSuffix(filename, ".mdx")) {
		return &InvalidFileError{}
	}

	data, readErr := os.ReadFile(filename)
	if readErr != nil {
		return readErr
	}

	lexer := lexer.New(string(data))
	parser := parser.New(lexer)
	elements, parseErr := parser.Parse(token.EOF)

	if parseErr != nil {
		return parseErr
	}

	outputFilename := strings.TrimSuffix(strings.TrimSuffix(filename, ".mdx"), ".md") + ".html"
	generator.GenerateDocument(outputFilename, elements)
	return nil
}
