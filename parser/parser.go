package parser

import (
	"fmt"
	"mdx/ast"
	"mdx/lexer"
	"mdx/token"
	"strconv"
	"strings"
)

type Parser struct {
	l          *lexer.Lexer
	currentTok token.Token
	nextTok    token.Token
}

func New(lex *lexer.Lexer) *Parser {
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

func (p *Parser) peekToken() token.Token {
	return p.nextTok
}

func (p *Parser) curTokenIs(tokenType token.TokenType) bool {
	return p.currentTok.Type == tokenType
}

func (p *Parser) peekTokenIs(tokenType token.TokenType) bool {
	return p.nextTok.Type == tokenType
}

func (p *Parser) Parse() ([]ast.Component, error) {
	elements := make([]ast.Component, 0)
	var properties []ast.Property
	var component ast.Component

	for p.currentTok.Type != token.EOF {
		if p.currentTok.Type == token.LSQUIRLY {
			var err error
			properties, err = p.parseProperties()
			if err != nil {
				return nil, err
			}
			continue
		} else {
			component = p.parseComponent(properties)
		}

		if component != nil {
			elements = append(elements, component)
			properties = nil
		}

		p.nextToken()
	}

	return elements, nil
}

func (p *Parser) parseComponent(properties []ast.Property) ast.Component {
	var component ast.Component

	switch p.currentTok.Type {
	case token.HASH:
		component = p.parseHeader(properties)
	case token.WORD:
		component = p.parseParagraph(properties)
	case token.BACKTICK:
		component = p.parseCode(properties)
	case token.ASTERISK:
		if p.peekTokenIs(token.ASTERISK) {
			component = p.parseStrong(properties)
		} else {
			component = p.parseEm(properties)
		}
	case token.GT:
		component = p.parseBlockQuote(properties)
	case token.LISTELEMENT:
		component = p.parseOrderedListElement(properties)
	case token.DASH:
		if p.peekTokenIs(token.SPACE) {
			component = p.parseUnorderedList(properties)
		} else if p.peekTokenIs(token.DASH) {
			p.nextToken()
			if p.peekTokenIs(token.DASH) {
				component = &ast.HorizontalRule{Properties: properties}
				p.nextToken()
			} else {
				component = p.parseFragment(properties)
				prefixFragment(component, "-")
			}
		} else {
			component = p.parseFragment(properties)
		}
	case token.BANG:
		if p.peekTokenIs(token.LBRACKET) {
			component = p.parseImage(properties)
		} else {
			component = p.parseFragment(properties)
		}
	case token.UNDERSCORE:
		if p.peekTokenIs(token.UNDERSCORE) {
			p.nextToken()
			if p.peekTokenIs(token.UNDERSCORE) {
				component = &ast.HorizontalRule{Properties: properties}
				p.nextToken()
			} else {
				component = p.parseFragment(properties)
				prefixFragment(component, "_")
			}
		} else {
			component = p.parseFragment(properties)
		}
	case token.LBRACKET:
		if p.peekTokenIs(token.SPACE) || p.peekTokenIs(token.NEWLINE) {
			component = p.parseDiv(properties)
		} else {
			component = p.parseLink(properties)
		}
	case token.LT:
		component = p.parseShortLink(properties)
	case token.TIDLE:
		component = p.parseButton(properties)
	case token.AT:
		component = p.parseNav(properties)
	case token.DOLLAR:
		component = p.parseSpan(properties)
	case token.CARET:
		if p.peekTokenIs(token.CARET) {
			component = p.parseCodeBlock(properties)
		} else {
			component = p.parseFragment(properties)
		}
	case token.NEWLINE:
		component = &ast.LineBreak{}
	case token.SLASH:
		if p.peekTokenIs(token.SLASH) {
			p.parseComment()
		} else {
			component = p.parseFragment(properties)
		}
	}

	return component
}

func (p *Parser) parseProperties() ([]ast.Property, error) {
	props := make([]ast.Property, 0)
	for !p.curTokenIs(token.RSQUIRLY) {
		if p.curTokenIs(token.DOT) {
			if !p.peekTokenIs(token.WORD) {
				return nil, &ParseError{errorReason: "Property formatted incorrectly. DOT must be followed by a WORD"}
			}

			p.nextToken()
			key := p.currentTok.Literal

			if !p.peekTokenIs(token.EQUALS) {
				return nil, &ParseError{errorReason: "Property formatted incorrectly. KEY must be follwed by EQUALS"}
			}

			p.nextToken()
			if !p.peekTokenIs(token.WORD) {
				return nil, &ParseError{errorReason: "Property formatted incorrectly. EQUALS must be followed by VALUE"}
			}

			p.nextToken()
			value := p.currentTok.Literal

			props = append(props, ast.Property{Name: key, Value: value})
		}

		p.nextToken()
	}

	p.nextToken()
	for p.curTokenIs(token.SPACE) || p.curTokenIs(token.NEWLINE) {
		p.nextToken()
	}

	return props, nil
}

func (p *Parser) parseFragment(properties []ast.Property) *ast.Fragment {
	content := p.parseTextLine()
	return &ast.Fragment{String: content}
}

func (p *Parser) parseTextLine() string {
	var lineString string
	for !p.curTokenIs(token.NEWLINE) {
		lineString += p.currentTok.Literal
		p.nextToken()
	}
	return lineString
}

func (p *Parser) parseTextBlock() string {
	var blockString string
	for !(p.curTokenIs(token.EOF) || p.curTokenIs(token.NEWLINE) && (p.peekTokenIs(token.NEWLINE) || p.peekTokenIs(token.EOF))) {
		blockString += p.currentTok.Literal
		p.nextToken()
	}
	return blockString
}

func (p *Parser) parseHeader(props []ast.Property) ast.Component {
	level := 0
	for p.curTokenIs(token.HASH) {
		level++
		p.nextToken()
	}

	// next token must be space to be a valid header, otherwise just return a <p>
	if !p.curTokenIs(token.SPACE) {
		return p.parseFragment(props)
	}

	p.nextToken()
	content := p.parseTextLine()
	return &ast.Header{Level: level, Text: content, Properties: props}
}

func (p *Parser) parseParagraph(props []ast.Property) ast.Component {
	content := strings.ReplaceAll(p.parseTextBlock(), "\\n", " ")
	if len(content) == 0 {
		return nil
	}

	return &ast.Paragraph{Text: content, Properties: props}
}

func prefixFragment(component ast.Component, prefix string) {
	switch c := (component).(type) {
	case *ast.Fragment:
		c.String = prefix + c.String
	}
}

func (p *Parser) parseCode(properties []ast.Property) ast.Component {
	p.nextToken()
	var codeString string

	for !p.curTokenIs(token.BACKTICK) {
		if p.curTokenIs(token.NEWLINE) || p.curTokenIs(token.EOF) {
			return &ast.Fragment{String: "`" + codeString}
		}
		codeString += p.currentTok.Literal
		p.nextToken()
	}

	return &ast.Code{Properties: properties, Text: codeString}
}

func (p *Parser) parseStrong(properties []ast.Property) ast.Component {
	p.nextToken()
	if p.peekTokenIs(token.SPACE) || p.peekTokenIs(token.NEWLINE) || p.peekTokenIs(token.EOF) {
		content := p.parseTextLine()
		return &ast.Fragment{String: "*" + content}
	}

	p.nextToken()
	var strongString string

	for !(p.curTokenIs(token.ASTERISK) && p.peekTokenIs(token.ASTERISK)) {
		strongString += p.currentTok.Literal
		p.nextToken()

		if p.curTokenIs(token.NEWLINE) || p.curTokenIs(token.EOF) {
			fragment := &ast.Fragment{String: strongString}
			prefixFragment(fragment, "**")
			return fragment
		}
	}

	p.nextToken()
	return &ast.Bold{Properties: properties, Text: strongString}
}

func (p *Parser) parseEm(properties []ast.Property) ast.Component {
	if p.peekTokenIs(token.SPACE) || p.peekTokenIs(token.NEWLINE) || p.peekTokenIs(token.EOF) {
		content := p.parseTextLine()
		return &ast.Fragment{String: content}
	}

	p.nextToken()
	var emString string

	for !p.curTokenIs(token.ASTERISK) {
		emString += p.currentTok.Literal
		p.nextToken()

		if p.curTokenIs(token.NEWLINE) || p.curTokenIs(token.EOF) {
			fragment := &ast.Fragment{String: emString}
			prefixFragment(fragment, "*")
			return fragment
		}
	}

	return &ast.Italic{Properties: properties, Text: emString}
}

func (p *Parser) parseBlockQuote(properties []ast.Property) ast.Component {
	p.nextToken()
	content := strings.ReplaceAll(p.parseTextBlock(), "\\n", "<br/>")
	content = strings.TrimSpace(content)
	if len(content) == 0 {
		return nil
	}

	return &ast.BlockQuote{Properties: properties, Text: content}
}

func (p *Parser) parseOrderedListElement(properties []ast.Property) ast.Component {
	start, parseErr := strconv.Atoi(strings.TrimSuffix(p.currentTok.Literal, "."))
	if parseErr != nil {
		start = 1
	}

	listElements := make([]ast.ListItem, 0)
	for !(p.curTokenIs(token.EOF) || (p.curTokenIs(token.NEWLINE) && !p.peekTokenIs(token.LISTELEMENT))) {
		p.nextToken()
		if p.curTokenIs(token.LISTELEMENT) {
			p.nextToken()
		}
		elementContent := strings.TrimSpace(p.parseTextLine())
		element := ast.ListItem{Component: &ast.Paragraph{Text: elementContent}}
		listElements = append(listElements, element)
	}

	return &ast.OrderedList{Properties: properties, ListItems: listElements, Start: start}
}

func (p *Parser) parseUnorderedList(properties []ast.Property) ast.Component {
	listElements := make([]ast.ListItem, 0)
	for !(p.curTokenIs(token.EOF) || (p.curTokenIs(token.NEWLINE) && !p.peekTokenIs(token.DASH))) {
		p.nextToken()
		if p.curTokenIs(token.DASH) {
			if !p.peekTokenIs(token.SPACE) {
				return &ast.UnorderedList{Properties: properties, ListItems: listElements}
			}

			p.nextToken()
		}

		elementContent := strings.TrimSpace(p.parseTextLine())
		element := ast.ListItem{Component: &ast.Paragraph{Text: elementContent}}
		listElements = append(listElements, element)
	}

	return &ast.UnorderedList{Properties: properties, ListItems: listElements}
}

func (p *Parser) parseImage(properties []ast.Property) ast.Component {
	p.nextToken()
	p.nextToken()

	var altText string
	for !p.curTokenIs(token.RBRACKET) {
		altText += p.currentTok.Literal
		p.nextToken()

		if p.curTokenIs(token.NEWLINE) || p.curTokenIs(token.EOF) {
			return &ast.Fragment{String: "![" + altText}
		}
	}

	if !p.peekTokenIs(token.LPAREN) {
		return &ast.Fragment{String: "![" + altText + "]"}
	}

	p.nextToken()
	p.nextToken()

	var urlString string
	for !p.curTokenIs(token.RPAREN) {
		urlString += p.currentTok.Literal
		p.nextToken()

		if p.curTokenIs(token.NEWLINE) || p.curTokenIs(token.EOF) {
			return &ast.Fragment{String: "![" + altText + "](" + urlString}
		}
	}

	return &ast.Image{Properties: properties, ImgUrl: urlString, AltText: altText}
}

func (p *Parser) parseDiv(properties []ast.Property) ast.Component {
	// var divString string

	// TODO: Do we need to refactor for child elements to write this one?

	return &ast.Fragment{}
}

func (p *Parser) parseLink(properties []ast.Property) ast.Component {

	return &ast.Fragment{}
}

func (p *Parser) parseShortLink(properties []ast.Property) ast.Component {

	return &ast.Fragment{}
}

func (p *Parser) parseButton(properties []ast.Property) ast.Component {
	p.nextToken()
	p.nextToken()

	var buttonLabel string
	for !p.curTokenIs(token.RBRACKET) {
		buttonLabel += p.currentTok.Literal
		p.nextToken()

		if p.curTokenIs(token.NEWLINE) || p.curTokenIs(token.EOF) {
			return &ast.Fragment{String: "~[" + buttonLabel}
		}
	}

	if !p.peekTokenIs(token.LPAREN) {
		return &ast.Fragment{String: "~[" + buttonLabel + "]"}
	}

	p.nextToken()
	p.nextToken()

	var onClick string
	for !p.curTokenIs(token.RPAREN) {
		onClick += p.currentTok.Literal
		p.nextToken()

		if p.curTokenIs(token.NEWLINE) || p.curTokenIs(token.EOF) {
			return &ast.Fragment{String: "~[" + buttonLabel + "](" + onClick}
		}
	}

	return &ast.Button{Properties: properties, OnClick: onClick, Child: &ast.Fragment{String: buttonLabel}}
}

func (p *Parser) parseNav(properties []ast.Property) ast.Component {

	return &ast.Fragment{}
}

func (p *Parser) parseSpan(properties []ast.Property) ast.Component {

	return &ast.Fragment{}
}

func (p *Parser) parseCodeBlock(properties []ast.Property) ast.Component {
	p.nextToken()
	p.nextToken()

	var codeBlockString string
	for !(p.curTokenIs(token.CARET) && p.peekTokenIs(token.CARET)) {
		codeBlockString += p.currentTok.Literal
		p.nextToken()

		if p.curTokenIs(token.EOF) {
			fragment := &ast.Fragment{String: "^^" + codeBlockString}
			return fragment
		}
	}
	p.nextToken()

	codeBlockString = strings.ReplaceAll(codeBlockString, "\\t", "    ")
	codeBlockString = strings.TrimPrefix(codeBlockString, "\\n")
	codeBlockString = strings.TrimSuffix(codeBlockString, "\\n")
	return &ast.CodeBlock{Properties: properties, Content: codeBlockString}
}

func (p *Parser) parseComment() {
	for !p.curTokenIs(token.NEWLINE) {
		p.nextToken()
	}
}
