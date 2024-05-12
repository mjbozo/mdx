package mdx

import (
	"testing"
)

func TestParseHeader1(t *testing.T) {
	input := "# Heading"
	lex := newLexer(input)
	parser := newParser(lex)
	elements, parseErr := parser.parse(eof)

	if parseErr != nil {
		t.Errorf("Header parse failed: %s", parseErr.Error())
	}

	if len(elements) != 1 {
		t.Errorf("Header parse failed: Expected 1 element, got=%d", len(elements))
	}

	element := elements[0]
	if head, ok := element.(*header); ok {
		if head.Level != 1 {
			t.Errorf("Header parse failed: Expected Header level 1, got=%d", head.Level)
		}

		if head.Text() != "Heading" {
			t.Errorf("Header parse failed: Expected text 'Heading', got=%s", head.Text())
		}
	} else {
		t.Errorf("Header parse failed: Expected Heading, got=%T", element)
	}
}

func TestParseHeader2(t *testing.T) {
	input := "## Heading Too"
	lex := newLexer(input)
	parser := newParser(lex)
	elements, parseErr := parser.parse(eof)

	if parseErr != nil {
		t.Errorf("Header2 parse failed: %s", parseErr.Error())
	}

	if len(elements) != 1 {
		t.Errorf("Header2 parse failed: Expected 1 element, got=%d", len(elements))
	}

	element := elements[0]
	if head, ok := element.(*header); ok {
		if head.Level != 2 {
			t.Errorf("Header2 parse failed: Expected Header level 1, got=%d", head.Level)
		}

		if head.Text() != "Heading Too" {
			t.Errorf("Header2 parse failed: Expected text 'Heading', got=%s", head.Text())
		}
	} else {
		t.Errorf("Header2 parse failed: Expected Heading, got=%T", element)
	}
}

func TestParseHeaderBetweenElements(t *testing.T) {
	input := `Hello
# World
MDX`
	lex := newLexer(input)
	parser := newParser(lex)
	elements, parseErr := parser.parse(eof)

	if parseErr != nil {
		t.Errorf("HeaderBetweenElements parse failed: %s", parseErr.Error())
	}

	if len(elements) != 3 {
		t.Errorf("HeaderBetweenElements parse failed: Expected 3 elements, got=%d", len(elements))
		t.FailNow()
	}

	element := elements[1]

	if head, ok := element.(*header); ok {
		if head.Level != 1 {
			t.Errorf("HeaderBetweenElements parse failed: Expected Header level 1, got=%d", head.Level)
		}

		if head.Text() != "World" {
			t.Errorf("HeaderBetweenElements parse failed: Expected text 'World', got=%s", head.Text())
		}
	} else {
		t.Errorf("HeaderBetweenElements parse failed: Expected Header, got=%T", element)
	}
}

func TestParseHeaderWithNestedElements(t *testing.T) {
	input := `# Hello **world**`
	lex := newLexer(input)
	parser := newParser(lex)
	elements, parseErr := parser.parse(eof)

	if parseErr != nil {
		t.Errorf("HeaderWithNestedElements parse failed: %s", parseErr)
	}

	if len(elements) != 1 {
		t.Errorf("HeaderWithNestedElements parse failed: Expected 1 element, got=%d", len(elements))
		t.FailNow()
	}

	element := elements[0]

	if head, ok := element.(*header); ok {
		if head.Level != 1 {
			t.Errorf("HeaderWithNestedElements parse failed: Expected Header level 1, got=%d", head.Level)
		}

		if len(head.Content) != 2 {
			t.Errorf("HeaderWithNestedElements parse failed: Expected 2 children, got=%d", len(head.Content))
			t.FailNow()
		}

		if frag, ok := head.Content[0].(*fragment); ok {
			if frag.String != "Hello " {
				t.Errorf("HeaderWithNestedElements parse failed: Expected 'Hello ', got=%s", frag.String)
			}
		} else {
			t.Errorf("HeaderWithNestedElements parse failed: Expected Fragment child, got=%T", element)
		}

		if strong, ok := head.Content[1].(*bold); ok {
			if strong.Text != "world" {
				t.Errorf("HeaderWithNestedElements parse failed: Expected 'world', got=%s", strong.Text)
			}
		} else {
			t.Errorf("HeaderWithNestedElements parse failed: Expected Strong child, got=%T", element)
		}
	} else {
		t.Errorf("HeaderWithNestedElements parse failed: Expected Header, got=%T", element)
	}
}

func TestParseParagraph(t *testing.T) {
	input := "Hello, world"
	lex := newLexer(input)
	parser := newParser(lex)
	elements, parseErr := parser.parse(eof)

	if parseErr != nil {
		t.Errorf("Paragraph parse failed: %s", parseErr.Error())
	}

	if len(elements) != 1 {
		t.Errorf("Paragraph parse failed: Expected 1 elements, got=%d", len(elements))
	}

	element := elements[0]

	if p, ok := element.(*paragraph); ok {
		if p.Text() != "Hello, world" {
			t.Errorf("Paragraph parse failed: Expected 'Hello, world', got=%s", p.Text())
		}
	} else {
		t.Errorf("Paragraph parse failed: Expected Paragraph, got=%T", element)
	}
}

func TestParseParagraphBetweenElements(t *testing.T) {
	input := `# Header 1
Paragraph test
# Header Too`
	lex := newLexer(input)
	parser := newParser(lex)
	elements, parseErr := parser.parse(eof)

	if parseErr != nil {
		t.Errorf("ParagraphBetweenElements parse failed: %s", parseErr.Error())
	}

	if len(elements) != 3 {
		t.Errorf("ParagraphBetweenElements parse failed: Expected 3 elements, got=%d", len(elements))
		t.FailNow()
	}

	element := elements[1]

	if p, ok := element.(*paragraph); ok {
		if p.Text() != "Paragraph test" {
			t.Errorf("ParagraphBetweenElements parse failed: Expected 'Paragraph test', got=%s", p.Text())
		}
	} else {
		t.Errorf("ParagraphBetweenElements parse failed: Expected Paragraph, got=%T", element)
	}
}

func TestParseParagraphWithNestedElements(t *testing.T) {
	input := "Hello, **world**"
	lex := newLexer(input)
	parser := newParser(lex)
	elements, parseErr := parser.parse(eof)

	if parseErr != nil {
		t.Errorf("ParagraphWithNestedElements parse failed: %s", parseErr.Error())
	}

	if len(elements) != 1 {
		t.Errorf("ParagraphWithNestedElements parse failed: Expected 1 elements, got=%d", len(elements))
	}

	element := elements[0]

	if p, ok := element.(*paragraph); ok {
		if len(p.Content) != 2 {
			t.Errorf("ParagraphWithNestedElements parse failed: Expected 2 children, got=%d", len(p.Content))
			t.FailNow()
		}

		if frag, ok := p.Content[0].(*fragment); ok {
			if frag.String != "Hello, " {
				t.Errorf("ParagraphWithNestedElements parse failed: Expected fragment 'Hello, ', got=%s", frag.String)
			}
		} else {
			t.Errorf("ParagraphWithNestedElements parse failed: Expected Fragment child, got=%T", element)
		}

		if strong, ok := p.Content[1].(*bold); ok {
			if strong.Text != "world" {
				t.Errorf("ParagraphWithNestedElements parse failed: Expected strong 'world', got=%s", strong.Text)
			}
		} else {
			t.Errorf("ParagraphWithNestedElements parse failed: Expected Fragment child, got=%T", element)
		}
	} else {
		t.Errorf("ParagraphWithNestedElements parse failed: Expected Paragraph, got=%T", element)
	}
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
			t.Errorf("Div parse 1 failed: Exected one child, got=%d", len(div.Children))
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
			t.Log(div.Html())
			t.Errorf("Div parse 2 failed: Exected one child, got=%d", len(div.Children))
		}
	} else {
		t.Errorf("Div parse failed: Expected Div type, got=%T", element)
	}
}
