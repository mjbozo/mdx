package mdx

import (
	"fmt"
	"testing"
)

func fail(t *testing.T, message string) {
	t.Helper()
	t.Errorf("%s failed: %s", t.Name(), message)
}

func TestParseHeader1(t *testing.T) {
	input := "# Heading"
	lex := newLexer(input)
	parser := newParser(lex)
	elements, parseErr := parser.parse(eof)

	if parseErr != nil {
		fail(t, parseErr.Error())
	}

	if len(elements) != 1 {
		fail(t, fmt.Sprintf("Expected 1 elements, got=%d", len(elements)))
	}

	element := elements[0]
	if head, ok := element.(*header); ok {
		if head.Level != 1 {
			fail(t, fmt.Sprintf("Expected Header level 1, got=%d", head.Level))
		}

		if head.Text() != "Heading" {
			fail(t, fmt.Sprintf("Expected text 'Heading', got=%s", head.Text()))
		}
	} else {
		fail(t, fmt.Sprintf("Expected Heading, got=%T", element))
	}
}

func TestParseHeader2(t *testing.T) {
	input := "## Heading Too"
	lex := newLexer(input)
	parser := newParser(lex)
	elements, parseErr := parser.parse(eof)

	if parseErr != nil {
		fail(t, fmt.Sprintf("%s", parseErr.Error()))
	}

	if len(elements) != 1 {
		fail(t, fmt.Sprintf("Expected 1 element, got=%d", len(elements)))
	}

	element := elements[0]
	if head, ok := element.(*header); ok {
		if head.Level != 2 {
			fail(t, fmt.Sprintf("Expected Header level 1, got=%d", head.Level))
		}

		if head.Text() != "Heading Too" {
			fail(t, fmt.Sprintf("Expected text 'Heading', got=%s", head.Text()))
		}
	} else {
		fail(t, fmt.Sprintf("Expected Heading, got=%T", element))
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
		fail(t, fmt.Sprintf("%s", parseErr.Error()))
	}

	if len(elements) != 3 {
		fail(t, fmt.Sprintf("Expected 3 elements, got=%d", len(elements)))
		t.FailNow()
	}

	element := elements[1]

	if head, ok := element.(*header); ok {
		if head.Level != 1 {
			fail(t, fmt.Sprintf("Expected Header level 1, got=%d", head.Level))
		}

		if head.Text() != "World" {
			fail(t, fmt.Sprintf("Expected text 'World', got=%s", head.Text()))
		}
	} else {
		fail(t, fmt.Sprintf("Expected Header, got=%T", element))
	}
}

func TestParseHeaderWithNestedElements(t *testing.T) {
	input := `# Hello **world**`
	lex := newLexer(input)
	parser := newParser(lex)
	elements, parseErr := parser.parse(eof)

	if parseErr != nil {
		fail(t, fmt.Sprintf("%s", parseErr))
	}

	if len(elements) != 1 {
		fail(t, fmt.Sprintf("Epected 1 element, got=%d", len(elements)))
		t.FailNow()
	}

	element := elements[0]

	if head, ok := element.(*header); ok {
		if head.Level != 1 {
			fail(t, fmt.Sprintf("Expected Header level 1, got=%d", head.Level))
		}

		if len(head.Content) != 2 {
			fail(t, fmt.Sprintf("Expected 2 children, got=%d", len(head.Content)))
			t.FailNow()
		}

		if frag, ok := head.Content[0].(*fragment); ok {
			if frag.String != "Hello " {
				fail(t, fmt.Sprintf("Expected 'Hello ', got=%s", frag.String))
			}
		} else {
			fail(t, fmt.Sprintf("Expected Fragment child, got=%T", element))
		}

		if strong, ok := head.Content[1].(*bold); ok {
			if strong.Text != "world" {
				fail(t, fmt.Sprintf("Expected 'world', got=%s", strong.Text))
			}
		} else {
			fail(t, fmt.Sprintf("Expected Strong child, got=%T", element))
		}
	} else {
		fail(t, fmt.Sprintf("Expected Header, got=%T", element))
	}
}

func TestParseParagraph(t *testing.T) {
	input := "Hello, world"
	lex := newLexer(input)
	parser := newParser(lex)
	elements, parseErr := parser.parse(eof)

	if parseErr != nil {
		fail(t, fmt.Sprintf("%s", parseErr.Error()))
	}

	if len(elements) != 1 {
		fail(t, fmt.Sprintf("Expected 1 elements, got=%d", len(elements)))
	}

	element := elements[0]

	if p, ok := element.(*paragraph); ok {
		if p.Text() != "Hello, world" {
			fail(t, fmt.Sprintf("Expected 'Hello, world', got=%s", p.Text()))
		}
	} else {
		fail(t, fmt.Sprintf("Expected Paragraph, got=%T", element))
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
		fail(t, fmt.Sprintf("%s", parseErr.Error()))
	}

	if len(elements) != 3 {
		fail(t, fmt.Sprintf("Expected 3 elements, got=%d", len(elements)))
		t.FailNow()
	}

	element := elements[1]

	if p, ok := element.(*paragraph); ok {
		if p.Text() != "Paragraph test" {
			fail(t, fmt.Sprintf("Expected 'Paragraph test', got=%s", p.Text()))
		}
	} else {
		fail(t, fmt.Sprintf("Expected Paragraph, got=%T", element))
	}
}

func TestParseParagraphWithNestedElements(t *testing.T) {
	input := "Hello, **world**"
	lex := newLexer(input)
	parser := newParser(lex)
	elements, parseErr := parser.parse(eof)

	if parseErr != nil {
		fail(t, fmt.Sprintf("%s", parseErr.Error()))
	}

	if len(elements) != 1 {
		fail(t, fmt.Sprintf("Expected 1 elements, got=%d", len(elements)))
	}

	element := elements[0]

	if p, ok := element.(*paragraph); ok {
		if len(p.Content) != 2 {
			fail(t, fmt.Sprintf("Expected 2 children, got=%d", len(p.Content)))
			t.FailNow()
		}

		if frag, ok := p.Content[0].(*fragment); ok {
			if frag.String != "Hello, " {
				fail(t, fmt.Sprintf("Expected fragment 'Hello, ', got=%s", frag.String))
			}
		} else {
			fail(t, fmt.Sprintf("Expected Fragment child, got=%T", element))
		}

		if strong, ok := p.Content[1].(*bold); ok {
			if strong.Text != "world" {
				fail(t, fmt.Sprintf("Expected strong 'world', got=%s", strong.Text))
			}
		} else {
			fail(t, fmt.Sprintf("Expected Fragment child, got=%T", element))
		}
	} else {
		fail(t, fmt.Sprintf("Expected Paragraph, got=%T", element))
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
		fail(t, fmt.Sprintf("%s", parseErr.Error()))
	}

	if len(elements) != 1 {
		fail(t, fmt.Sprintf("Expected 1 element, got=%d", len(elements)))
	}

	element := elements[0]
	if code, ok := element.(*code); ok {
		if code.Text != "print('Hello, world!')" {
			fail(t, fmt.Sprintf("Content incorrect, got=%s", code.Text))
		}
	} else {
		fail(t, fmt.Sprintf("Expected Code type, got=%T", element))
	}

	input = "`print('Hello, world!')"
	lex = newLexer(input)
	parser = newParser(lex)
	elements, parseErr = parser.parse(eof)
	if parseErr != nil {
		fail(t, fmt.Sprintf("%s", parseErr.Error()))
	}

	if len(elements) != 1 {
		fail(t, fmt.Sprintf("Expected 1 element, got=%d", len(elements)))
	}

	element = elements[0]
	if frag, ok := element.(*fragment); ok {
		if frag.String != input {
			fail(t, fmt.Sprintf("Got fragment with string %s", frag.String))
		}
	} else {
		fail(t, fmt.Sprintf("Expected Fragment, got=%T", element))
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
		fail(t, fmt.Sprintf("%s", parseErr.Error()))
	}

	if len(elements) != 1 {
		fail(t, fmt.Sprintf("Expected 1 element, got=%d", len(elements)))
	}

	element := elements[0]
	if div, ok := element.(*div); ok {
		if len(div.Children) != 1 {
			fail(t, fmt.Sprintf("Exected one child, got=%d", len(div.Children)))
		}
	} else {
		fail(t, fmt.Sprintf("Expected Div type, got=%T", element))
	}

	input = `[
    Hello
]

`
	lex = newLexer(input)
	parser = newParser(lex)
	elements, parseErr = parser.parse(eof)
	if parseErr != nil {
		fail(t, fmt.Sprintf("%s", parseErr.Error()))
	}

	if len(elements) != 1 {
		fail(t, fmt.Sprintf("Expected 1 element, got=%d", len(elements)))
	}

	element = elements[0]
	if div, ok := element.(*div); ok {
		if len(div.Children) != 1 {
			fail(t, fmt.Sprintf("Exected one child, got=%d", len(div.Children)))
		}
	} else {
		fail(t, fmt.Sprintf("Expected Div type, got=%T", element))
	}
}