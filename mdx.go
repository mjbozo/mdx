package mdx

import (
	"os"
	"strings"

	"github.com/matt-bourke/mdx/generator"
	"github.com/matt-bourke/mdx/lexer"
	"github.com/matt-bourke/mdx/parser"
	"github.com/matt-bourke/mdx/token"
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
// TODO: Need a way to put data in <head> section. UPDATE: Maybe not
// TODO: Div doesn't return fragment when no closing bracket exists
// TODO: Newlines inserting line break elements after blocks finish

func Transform(inputFilename string) (string, error) {
	if !(strings.HasSuffix(inputFilename, ".md") || strings.HasSuffix(inputFilename, ".mdx")) {
		return "", &InvalidFileError{}
	}

	data, readErr := os.ReadFile(inputFilename)
	if readErr != nil {
		return "", readErr
	}

	lexer := lexer.New(string(data))
	parser := parser.New(lexer)
	elements, parseErr := parser.Parse(token.EOF)

	if parseErr != nil {
		return "", parseErr
	}

	htmlString := generator.TransformMDX(elements)

	return htmlString, nil
}

func Generate(config *generator.GeneratorConfig) error {
	if !(strings.HasSuffix(config.InputFilename, ".md") || strings.HasSuffix(config.InputFilename, ".mdx")) {
		return &InvalidFileError{}
	}

	data, readErr := os.ReadFile(config.InputFilename)
	if readErr != nil {
		return readErr
	}

	lexer := lexer.New(string(data))
	parser := parser.New(lexer)
	elements, parseErr := parser.Parse(token.EOF)

	if parseErr != nil {
		return parseErr
	}

	err := generator.GenerateHtml(elements, config)

	return err
}
