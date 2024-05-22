package mdx

import (
	"fmt"
	"strings"
)

type component interface {
	// Convert element to HTML representation
	Html() string
}

type property struct {
	Name  string
	Value string
}

type fragment struct {
	String string
}

func (f *fragment) Html() string {
	return fmt.Sprintf("%s", f.String)
}

type lineBreak struct{}

func (lb *lineBreak) Html() string {
	return "<br/>"
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

type paragraph struct {
	Properties []property
	Content    []component
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

type link struct {
	Properties []property
	Url        string
	Text       string
}

func (l *link) Html() string {
	var propertyString string
	for _, property := range l.Properties {
		propertyString += fmt.Sprintf(" %s=\"%s\"", property.Name, property.Value)
	}
	return fmt.Sprintf("<a%s href=\"%s\">%s</a>", propertyString, l.Url, l.Text)
}

type button struct {
	Properties []property
	Child      component
	OnClick    string
}

func (b *button) Html() string {
	var propertyString string
	for _, property := range b.Properties {
		propertyString += fmt.Sprintf(" %s=\"%s\"", property.Name, property.Value)
	}
	return fmt.Sprintf("<button%s onclick=\"%s()\">\n    %s\n</button>", propertyString, b.OnClick, b.Child.Html())
}

type div struct {
	Properties []property
	Children   []component
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

type span struct {
	Properties []property
	Children   []component
}

func (s *span) Html() string {
	var spanString string
	var propertyString string
	for _, property := range s.Properties {
		propertyString += fmt.Sprintf(" %s=\"%s\"", property.Name, property.Value)
	}

	if len(s.Children) == 0 {
		return fmt.Sprintf("<span%s/>", propertyString)
	}

	spanString += fmt.Sprintf("<span%s>", propertyString)
	for _, child := range s.Children {
		spanString += fmt.Sprintf("%s", child.Html())
	}
	spanString += fmt.Sprintf("</span>")

	return spanString
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
