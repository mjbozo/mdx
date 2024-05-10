package mdx

import (
	"testing"
)

func TestHeading(t *testing.T) {

}

func TestParseParagraph(t *testing.T) {

}

func TestParseProperties(t *testing.T) {

}

func TestParseCode(t *testing.T) {
	input := "`print('Hello, world!')`"
	lex := newLexer(input)
	parser := newParser(lex)
	elements, parseErr := parser.parse(eof)
	if parseErr != nil {
		t.Errorf("Code parse failed: %s", parseErr.Error())
	}

	if len(elements) != 1 {
		t.Errorf("Code parse failed: Expected 1 element, got=%d", len(elements))
	}

	element := elements[0]
	if code, ok := element.(*code); ok {
		if code.Text != "print('Hello, world!')" {
			t.Errorf("Code parse failed: Content incorrect, got=%s", code.Text)
		}
	} else {
		t.Errorf("Code parse failed: Expected Code type, got=%T", element)
	}

	input = "`print('Hello, world!')"
	lex = newLexer(input)
	parser = newParser(lex)
	elements, parseErr = parser.parse(eof)
	if parseErr != nil {
		t.Errorf("Code parse failed: %s", parseErr.Error())
	}

	if len(elements) != 1 {
		t.Errorf("Code parse failed: Expected 1 element, got=%d", len(elements))
	}

	element = elements[0]
	if frag, ok := element.(*fragment); ok {
		if frag.String != input {
			t.Errorf("Invalid code parse failed: Got fragment with string %s", frag.String)
		}
	} else {
		t.Errorf("Invalid code parse failed: Expected Fragment, got=%T", element)
	}
}

func TestParseDiv(t *testing.T) {
	input := `[
    Hello
]`
	lex := newLexer(input)
	parser := newParser(lex)
	elements, parseErr := parser.parse(eof)
	if parseErr != nil {
		t.Errorf("Div parse failed: %s", parseErr.Error())
	}

	if len(elements) != 1 {
		t.Errorf("Div parse failed: Expected 1 element, got=%d", len(elements))
	}

	element := elements[0]
	if div, ok := element.(*div); ok {
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
	lex = newLexer(input)
	parser = newParser(lex)
	elements, parseErr = parser.parse(eof)
	if parseErr != nil {
		t.Errorf("Div parse failed: %s", parseErr.Error())
	}

	if len(elements) != 1 {
		t.Errorf("Div parse failed: Expected 1 element, got=%d", len(elements))
	}

	element = elements[0]
	if div, ok := element.(*div); ok {
		if len(div.Children) != 1 {
			t.Errorf("Div parse failed: Exected one child, got=%d", len(div.Children))
		}
	} else {
		t.Errorf("Div parse failed: Expected Div type, got=%T", element)
	}
}
