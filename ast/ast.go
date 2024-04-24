package ast

import (
	"fmt"
	"strings"
)

type Component interface {
	Html() string
}

type Property struct {
	Name  string
	Value string
}

type Fragment struct {
	String string
}

func (f *Fragment) Html() string {
	return fmt.Sprintf("%s", f.String)
}

type LineBreak struct{}

func (lb *LineBreak) Html() string {
	return "<br/>"
}

type Header struct {
	Properties []Property
	Level      int
	Text       string
}

func (h *Header) Html() string {
	if h.Level == 0 {
		h.Level = 1
	}

	var propertyString string
	for _, property := range h.Properties {
		propertyString += fmt.Sprintf(" %s=\"%s\"", property.Name, property.Value)
	}
	return fmt.Sprintf("<h%d%s>%s</h%d>", h.Level, propertyString, h.Text, h.Level)
}

type Paragraph struct {
	Properties []Property
	Text       string
}

func (p *Paragraph) Html() string {
	var propertyString string
	for _, property := range p.Properties {
		propertyString += fmt.Sprintf(" %s=\"%s\"", property.Name, property.Value)
	}
	return fmt.Sprintf("<p%s>%s</p>", propertyString, p.Text)
}

type Code struct {
	Properties []Property
	Text       string
}

func (c *Code) Html() string {
	var propertyString string
	for _, property := range c.Properties {
		propertyString += fmt.Sprintf(" %s=\"%s\"", property.Name, property.Value)
	}
	return fmt.Sprintf("<code%s>%s</code>", propertyString, c.Text)
}

type Bold struct {
	Properties []Property
	Text       string
}

func (b *Bold) Html() string {
	var propertyString string
	for _, property := range b.Properties {
		propertyString += fmt.Sprintf(" %s=\"%s\"", property.Name, property.Value)
	}
	return fmt.Sprintf("<strong%s>%s</strong>", propertyString, b.Text)
}

type Italic struct {
	Properties []Property
	Text       string
}

func (i *Italic) Html() string {
	var propertyString string
	for _, property := range i.Properties {
		propertyString += fmt.Sprintf(" %s=\"%s\"", property.Name, property.Value)
	}
	return fmt.Sprintf("<em%s>%s</em>", propertyString, i.Text)
}

type BlockQuote struct {
	Properties []Property
	Text       string
}

func (bq *BlockQuote) Html() string {
	var propertyString string
	for _, property := range bq.Properties {
		propertyString += fmt.Sprintf(" %s=\"%s\"", property.Name, property.Value)
	}
	return fmt.Sprintf("<blockquote%s>%s</blockquote>", propertyString, bq.Text)
}

type ListItem struct {
	Properties []Property
	Component  Component
}

func (li *ListItem) Html() string {
	var propertyString string
	for _, property := range li.Properties {
		propertyString += fmt.Sprintf(" %s=\"%s\"", property.Name, property.Value)
	}
	return fmt.Sprintf("<li%s>%s</li>", propertyString, li.Component.Html())
}

type OrderedList struct {
	Properties []Property
	ListItems  []ListItem
	Start      int
}

func (ol *OrderedList) Html() string {
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

type UnorderedList struct {
	Properties []Property
	ListItems  []ListItem
}

func (ul *UnorderedList) Html() string {
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

type Image struct {
	Properties []Property
	ImgUrl     string
	AltText    string
}

func (img *Image) Html() string {
	var propertyString string
	for _, property := range img.Properties {
		propertyString += fmt.Sprintf(" %s=\"%s\"", property.Name, property.Value)
	}
	return fmt.Sprintf("<img%s src=\"%s\" alt=\"%s\"/>", propertyString, img.ImgUrl, img.AltText)
}

type HorizontalRule struct {
	Properties []Property
}

func (hr *HorizontalRule) Html() string {
	var propertyString string
	for _, property := range hr.Properties {
		propertyString += fmt.Sprintf(" %s=\"%s\"", property.Name, property.Value)
	}
	return fmt.Sprintf("<hr%s/>", propertyString)
}

type Link struct {
	Properties []Property
	Url        string
	Text       string
}

func (l *Link) Html() string {
	var propertyString string
	for _, property := range l.Properties {
		propertyString += fmt.Sprintf(" %s=\"%s\"", property.Name, property.Value)
	}
	return fmt.Sprintf("<a%s href=\"%s\">%s</a>", propertyString, l.Url, l.Text)
}

type Button struct {
	Properties []Property
	Child      Component
	OnClick    string
}

func (b *Button) Html() string {
	var propertyString string
	for _, property := range b.Properties {
		propertyString += fmt.Sprintf(" %s=\"%s\"", property.Name, property.Value)
	}
	return fmt.Sprintf("<button%s onclick=\"%s()\">\n    %s\n</button>", propertyString, b.OnClick, b.Child.Html())
}

type Div struct {
	Properties []Property
	Children   []Component
}

func (d *Div) Html() string {
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

type Nav struct {
	Properties []Property
	Children   []Component
}

func (n *Nav) Html() string {
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

type Span struct {
	Properties []Property
	Children   []Component
}

func (s *Span) Html() string {
	var spanString string
	var propertyString string
	for _, property := range s.Properties {
		propertyString += fmt.Sprintf(" %s=\"%s\"", property.Name, property.Value)
	}

	if len(s.Children) == 0 {
		return fmt.Sprintf("<span%s/>", propertyString)
	}

	spanString += fmt.Sprintf("<span%s>\n", propertyString)
	for _, child := range s.Children {
		spanString += fmt.Sprintf("    %s\n", child.Html())
	}
	spanString += fmt.Sprintf("</span>")

	return spanString
}

type CodeBlock struct {
	Properties []Property
	Content    string
}

func (cb *CodeBlock) Html() string {
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

type Body struct {
	Children []Component
}

func (b *Body) Html() string {
	bodyString := "<body>\n"
	for _, child := range b.Children {
		bodyString += fmt.Sprintf("    %s\n", child.Html())
	}
	bodyString += "</body>"
	return bodyString
}
