package mdx

import (
	"fmt"
	"strconv"
	"strings"
)

type Parser struct {
	l          *Lexer
	currentTok Token
	nextTok    Token
}

func NewParser(lex *Lexer) *Parser {
	parser := &Parser{l: lex}
	parser.nextToken()
	parser.nextToken()
	return parser
}

type ParseError struct {
	error
	errorReason string
}

func (p *ParseError) Error() string {
	return fmt.Sprintf("ParseError occurred: %s", p.errorReason)
}

func (p *Parser) nextToken() {
	p.currentTok = p.nextTok
	p.nextTok = p.l.NextToken()
}

func (p *Parser) peekToken() Token {
	return p.nextTok
}

func (p *Parser) curTokenIs(tokenType TokenType) bool {
	return p.currentTok.Type == tokenType
}

func (p *Parser) peekTokenIs(tokenType TokenType) bool {
	return p.nextTok.Type == tokenType
}

func (p *Parser) isDoubleBreak() bool {
	return p.curTokenIs(NEWLINE) && (p.peekTokenIs(NEWLINE) || p.peekTokenIs(EOF))
}

func (p *Parser) Parse(delim TokenType) ([]Component, error) {
	elements := make([]Component, 0)
	var properties []Property
	var component Component

	for p.currentTok.Type != delim && p.currentTok.Type != EOF {
		if p.currentTok.Type == LSQUIRLY {
			var err error
			properties, err = p.parseProperties()
			if err != nil {
				return nil, err
			}
			continue
		} else {
			component = p.parseComponent(properties, delim)
		}

		if component != nil {
			elements = append(elements, component)
			properties = nil
		}

		p.nextToken()
	}

	return elements, nil
}

func (p *Parser) parseComponent(properties []Property, closing TokenType) Component {
	var component Component

	switch p.currentTok.Type {
	case HASH:
		component = p.parseHeader(properties, closing)
	case WORD:
		component = p.parseParagraph(properties, closing)
	case BACKTICK:
		component = p.parseCode(properties, closing)
	case ASTERISK:
		if p.peekTokenIs(ASTERISK) {
			component = p.parseStrong(properties, closing)
		} else {
			component = p.parseEm(properties, closing)
		}
	case GT:
		component = p.parseBlockQuote(properties, closing)
	case LISTELEMENT:
		component = p.parseOrderedListElement(properties, closing)
	case DASH:
		if p.peekTokenIs(SPACE) {
			component = p.parseUnorderedList(properties, closing)
		} else if p.peekTokenIs(DASH) {
			p.nextToken()
			if p.peekTokenIs(DASH) {
				component = &HorizontalRule{Properties: properties}
				p.nextToken()
			} else {
				component = p.parseFragment(properties, closing)
				prefixFragment(component, "-", closing)
			}
		} else {
			component = p.parseFragment(properties, closing)
		}
	case BANG:
		if p.peekTokenIs(LBRACKET) {
			component = p.parseImage(properties, closing)
		} else {
			component = p.parseFragment(properties, closing)
		}
	case UNDERSCORE:
		if p.peekTokenIs(UNDERSCORE) {
			p.nextToken()
			if p.peekTokenIs(UNDERSCORE) {
				component = &HorizontalRule{Properties: properties}
				p.nextToken()
			} else {
				component = p.parseFragment(properties, closing)
				prefixFragment(component, "_", closing)
			}
		} else {
			component = p.parseFragment(properties, closing)
		}
	case LBRACKET:
		if p.peekTokenIs(SPACE) || p.peekTokenIs(NEWLINE) {
			component = p.parseDiv(properties, closing)
		} else {
			component = p.parseLink(properties, closing)
		}
	case LT:
		component = p.parseShortLink(properties, closing)
	case TIDLE:
		component = p.parseButton(properties, closing)
	case AT:
		component = p.parseNav(properties, closing)
	case DOLLAR:
		component = p.parseSpan(properties, closing)
	case CARET:
		if p.peekTokenIs(CARET) {
			component = p.parseCodeBlock(properties, closing)
		} else {
			component = p.parseFragment(properties, closing)
		}
	case NEWLINE:
		component = &LineBreak{}
	case SLASH:
		if p.peekTokenIs(SLASH) {
			p.parseComment()
		} else {
			component = p.parseFragment(properties, closing)
		}
	}

	// if block component, skip newlines
	if component != nil && isBlockElement(component) {
		for p.peekTokenIs(NEWLINE) {
			p.nextToken()
		}
	}

	return component
}

func isBlockElement(component Component) bool {
	switch component.(type) {
	case *Div,
		*CodeBlock,
		*HorizontalRule,
		*Image,
		*Button,
		*Nav:
		return true
	}
	return false
}

func (p *Parser) parseProperties() ([]Property, error) {
	props := make([]Property, 0)
	for !p.curTokenIs(RSQUIRLY) {
		if p.curTokenIs(DOT) {
			if !p.peekTokenIs(WORD) {
				return nil, &ParseError{errorReason: "Property formatted incorrectly. DOT must be followed by a WORD"}
			}

			p.nextToken()
			key := p.currentTok.Literal

			if !p.peekTokenIs(EQUALS) {
				return nil, &ParseError{errorReason: "Property formatted incorrectly. KEY must be follwed by EQUALS"}
			}

			p.nextToken()
			if !p.peekTokenIs(WORD) {
				return nil, &ParseError{errorReason: "Property formatted incorrectly. EQUALS must be followed by VALUE"}
			}

			p.nextToken()
			value := p.currentTok.Literal

			props = append(props, Property{Name: key, Value: value})
		}

		p.nextToken()
	}

	p.nextToken()
	for p.curTokenIs(SPACE) || p.curTokenIs(NEWLINE) {
		p.nextToken()
	}

	return props, nil
}

func (p *Parser) parseFragment(properties []Property, closing TokenType) *Fragment {
	content := p.parseTextLine(closing)
	return &Fragment{String: content}
}

func (p *Parser) parseTextLine(closing TokenType) string {
	var lineString string
	for !(p.curTokenIs(NEWLINE) || p.curTokenIs(closing)) {
		lineString += p.currentTok.Literal
		p.nextToken()
	}
	return lineString
}

func (p *Parser) parseTextBlock(closing TokenType) string {
	var blockString string
	for !(p.curTokenIs(EOF) || p.curTokenIs(closing) || p.peekTokenIs(closing) || p.isDoubleBreak()) {
		blockString += p.currentTok.Literal
		p.nextToken()
	}

	return blockString
}

func (p *Parser) parseHeader(props []Property, closing TokenType) Component {
	level := 0
	for p.curTokenIs(HASH) {
		level++
		p.nextToken()
	}

	// next token must be space to be a valid header, otherwise just return a <p>
	if !p.curTokenIs(SPACE) {
		return p.parseFragment(props, closing)
	}

	p.nextToken()
	content := p.parseTextLine(closing)
	return &Header{Level: level, Text: content, Properties: props}
}

func (p *Parser) parseParagraph(props []Property, closing TokenType) Component {
	content := strings.ReplaceAll(p.parseTextBlock(closing), "\\n", " ")
	if len(content) == 0 {
		return nil
	}

	return &Paragraph{Text: content, Properties: props}
}

func prefixFragment(component Component, prefix string, closing TokenType) {
	switch c := (component).(type) {
	case *Fragment:
		c.String = prefix + c.String
	}
}

func (p *Parser) parseCode(properties []Property, closing TokenType) Component {
	p.nextToken()
	var codeString string

	for !p.curTokenIs(BACKTICK) {
		if p.curTokenIs(NEWLINE) || p.curTokenIs(EOF) {
			return &Fragment{String: "`" + codeString}
		}
		codeString += p.currentTok.Literal
		p.nextToken()
	}

	return &Code{Properties: properties, Text: codeString}
}

func (p *Parser) parseStrong(properties []Property, closing TokenType) Component {
	p.nextToken()
	if p.peekTokenIs(SPACE) || p.peekTokenIs(NEWLINE) || p.peekTokenIs(EOF) {
		content := p.parseTextLine(closing)
		return &Fragment{String: "*" + content}
	}

	p.nextToken()
	var strongString string

	for !(p.curTokenIs(ASTERISK) && p.peekTokenIs(ASTERISK)) {
		strongString += p.currentTok.Literal
		p.nextToken()

		if p.curTokenIs(NEWLINE) || p.curTokenIs(EOF) {
			fragment := &Fragment{String: strongString}
			prefixFragment(fragment, "**", closing)
			return fragment
		}
	}

	p.nextToken()
	return &Bold{Properties: properties, Text: strongString}
}

func (p *Parser) parseEm(properties []Property, closing TokenType) Component {
	if p.peekTokenIs(SPACE) || p.peekTokenIs(NEWLINE) || p.peekTokenIs(EOF) {
		content := p.parseTextLine(closing)
		return &Fragment{String: content}
	}

	p.nextToken()
	var emString string

	for !p.curTokenIs(ASTERISK) {
		emString += p.currentTok.Literal
		p.nextToken()

		if p.curTokenIs(NEWLINE) || p.curTokenIs(EOF) {
			fragment := &Fragment{String: emString}
			prefixFragment(fragment, "*", closing)
			return fragment
		}
	}

	return &Italic{Properties: properties, Text: emString}
}

func (p *Parser) parseBlockQuote(properties []Property, closing TokenType) Component {
	p.nextToken()
	content := strings.ReplaceAll(p.parseTextBlock(closing), "\\n", "<br/>")
	content = strings.TrimSpace(content)
	if len(content) == 0 {
		return nil
	}

	return &BlockQuote{Properties: properties, Text: content}
}

func (p *Parser) parseOrderedListElement(properties []Property, closing TokenType) Component {
	start, parseErr := strconv.Atoi(strings.TrimSuffix(p.currentTok.Literal, "."))
	if parseErr != nil {
		start = 1
	}

	listElements := make([]ListItem, 0)
	for !(p.curTokenIs(EOF) || (p.curTokenIs(NEWLINE) && !p.peekTokenIs(LISTELEMENT))) {
		p.nextToken()
		if p.curTokenIs(LISTELEMENT) {
			p.nextToken()
		}
		elementContent := strings.TrimSpace(p.parseTextLine(closing))
		element := ListItem{Component: &Paragraph{Text: elementContent}}
		listElements = append(listElements, element)
	}

	return &OrderedList{Properties: properties, ListItems: listElements, Start: start}
}

func (p *Parser) parseUnorderedList(properties []Property, closing TokenType) Component {
	listElements := make([]ListItem, 0)
	for !(p.curTokenIs(EOF) || (p.curTokenIs(NEWLINE) && !p.peekTokenIs(DASH))) {
		p.nextToken()
		if p.curTokenIs(DASH) {
			if !p.peekTokenIs(SPACE) {
				return &UnorderedList{Properties: properties, ListItems: listElements}
			}

			p.nextToken()
		}

		elementContent := strings.TrimSpace(p.parseTextLine(closing))
		element := ListItem{Component: &Paragraph{Text: elementContent}}
		listElements = append(listElements, element)
	}

	return &UnorderedList{Properties: properties, ListItems: listElements}
}

func (p *Parser) parseImage(properties []Property, closing TokenType) Component {
	p.nextToken()
	p.nextToken()

	var altText string
	for !p.curTokenIs(RBRACKET) {
		altText += p.currentTok.Literal
		p.nextToken()

		if p.curTokenIs(NEWLINE) || p.curTokenIs(EOF) {
			return &Fragment{String: "![" + altText}
		}
	}

	if !p.peekTokenIs(LPAREN) {
		return &Fragment{String: "![" + altText + "]"}
	}

	p.nextToken()
	p.nextToken()

	var urlString string
	for !p.curTokenIs(RPAREN) {
		urlString += p.currentTok.Literal
		p.nextToken()

		if p.curTokenIs(NEWLINE) || p.curTokenIs(EOF) {
			return &Fragment{String: "![" + altText + "](" + urlString}
		}
	}

	return &Image{Properties: properties, ImgUrl: urlString, AltText: altText}
}

func (p *Parser) parseDiv(properties []Property, closing TokenType) Component {
	children := make([]Component, 0)

	p.nextToken()
	for p.curTokenIs(NEWLINE) {
		p.nextToken()
	}

	components, err := p.Parse(RBRACKET)
	if err != nil {
		panic(err.Error())
	}

	for _, component := range components {
		children = append(children, component)
	}

	if p.peekTokenIs(NEWLINE) {
		p.nextToken()
	}

	return &Div{Properties: properties, Children: children}
}

func (p *Parser) parseLink(properties []Property, closing TokenType) Component {
	p.nextToken()

	var displayText string
	for !p.curTokenIs(RBRACKET) {
		displayText += p.currentTok.Literal
		p.nextToken()

		if p.curTokenIs(NEWLINE) || p.curTokenIs(EOF) {
			return &Fragment{String: "[" + displayText}
		}
	}

	if !p.peekTokenIs(LPAREN) {
		return &Fragment{String: "[" + displayText + "]"}
	}

	p.nextToken()
	p.nextToken()

	var urlString string
	for !p.curTokenIs(RPAREN) {
		urlString += p.currentTok.Literal
		p.nextToken()

		if p.curTokenIs(NEWLINE) || p.curTokenIs(EOF) {
			return &Fragment{String: "[" + displayText + "](" + urlString}
		}
	}

	return &Link{Properties: properties, Url: urlString, Text: displayText}
}

func (p *Parser) parseShortLink(properties []Property, closing TokenType) Component {
	p.nextToken()

	var urlString string
	for !p.curTokenIs(GT) {
		urlString += p.currentTok.Literal
		p.nextToken()

		if p.curTokenIs(NEWLINE) || p.curTokenIs(EOF) {
			return &Fragment{String: "<" + urlString}
		}
	}

	return &Link{Properties: properties, Url: urlString, Text: urlString}
}

func (p *Parser) parseButton(properties []Property, closing TokenType) Component {
	p.nextToken()
	p.nextToken()

	var buttonLabel string
	for !p.curTokenIs(RBRACKET) {
		buttonLabel += p.currentTok.Literal
		p.nextToken()

		if p.curTokenIs(NEWLINE) || p.curTokenIs(EOF) {
			return &Fragment{String: "~[" + buttonLabel}
		}
	}

	if !p.peekTokenIs(LPAREN) {
		return &Fragment{String: "~[" + buttonLabel + "]"}
	}

	p.nextToken()
	p.nextToken()

	var onClick string
	for !p.curTokenIs(RPAREN) {
		onClick += p.currentTok.Literal
		p.nextToken()

		if p.curTokenIs(NEWLINE) || p.curTokenIs(EOF) {
			return &Fragment{String: "~[" + buttonLabel + "](" + onClick}
		}
	}

	return &Button{Properties: properties, OnClick: onClick, Child: &Fragment{String: buttonLabel}}
}

func (p *Parser) parseNav(properties []Property, closing TokenType) Component {
	children := make([]Component, 0)

	p.nextToken()
	components, err := p.Parse(AT)
	if err != nil {
		panic(err.Error())
	}

	for _, component := range components {
		// don't put line breaks in nav element
		if _, ok := component.(*LineBreak); !ok {
			children = append(children, component)
		}
	}

	return &Nav{Properties: properties, Children: children}
}

func (p *Parser) parseSpan(properties []Property, closing TokenType) Component {
	if p.peekTokenIs(NEWLINE) || p.peekTokenIs(EOF) {
		content := p.parseTextLine(closing)
		return &Fragment{String: content}
	}

	p.nextToken()
	var spanString string

	for !p.curTokenIs(DOLLAR) {
		spanString += p.currentTok.Literal
		p.nextToken()

		if p.curTokenIs(NEWLINE) || p.curTokenIs(EOF) {
			fragment := &Fragment{String: spanString}
			prefixFragment(fragment, "$", closing)
			return fragment
		}
	}

	children := []Component{&Fragment{String: strings.TrimSpace(spanString)}}

	return &Span{Properties: properties, Children: children}
}

func (p *Parser) parseCodeBlock(properties []Property, closing TokenType) Component {
	p.nextToken()
	p.nextToken()

	var codeBlockString string
	for !(p.curTokenIs(CARET) && p.peekTokenIs(CARET)) {
		codeBlockString += p.currentTok.Literal
		p.nextToken()

		if p.curTokenIs(EOF) {
			fragment := &Fragment{String: "^^" + codeBlockString}
			return fragment
		}
	}
	p.nextToken()

	codeBlockString = strings.ReplaceAll(codeBlockString, "\\t", "    ")
	codeBlockString = strings.TrimPrefix(codeBlockString, "\\n")
	codeBlockString = strings.TrimSuffix(codeBlockString, "\\n")
	return &CodeBlock{Properties: properties, Content: codeBlockString}
}

func (p *Parser) parseComment() {
	for !p.curTokenIs(NEWLINE) {
		p.nextToken()
	}
}
