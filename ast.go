package mdx

import (
	"fmt"
	"slices"
	"strings"
)

type ComponentType int

const (
	Block = iota
	Inline
)

type component interface {
	// Convert element to HTML representation
	Html() string
	Type() ComponentType

	// Format method rendering html but indented appropriately
	// Hoping for this to replace the Html method and make everything cleaner
	Format(int) string
}

type property struct {
	Name  string
	Value string
}

type fragment struct {
	Value string
}

func (f *fragment) String() string {
	return fmt.Sprintf("Fragment{%s}", f.Value)
}

func (f *fragment) Html() string {
	return fmt.Sprintf("%s", f.Value)
}

func (f *fragment) Type() ComponentType {
	return Inline
}

func (f *fragment) Format(indentLevel int) string {
	indentPrefix := strings.Repeat(INDENT, indentLevel)
	formattedOutput := indentPrefix + f.Value
	return formattedOutput
}

type lineBreak struct{}

func (lb *lineBreak) Html() string {
	return "<br/>"
}

func (lb *lineBreak) Type() ComponentType {
	return Block
}

func (lb *lineBreak) Format(indentLevel int) string {
	tag := "<br/>"
	indentPrefix := strings.Repeat(INDENT, indentLevel)
	formattedOutput := "\n" + indentPrefix + tag + "\n"
	return formattedOutput
}

type header struct {
	Properties []property
	Level      int
	Content    []component
}

func (h *header) InnerHtml() string {
	var contentString string
	for _, child := range h.Content {
		contentString += child.Html()
	}
	return contentString
}

func (h *header) Html() string {
	if h.Level == 0 {
		h.Level = 1
	}

	var propertyString string
	for _, property := range h.Properties {
		propertyString += fmt.Sprintf(" %s=\"%s\"", property.Name, property.Value)
	}

	return fmt.Sprintf("<h%d%s>%s</h%d>", h.Level, propertyString, h.InnerHtml(), h.Level)
}

func (h *header) Type() ComponentType {
	return Block
}

func (h *header) Format(indentLevel int) string {
	var propertyString string
	for _, property := range h.Properties {
		propertyString += fmt.Sprintf(" %s=\"%s\"", property.Name, property.Value)
	}

	openingTag := fmt.Sprintf("<h%d%s>", h.Level, propertyString)
	closingTag := fmt.Sprintf("</h%d>", h.Level)
	indentPrefix := strings.Repeat(INDENT, indentLevel)

	formattedOutput := "\n" + indentPrefix + openingTag

	containsBlockElement := slices.ContainsFunc(h.Content, func(c component) bool {
		return c.Type() == Block
	})

	if containsBlockElement {
		// put child component on new line and indented + 1
		if len(h.Content) > 0 && h.Content[0].Type() == Inline {
			formattedOutput += "\n"
		}

		var inlineString string
		for _, child := range h.Content {
			if child.Type() == Inline {
				inlineString += indentPrefix + INDENT + child.Html()
			} else {
				if len(inlineString) > 0 {
					// split the inline string and append each
					lineLength := len(indentPrefix) + len(openingTag) + len(inlineString) + len(closingTag)
					appendInlineString(inlineString, indentPrefix, lineLength, &formattedOutput)
					inlineString = ""
				}

				formattedOutput += child.Format(indentLevel + 1)
			}
		}
		if len(inlineString) > 0 {
			// split the inline string and append each
			lineLength := len(indentPrefix) + len(openingTag) + len(inlineString) + len(closingTag)
			appendInlineString(inlineString, indentPrefix, lineLength, &formattedOutput)
			inlineString = ""
		}
		formattedOutput += "\n" + indentPrefix
	} else {
		// check if everything can fit on one line. if not, figure it out
		inlineString := h.InnerHtml()
		lineLength := len(indentPrefix) + len(openingTag) + len(inlineString) + len(closingTag)
		appendInlineString(inlineString, indentPrefix, lineLength, &formattedOutput)
	}

	formattedOutput += closingTag + "\n"

	return formattedOutput
}

type paragraph struct {
	Properties []property
	Content    []component
}

func (p *paragraph) String() string {
	var contentString string
	for _, child := range p.Content {
		contentString += fmt.Sprintf("%s ", child)
	}
	return fmt.Sprintf("Paragraph{Content=[%s]}", strings.TrimSpace(contentString))
}

func (p *paragraph) InnerHtml() string {
	var contentString string
	for _, child := range p.Content {
		contentString += child.Html()
	}
	return contentString
}

func (p *paragraph) Html() string {
	var propertyString string
	for _, property := range p.Properties {
		propertyString += fmt.Sprintf(" %s=\"%s\"", property.Name, property.Value)
	}

	return fmt.Sprintf("<p%s>%s</p>", propertyString, p.InnerHtml())
}

func (p *paragraph) Type() ComponentType {
	return Inline
}

func (p *paragraph) Format(indentLevel int) string {
	var propertyString string
	for _, property := range p.Properties {
		propertyString += fmt.Sprintf(" %s=\"%s\"", property.Name, property.Value)
	}

	openingTag := fmt.Sprintf("<p%s>", propertyString)
	closingTag := "</p>"

	indentPrefix := strings.Repeat(INDENT, indentLevel)
	formattedOutput := openingTag

	containsBlockElement := slices.ContainsFunc(p.Content, func(c component) bool {
		return c.Type() == Block
	})

	if containsBlockElement {
		// put child component on new line and indented + 1
		if len(p.Content) > 0 && p.Content[0].Type() == Inline {
			formattedOutput += "\n"
		}

		var inlineString string
		for _, child := range p.Content {
			if child.Type() == Inline {
				inlineString += indentPrefix + INDENT + child.Html()
			} else {
				if len(inlineString) > 0 {
					// split the inline string and append each
					lineLength := len(indentPrefix) + len(openingTag) + len(inlineString) + len(closingTag)
					appendInlineString(inlineString, indentPrefix, lineLength, &formattedOutput)
					inlineString = ""
				}

				formattedOutput += child.Format(indentLevel + 1)
			}
		}
		if len(inlineString) > 0 {
			// split the inline string and append each
			lineLength := len(indentPrefix) + len(openingTag) + len(inlineString) + len(closingTag)
			appendInlineString(inlineString, indentPrefix, lineLength, &formattedOutput)
			inlineString = ""
		} else {
			formattedOutput += indentPrefix
		}
	} else {
		// check if everything can fit on one line. if not, figure it out
		inlineString := p.InnerHtml()
		lineLength := len(indentPrefix) + len(openingTag) + len(inlineString) + len(closingTag)
		appendInlineString(inlineString, indentPrefix, lineLength, &formattedOutput)
	}

	formattedOutput += closingTag + "\n"

	return formattedOutput
}

type code struct {
	Properties []property
	Text       string
}

func (c *code) Html() string {
	var propertyString string
	for _, property := range c.Properties {
		propertyString += fmt.Sprintf(" %s=\"%s\"", property.Name, property.Value)
	}
	return fmt.Sprintf("<code%s>%s</code>", propertyString, c.Text)
}

func (c *code) Type() ComponentType {
	return Inline
}

func (c *code) Format(indentLevel int) string {
	var propertyString string
	for _, property := range c.Properties {
		propertyString += fmt.Sprintf(" %s=\"%s\"", property.Name, property.Value)
	}

	formattedOutput := fmt.Sprintf("<code%s>%s</code", propertyString, c.Text)

	return formattedOutput
}

type bold struct {
	Properties []property
	Content    []component
}

func (b *bold) InnerHtml() string {
	var contentString string
	for _, child := range b.Content {
		contentString += child.Html()
	}
	return contentString
}

func (b *bold) Html() string {
	var propertyString string
	for _, property := range b.Properties {
		propertyString += fmt.Sprintf(" %s=\"%s\"", property.Name, property.Value)
	}
	return fmt.Sprintf("<strong%s>%s</strong>", propertyString, b.InnerHtml())
}

func (b *bold) Type() ComponentType {
	return Inline
}

func (b *bold) Format(indentLevel int) string {
	var propertyString string
	for _, property := range b.Properties {
		propertyString += fmt.Sprintf(" %s=\"%s\"", property.Name, property.Value)
	}

	openingTag := fmt.Sprintf("<strong%s>", propertyString)
	closingTag := "</strong>"

	indentPrefix := strings.Repeat(INDENT, indentLevel)
	formattedOutput := openingTag

	containsBlockElement := slices.ContainsFunc(b.Content, func(c component) bool {
		return c.Type() == Block
	})

	if containsBlockElement {
		// put child component on new line and indented + 1
		if len(b.Content) > 0 && b.Content[0].Type() == Inline {
			formattedOutput += "\n"
		}

		var inlineString string
		for _, child := range b.Content {
			if child.Type() == Inline {
				inlineString += indentPrefix + INDENT + child.Html()
			} else {
				if len(inlineString) > 0 {
					// split the inline string and append each
					lineLength := len(indentPrefix) + len(openingTag) + len(inlineString) + len(closingTag)
					appendInlineString(inlineString, indentPrefix, lineLength, &formattedOutput)
					inlineString = ""
				}

				formattedOutput += child.Format(indentLevel + 1)
			}
		}
		if len(inlineString) > 0 {
			// split the inline string and append each
			lineLength := len(indentPrefix) + len(openingTag) + len(inlineString) + len(closingTag)
			appendInlineString(inlineString, indentPrefix, lineLength, &formattedOutput)
			inlineString = ""
		}
	} else {
		// check if everything can fit on one line. if not, figure it out
		inlineString := b.InnerHtml()
		lineLength := len(indentPrefix) + len(openingTag) + len(inlineString) + len(closingTag)
		appendInlineString(inlineString, indentPrefix, lineLength, &formattedOutput)
	}

	formattedOutput += closingTag + "\n"

	return formattedOutput
}

type italic struct {
	Properties []property
	Content    []component
}

func (i *italic) InnerHtml() string {
	var contentString string
	for _, child := range i.Content {
		contentString += child.Html()
	}
	return contentString
}

func (i *italic) Html() string {
	var propertyString string
	for _, property := range i.Properties {
		propertyString += fmt.Sprintf(" %s=\"%s\"", property.Name, property.Value)
	}
	return fmt.Sprintf("<em%s>%s</em>", propertyString, i.InnerHtml())
}

func (i *italic) Type() ComponentType {
	return Inline
}

func (i *italic) Format(indentLevel int) string {
	var propertyString string
	for _, property := range i.Properties {
		propertyString += fmt.Sprintf(" %s=\"%s\"", property.Name, property.Value)
	}

	openingTag := fmt.Sprintf("<em%s>", propertyString)
	closingTag := "</em>"

	indentPrefix := strings.Repeat(INDENT, indentLevel)
	formattedOutput := openingTag

	containsBlockElement := slices.ContainsFunc(i.Content, func(c component) bool {
		return c.Type() == Block
	})

	if containsBlockElement {
		// put child component on new line and indented + 1
		if len(i.Content) > 0 && i.Content[0].Type() == Inline {
			formattedOutput += "\n"
		}

		var inlineString string
		for _, child := range i.Content {
			if child.Type() == Inline {
				inlineString += indentPrefix + INDENT + child.Html()
			} else {
				if len(inlineString) > 0 {
					// split the inline string and append each
					lineLength := len(indentPrefix) + len(openingTag) + len(inlineString) + len(closingTag)
					appendInlineString(inlineString, indentPrefix, lineLength, &formattedOutput)
					inlineString = ""
				}

				formattedOutput += child.Format(indentLevel + 1)
			}
		}
		if len(inlineString) > 0 {
			// split the inline string and append each
			lineLength := len(indentPrefix) + len(openingTag) + len(inlineString) + len(closingTag)
			appendInlineString(inlineString, indentPrefix, lineLength, &formattedOutput)
			inlineString = ""
		}
	} else {
		// check if everything can fit on one line. if not, figure it out
		inlineString := i.InnerHtml()
		lineLength := len(indentPrefix) + len(openingTag) + len(inlineString) + len(closingTag)
		appendInlineString(inlineString, indentPrefix, lineLength, &formattedOutput)
	}

	formattedOutput += closingTag + "\n"

	return formattedOutput
}

type blockQuote struct {
	Properties []property
	Content    []component
}

func (bq *blockQuote) InnerHtml() string {
	var contentString string
	for _, child := range bq.Content {
		contentString += child.Html()
	}
	return contentString
}

func (bq *blockQuote) Html() string {
	var propertyString string
	for _, property := range bq.Properties {
		propertyString += fmt.Sprintf(" %s=\"%s\"", property.Name, property.Value)
	}
	return fmt.Sprintf("<blockquote%s>%s</blockquote>", propertyString, bq.InnerHtml())
}

func (bq *blockQuote) Type() ComponentType {
	return Block
}

func (bq *blockQuote) Format(indentLevel int) string {
	var propertyString string
	for _, property := range bq.Properties {
		propertyString += fmt.Sprintf(" %s=\"%s\"", property.Name, property.Value)
	}

	openingTag := fmt.Sprintf("<blockquote%s>", propertyString)
	closingTag := "</blockquote>"
	indentPrefix := strings.Repeat(INDENT, indentLevel)

	formattedOutput := "\n" + indentPrefix + openingTag

	containsBlockElement := slices.ContainsFunc(bq.Content, func(c component) bool {
		return c.Type() == Block
	})

	if containsBlockElement {
		// put child component on new line and indented + 1
		if len(bq.Content) > 0 && bq.Content[0].Type() == Inline {
			formattedOutput += "\n"
		}

		var inlineString string
		for _, child := range bq.Content {
			if child.Type() == Inline {
				inlineString += indentPrefix + INDENT + child.Html()
			} else {
				if len(inlineString) > 0 {
					// split the inline string and append each
					lineLength := len(indentPrefix) + len(openingTag) + len(inlineString) + len(closingTag)
					appendInlineString(inlineString, indentPrefix, lineLength, &formattedOutput)
					inlineString = ""
				}

				formattedOutput += child.Format(indentLevel + 1)
			}
		}
		if len(inlineString) > 0 {
			// split the inline string and append each
			lineLength := len(indentPrefix) + len(openingTag) + len(inlineString) + len(closingTag)
			appendInlineString(inlineString, indentPrefix, lineLength, &formattedOutput)
			inlineString = ""
		}
		formattedOutput += "\n" + indentPrefix
	} else {
		// check if everything can fit on one line. if not, figure it out
		inlineString := bq.InnerHtml()
		lineLength := len(indentPrefix) + len(openingTag) + len(inlineString) + len(closingTag)
		appendInlineString(inlineString, indentPrefix, lineLength, &formattedOutput)
	}

	formattedOutput += closingTag + "\n"

	return formattedOutput
}

type listItem struct {
	Properties []property
	Component  component
}

func (li *listItem) Html() string {
	var propertyString string
	for _, property := range li.Properties {
		propertyString += fmt.Sprintf(" %s=\"%s\"", property.Name, property.Value)
	}
	return fmt.Sprintf("<li%s>%s</li>", propertyString, li.Component.Html())
}

func (li *listItem) Type() ComponentType {
	return Block
}

func (li *listItem) Format(indentLevel int) string {
	var propertyString string
	for _, property := range li.Properties {
		propertyString += fmt.Sprintf(" %s=\"%s\"", property.Name, property.Value)
	}

	openingTag := fmt.Sprintf("<li%s>\n", propertyString)
	closingTag := "</li>"
	indentPrefix := strings.Repeat(INDENT, indentLevel)

	formattedOutput := indentPrefix + openingTag

	if li.Component.Type() == Block {
		formattedOutput += li.Component.Format(indentLevel + 1)
		formattedOutput += indentPrefix
	} else {
		inlineString := strings.Repeat(INDENT, indentLevel+1) + li.Component.Format(indentLevel+1)
		lineLength := len(indentPrefix) + len(openingTag) + len(inlineString) + len(closingTag)
		appendInlineString(inlineString, indentPrefix, lineLength, &formattedOutput)
	}

	formattedOutput += "\n" + indentPrefix + closingTag + "\n"

	return formattedOutput
}

type orderedList struct {
	Properties []property
	ListItems  []listItem
	Start      int
}

func (ol *orderedList) Html() string {
	var propertyString string
	for _, property := range ol.Properties {
		propertyString += fmt.Sprintf(" %s=\"%s\"", property.Name, property.Value)
	}

	var listItemString string
	for _, item := range ol.ListItems {
		listItemString += fmt.Sprintf("    %s\n", item.Html())
	}

	return fmt.Sprintf("<ol start=\"%d\"%s>\n%s</ol>", ol.Start, propertyString, listItemString)
}

func (ol *orderedList) Type() ComponentType {
	return Block
}

func (ol *orderedList) Format(indentLevel int) string {
	var propertyString string
	for _, property := range ol.Properties {
		propertyString += fmt.Sprintf(" %s=\"%s\"", property.Name, property.Value)
	}

	openingTag := fmt.Sprintf("<ol%s>", propertyString)
	closingTag := "</ol>"
	indentPrefix := strings.Repeat(INDENT, indentLevel)

	formattedOutput := "\n" + indentPrefix + openingTag + "\n"
	for _, item := range ol.ListItems {
		formattedOutput += item.Format(indentLevel + 1)
	}

	formattedOutput += indentPrefix + closingTag + "\n"

	return formattedOutput
}

type unorderedList struct {
	Properties []property
	ListItems  []listItem
}

func (ul *unorderedList) Html() string {
	var propertyString string
	for _, property := range ul.Properties {
		propertyString += fmt.Sprintf(" %s=\"%s\"", property.Name, property.Value)
	}

	var listItemString string
	for _, item := range ul.ListItems {
		listItemString += fmt.Sprintf("    %s\n", item.Html())
	}

	return fmt.Sprintf("<ul%s>\n%s</ul>", propertyString, listItemString)
}

func (ul *unorderedList) Type() ComponentType {
	return Block
}

func (ul *unorderedList) Format(indentLevel int) string {
	var propertyString string
	for _, property := range ul.Properties {
		propertyString += fmt.Sprintf(" %s=\"%s\"", property.Name, property.Value)
	}

	openingTag := fmt.Sprintf("<ul%s>", propertyString)
	closingTag := "</ul>"
	indentPrefix := strings.Repeat(INDENT, indentLevel)

	formattedOutput := "\n" + indentPrefix + openingTag + "\n"
	for _, item := range ul.ListItems {
		formattedOutput += item.Format(indentLevel + 1)
	}

	formattedOutput += indentPrefix + closingTag + "\n"

	return formattedOutput
}

type image struct {
	Properties []property
	ImgUrl     string
	AltText    string
}

func (img *image) Html() string {
	var propertyString string
	for _, property := range img.Properties {
		propertyString += fmt.Sprintf(" %s=\"%s\"", property.Name, property.Value)
	}
	return fmt.Sprintf("<img%s src=\"%s\" alt=\"%s\"/>", propertyString, img.ImgUrl, img.AltText)
}

func (img *image) Type() ComponentType {
	return Block
}

func (img *image) Format(indentLevel int) string {
	var propertyString string
	for _, property := range img.Properties {
		propertyString += fmt.Sprintf(" %s=\"%s\"", property.Name, property.Value)
	}

	indentPrefix := strings.Repeat(INDENT, indentLevel)
	tag := fmt.Sprintf("<img%s src=\"%s\" alt=\"%s\"/>", propertyString, img.ImgUrl, img.AltText)
	formattedOutput := indentPrefix + tag + "\n"

	return formattedOutput
}

type horizontalRule struct {
	Properties []property
}

func (hr *horizontalRule) Html() string {
	var propertyString string
	for _, property := range hr.Properties {
		propertyString += fmt.Sprintf(" %s=\"%s\"", property.Name, property.Value)
	}
	return fmt.Sprintf("<hr%s/>", propertyString)
}

func (hr *horizontalRule) Type() ComponentType {
	return Block
}

func (hr *horizontalRule) Format(indentLevel int) string {
	var propertyString string
	for _, property := range hr.Properties {
		propertyString += fmt.Sprintf(" %s=\"%s\"", property.Name, property.Value)
	}

	indentPrefix := strings.Repeat(INDENT, indentLevel)
	tag := fmt.Sprintf("<hr%s/>", propertyString)
	formattedOutput := indentPrefix + tag + "\n"

	return formattedOutput
}

type link struct {
	Properties []property
	Url        string
	Content    []component
}

func (l *link) InnerHtml() string {
	var contentString string
	for _, child := range l.Content {
		contentString += child.Html()
	}
	return contentString
}

func (l *link) Html() string {
	var propertyString string
	for _, property := range l.Properties {
		propertyString += fmt.Sprintf(" %s=\"%s\"", property.Name, property.Value)
	}
	return fmt.Sprintf("<a%s href=\"%s\" target=_blank>%s</a>", propertyString, l.Url, l.InnerHtml())
}

func (l *link) Type() ComponentType {
	return Inline
}

func (l *link) Format(indentLevel int) string {
	var propertyString string
	for _, property := range l.Properties {
		propertyString += fmt.Sprintf(" %s=\"%s\"", property.Name, property.Value)
	}

	openingTag := fmt.Sprintf("<a%s href=\"%s\" target=_blank>", propertyString, l.Url)
	closingTag := "</a>"
	indentPrefix := strings.Repeat(INDENT, indentLevel)

	formattedOutput := openingTag
	containsBlockElement := slices.ContainsFunc(l.Content, func(c component) bool {
		return c.Type() == Block
	})

	if containsBlockElement {
		// put child component on new line and indented + 1
		if len(l.Content) > 0 && l.Content[0].Type() == Inline {
			formattedOutput += "\n"
		}

		var inlineString string
		for _, child := range l.Content {
			if child.Type() == Inline {
				inlineString += indentPrefix + INDENT + child.Html()
			} else {
				if len(inlineString) > 0 {
					// split the inline string and append each
					lineLength := len(indentPrefix) + len(openingTag) + len(inlineString) + len(closingTag)
					appendInlineString(inlineString, indentPrefix, lineLength, &formattedOutput)
					inlineString = ""
				}

				formattedOutput += child.Format(indentLevel + 1)
			}
		}
		if len(inlineString) > 0 {
			// split the inline string and append each
			lineLength := len(indentPrefix) + len(openingTag) + len(inlineString) + len(closingTag)
			appendInlineString(inlineString, indentPrefix, lineLength, &formattedOutput)
			inlineString = ""
		}
	} else {
		// check if everything can fit on one line. if not, figure it out
		inlineString := l.InnerHtml()
		lineLength := len(indentPrefix) + len(openingTag) + len(inlineString) + len(closingTag)
		appendInlineString(inlineString, indentPrefix, lineLength, &formattedOutput)
	}

	formattedOutput += closingTag + "\n"

	return formattedOutput
}

type button struct {
	Properties []property
	Content    []component
	OnClick    string
}

func (b *button) String() string {
	var contentString string
	for _, child := range b.Content {
		contentString += fmt.Sprintf("%s ", child)
	}
	return fmt.Sprintf("Button{OnClick='%s', Content=[%s]}", b.OnClick, strings.TrimSpace(contentString))
}

func (b *button) InnerHtml() string {
	var contentString string
	for _, child := range b.Content {
		contentString += child.Html()
	}
	return contentString
}

func (b *button) Html() string {
	var propertyString string
	for _, property := range b.Properties {
		propertyString += fmt.Sprintf(" %s=\"%s\"", property.Name, property.Value)
	}
	return fmt.Sprintf("<button%s onclick=\"%s()\">\n    %s\n</button>", propertyString, b.OnClick, b.InnerHtml())
}

func (b *button) Type() ComponentType {
	return Block
}

func (b *button) Format(indentLevel int) string {
	var propertyString string
	for _, property := range b.Properties {
		propertyString += fmt.Sprintf(" %s=\"%s\"", property.Name, property.Value)
	}

	var formattedOutput = "\n"
	var indentPrefix string
	for range indentLevel {
		indentPrefix += INDENT
	}

	openingTag := fmt.Sprintf("<button%s onclick=\"%s()\">", propertyString, b.OnClick)
	closingTag := "</button>"

	formattedOutput += indentPrefix + openingTag

	containsBlockElement := slices.ContainsFunc(b.Content, func(c component) bool {
		return c.Type() == Block
	})

	if containsBlockElement {
		// put child component on new line and indented + 1
		if len(b.Content) > 0 && b.Content[0].Type() == Inline {
			formattedOutput += "\n"
		}

		var inlineString string
		for _, child := range b.Content {
			if child.Type() == Inline {
				inlineString += indentPrefix + INDENT + child.Html()
			} else {
				if len(inlineString) > 0 {
					// split the inline string and append each
					lineLength := len(indentPrefix) + len(openingTag) + len(inlineString) + len(closingTag)
					appendInlineString(inlineString, indentPrefix, lineLength, &formattedOutput)
					inlineString = ""
				}

				formattedOutput += child.Format(indentLevel + 1)
			}
		}
		if len(inlineString) > 0 {
			// split the inline string and append each
			lineLength := len(indentPrefix) + len(openingTag) + len(inlineString) + len(closingTag)
			appendInlineString(inlineString, indentPrefix, lineLength, &formattedOutput)
			inlineString = ""
		}
		formattedOutput += "\n" + indentPrefix
	} else {
		// check if everything can fit on one line. if not, figure it out
		inlineString := b.InnerHtml()
		lineLength := len(indentPrefix) + len(openingTag) + len(inlineString) + len(closingTag)
		appendInlineString(inlineString, indentPrefix, lineLength, &formattedOutput)
	}

	formattedOutput += closingTag + "\n"

	return formattedOutput
}

type div struct {
	Properties []property
	Children   []component
}

func (d *div) String() string {
	var contentString string
	for _, child := range d.Children {
		contentString += fmt.Sprintf("%s ", child)
	}
	return fmt.Sprintf("Div{Children=[%s]}", strings.TrimSpace(contentString))
}

func (d *div) Html() string {
	var divString string
	var propertyString string
	for _, property := range d.Properties {
		propertyString += fmt.Sprintf(" %s=\"%s\"", property.Name, property.Value)
	}

	if len(d.Children) == 0 {
		return fmt.Sprintf("<div%s/>", propertyString)
	}

	divString += fmt.Sprintf("<div%s>\n", propertyString)
	for _, child := range d.Children {
		divString += fmt.Sprintf("    %s\n", child.Html())
	}
	divString += "</div>"

	return divString
}

func (d *div) Type() ComponentType {
	return Block
}

func (d *div) Format(indentLevel int) string {
	var propertyString string
	for _, property := range d.Properties {
		propertyString += fmt.Sprintf(" %s=\"%s\"", property.Name, property.Value)
	}

	openingTag := fmt.Sprintf("<div%s>", propertyString)
	closingTag := "</div>"
	indentPrefix := strings.Repeat(INDENT, indentLevel)

	formattedOutput := "\n" + indentPrefix + openingTag + "\n"

	for _, child := range d.Children {
		if child.Type() == Inline {
			formattedOutput += strings.Repeat(INDENT, indentLevel+1)
		}
		formattedOutput += child.Format(indentLevel + 1)
	}

	formattedOutput += "\n" + indentPrefix + closingTag + "\n"
	return formattedOutput
}

type nav struct {
	Properties []property
	Children   []component
}

func (n *nav) Html() string {
	var navString string
	var propertyString string
	for _, property := range n.Properties {
		propertyString += fmt.Sprintf(" %s=\"%s\"", property.Name, property.Value)
	}

	if len(n.Children) == 0 {
		return fmt.Sprintf("<nav%s/>", propertyString)
	}

	navString += fmt.Sprintf("<nav%s>\n", propertyString)
	for _, child := range n.Children {
		navString += fmt.Sprintf("    %s\n", child.Html())
	}
	navString += "</nav>"

	return navString
}

func (n *nav) Type() ComponentType {
	return Block
}

func (n *nav) Format(indentLevel int) string {
	var propertyString string
	for _, property := range n.Properties {
		propertyString += fmt.Sprintf(" %s=\"%s\"", property.Name, property.Value)
	}

	openingTag := fmt.Sprintf("<nav%s>", propertyString)
	closingTag := "</nav>"
	indentPrefix := strings.Repeat(INDENT, indentLevel)

	formattedOutput := "\n" + indentPrefix + openingTag + "\n"

	for _, child := range n.Children {
		if child.Type() == Inline {
			formattedOutput += strings.Repeat(INDENT, indentLevel+1)
		}
		formattedOutput += child.Format(indentLevel + 1)
	}

	formattedOutput += "\n" + indentPrefix + closingTag + "\n"
	return formattedOutput
}

type span struct {
	Properties []property
	Content    []component
}

func (s *span) InnerHtml() string {
	var contentString string
	for _, child := range s.Content {
		contentString += child.Html()
	}
	return contentString
}

func (s *span) Html() string {
	var propertyString string
	for _, property := range s.Properties {
		propertyString += fmt.Sprintf(" %s=\"%s\"", property.Name, property.Value)
	}

	if len(s.Content) == 0 {
		return fmt.Sprintf("<span%s/>", propertyString)
	}

	return fmt.Sprintf("<span%s>%s</span>", propertyString, s.InnerHtml())
}

func (s *span) Type() ComponentType {
	return Inline
}

func (s *span) Format(indentLevel int) string {
	var propertyString string
	for _, property := range s.Properties {
		propertyString += fmt.Sprintf(" %s=\"%s\"", property.Name, property.Value)
	}

	openingTag := fmt.Sprintf("<span%s>", propertyString)
	closingTag := "</span>"

	indentPrefix := strings.Repeat(INDENT, indentLevel)
	formattedOutput := openingTag

	containsBlockElement := slices.ContainsFunc(s.Content, func(c component) bool {
		return c.Type() == Block
	})

	if containsBlockElement {
		// put child component on new line and indented + 1
		if len(s.Content) > 0 && s.Content[0].Type() == Inline {
			formattedOutput += "\n"
		}

		var inlineString string
		for _, child := range s.Content {
			if child.Type() == Inline {
				inlineString += indentPrefix + INDENT + child.Html()
			} else {
				if len(inlineString) > 0 {
					// split the inline string and append each
					lineLength := len(indentPrefix) + len(openingTag) + len(inlineString) + len(closingTag)
					appendInlineString(inlineString, indentPrefix, lineLength, &formattedOutput)
					inlineString = ""
				}

				formattedOutput += child.Format(indentLevel + 1)
			}
		}
		if len(inlineString) > 0 {
			// split the inline string and append each
			lineLength := len(indentPrefix) + len(openingTag) + len(inlineString) + len(closingTag)
			appendInlineString(inlineString, indentPrefix, lineLength, &formattedOutput)
			inlineString = ""
		}
	} else {
		// check if everything can fit on one line. if not, figure it out
		inlineString := s.InnerHtml()
		lineLength := len(indentPrefix) + len(openingTag) + len(inlineString) + len(closingTag)
		appendInlineString(inlineString, indentPrefix, lineLength, &formattedOutput)
	}

	formattedOutput += closingTag + "\n"

	return formattedOutput
}

type codeBlock struct {
	Properties []property
	Content    string
}

func (cb *codeBlock) Html() string {
	var codeBlockString string
	var propertiesString string
	for _, property := range cb.Properties {
		propertiesString += fmt.Sprintf(" %s=\"%s\"", property.Name, property.Value)
	}

	lines := strings.Split(cb.Content, "\\n")
	codeBlockString += fmt.Sprintf("<div class=\"code-block\"%s>\n", propertiesString)
	for _, line := range lines {
		codeBlockString += fmt.Sprintf("    <pre>%s</pre>\n", line)
	}
	codeBlockString += "</div>"

	return codeBlockString
}

func (cb *codeBlock) Type() ComponentType {
	return Block
}

func (cb *codeBlock) Format(indentLevel int) string {
	var propertyString string
	for _, property := range cb.Properties {
		propertyString += fmt.Sprintf(" %s=\"%s\"", property.Name, property.Value)
	}

	openingTag := fmt.Sprintf("<div class=\"code-block\"%s>", propertyString)
	closingTag := "</div>"
	indentPrefix := strings.Repeat(INDENT, indentLevel)

	formattedOutput := "\n" + indentPrefix + openingTag + "\n"
	lines := strings.Split(cb.Content, "\\n")
	for _, line := range lines {
		formattedOutput += indentPrefix + INDENT
		formattedOutput += fmt.Sprintf("<pre>%s</pre>\n", line)
	}

	formattedOutput += indentPrefix + closingTag + "\n"

	return formattedOutput
}

type body struct {
	Children []component
}

func (b *body) Html() string {
	bodyString := "<body>\n"
	for _, child := range b.Children {
		bodyString += fmt.Sprintf("    %s\n", child.Html())
	}
	bodyString += "</body>"
	return bodyString
}

func (b *body) Type() ComponentType {
	return Block
}

func (b *body) Format(indentLevel int) string {
	formattedOutput := "\n"
	var indentPrefix string
	for range indentLevel {
		indentPrefix += INDENT
	}
	formattedOutput += indentPrefix + "<body>"

	if len(b.Children) > 0 && b.Children[0].Type() == Inline {
		formattedOutput += "\n"
	}

	for _, child := range b.Children {
		if child.Type() == Inline {
			formattedOutput += strings.Repeat(INDENT, indentLevel+1)
		}

		formattedOutput += child.Format(indentLevel + 1)
	}

	if len(b.Children) > 0 && b.Children[len(b.Children)-1].Type() == Inline {
		formattedOutput += "\n"
	}

	formattedOutput += indentPrefix + "</body>"

	return formattedOutput
}

func appendInlineString(inlineString string, indentPrefix string, lineLength int, formattedOutput *string) {
	if lineLength < MAX_LENGTH {
		*formattedOutput += inlineString
	} else {
		*formattedOutput += "\n"
		indentString := indentPrefix + INDENT

		words := strings.Split(inlineString, " ")
		currentLine := words[0]
		words = words[1:]

		for len(words) > 0 {
			canAdd := MAX_LENGTH - len(currentLine) - len(indentString)

			// this is not what I originally intended, but I think it might look better like this?
			// to make it as originally intended, remove `currentLine` from below length check
			if len(currentLine+words[0]) > canAdd {
				*formattedOutput += indentString + currentLine + "\n"
				currentLine = ""
			}

			if len(currentLine) > 0 {
				currentLine += " "
			}
			currentLine += words[0]
			words = words[1:]
		}

		if len(currentLine) > 0 {
			*formattedOutput += indentString + currentLine + "\n"
		}

		*formattedOutput += indentPrefix
	}
}
