package ast

import (
	"testing"
)

func defaultProps(t *testing.T) []Property {
	t.Helper()

	properties := make([]Property, 0)
	properties = append(properties, Property{Name: "class", Value: "test"})
	properties = append(properties, Property{Name: "style", Value: "background-color: red"})

	return properties
}

func TestFragment(t *testing.T) {
	fragment := Fragment{}
	fragmentHtml := fragment.Html()
	expected := ""
	if fragmentHtml != expected {
		t.Errorf("Fragment wrong, got=%q", fragmentHtml)
	}

	fragment.String = "Hello"
	fragmentHtml = fragment.Html()
	expected = "Hello"
	if fragmentHtml != expected {
		t.Errorf("Fragment wrong, got=%q", fragmentHtml)
	}
}

func TestHeader(t *testing.T) {
	header := Header{Level: 1, Text: "Test"}
	headerHtml := header.Html()
	expected := "<h1>Test</h1>"
	if headerHtml != expected {
		t.Errorf("Header wrong, got=%q", headerHtml)
	}

	header = Header{Level: 2, Text: "Test2"}
	headerHtml = header.Html()
	expected = "<h2>Test2</h2>"
	if headerHtml != expected {
		t.Errorf("Header wrong, got=%q", headerHtml)
	}

	properties := defaultProps(t)
	header = Header{Level: 6, Text: "Test with Props", Properties: properties}
	headerHtml = header.Html()
	expected = "<h6 class=\"test\" style=\"background-color: red\">Test with Props</h6>"
	if headerHtml != expected {
		t.Errorf("Header wrong, got=%q", headerHtml)
	}
}

func TestParagraph(t *testing.T) {
	paragraph := Paragraph{Text: "Paragraph test"}
	paragraphHtml := paragraph.Html()
	expected := "<p>Paragraph test</p>"
	if paragraphHtml != expected {
		t.Errorf("Paragraph wrong, got=%q", paragraphHtml)
	}
}

func TestCode(t *testing.T) {
	code := Code{Text: "fmt.Printf(\"Hello, world!\n\")"}
	codeHtml := code.Html()
	expected := "<code>fmt.Printf(\"Hello, world!\n\")</code>"
	if codeHtml != expected {
		t.Errorf("Code wrong, got=%q", codeHtml)
	}

	code.Properties = defaultProps(t)
	codeHtml = code.Html()
	expected = "<code class=\"test\" style=\"background-color: red\">fmt.Printf(\"Hello, world!\n\")</code>"
	if codeHtml != expected {
		t.Errorf("Code properties wrong, got=%q", codeHtml)
	}
}

func TestBold(t *testing.T) {
	bold := Bold{Text: "stronk"}
	boldHtml := bold.Html()
	expected := "<strong>stronk</strong>"
	if boldHtml != expected {
		t.Errorf("Bold wrong, got=%q", boldHtml)
	}

	bold.Properties = defaultProps(t)
	boldHtml = bold.Html()
	expected = "<strong class=\"test\" style=\"background-color: red\">stronk</strong>"
	if boldHtml != expected {
		t.Errorf("Bold properties wrong, got=%q", boldHtml)
	}
}

func TestItalic(t *testing.T) {
	italic := Italic{Text: "italian"}
	italicHtml := italic.Html()
	expected := "<em>italian</em>"
	if italicHtml != expected {
		t.Errorf("Italic wrong, got=%q", italicHtml)
	}

	italic.Properties = defaultProps(t)
	italicHtml = italic.Html()
	expected = "<em class=\"test\" style=\"background-color: red\">italian</em>"
	if italicHtml != expected {
		t.Errorf("Italic properties wrong, got=%q", italicHtml)
	}
}

func TestBlockQuote(t *testing.T) {
	blockquote := BlockQuote{Text: "quote"}
	blockquoteHtml := blockquote.Html()
	expected := "<blockquote>quote</blockquote>"
	if blockquoteHtml != expected {
		t.Errorf("Blockquote wrong, got=%q", blockquoteHtml)
	}

	blockquote.Properties = defaultProps(t)
	blockquoteHtml = blockquote.Html()
	expected = "<blockquote class=\"test\" style=\"background-color: red\">quote</blockquote>"
	if blockquoteHtml != expected {
		t.Errorf("Blockquote properties wrong, got=%q", blockquoteHtml)
	}
}

func TestListItem(t *testing.T) {
	listItem := ListItem{Component: &Paragraph{Text: "Item #1"}}
	listItemHtml := listItem.Html()
	expected := "<li><p>Item #1</p></li>"
	if listItemHtml != expected {
		t.Errorf("ListItem wrong, got=%q", listItemHtml)
	}

	listItem.Properties = defaultProps(t)
	listItemHtml = listItem.Html()
	expected = "<li class=\"test\" style=\"background-color: red\"><p>Item #1</p></li>"
	if listItemHtml != expected {
		t.Errorf("ListItem properties wrong, got=%q", listItemHtml)
	}
}

func TestOrderedList(t *testing.T) {
	listItem1 := ListItem{Component: &Paragraph{Text: "Item #1"}}
	listItem2 := ListItem{Component: &Paragraph{Text: "Item #2"}}
	listItems := []ListItem{listItem1, listItem2}
	list := OrderedList{ListItems: listItems, Start: 1}
	listHtml := list.Html()
	expected := "<ol start=\"1\">\n    <li><p>Item #1</p></li>\n    <li><p>Item #2</p></li>\n</ol>"
	if listHtml != expected {
		t.Errorf("OrderedList wrong, got=%q", listHtml)
	}

	list.Properties = defaultProps(t)
	list.Start = 5
	listHtml = list.Html()
	expected = "<ol start=\"5\" class=\"test\" style=\"background-color: red\">\n    <li><p>Item #1</p></li>\n    <li><p>Item #2</p></li>\n</ol>"
	if listHtml != expected {
		t.Errorf("OrderedList properties wrong, got=%q", listHtml)
	}
}

func TestUnorderedList(t *testing.T) {
	listItem1 := ListItem{Component: &Paragraph{Text: "Item #1"}}
	listItem2 := ListItem{Component: &Paragraph{Text: "Item #2"}}
	listItems := []ListItem{listItem1, listItem2}
	list := UnorderedList{ListItems: listItems}
	listHtml := list.Html()
	expected := "<ul>\n    <li><p>Item #1</p></li>\n    <li><p>Item #2</p></li>\n</ul>"
	if listHtml != expected {
		t.Errorf("UnorderedList wrong, got=%q", listHtml)
	}

	list.Properties = defaultProps(t)
	listHtml = list.Html()
	expected = "<ul class=\"test\" style=\"background-color: red\">\n    <li><p>Item #1</p></li>\n    <li><p>Item #2</p></li>\n</ul>"
	if listHtml != expected {
		t.Errorf("UnorderedList properties wrong, got=%q", listHtml)
	}
}

func TestImage(t *testing.T) {
	img := Image{ImgUrl: "https://img.pokemondb.net/artwork/avif/regirock.avif", AltText: "Reginald"}
	imgHtml := img.Html()
	expected := "<img src=\"https://img.pokemondb.net/artwork/avif/regirock.avif\" alt=\"Reginald\"/>"
	if imgHtml != expected {
		t.Errorf("Image wrong, got=%q", imgHtml)
	}

	img.Properties = defaultProps(t)
	imgHtml = img.Html()
	expected = "<img class=\"test\" style=\"background-color: red\" src=\"https://img.pokemondb.net/artwork/avif/regirock.avif\" alt=\"Reginald\"/>"
}

func TestHorizontalRule(t *testing.T) {
	rule := HorizontalRule{}
	ruleHtml := rule.Html()
	expected := "<hr/>"
	if ruleHtml != expected {
		t.Errorf("HorizontalRule wrong, got=%q", ruleHtml)
	}

	rule.Properties = defaultProps(t)
	ruleHtml = rule.Html()
	expected = "<hr class=\"test\" style=\"background-color: red\"/>"
	if ruleHtml != expected {
		t.Errorf("HorizontalRule properties wrong, got=%q", ruleHtml)
	}
}

func TestLink(t *testing.T) {
	link := Link{Url: "https://google.com", Text: "Google"}
	linkHtml := link.Html()
	expected := "<a href=\"https://google.com\">Google</a>"
	if linkHtml != expected {
		t.Errorf("Link wrong, got=%q", linkHtml)
	}

	link.Properties = defaultProps(t)
	linkHtml = link.Html()
	expected = "<a class=\"test\" style=\"background-color: red\" href=\"https://google.com\">Google</a>"
	if linkHtml != expected {
		t.Errorf("Link properties wrong, got=%q", linkHtml)
	}
}

func TestButton(t *testing.T) {
	button := Button{OnClick: "handleClick", Child: &Paragraph{Text: "Click Me"}}
	buttonHtml := button.Html()
	expected := "<button onclick=\"handleClick()\">\n    <p>Click Me</p>\n</button>"
	if buttonHtml != expected {
		t.Errorf("Button wrong, got=%q", buttonHtml)
	}

	button.Properties = defaultProps(t)
	buttonHtml = button.Html()
	expected = "<button class=\"test\" style=\"background-color: red\" onclick=\"handleClick()\">\n    <p>Click Me</p>\n</button>"
	if buttonHtml != expected {
		t.Errorf("Button properties wrong, got=%q", buttonHtml)
	}
}

func TestDiv(t *testing.T) {
	emptyDiv := Div{}
	divHtml := emptyDiv.Html()
	expected := "<div/>"
	if divHtml != expected {
		t.Errorf("Empty Div wrong, got=%q", divHtml)
	}

	properties := defaultProps(t)
	propertyDiv := Div{Properties: properties}
	divHtml = propertyDiv.Html()
	expected = "<div class=\"test\" style=\"background-color: red\"/>"
	if divHtml != expected {
		t.Errorf("Property Div wrong, got=%q", divHtml)
	}

	p := &Paragraph{Text: "child"}
	childDiv := Div{Children: []Component{p}}
	divHtml = childDiv.Html()
	expected = "<div>\n    <p>child</p>\n</div>"
	if divHtml != expected {
		t.Errorf("Child div wrong, got=%q", divHtml)
	}
}

func TestNav(t *testing.T) {
	nav := Nav{}
	navHtml := nav.Html()
	expected := "<nav/>"
	if navHtml != expected {
		t.Errorf("Empty nav wrong, got=%q", navHtml)
	}

	nav.Properties = defaultProps(t)
	navHtml = nav.Html()
	expected = "<nav class=\"test\" style=\"background-color: red\"/>"
	if navHtml != expected {
		t.Errorf("Nav properties wrong, got=%q", navHtml)
	}

	nav.Children = []Component{&Link{Url: "https://test.com", Text: "Test"}}
	navHtml = nav.Html()
	expected = "<nav class=\"test\" style=\"background-color: red\">\n    <a href=\"https://test.com\">Test</a>\n</nav>"
	if navHtml != expected {
		t.Errorf("Nav children wrong, got=%q", navHtml)
	}
}

func TestSpan(t *testing.T) {
	span := Span{}
	spanHtml := span.Html()
	expected := "<span/>"
	if spanHtml != expected {
		t.Errorf("Span wrong, got=%q", spanHtml)
	}

	span.Children = []Component{&Paragraph{Text: "Hello"}}
	spanHtml = span.Html()
	expected = "<span>\n    <p>Hello</p>\n</span>"
	if spanHtml != expected {
		t.Errorf("Span children wrong, got=%q", spanHtml)
	}
}

func TestCodeBlock(t *testing.T) {
	content := `package main

import "fmt"

func main() {
    fmt.Println("Hello, world!")
}`
	codeBlock := CodeBlock{Content: content}
	codeBlockHtml := codeBlock.Html()
	expected := `<div class="code-block">
    <pre>package main</pre>
    <pre></pre>
    <pre>import "fmt"</pre>
    <pre></pre>
    <pre>func main() {</pre>
    <pre>    fmt.Println("Hello, world!")</pre>
    <pre>}</pre>
</div>`

	if codeBlockHtml != expected {
		t.Errorf("CodeBlock wrong\ngot=     %q\nexpected=%q", codeBlockHtml, expected)
	}
}