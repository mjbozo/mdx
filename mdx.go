package mdx

import (
	"os"
	"strings"
)

type InvalidFileError struct {
	error
}

func (e *InvalidFileError) Error() string {
	return "Invalid file type. File must have .md or .mdx extension"
}

// TODO: HTML Formatting
// TODO: Code block formatting
// TODO: Ordered list only parses when list starts with single digit number
// TODO: Format HTML correctly when generating file
// TODO: Blockquote are not able to be nested at the moment
// TODO: Ordered list elements only render <p> and only parse continuous lists (doesn't continue over empty lines)
// TODO: Unordered list can't be immediately followed by horizontal rule
// TODO: Button child is hard coded as fragment, doesn't read child components yet
// TODO: Need a way to put data in <head> section. UPDATE: Maybe not
// TODO: Div doesn't return fragment when no closing bracket exists
// TODO: Should be able to have <code> blocks inside <p> tags
// TODO: Unordered lists should ignore all whitespace between list components

func Transform(inputFilename string) (string, error) {
	if !(strings.HasSuffix(inputFilename, ".md") || strings.HasSuffix(inputFilename, ".mdx")) {
		return "", &InvalidFileError{}
	}

	data, readErr := os.ReadFile(inputFilename)
	if readErr != nil {
		return "", readErr
	}

	lexer := NewLexer(string(data))
	parser := NewParser(lexer)
	elements, parseErr := parser.Parse(EOF)

	if parseErr != nil {
		return "", parseErr
	}

	htmlString := TransformMDX(elements)

	return htmlString, nil
}

func Generate(config *GeneratorConfig) (int, error) {
	if !(strings.HasSuffix(config.InputFilename, ".md") || strings.HasSuffix(config.InputFilename, ".mdx")) {
		return 0, &InvalidFileError{}
	}

	data, readErr := os.ReadFile(config.InputFilename)
	if readErr != nil {
		return 0, readErr
	}

	lexer := NewLexer(string(data))
	parser := NewParser(lexer)
	elements, parseErr := parser.Parse(EOF)

	if parseErr != nil {
		return 0, parseErr
	}

	n, err := GenerateHtml(elements, config)

	return n, err
}
