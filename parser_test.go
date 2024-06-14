package mdx

import (
	"fmt"
	"reflect"
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

func validateLength(t *testing.T, actual, expected int) {
	t.Helper()
	if actual != expected {
		fail(t, fmt.Sprintf("Expected %d element(s), got=%d", expected, actual))
		t.FailNow()
	}
}

func TestParseProperties(t *testing.T) {
	input := "{ .class=test } Hello, world"
	elements := execute(t, input)
	validateLength(t, len(elements), 1)
	element := elements[0]

	if p, ok := element.(*paragraph); ok {
		properties := p.Properties

		validateLength(t, len(properties), 1)
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
	validateLength(t, len(elements), 1)
	element := elements[0]

	if p, ok := element.(*paragraph); ok {
		validateLength(t, len(p.Content), 2)

		if frag, ok := p.Content[0].(*fragment); ok {
			if frag.Value != "Hello, " {
				fail(t, fmt.Sprintf("Expected 'Hello, ', got=%s", frag.Value))
			}
		} else {
			fail(t, fmt.Sprintf("Expected Fragment, got=%T", p.Content[0]))
		}

		if s, ok := p.Content[1].(*span); ok {
			validateLength(t, len(s.Content), 1)
			validateLength(t, len(s.Properties), 1)

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
	validateLength(t, len(elements), 2)
	element := elements[1]

	if div, ok := element.(*div); ok {
		validateLength(t, len(div.Properties), 1)

		divProperty := div.Properties[0]
		if divProperty.Name != "class" {
			fail(t, fmt.Sprintf("Expected Div property Name=class, got=%s", divProperty.Name))
		}

		if divProperty.Value != "container" {
			fail(t, fmt.Sprintf("Expected Div property Value=container, got=%s", divProperty.Value))
		}

		validateLength(t, len(div.Children), 2)

		child := div.Children[1]
		if p, ok := child.(*paragraph); ok {
			validateLength(t, len(p.Properties), 2)

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

func TestParseHeader1(t *testing.T) {
	input := "# Heading"
	elements := execute(t, input)
	validateLength(t, len(elements), 1)
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
	validateLength(t, len(elements), 1)
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
	validateLength(t, len(elements), 3)
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
	validateLength(t, len(elements), 1)
	element := elements[0]

	if head, ok := element.(*header); ok {
		if head.Level != 1 {
			fail(t, fmt.Sprintf("Expected Header level 1, got=%d", head.Level))
		}

		validateLength(t, len(head.Content), 2)

		if frag, ok := head.Content[0].(*fragment); ok {
			if frag.Value != "Hello " {
				fail(t, fmt.Sprintf("Expected 'Hello ', got=%s", frag.Value))
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
	validateLength(t, len(elements), 1)
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
	validateLength(t, len(elements), 3)

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
	validateLength(t, len(elements), 1)
	element := elements[0]

	if p, ok := element.(*paragraph); ok {
		validateLength(t, len(p.Content), 2)

		if frag, ok := p.Content[0].(*fragment); ok {
			if frag.Value != "Hello, " {
				fail(t, fmt.Sprintf("Expected fragment 'Hello, ', got=%s", frag.Value))
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

	validateLength(t, len(elements), 1)

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

	validateLength(t, len(elements), 1)

	element = elements[0]
	if p, ok := element.(*paragraph); ok {
		validateLength(t, len(p.Content), 1)
		if frag, ok := p.Content[0].(*fragment); ok {
			if frag.Value != input {
				fail(t, fmt.Sprintf("Got fragment with string %s", frag.Value))
			}
		} else {
			fail(t, fmt.Sprintf("Expected Fragment, got=%T", element))
		}
	} else {
		fail(t, fmt.Sprintf("Expected paragraph, got=%T", element))
	}
}

func TestParseCodeBetweenElements(t *testing.T) {
	input := "# Header\n`hello, world` goodbye, code"
	elements := execute(t, input)
	validateLength(t, len(elements), 2)
	element := elements[1]

	if p, ok := element.(*paragraph); ok {
		validateLength(t, len(p.Content), 2)

		if code, ok := p.Content[0].(*code); ok {
			if code.Text != "hello, world" {
				fail(t, fmt.Sprintf("Expected 'hello, world', got=%s", code.Text))
			}
		} else {
			fail(t, fmt.Sprintf("Expected Code element, got=%T", element))
		}
	} else {
		fail(t, fmt.Sprintf("Expected Paragraph element, got=%T", element))
	}
}

func TestParseCodeWithinElement(t *testing.T) {
	input := "Code: `hello, world` goodbye, code"
	elements := execute(t, input)
	validateLength(t, len(elements), 1)

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
	validateLength(t, len(elements), 1)
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
	validateLength(t, len(elements), 3)
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
	validateLength(t, len(elements), 1)

	if strong, ok := elements[0].(*bold); ok {
		validateLength(t, len(strong.Content), 3)

		if frag, ok := strong.Content[0].(*fragment); ok {
			if frag.Value != "Extreme " {
				fail(t, fmt.Sprintf("Expected Strong Fragment 1 text='Extreme ', got=%s", frag.Value))
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
			if frag.Value != " time" {
				fail(t, fmt.Sprintf("Expected Strong Fragment 1 text=' time', got=%s", frag.Value))
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
	validateLength(t, len(elements), 1)
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
	validateLength(t, len(elements), 3)
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
	validateLength(t, len(elements), 1)
	element := elements[0]

	if em, ok := element.(*italic); ok {
		validateLength(t, len(em.Content), 3)

		if f, ok := em.Content[0].(*fragment); ok {
			if f.Value != "speedy " {
				fail(t, fmt.Sprintf("Expected Em Fragment text='speedy ', got=%s", f.Value))
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
			if f.Value != " session" {
				fail(t, fmt.Sprintf("Expected Em Fragment text=' session', got=%s", f.Value))
			}
		} else {
			fail(t, fmt.Sprintf("Expected Em Fragment, got=%T", em.Content[2]))
		}
	} else {
		fail(t, fmt.Sprintf("Expected Em, got=%T", element))
	}
}

func TestParseBlockQuote(t *testing.T) {
	input := "> Quote me"
	elements := execute(t, input)
	validateLength(t, len(elements), 1)
	element := elements[0]

	if quote, ok := element.(*blockQuote); ok {
		if quote.InnerHtml() != "Quote me" {
			fail(t, fmt.Sprintf("Expected blockQuote text='Quote me', got='%s'", quote.InnerHtml()))
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
	validateLength(t, len(elements), 3)
	element := elements[1]

	if quote, ok := element.(*blockQuote); ok {
		if quote.InnerHtml() != "Quote" {
			fail(t, fmt.Sprintf("Expected blockQuote text='Quote', got=%s", quote.InnerHtml()))
		}
	} else {
		fail(t, fmt.Sprintf("Expected BlockQuote, got=%T", element))
	}
}

func TestParseMultiLayerBlockQuote(t *testing.T) {
	input := "> > Quote"
	elements := execute(t, input)
	validateLength(t, len(elements), 1)
	first := elements[0]

	if quote, ok := first.(*blockQuote); ok {
		validateLength(t, len(quote.Content), 1)

		second := quote.Content[0]

		if innerQuote, ok := second.(*blockQuote); ok {
			validateLength(t, len(innerQuote.Content), 1)

			if frag, ok := innerQuote.Content[0].(*fragment); ok {
				if frag.Value != "Quote" {
					fail(t, fmt.Sprintf("Expected fragment text='Quote', got='%s'", frag.Value))
				}
			}
		} else {
			fail(t, fmt.Sprintf("Expected BlockQuote child, got=%T", second))
		}
	} else {
		fail(t, fmt.Sprintf("Expected BlockQuote, got=%T", first))
	}
}

func TestParseBlockQuoteWithNestedElements(t *testing.T) {
	input := `> Quote
**stronk**
> > Nested
> >
> > Separated`
	elements := execute(t, input)

	if len(elements) != 1 {
		fail(t, fmt.Sprintf("Expected 1 element, got=%d", len(elements)))
		t.FailNow()
	}

	element := elements[0]

	if first, ok := element.(*blockQuote); ok {
		if len(first.Content) != 4 {
			fail(t, fmt.Sprintf("Expected 4 elements, got=%d", len(first.Content)))
			t.FailNow()
		}

		if firstFrag, ok := first.Content[0].(*fragment); ok {
			if firstFrag.Value != "Quote" {
				fail(t, fmt.Sprintf("Expected fragment text='Quote ', got='%s'", firstFrag.Value))
			}
		} else {
			fail(t, fmt.Sprintf("Expected Fragment, got=%T", first.Content[0]))
		}

		if space, ok := first.Content[1].(*fragment); ok {
			if space.Value != " " {
				fail(t, fmt.Sprintf("Expected fragment text=' ', got='%s'", space.Value))
			}
		} else {
			fail(t, fmt.Sprintf("Expected Fragment, got=%T", first.Content[1]))
		}

		if stronk, ok := first.Content[2].(*bold); ok {
			if len(stronk.Content) != 1 {
				fail(t, fmt.Sprintf("Expected 1 bold child, got=%d", len(stronk.Content)))
				t.FailNow()
			}

			if stronkFrag, ok := stronk.Content[0].(*fragment); ok {
				if stronkFrag.Value != "stronk" {
					fail(t, fmt.Sprintf("Expected bold text='stronk', got='%s'", stronkFrag.Value))
				}
			} else {
				fail(t, fmt.Sprintf("Expected Fragment, got=%T", stronk.Content[0]))
			}
		} else {
			fail(t, fmt.Sprintf("Expected Fragment, got=%T", first.Content[3]))
		}

		if second, ok := first.Content[3].(*blockQuote); ok {
			if len(second.Content) != 3 {
				fail(t, fmt.Sprintf("Expected 3 second children, got=%d", len(second.Content)))
				t.FailNow()
			}

			if firstInner, ok := second.Content[0].(*fragment); ok {
				if firstInner.Value != "Nested" {
					fail(t, fmt.Sprintf("Expected second text='Nested', got='%s'", firstInner.Value))
				}
			} else {
				fail(t, fmt.Sprintf("Expected Fragment, got=%T", second.Content[0]))
			}

			if _, ok := second.Content[1].(*lineBreak); !ok {
				fail(t, fmt.Sprintf("Expected LineBreak, got=%T", second.Content[1]))
			}

			if thirdInner, ok := second.Content[2].(*fragment); ok {
				if thirdInner.Value != "Separated" {
					fail(t, fmt.Sprintf("Expected second nested text='Separated', got='%s'", thirdInner.Value))
				}
			} else {
				fail(t, fmt.Sprintf("Expected Fragment, got=%T", second.Content[2]))
			}
		} else {
			fail(t, fmt.Sprintf("Expected BlockQuote, got=%T", first))
		}
	} else {
		fail(t, fmt.Sprintf("Expected BlockQuote, got=%T", first))
	}
}

func TestParseOrderedList(t *testing.T) {
	input := `1. First
2. Second`
	elements := execute(t, input)
	validateLength(t, len(elements), 1)
	element := elements[0]

	if list, ok := element.(*orderedList); ok {
		validateLength(t, len(list.ListItems), 2)

		first := list.ListItems[0]
		if p, ok := first.Component.(*paragraph); ok {
			validateLength(t, len(p.Content), 1)

			if pText, ok := p.Content[0].(*fragment); ok {
				if pText.Value != "First" {
					fail(t, fmt.Sprintf("Expected item text='First', got='%s'", pText.Value))
				}
			} else {
				fail(t, fmt.Sprintf("Expected fragment, got=%T", p.Content[0]))
			}
		} else {
			fail(t, fmt.Sprintf("Expected paragraph list item, got=%T", first.Component))
		}

		second := list.ListItems[1]
		if p, ok := second.Component.(*paragraph); ok {
			validateLength(t, len(p.Content), 1)

			if pText, ok := p.Content[0].(*fragment); ok {
				if pText.Value != "Second" {
					fail(t, fmt.Sprintf("Expected item text='Second', got='%s'", pText.Value))
				}
			} else {
				fail(t, fmt.Sprintf("Expected fragment, got=%T", p.Content[0]))
			}
		} else {
			fail(t, fmt.Sprintf("Expected paragraph list item, got=%T", second.Component))
		}
	} else {
		fail(t, fmt.Sprintf("Expected OrderedList, got=%T", element))
	}
}

func TestParseOrderedListBetweenElements(t *testing.T) {
	input := `# Header
1. First
2. Second
# Another Header`
	elements := execute(t, input)
	validateLength(t, len(elements), 3)
	element := elements[1]

	if list, ok := element.(*orderedList); ok {
		validateLength(t, len(list.ListItems), 2)
	} else {
		fail(t, fmt.Sprintf("Expected OrderedList, got=%T", element))
	}
}

func TestParseUnorderedList(t *testing.T) {
	input := `- First
- Second`

	elements := execute(t, input)
	validateLength(t, len(elements), 1)
	element := elements[0]

	if list, ok := element.(*unorderedList); ok {
		validateLength(t, len(list.ListItems), 2)

		first := list.ListItems[0]
		if p, ok := first.Component.(*paragraph); ok {
			validateLength(t, len(p.Content), 1)

			if pText, ok := p.Content[0].(*fragment); ok {
				if pText.Value != "First" {
					fail(t, fmt.Sprintf("Expected item text='First', got='%s'", pText.Value))
				}
			} else {
				fail(t, fmt.Sprintf("Expected fragment, got=%T", p.Content[0]))
			}
		} else {
			fail(t, fmt.Sprintf("Expected paragraph list item, got=%T", first.Component))
		}

		second := list.ListItems[1]
		if p, ok := second.Component.(*paragraph); ok {
			validateLength(t, len(p.Content), 1)

			if pText, ok := p.Content[0].(*fragment); ok {
				if pText.Value != "Second" {
					fail(t, fmt.Sprintf("Expected item text='Second', got='%s'", pText.Value))
				}
			} else {
				fail(t, fmt.Sprintf("Expected fragment, got=%T", p.Content[0]))
			}
		} else {
			fail(t, fmt.Sprintf("Expected paragraph list item, got=%T", second.Component))
		}
	} else {
		fail(t, fmt.Sprintf("Expected OrderedList, got=%T", element))
	}
}

func TestParseUnorderedListBetweenElements(t *testing.T) {
	input := `# Header
- First
- Second
# Another Header`
	elements := execute(t, input)
	validateLength(t, len(elements), 3)
	element := elements[1]

	if list, ok := element.(*unorderedList); ok {
		validateLength(t, len(list.ListItems), 2)
	} else {
		fail(t, fmt.Sprintf("Expected UnorderedList, got=%T", element))
	}
}

func TesParseImage(t *testing.T) {
	input := "![Image](https://imgurl.com)"
	elements := execute(t, input)
	validateLength(t, len(elements), 1)
	element := elements[0]

	if img, ok := element.(*image); ok {
		if img.AltText != "Image" {
			fail(t, fmt.Sprintf("Expected Image AltText='Image', got='%s'", img.AltText))
		}

		if img.ImgUrl != "https://imgurl.com" {
			fail(t, fmt.Sprintf("Expected Image AltText='https://imgurl.com', got='%s'", img.ImgUrl))
		}
	} else {
		fail(t, fmt.Sprintf("Expected Image, got=%T", element))
	}
}

func TestParseImageBetweenElements(t *testing.T) {
	input := `$ span $
![Image](https://imgurl.com)
Hello`
	elements := execute(t, input)
	validateLength(t, len(elements), 3)
	element := elements[1]

	if img, ok := element.(*image); ok {
		if img.AltText != "Image" {
			fail(t, fmt.Sprintf("Expected Image text='Image', got='%s'", img.AltText))
		}

		if img.ImgUrl != "https://imgurl.com" {
			fail(t, fmt.Sprintf("Expected Image url='https://imgurl.com', got='%s'", img.ImgUrl))
		}
	} else {
		fail(t, fmt.Sprintf("Expected Image, got=%T", element))
	}
}

func TestParseHorizontalRule(t *testing.T) {
	input := "---"
	elements := execute(t, input)
	validateLength(t, len(elements), 1)
	element := elements[0]

	if _, ok := element.(*horizontalRule); !ok {
		fail(t, fmt.Sprintf("Expected HorizontalRule, got=%T", element))
	}

	input = "___"
	elements = execute(t, input)
	validateLength(t, len(elements), 1)
	element = elements[0]

	if _, ok := element.(*horizontalRule); !ok {
		fail(t, fmt.Sprintf("Expected HorizontalRule, got=%T", element))
	}
}

func TestParseHorizontalRuleBetweenElements(t *testing.T) {
	input := `$ span $
---
Hello`
	elements := execute(t, input)
	validateLength(t, len(elements), 3)
	element := elements[1]

	if _, ok := element.(*horizontalRule); !ok {
		fail(t, fmt.Sprintf("Expected HorizontalRule, got=%T", element))
	}

	input = `$ span $
___
Hello`
	elements = execute(t, input)
	validateLength(t, len(elements), 3)
	element = elements[1]

	if _, ok := element.(*horizontalRule); !ok {
		fail(t, fmt.Sprintf("Expected HorizontalRule, got=%T", element))
	}
}

func TestParseLink(t *testing.T) {
	input := "[Text](https://linkurl.com)"
	elements := execute(t, input)
	validateLength(t, len(elements), 1)
	element := elements[0]

	if link, ok := element.(*link); ok {
		if link.InnerHtml() != "Text" {
			fail(t, fmt.Sprintf("Expected Link text='Text', got='%s'", link.InnerHtml()))
		}

		if link.Url != "https://linkurl.com" {
			fail(t, fmt.Sprintf("Expected Link Url='https://linkurl.com', got='%s'", link.Url))
		}
	} else {
		fail(t, fmt.Sprintf("Expected Link, got=%T", element))
	}

	input = "<https://linkurl.com>"
	elements = execute(t, input)
	validateLength(t, len(elements), 1)
	element = elements[0]

	if link, ok := element.(*link); ok {
		if link.InnerHtml() != "https://linkurl.com" {
			fail(t, fmt.Sprintf("Expected Link text='https://linkurl.com', got='%s'", link.InnerHtml()))
		}

		if link.Url != "https://linkurl.com" {
			fail(t, fmt.Sprintf("Expected Link Url='https://linkurl.com', got='%s'", link.Url))
		}
	} else {
		fail(t, fmt.Sprintf("Expected Link, got=%T", element))
	}
}

func TestParseLinkBetweenElements(t *testing.T) {
	input := `# Header
[Link](https://linkurl.com)
> Quote`
	elements := execute(t, input)
	validateLength(t, len(elements), 3)
	element := elements[1]

	if link, ok := element.(*link); ok {
		if link.InnerHtml() != "Link" {
			fail(t, fmt.Sprintf("Expected Link text='Link', got='%s'", link.InnerHtml()))
		}

		if link.Url != "https://linkurl.com" {
			fail(t, fmt.Sprintf("Expected Link Url='https://linkurl.com', got='%s'", link.Url))
		}
	} else {
		fail(t, fmt.Sprintf("Expected Link, got=%T", element))
	}

	input = `# Header
<https://linkurl.com>
> Quote
`
	elements = execute(t, input)
	validateLength(t, len(elements), 3)
	element = elements[1]

	if link, ok := element.(*link); ok {
		if link.InnerHtml() != "https://linkurl.com" {
			fail(t, fmt.Sprintf("Expected Link text='https://linkurl.com', got='%s'", link.InnerHtml()))
		}

		if link.Url != "https://linkurl.com" {
			fail(t, fmt.Sprintf("Expected Link Url='https://linkurl.com', got='%s'", link.Url))
		}
	} else {
		fail(t, fmt.Sprintf("Expected Link, got=%T", element))
	}
}

func TestParseLinkWithNestedElements(t *testing.T) {
	input := "[[ Link ]](https://linkurl.com)"
	elements := execute(t, input)
	validateLength(t, len(elements), 1)
	element := elements[0]

	if link, ok := element.(*link); ok {
		validateLength(t, len(link.Content), 1)
		if div, ok := link.Content[0].(*div); ok {
			validateLength(t, len(div.Children), 1)
			if p, ok := div.Children[0].(*paragraph); ok {
				validateLength(t, len(p.Content), 1)
				if frag, ok := p.Content[0].(*fragment); ok {
					if frag.Value != "Link" {
						fail(t, fmt.Sprintf("Expected Fragment text='Link', got='%s'", frag.Value))
					}
				} else {
					fail(t, fmt.Sprintf("Expected Fragment, got=%T", p.Content[0]))
				}
			} else {
				fail(t, fmt.Sprintf("Expected Paragraph, got=%T", div.Children[0]))
			}
		} else {
			fail(t, fmt.Sprintf("Expected Div, got=%T", link.Content[0]))
		}

		if link.Url != "https://linkurl.com" {
			fail(t, fmt.Sprintf("Expected Link Url='https://linkurl.com', got='%s'", link.Url))
		}
	} else {
		fail(t, fmt.Sprintf("Expected Link, got=%T", element))
	}
}

func TestParseButton(t *testing.T) {
	input := "~[Click Me](handleClick)"
	elements := execute(t, input)
	validateLength(t, len(elements), 1)
	element := elements[0]

	if button, ok := element.(*button); ok {
		if button.OnClick != "handleClick" {
			fail(t, fmt.Sprintf("Expected OnClick='handleClick', got='%s'", button.OnClick))
		}

		validateLength(t, len(button.Content), 1)

		if p, ok := button.Content[0].(*paragraph); ok {
			validateLength(t, len(p.Content), 1)
			if frag, ok := p.Content[0].(*fragment); ok {
				if frag.Value != "Click Me" {
					fail(t, fmt.Sprintf("Expected Fragment text='Click Me', got='%s'", frag.Value))
				}
			} else {
				fail(t, fmt.Sprintf("Expected Fragment, got=%T", p.Content[0]))
			}
		} else {
			fail(t, fmt.Sprintf("Expected Paragraph, got=%T", button.Content[0]))
		}
	}
}

func TestParseButtonBetweenElements(t *testing.T) {
	input := `# Header
~[Click Me](handleClick)
> Quote Me`
	elements := execute(t, input)
	validateLength(t, len(elements), 3)
	element := elements[1]

	if button, ok := element.(*button); ok {
		validateLength(t, len(button.Content), 1)
		if p, ok := button.Content[0].(*paragraph); ok {
			validateLength(t, len(p.Content), 1)
			if frag, ok := p.Content[0].(*fragment); ok {
				if frag.Value != "Click Me" {
					fail(t, fmt.Sprintf("Expected Paragraph text='Click Me', got='%s'", frag.Value))
				}
			} else {
				fail(t, fmt.Sprintf("Expected Fragment, got=%T", p.Content[0]))
			}
		} else {
			fail(t, fmt.Sprintf("Expected Paragraph, got=%T", button.Content[0]))
		}
	} else {
		fail(t, fmt.Sprintf("Expected Button, got=%T", element))
	}
}

func TestParseButtonWithNestedElements(t *testing.T) {
	input := "~[[ Click Me ]](handleClick)"
	elements := execute(t, input)
	validateLength(t, len(elements), 1)
	element := elements[0]

	if button, ok := element.(*button); ok {
		if button.OnClick != "handleClick" {
			fail(t, fmt.Sprintf("Expected OnClick='handleClick', got='%s'", button.OnClick))
		}

		validateLength(t, len(button.Content), 1)
		if div, ok := button.Content[0].(*div); ok {
			validateLength(t, len(div.Children), 1)
			if p, ok := div.Children[0].(*paragraph); ok {
				validateLength(t, len(p.Content), 1)
				if frag, ok := p.Content[0].(*fragment); ok {
					if frag.Value != "Click Me" {
						fail(t, fmt.Sprintf("Expected fragment value='Click Me', got='%s'", frag.Value))
					}
				} else {
					fail(t, fmt.Sprintf("Expected fragment child, got=%T", p.Content[0]))
				}
			} else {
				fail(t, fmt.Sprintf("Expected div child Paragraph, got=%T", div.Children[0]))
			}
		} else {
			fail(t, fmt.Sprintf("Expected child Div, got=%T", button.Content[0]))
		}
	} else {
		fail(t, fmt.Sprintf("Expected Button, got=%T", element))
	}
}

func TestParseDiv(t *testing.T) {
	input := `[
    Hello
]`
	elements := execute(t, input)
	validateLength(t, len(elements), 1)
	element := elements[0]

	if div, ok := element.(*div); ok {
		validateLength(t, len(div.Children), 1)
	} else {
		fail(t, fmt.Sprintf("Expected Div type, got=%T", element))
	}

	input = `[
    Hello
]

`
	elements = execute(t, input)
	validateLength(t, len(elements), 1)
	element = elements[0]

	if div, ok := element.(*div); ok {
		validateLength(t, len(div.Children), 1)
	} else {
		fail(t, fmt.Sprintf("Expected Div type, got=%T", element))
	}
}

func TestParseNav(t *testing.T) {
	input := "@ Navigate @"
	elements := execute(t, input)

	validateLength(t, len(elements), 1)

	element := elements[0]

	if nav, ok := element.(*nav); ok {
		validateLength(t, len(nav.Children), 1)

		child := nav.Children[0]
		if p, ok := child.(*paragraph); ok {
			validateLength(t, len(p.Content), 1)

			if frag, ok := p.Content[0].(*fragment); ok {
				if frag.Value != "Navigate" {
					fail(t, fmt.Sprintf("Expected Nav text='Navigate', got='%s'", frag.Value))
				}
			} else {
				fail(t, fmt.Sprintf("Expected fragment, got=%T", child))
			}
		} else {
			fail(t, fmt.Sprintf("Expected fragment, got=%T", child))
		}
	} else {
		fail(t, fmt.Sprintf("Expected Nav, got=%T", element))
	}
}

func TestParseNavBetweenElements(t *testing.T) {
	input := `# Header
@ Nav @
> Quote`
	elements := execute(t, input)
	validateLength(t, len(elements), 3)
	element := elements[1]

	if nav, ok := element.(*nav); ok {
		validateLength(t, len(nav.Children), 1)
		if p, ok := nav.Children[0].(*paragraph); ok {
			validateLength(t, len(p.Content), 1)
			if frag, ok := p.Content[0].(*fragment); ok {
				if frag.Value != "Nav" {
					fail(t, fmt.Sprintf("Expected fragment text='Nav', got=''%s'", frag.Value))
				}
			} else {
				fail(t, fmt.Sprintf("Expected Fragment, got=%T", p.Content[0]))
			}
		} else {
			fail(t, fmt.Sprintf("Expected Paragraph, got=%T", nav.Children[0]))
		}
	} else {
		fail(t, fmt.Sprintf("Expected Nav, got=%T", element))
	}
}

func TestParseNavWithNestedElements(t *testing.T) {
	input := `@
[ # Title ]
[
	> Quote
]
@
`
	elements := execute(t, input)
	validateLength(t, len(elements), 1)
	element := elements[0]

	if nav, ok := element.(*nav); ok {
		validateLength(t, len(nav.Children), 2)
		if div1, ok := nav.Children[0].(*div); ok {
			validateLength(t, len(div1.Children), 1)
		} else {
			fail(t, fmt.Sprintf("Expected Div, got=%T", nav.Children[0]))
		}

		if div2, ok := nav.Children[1].(*div); ok {
			validateLength(t, len(div2.Children), 1)
		} else {
			fail(t, fmt.Sprintf("Expected Div, got=%T", nav.Children[1]))
		}
	} else {
		fail(t, fmt.Sprintf("Expected Nav, got=%T", element))
	}
}

func TestParseSpan(t *testing.T) {
	input := " $ span $"
	elements := execute(t, input)
	validateLength(t, len(elements), 1)
	element := elements[0]

	if span, ok := element.(*span); ok {
		validateLength(t, len(span.Content), 1)

		child := span.Content[0]
		if frag, ok := child.(*fragment); ok {
			if frag.Value != "span" {
				fail(t, fmt.Sprintf("Expected text='span', got='%s'", frag.Value))
			}
		} else {
			fail(t, fmt.Sprintf("Expected fragment, got=%T", child))
		}
	} else {
		fail(t, fmt.Sprintf("Expected Span, got=%T", element))
	}
}

func TestParseSpanBetweenElements(t *testing.T) {
	input := `# Header
$ span $
> Quote`
	elements := execute(t, input)
	validateLength(t, len(elements), 3)
	element := elements[1]

	if span, ok := element.(*span); ok {
		validateLength(t, len(span.Content), 1)
		if frag, ok := span.Content[0].(*fragment); ok {
			if frag.Value != "span" {
				fail(t, fmt.Sprintf("Expected span text='span', got='%s'", frag.Value))
			}
		} else {
			fail(t, fmt.Sprintf("Expected Fragment, got=%T", span.Content[0]))
		}
	} else {
		fail(t, fmt.Sprintf("Expected Span, got=%T", element))
	}
}

func TestParseSpanWithNestedElements(t *testing.T) {
	input := "$ Hello, `world`! $"
	elements := execute(t, input)
	validateLength(t, len(elements), 1)
	element := elements[0]

	if span, ok := element.(*span); ok {
		validateLength(t, len(span.Content), 3)
		if frag1, ok := span.Content[0].(*fragment); ok {
			if frag1.Value != "Hello, " {
				fail(t, fmt.Sprintf("Expected Fragment text='Hello, ', got='%s'", frag1.Value))
			}
		}

		if code, ok := span.Content[1].(*code); ok {
			if code.Text != "world" {
				fail(t, fmt.Sprintf("Expected Fragment text='world', got='%s'", code.Text))
			}
		}

		if frag2, ok := span.Content[2].(*fragment); ok {
			if frag2.Value != "!" {
				fail(t, fmt.Sprintf("Expected Fragment text='!', got='%s'", frag2.Value))
			}
		}
	} else {
		fail(t, fmt.Sprintf("Expected Span, got=%T", element))
	}
}

func TestParseCodeBlock(t *testing.T) {
	input := `^^
func main() {
	fmt.Println("Hello, world!")
}
^^`
	elements := execute(t, input)
	validateLength(t, len(elements), 1)
	element := elements[0]

	if codeBlock, ok := element.(*codeBlock); ok {
		expectedCode := "func main() {\\n    fmt.Println(\"Hello, world!\")\\n}"
		if codeBlock.Content != expectedCode {
			fail(t, fmt.Sprintf("Expected content='%s', got='%s'", expectedCode, codeBlock.Content))
		}
	} else {
		fail(t, fmt.Sprintf("Expected CodeBlock, got=%T", element))
	}
}

func TestParseCodeBlockBetweenElements(t *testing.T) {
	input := `# Heading
^^
func main() {
	fmt.Println("Hello, world!")
}
^^
> Quote`
	elements := execute(t, input)
	validateLength(t, len(elements), 3)
	element := elements[1]

	if codeBlock, ok := element.(*codeBlock); ok {
		expectedCode := "func main() {\\n    fmt.Println(\"Hello, world!\")\\n}"
		if codeBlock.Content != expectedCode {
			fail(t, fmt.Sprintf("Expected content='%s', got='%s'", expectedCode, codeBlock.Content))
		}
	} else {
		fail(t, fmt.Sprintf("Expected CodeBlock, got=%T", element))
	}
}

func TestParser(t *testing.T) {
	inputs := map[string][]component{
		"test\ntest": {
			&paragraph{Content: []component{&fragment{Value: "test test"}}},
		},
		"test\n\ntest": {
			&paragraph{Content: []component{&fragment{Value: "test"}}},
			&paragraph{Content: []component{&fragment{Value: "test"}}},
		},
	}

	for test, expected := range inputs {
		actual := execute(t, test)
		if !reflect.DeepEqual(actual, expected) {
			fail(t, fmt.Sprintf("Expected %q, got=%q", expected, actual))
		}
	}
}
