package mdx

import (
	"fmt"
	"testing"
)

func execute(t *testing.T, input string) []component {
	t.Helper()
	elements, parseErr := newParser(newLexer(input)).parse(eof)
	if parseErr != nil {
		fail(t, parseErr.Error())
	}
	return elements
}

func fail(t *testing.T, message string) {
	t.Helper()
	t.Errorf("%s failed: %s", t.Name(), message)
}

func TestParseProperties(t *testing.T) {
	input := "{ .class=test } Hello, world"
	elements := execute(t, input)

	if len(elements) != 1 {
		fail(t, fmt.Sprintf("Expected 1 element, got=%d", len(elements)))
		t.FailNow()
	}

	element := elements[0]

	if p, ok := element.(*paragraph); ok {
		properties := p.Properties

		if len(properties) != 1 {
			fail(t, fmt.Sprintf("Expected 1 property, got=%d", len(properties)))
			t.FailNow()
		}

		property := properties[0]

		if property.Name != "class" {
			fail(t, fmt.Sprintf("Expected Name 'class', got=%s", property.Name))
		}

		if property.Value != "test" {
			fail(t, fmt.Sprintf("Expected Value 'test', got=%s", property.Name))
		}
	} else {
		fail(t, fmt.Sprintf("Expected Paragraph, got=%T", element))
	}
}

func TestParsePropertiesInline(t *testing.T) {
	input := "Hello, { .class=groovy } $ world $"
	elements := execute(t, input)

	if len(elements) != 1 {
		fail(t, fmt.Sprintf("Expected 1 element, got=%d", len(elements)))
		t.FailNow()
	}

	element := elements[0]
	if p, ok := element.(*paragraph); ok {
		if len(p.Content) != 2 {
			fail(t, fmt.Sprintf("Expected 2 child elements, got=%d", len(p.Content)))
			t.FailNow()
		}

		if frag, ok := p.Content[0].(*fragment); ok {
			if frag.String != "Hello, " {
				fail(t, fmt.Sprintf("Expected 'Hello, ', got=%s", frag.String))
			}
		} else {
			fail(t, fmt.Sprintf("Expected Fragment, got=%T", p.Content[0]))
		}

		if s, ok := p.Content[1].(*span); ok {
			if len(s.Children) != 1 {
				fail(t, fmt.Sprintf("Expected 1 child, got=%d", len(s.Children)))
			}

			if len(s.Properties) != 1 {
				fail(t, fmt.Sprintf("Expected 1 property, got=%d", len(s.Properties)))
			}

			properties := s.Properties[0]

			if properties.Name != "class" {
				fail(t, fmt.Sprintf("Expected Name 'class', got=%s", properties.Name))
			}

			if properties.Value != "groovy" {
				fail(t, fmt.Sprintf("Expected Value 'groovy', got=%s", properties.Name))
			}
		} else {
			fail(t, fmt.Sprintf("Expected paragraph (2), got=%T", elements[1]))
		}
	} else {
		fail(t, fmt.Sprintf("Expected paragraph (1), got=%T", elements[0]))
	}

}

func TestParseNestedProperties(t *testing.T) {
	input := `# Hello
{ .class=container }
[
	## Section
	{ .class=content .data-parent=container }
	Zuzzy
]`
	elements := execute(t, input)

	if len(elements) != 2 {
		fail(t, fmt.Sprintf("Expected 2 elements, got=%d", len(elements)))
	}

	element := elements[1]

	if div, ok := element.(*div); ok {
		if len(div.Properties) != 1 {
			fail(t, fmt.Sprintf("Expected 1 Div property,, got=%d", len(div.Properties)))
			t.FailNow()
		}

		divProperty := div.Properties[0]
		if divProperty.Name != "class" {
			fail(t, fmt.Sprintf("Expected Div property Name=class, got=%s", divProperty.Name))
		}

		if divProperty.Value != "container" {
			fail(t, fmt.Sprintf("Expected Div property Value=container, got=%s", divProperty.Value))
		}

		if len(div.Children) != 2 {
			fail(t, fmt.Sprintf("Expected 2 Div children, got=%d", len(div.Children)))
		}

		child := div.Children[1]
		if p, ok := child.(*paragraph); ok {
			if len(p.Properties) != 2 {
				fail(t, fmt.Sprintf("Expected 2 Paragraph properties, got=%d", len(p.Properties)))
			}

			if p.Properties[0].Name != "class" {
				fail(t, fmt.Sprintf("Expected Paragraph property Name=class, got=%s", p.Properties[0].Name))

			}

			if p.Properties[0].Value != "content" {
				fail(t, fmt.Sprintf("Expected Paragraph property Value=content, got=%s", p.Properties[0].Name))
			}

			if p.Properties[1].Name != "data-parent" {
				fail(t, fmt.Sprintf("Expected Paragraph property Name=data-parent, got=%s", p.Properties[1].Name))

			}

			if p.Properties[1].Value != "container" {
				fail(t, fmt.Sprintf("Expected Paragraph property Value=container, got=%s", p.Properties[1].Value))
			}
		} else {
			fail(t, fmt.Sprintf("Expected Paragraph child, got=%T", child))
		}
	} else {
		fail(t, fmt.Sprintf("Expected Div, got=%T", element))
	}
}

// TODO: Fragment tests
func TestParseFragment(t *testing.T) {

}

func TestParseHeader1(t *testing.T) {
	input := "# Heading"
	elements := execute(t, input)

	if len(elements) != 1 {
		fail(t, fmt.Sprintf("Expected 1 elements, got=%d", len(elements)))
	}

	element := elements[0]
	if head, ok := element.(*header); ok {
		if head.Level != 1 {
			fail(t, fmt.Sprintf("Expected Header level 1, got=%d", head.Level))
		}

		if head.InnerHtml() != "Heading" {
			fail(t, fmt.Sprintf("Expected text 'Heading', got=%s", head.InnerHtml()))
		}
	} else {
		fail(t, fmt.Sprintf("Expected Heading, got=%T", element))
	}
}

func TestParseHeader2(t *testing.T) {
	input := "## Heading Too"
	elements := execute(t, input)

	if len(elements) != 1 {
		fail(t, fmt.Sprintf("Expected 1 element, got=%d", len(elements)))
	}

	element := elements[0]
	if head, ok := element.(*header); ok {
		if head.Level != 2 {
			fail(t, fmt.Sprintf("Expected Header level 1, got=%d", head.Level))
		}

		if head.InnerHtml() != "Heading Too" {
			fail(t, fmt.Sprintf("Expected text 'Heading', got=%s", head.InnerHtml()))
		}
	} else {
		fail(t, fmt.Sprintf("Expected Heading, got=%T", element))
	}
}

func TestParseHeaderBetweenElements(t *testing.T) {
	input := `Hello
# World
MDX`
	elements := execute(t, input)

	if len(elements) != 3 {
		fail(t, fmt.Sprintf("Expected 3 elements, got=%d", len(elements)))
		t.FailNow()
	}

	element := elements[1]

	if head, ok := element.(*header); ok {
		if head.Level != 1 {
			fail(t, fmt.Sprintf("Expected Header level 1, got=%d", head.Level))
		}

		if head.InnerHtml() != "World" {
			fail(t, fmt.Sprintf("Expected text 'World', got=%s", head.InnerHtml()))
		}
	} else {
		fail(t, fmt.Sprintf("Expected Header, got=%T", element))
	}
}

func TestParseHeaderWithNestedElements(t *testing.T) {
	input := `# Hello **world**`
	elements := execute(t, input)

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
			if strong.InnerHtml() != "world" {
				fail(t, fmt.Sprintf("Expected 'world', got=%s", strong.InnerHtml()))
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
	elements := execute(t, input)

	if len(elements) != 1 {
		fail(t, fmt.Sprintf("Expected 1 elements, got=%d", len(elements)))
	}

	element := elements[0]

	if p, ok := element.(*paragraph); ok {
		if p.InnerHtml() != "Hello, world" {
			fail(t, fmt.Sprintf("Expected 'Hello, world', got=%s", p.InnerHtml()))
		}
	} else {
		fail(t, fmt.Sprintf("Expected Paragraph, got=%T", element))
	}
}

func TestParseParagraphBetweenElements(t *testing.T) {
	input := `# Header 1
Paragraph test
# Header Too`
	elements := execute(t, input)

	if len(elements) != 3 {
		fail(t, fmt.Sprintf("Expected 3 elements, got=%d", len(elements)))
		t.FailNow()
	}

	element := elements[1]

	if p, ok := element.(*paragraph); ok {
		if p.InnerHtml() != "Paragraph test" {
			fail(t, fmt.Sprintf("Expected 'Paragraph test', got=%s", p.InnerHtml()))
		}
	} else {
		fail(t, fmt.Sprintf("Expected Paragraph, got=%T", element))
	}
}

func TestParseParagraphWithNestedElements(t *testing.T) {
	input := "Hello, **world**"
	elements := execute(t, input)

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
			if strong.InnerHtml() != "world" {
				fail(t, fmt.Sprintf("Expected strong 'world', got=%s", strong.InnerHtml()))
			}
		} else {
			fail(t, fmt.Sprintf("Expected Fragment child, got=%T", element))
		}
	} else {
		fail(t, fmt.Sprintf("Expected Paragraph, got=%T", element))
	}
}

func TestParseCode(t *testing.T) {
	input := "`print('Hello, world!')`"
	elements := execute(t, input)

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
	elements = execute(t, input)

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

func TestParseCodeBetweenElements(t *testing.T) {
	input := "# Header\n"
	input += "`hello, world` "
	input += "goodbye, code"
	elements := execute(t, input)

	if len(elements) != 3 {
		fail(t, fmt.Sprintf("Expected 3 elements, got=%d", len(elements)))
		t.FailNow()
	}

	element := elements[1]

	if code, ok := element.(*code); ok {
		if code.Text != "hello, world" {
			fail(t, fmt.Sprintf("Expected 'hello, world', got=%s", code.Text))
		}
	} else {
		fail(t, fmt.Sprintf("Expected Code element, got=%T", element))
	}
}

func TestParseCodeWithinElement(t *testing.T) {
	input := "Code: `hello, world` goodbye, code"
	elements := execute(t, input)

	if len(elements) != 1 {
		fail(t, fmt.Sprintf("Expected 1 element, got=%d", len(elements)))
		t.FailNow()
	}

	if p, ok := elements[0].(*paragraph); ok {
		if len(p.Content) != 3 {
			fail(t, fmt.Sprintf("Expected 3 Paragraph children, got=%d", len(p.Content)))
			t.FailNow()
		}

		if code, ok := p.Content[1].(*code); ok {
			if code.Text != "hello, world" {
				fail(t, fmt.Sprintf("Expected 'hello, world', got=%s", code.Text))
			}
		} else {
			fail(t, fmt.Sprintf("Expected Code element, got=%T", p.Content[1]))
		}
	}
}

func TestParseStrong(t *testing.T) {
	input := "**stronk**"
	elements := execute(t, input)

	if len(elements) != 1 {
		fail(t, fmt.Sprintf("Expected 1 element, got=%d", len(elements)))
		t.FailNow()
	}

	element := elements[0]
	if strong, ok := element.(*bold); ok {
		if strong.InnerHtml() != "stronk" {
			fail(t, fmt.Sprintf("Expected text 'stronk', got=%s", strong.InnerHtml()))
		}
	} else {
		fail(t, fmt.Sprintf("Expected Bold, got=%T", element))
	}
}

func TestParseStrongBetweenElements(t *testing.T) {
	input := `# Header
**stronk**
MDX`
	elements := execute(t, input)

	if len(elements) != 3 {
		fail(t, fmt.Sprintf("Expected 3 elements, got=%d", len(elements)))
		t.FailNow()
	}

	element := elements[1]

	if strong, ok := element.(*bold); ok {
		if strong.InnerHtml() != "stronk" {
			fail(t, fmt.Sprintf("Expected text 'stronk', got=%s", strong.InnerHtml()))
		}
	} else {
		fail(t, fmt.Sprintf("Expected Bold, got=%T", element))
	}
}

func TestParseStrongWithNestedElements(t *testing.T) {
	input := "**Extreme `coding` time**"
	elements := execute(t, input)

	if len(elements) != 1 {
		fail(t, fmt.Sprintf("Expected 1 element, got=%d", len(elements)))
		t.FailNow()
	}

	if strong, ok := elements[0].(*bold); ok {
		if len(strong.Content) != 3 {
			fail(t, fmt.Sprintf("Expected 3 bold children, got=%d", len(strong.Content)))
			t.FailNow()
		}

		if frag, ok := strong.Content[0].(*fragment); ok {
			if frag.String != "Extreme " {
				fail(t, fmt.Sprintf("Expected Strong Fragment 1 text='Extreme ', got=%s", frag.String))
			}
		} else {
			fail(t, fmt.Sprintf("Expected Fragment, got=%T", strong.Content[0]))
		}

		if code, ok := strong.Content[1].(*code); ok {
			if code.Text != "coding" {
				fail(t, fmt.Sprintf("Expected Code text='coding', got=%s", code.Text))
			}
		} else {
			fail(t, fmt.Sprintf("Expected Code child, got=%T", strong.Content[1]))
		}

		if frag, ok := strong.Content[2].(*fragment); ok {
			if frag.String != " time" {
				fail(t, fmt.Sprintf("Expected Strong Fragment 1 text=' time', got=%s", frag.String))
			}
		} else {
			fail(t, fmt.Sprintf("Expected Fragment, got=%T", strong.Content[2]))
		}
	} else {
		fail(t, fmt.Sprintf("Expected Paragraph, got=%T", elements[0]))
	}
}

func TestParseEm(t *testing.T) {
	input := "*slinky*"
	elements := execute(t, input)

	if len(elements) != 1 {
		fail(t, fmt.Sprintf("Expected 1 element, got=%d", len(elements)))
		t.FailNow()
	}

	element := elements[0]

	if em, ok := element.(*italic); ok {
		if em.InnerHtml() != "slinky" {
			fail(t, fmt.Sprintf("Expected Em text='slinky', got=%s", em.InnerHtml()))
		}
	} else {
		fail(t, fmt.Sprintf("Expected Em, got=%T", element))
	}
}

func TestParseEmBetweenElements(t *testing.T) {
	input := `# Header
*slinky*
MDX`
	elements := execute(t, input)

	if len(elements) != 3 {
		fail(t, fmt.Sprintf("Expected 3 elements, got=%d", len(elements)))
		t.FailNow()
	}

	element := elements[1]

	if em, ok := element.(*italic); ok {
		if em.InnerHtml() != "slinky" {
			fail(t, fmt.Sprintf("Expected Em text='slinky', got=%s", em.InnerHtml()))
		}
	} else {
		fail(t, fmt.Sprintf("Expected Em, got=%T", element))
	}
}

func TestParseEmWithNestedElements(t *testing.T) {
	input := "*speedy `coding` session*"
	elements := execute(t, input)

	if len(elements) != 1 {
		fail(t, fmt.Sprintf("Expected 1 element, got=%d", len(elements)))
		t.FailNow()
	}

	element := elements[0]

	if em, ok := element.(*italic); ok {
		if len(em.Content) != 3 {
			fail(t, fmt.Sprintf("Expected 3 Em children, got=%d", len(em.Content)))
		}

		if f, ok := em.Content[0].(*fragment); ok {
			if f.String != "speedy " {
				fail(t, fmt.Sprintf("Expected Em Fragment text='speedy ', got=%s", f.String))
			}
		} else {
			fail(t, fmt.Sprintf("Expected Em Fragment, got=%T", em.Content[0]))
		}

		if c, ok := em.Content[1].(*code); ok {
			if c.Text != "coding" {
				fail(t, fmt.Sprintf("Expected Em Code text='coding', got=%s", c.Text))
			}
		} else {
			fail(t, fmt.Sprintf("Expected Em Code, got=%T", em.Content[1]))
		}

		if f, ok := em.Content[2].(*fragment); ok {
			if f.String != " session" {
				fail(t, fmt.Sprintf("Expected Em Fragment text=' session', got=%s", f.String))
			}
		} else {
			fail(t, fmt.Sprintf("Expected Em Fragment, got=%T", em.Content[2]))
		}
	} else {
		fail(t, fmt.Sprintf("Expected Em, got=%T", element))
	}
}

// TODO: Block Quote tests
func TestParseBlockQuote(t *testing.T) {
	input := "> Quote me"
	elements := execute(t, input)

	if len(elements) != 1 {
		fail(t, fmt.Sprintf("Expected 1 element, got=%d", len(elements)))
		t.FailNow()
	}

	element := elements[0]

	if quote, ok := element.(*blockQuote); ok {
		if quote.Text != "Quote me" {
			fail(t, fmt.Sprintf("Expected blockQuote text='Quote me', got=%s", quote.Text))
		}
	} else {
		fail(t, fmt.Sprintf("Expected BlockQuote, got=%T", element))
	}
}

func TestParseBlockQuoteBetweenElements(t *testing.T) {
	input := `# Head
> Quote

MDX`
	elements := execute(t, input)

	if len(elements) != 4 {
		fail(t, fmt.Sprintf("Expected 4 elements, got=%d", len(elements)))
		t.FailNow()
	}

	element := elements[1]

	if quote, ok := element.(*blockQuote); ok {
		if quote.Text != "Quote" {
			fail(t, fmt.Sprintf("Expected blockQuote text='Quote', got=%s", quote.Text))
		}
	} else {
		fail(t, fmt.Sprintf("Expected BlockQuote, got=%T", element))
	}
}

func TestParseBlockQuoteWithNestedElements(t *testing.T) {
	// 	input := `> Quote
	// **stronk**
	// > > Nested`
	// 	elements := execute(t, input)

	//	if len(elements) != 1 {
	//		fail(t, fmt.Sprintf("Expected 1 element, got=%d", len(elements)))
	//	}
}

// TODO: List Item tests
func TestParseListItem(t *testing.T) {

}

func TestParseListItemBetweenElements(t *testing.T) {

}

func TestParseListItemWithNestedElements(t *testing.T) {

}

// TODO: Ordered List tests
func TestParseOrderedList(t *testing.T) {

}

func TestParseOrderedListBetweenElements(t *testing.T) {

}

// TODO: Unordered List tests
func TestParseUnorderedList(t *testing.T) {

}

func TestParseUnorderedListBetweenElements(t *testing.T) {

}

// TODO: Image tests
func TesParsetImage(t *testing.T) {

}

func TestParseImageBetweenElements(t *testing.T) {

}

// TODO: Horizontal Rule tests
func TestParseHorizontalRule(t *testing.T) {

}

func TestParseHorizontalRuleBetweenElements(t *testing.T) {

}

// TODO: Link tests
func TestParseLink(t *testing.T) {

}

func TestParseLinkBetweenElements(t *testing.T) {

}

func TestParseLinkWithNestedElements(t *testing.T) {

}

// TODO: Button tests
func TestParseButton(t *testing.T) {

}

func TestParseButtonBetweenElements(t *testing.T) {

}

func TestParseButtonWithNestedElements(t *testing.T) {

}

func TestParseDiv(t *testing.T) {
	input := `[
    Hello
]`
	lex := newLexer(input)
	parser := newParser(lex)
	elements, parseErr := parser.parse(eof)
	if parseErr != nil {
		fail(t, parseErr.Error())
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
		fail(t, parseErr.Error())
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

// TODO: Nav tests
func TestParseNav(t *testing.T) {

}

func TestParseNavBetweenElements(t *testing.T) {

}

func TestParseNavWithNestedElements(t *testing.T) {

}

// TODO: Span tests
func TestParseSpan(t *testing.T) {

}

func TestParseSpanBetweenElements(t *testing.T) {

}

func TestParseSpanWithNestedElements(t *testing.T) {

}

// TODO: Code Block tests
func TestParseCodeBlock(t *testing.T) {

}

func TestParseCodeBlockBetweenElements(t *testing.T) {

}

func TestParseCodeBlockWithNestedElements(t *testing.T) {

}

// TODO: Body tests (needed?)
func TestParseBody(t *testing.T) {

}

func TestParseBodyBetweenElements(t *testing.T) {

}

func TestParseBodyWithNestedElements(t *testing.T) {

}
