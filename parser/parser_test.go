package parser

import (
	"github.com/matt-bourke/mdx/ast"
	"github.com/matt-bourke/mdx/lexer"
	"github.com/matt-bourke/mdx/token"
	"testing"
)

func TestHeading(t *testing.T) {

}

func TestParagraph(t *testing.T) {

}

func TestProperties(t *testing.T) {

}

func TestCode(t *testing.T) {
	input := "`print('Hello, world!')`"
	lex := lexer.New(input)
	parser := New(lex)
	elements, parseErr := parser.Parse(token.EOF)
	if parseErr != nil {
		t.Errorf("Code parse failed: %s", parseErr.Error())
	}

	if len(elements) != 1 {
		t.Errorf("Code parse failed: Expected 1 element, got=%d", len(elements))
	}

	element := elements[0]
	if code, ok := element.(*ast.Code); ok {
		if code.Text != "print('Hello, world!')" {
			t.Errorf("Code parse failed: Content incorrect, got=%s", code.Text)
		}
	} else {
		t.Errorf("Code parse failed: Expected Code type, got=%T", element)
	}

	input = "`print('Hello, world!')"
	lex = lexer.New(input)
	parser = New(lex)
	elements, parseErr = parser.Parse(token.EOF)
	if parseErr != nil {
		t.Errorf("Code parse failed: %s", parseErr.Error())
	}

	if len(elements) != 1 {
		t.Errorf("Code parse failed: Expected 1 element, got=%d", len(elements))
	}

	element = elements[0]
	if frag, ok := element.(*ast.Fragment); ok {
		if frag.String != input {
			t.Errorf("Invalid code parse failed: Got fragment with string %s", frag.String)
		}
	} else {
		t.Errorf("Invalid code parse failed: Expected Fragment, got=%T", element)
	}
}

func TestDiv(t *testing.T) {
	input := `[
    Hello
]`
	lex := lexer.New(input)
	parser := New(lex)
	elements, parseErr := parser.Parse(token.EOF)
	if parseErr != nil {
		t.Errorf("Div parse failed: %s", parseErr.Error())
	}

	if len(elements) != 1 {
		t.Errorf("Div parse failed: Expected 1 element, got=%d", len(elements))
	}

	element := elements[0]
	if div, ok := element.(*ast.Div); ok {
		if len(div.Children) != 1 {
			t.Errorf("Div parse failed: Exected one child, got=%d", len(div.Children))
		}
	} else {
		t.Errorf("Div parse failed: Expected Div type, got=%T", element)
	}

	input = `[
    Hello
]

`
	lex = lexer.New(input)
	parser = New(lex)
	elements, parseErr = parser.Parse(token.EOF)
	if parseErr != nil {
		t.Errorf("Div parse failed: %s", parseErr.Error())
	}

	if len(elements) != 1 {
		t.Errorf("Div parse failed: Expected 1 element, got=%d", len(elements))
	}

	element = elements[0]
	if div, ok := element.(*ast.Div); ok {
		if len(div.Children) != 1 {
			t.Errorf("Div parse failed: Exected one child, got=%d", len(div.Children))
		}
	} else {
		t.Errorf("Div parse failed: Expected Div type, got=%T", element)
	}
}
