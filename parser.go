package mdx

import (
	"fmt"
	"strconv"
	"strings"
)

// TODO: Properly handle nested components

type parser struct {
	l          *lexer
	currentTok token
	nextTok    token
}

func newParser(lex *lexer) *parser {
	parser := &parser{l: lex}
	parser.nextToken()
	parser.nextToken()
	return parser
}

type parseError struct {
	error
	errorReason string
}

func (p *parseError) Error() string {
	return fmt.Sprintf("ParseError occurred: %s", p.errorReason)
}

func (p *parser) nextToken() {
	p.currentTok = p.nextTok
	p.nextTok = p.l.nextToken()
}

func (p *parser) peekToken() *token {
	return &p.nextTok
}

func (p *parser) curTokenIs(tokType tokenType) bool {
	return p.currentTok.Type == tokType
}

func (p *parser) peekTokenIs(tokType tokenType) bool {
	return p.nextTok.Type == tokType
}

func (p *parser) isDoubleBreak() bool {
	return p.curTokenIs(newline) && (p.peekTokenIs(newline) || p.peekTokenIs(eof))
}

func (p *parser) isAfterNewline(tokType tokenType) bool {
	return p.curTokenIs(newline) && p.peekTokenIs(tokType)
}

func (p *parser) isNextLineElement() bool {
	return p.curTokenIs(newline) && p.peekToken().IsElementToken()
}

func (p *parser) parse(delim tokenType) ([]component, error) {
	elements := make([]component, 0)
	var properties []property
	var component component

	for p.currentTok.Type != delim && p.currentTok.Type != eof {
		if p.currentTok.Type == lsquirly {
			var err error
			properties, err, _ = p.parseProperties()
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

func (p *parser) parseComponent(properties []property, closing tokenType) component {
	var component component

	switch p.currentTok.Type {
	case hash:
		component = p.parseHeader(properties, closing)
	case word:
		component = p.parseParagraph(properties, closing)
	case backtick:
		component = p.parseCode(properties, closing)
	case asterisk:
		if p.peekTokenIs(asterisk) {
			component = p.parseStrong(properties, closing)
		} else {
			component = p.parseEm(properties, closing)
		}
	case gt:
		component = p.parseBlockQuote(properties, closing)
	case listelement:
		component = p.parseOrderedListElement(properties, closing)
	case dash:
		if p.peekTokenIs(space) {
			component = p.parseUnorderedList(properties, closing)
		} else if p.peekTokenIs(dash) {
			p.nextToken()
			if p.peekTokenIs(dash) {
				component = &horizontalRule{Properties: properties}
				p.nextToken()
			} else {
				component = p.parseFragment(properties, closing)
				prefixFragment(component, "-", closing)
			}
		} else {
			component = p.parseFragment(properties, closing)
		}
	case bang:
		if p.peekTokenIs(lbracket) {
			component = p.parseImage(properties, closing)
		} else {
			component = p.parseFragment(properties, closing)
		}
	case underscore:
		if p.peekTokenIs(underscore) {
			p.nextToken()
			if p.peekTokenIs(underscore) {
				component = &horizontalRule{Properties: properties}
				p.nextToken()
			} else {
				component = p.parseFragment(properties, closing)
				prefixFragment(component, "_", closing)
			}
		} else {
			component = p.parseFragment(properties, closing)
		}
	case lbracket:
		if p.peekTokenIs(space) || p.peekTokenIs(newline) {
			component = p.parseDiv(properties, closing)
		} else {
			component = p.parseLink(properties, closing)
		}
	case lt:
		component = p.parseShortLink(properties, closing)
	case tidle:
		component = p.parseButton(properties, closing)
	case at:
		component = p.parseNav(properties, closing)
	case dollar:
		component = p.parseSpan(properties, closing)
	case caret:
		if p.peekTokenIs(caret) {
			component = p.parseCodeBlock(properties, closing)
		} else {
			component = p.parseFragment(properties, closing)
		}
	case newline:
		component = &lineBreak{}
	case slash:
		if p.peekTokenIs(slash) {
			p.parseComment()
		} else {
			component = p.parseFragment(properties, closing)
		}
	}

	// if block component, skip newlines
	if component != nil && isBlockElement(component) {
		for p.peekTokenIs(newline) {
			p.nextToken()
		}
	}

	return component
}

func isBlockElement(comp component) bool {
	switch comp.(type) {
	case *div,
		*codeBlock,
		*horizontalRule,
		*image,
		*button,
		*nav:
		return true
	}
	return false
}

func (p *parser) parseProperties() ([]property, error, string) {
	props := make([]property, 0)
	propsString := "{"
	for !p.curTokenIs(rsquirly) {
		if p.curTokenIs(dot) {
			if !p.peekTokenIs(word) {
				return nil, &parseError{errorReason: "Property formatted incorrectly. DOT must be followed by a WORD"}, propsString
			}

			p.nextToken()
			propsString += p.currentTok.Literal
			key := p.currentTok.Literal

			if !p.peekTokenIs(equals) {
				return nil, &parseError{errorReason: "Property formatted incorrectly. KEY must be follwed by EQUALS"}, propsString
			}

			p.nextToken()
			propsString += p.currentTok.Literal
			if !p.peekTokenIs(word) {
				return nil, &parseError{errorReason: "Property formatted incorrectly. EQUALS must be followed by VALUE"}, propsString
			}

			p.nextToken()
			// propsString += p.currentTok.Literal
			value := p.currentTok.Literal

			props = append(props, property{Name: key, Value: value})
		}

		p.nextToken()
		propsString += p.currentTok.Literal
	}

	p.nextToken()
	for p.curTokenIs(space) || p.curTokenIs(newline) {
		p.nextToken()
	}

	return props, nil, ""
}

func (p *parser) parseFragment(properties []property, closing tokenType) *fragment {
	content := p.parseTextLine(closing)
	return &fragment{String: content}
}

func (p *parser) parseTextLine(closing tokenType) string {
	var lineString string
	for !(p.curTokenIs(newline) || p.curTokenIs(closing)) {
		lineString += p.currentTok.Literal
		p.nextToken()
	}
	return lineString
}

func (p *parser) parseLine(closing tokenType) []component {
	lineElements := make([]component, 0)
	var lineString string

	for !(p.curTokenIs(newline) || p.curTokenIs(closing)) {
		if p.currentTok.IsElementToken() {
			if len(lineString) > 0 {
				lineElements = append(lineElements, &fragment{String: lineString})
				lineString = ""
			}

			lineElements = append(lineElements, p.parseComponent(nil, closing))
		} else if p.curTokenIs(lsquirly) {
			properties, parseErr, propsText := p.parseProperties()
			if parseErr != nil {
				lineString += propsText
				p.nextToken()
			} else {
				if len(lineString) > 0 {
					lineElements = append(lineElements, &fragment{String: lineString})
					lineString = ""
				}
				lineElements = append(lineElements, p.parseComponent(properties, closing))
			}
		} else {
			lineString += p.currentTok.Literal
			p.nextToken()
		}
	}

	if len(lineString) > 0 {
		lineElements = append(lineElements, &fragment{String: lineString})
	}

	return lineElements
}

func (p *parser) parseLineDoubleClose(closing tokenType) []component {
	lineElements := make([]component, 0)
	var lineString string

	for !(p.curTokenIs(newline) || (p.curTokenIs(closing) && p.peekTokenIs(closing))) {
		if p.currentTok.IsElementToken() {
			if len(lineString) > 0 {
				lineElements = append(lineElements, &fragment{String: lineString})
				lineString = ""
			}

			lineElements = append(lineElements, p.parseComponent(nil, closing))
		} else if p.curTokenIs(lsquirly) {
			properties, parseErr, propsText := p.parseProperties()
			if parseErr != nil {
				lineString += propsText
				p.nextToken()
			} else {
				if len(lineString) > 0 {
					lineElements = append(lineElements, &fragment{String: lineString})
					lineString = ""
				}
				lineElements = append(lineElements, p.parseComponent(properties, closing))
			}
		} else {
			lineString += p.currentTok.Literal
			p.nextToken()
		}
	}

	if len(lineString) > 0 {
		lineElements = append(lineElements, &fragment{String: lineString})
	}

	return lineElements
}

func (p *parser) parseTextBlock(closing tokenType) string {
	var blockString string
	for !(p.curTokenIs(eof) || p.curTokenIs(closing) || p.isDoubleBreak() || p.isAfterNewline(closing) || p.isNextLineElement()) {
		blockString += p.currentTok.Literal
		p.nextToken()
	}

	return blockString
}

func (p *parser) parseBlock(closing tokenType) []component {
	blockElements := make([]component, 0)
	var blockString string

	for !(p.curTokenIs(eof) || p.curTokenIs(closing) || p.isDoubleBreak() || p.isAfterNewline(closing) || p.isNextLineElement()) {
		if p.currentTok.IsElementToken() {
			if len(blockString) > 0 {
				blockElements = append(blockElements, &fragment{String: blockString})
				blockString = ""
			}

			blockElements = append(blockElements, p.parseComponent(nil, closing))
		} else if p.curTokenIs(lsquirly) {
			properties, parseErr, propsText := p.parseProperties()
			if parseErr != nil {
				blockString += propsText
				p.nextToken()
			} else {
				if len(blockString) > 0 {
					blockElements = append(blockElements, &fragment{String: blockString})
					blockString = ""
				}
				blockElements = append(blockElements, p.parseComponent(properties, closing))
			}
		} else {
			blockString += p.currentTok.Literal
			p.nextToken()
		}
	}

	if len(blockString) > 0 {
		blockElements = append(blockElements, &fragment{String: blockString})
	}

	return blockElements
}

func (p *parser) parseHeader(props []property, closing tokenType) component {
	level := 0
	for p.curTokenIs(hash) {
		level++
		p.nextToken()
	}

	// next token must be space to be a valid header, otherwise just return a <p>
	if !p.curTokenIs(space) {
		return p.parseFragment(props, closing)
	}

	p.nextToken()
	contentElements := p.parseLine(closing)
	return &header{Level: level, Content: contentElements, Properties: props}
}

func (p *parser) parseParagraph(props []property, closing tokenType) component {
	// content := strings.ReplaceAll(p.parseTextBlock(closing), "\\n", " ")
	contentElements := p.parseBlock(closing)
	if len(contentElements) == 0 {
		return nil
	}

	return &paragraph{Content: contentElements, Properties: props}
}

func prefixFragment(comp component, prefix string, closing tokenType) {
	switch c := (comp).(type) {
	case *fragment:
		c.String = prefix + c.String
	}
}

func (p *parser) parseCode(properties []property, closing tokenType) component {
	p.nextToken()
	var codeString string

	for !p.curTokenIs(backtick) {
		if p.curTokenIs(newline) || p.curTokenIs(eof) {
			return &fragment{String: "`" + codeString}
		}
		codeString += p.currentTok.Literal
		p.nextToken()
	}

	p.nextToken()
	return &code{Properties: properties, Text: codeString}
}

func (p *parser) parseStrong(properties []property, closing tokenType) component {
	p.nextToken()
	if p.peekTokenIs(space) || p.peekTokenIs(newline) || p.peekTokenIs(eof) {
		content := p.parseTextLine(closing)
		return &fragment{String: "*" + content}
	}

	p.nextToken()

	content := p.parseLineDoubleClose(asterisk)

	p.nextToken()
	p.nextToken()

	return &bold{Properties: properties, Content: content}
}

func (p *parser) parseEm(properties []property, closing tokenType) component {
	if p.peekTokenIs(space) || p.peekTokenIs(newline) || p.peekTokenIs(eof) {
		content := p.parseTextLine(closing)
		p.nextToken()
		return &fragment{String: content}
	}

	p.nextToken()
	content := p.parseLine(asterisk)

	p.nextToken()
	return &italic{Properties: properties, Content: content}
}

func (p *parser) parseBlockQuote(properties []property, closing tokenType) component {
	p.nextToken()
	content := strings.ReplaceAll(p.parseTextBlock(closing), "\\n", "<br/>")
	content = strings.TrimSpace(content)
	if len(content) == 0 {
		return nil
	}

	return &blockQuote{Properties: properties, Text: content}
}

func (p *parser) parseOrderedListElement(properties []property, closing tokenType) component {
	start, parseErr := strconv.Atoi(strings.TrimSuffix(p.currentTok.Literal, "."))
	if parseErr != nil {
		start = 1
	}

	listElements := make([]listItem, 0)
	for !(p.curTokenIs(eof) || (p.curTokenIs(newline) && !p.peekTokenIs(listelement))) {
		p.nextToken()
		if p.curTokenIs(listelement) {
			p.nextToken()
		}
		elementContent := strings.TrimSpace(p.parseTextLine(closing))
		element := listItem{Component: &paragraph{Content: []component{&fragment{String: elementContent}}}}
		listElements = append(listElements, element)
	}

	return &orderedList{Properties: properties, ListItems: listElements, Start: start}
}

func (p *parser) parseUnorderedList(properties []property, closing tokenType) component {
	listElements := make([]listItem, 0)
	for !(p.curTokenIs(eof) || (p.curTokenIs(newline) && !p.peekTokenIs(dash))) {
		p.nextToken()
		if p.curTokenIs(dash) {
			if !p.peekTokenIs(space) {
				return &unorderedList{Properties: properties, ListItems: listElements}
			}

			p.nextToken()
		}

		elementContent := strings.TrimSpace(p.parseTextLine(closing))
		element := listItem{Component: &paragraph{Content: []component{&fragment{String: elementContent}}}}
		listElements = append(listElements, element)
	}

	return &unorderedList{Properties: properties, ListItems: listElements}
}

func (p *parser) parseImage(properties []property, closing tokenType) component {
	p.nextToken()
	p.nextToken()

	var altText string
	for !p.curTokenIs(rbracket) {
		altText += p.currentTok.Literal
		p.nextToken()

		if p.curTokenIs(newline) || p.curTokenIs(eof) {
			return &fragment{String: "![" + altText}
		}
	}

	if !p.peekTokenIs(lparen) {
		return &fragment{String: "![" + altText + "]"}
	}

	p.nextToken()
	p.nextToken()

	var urlString string
	for !p.curTokenIs(rparen) {
		urlString += p.currentTok.Literal
		p.nextToken()

		if p.curTokenIs(newline) || p.curTokenIs(eof) {
			return &fragment{String: "![" + altText + "](" + urlString}
		}
	}

	return &image{Properties: properties, ImgUrl: urlString, AltText: altText}
}

func (p *parser) parseDiv(properties []property, closing tokenType) component {
	children := make([]component, 0)

	p.nextToken()
	for p.curTokenIs(newline) {
		p.nextToken()
	}

	components, err := p.parse(rbracket)
	if err != nil {
		panic(err.Error())
	}

	for _, component := range components {
		children = append(children, component)
	}

	if p.peekTokenIs(newline) {
		p.nextToken()
	}

	return &div{Properties: properties, Children: children}
}

func (p *parser) parseLink(properties []property, closing tokenType) component {
	p.nextToken()

	var displayText string
	for !p.curTokenIs(rbracket) {
		displayText += p.currentTok.Literal
		p.nextToken()

		if p.curTokenIs(newline) || p.curTokenIs(eof) {
			return &fragment{String: "[" + displayText}
		}
	}

	if !p.peekTokenIs(lparen) {
		return &fragment{String: "[" + displayText + "]"}
	}

	p.nextToken()
	p.nextToken()

	var urlString string
	for !p.curTokenIs(rparen) {
		urlString += p.currentTok.Literal
		p.nextToken()

		if p.curTokenIs(newline) || p.curTokenIs(eof) {
			return &fragment{String: "[" + displayText + "](" + urlString}
		}
	}

	return &link{Properties: properties, Url: urlString, Text: displayText}
}

func (p *parser) parseShortLink(properties []property, closing tokenType) component {
	p.nextToken()

	var urlString string
	for !p.curTokenIs(gt) {
		urlString += p.currentTok.Literal
		p.nextToken()

		if p.curTokenIs(newline) || p.curTokenIs(eof) {
			return &fragment{String: "<" + urlString}
		}
	}

	return &link{Properties: properties, Url: urlString, Text: urlString}
}

func (p *parser) parseButton(properties []property, closing tokenType) component {
	p.nextToken()
	p.nextToken()

	var buttonLabel string
	for !p.curTokenIs(rbracket) {
		buttonLabel += p.currentTok.Literal
		p.nextToken()

		if p.curTokenIs(newline) || p.curTokenIs(eof) {
			return &fragment{String: "~[" + buttonLabel}
		}
	}

	if !p.peekTokenIs(lparen) {
		return &fragment{String: "~[" + buttonLabel + "]"}
	}

	p.nextToken()
	p.nextToken()

	var onClick string
	for !p.curTokenIs(rparen) {
		onClick += p.currentTok.Literal
		p.nextToken()

		if p.curTokenIs(newline) || p.curTokenIs(eof) {
			return &fragment{String: "~[" + buttonLabel + "](" + onClick}
		}
	}

	return &button{Properties: properties, OnClick: onClick, Child: &fragment{String: buttonLabel}}
}

func (p *parser) parseNav(properties []property, closing tokenType) component {
	children := make([]component, 0)

	p.nextToken()
	components, err := p.parse(at)
	if err != nil {
		panic(err.Error())
	}

	for _, component := range components {
		// don't put line breaks in nav element
		if _, ok := component.(*lineBreak); !ok {
			children = append(children, component)
		}
	}

	return &nav{Properties: properties, Children: children}
}

func (p *parser) parseSpan(properties []property, closing tokenType) component {
	if p.peekTokenIs(newline) || p.peekTokenIs(eof) {
		content := p.parseTextLine(closing)
		return &fragment{String: content}
	}

	p.nextToken()
	var spanString string

	for !p.curTokenIs(dollar) {
		spanString += p.currentTok.Literal
		p.nextToken()

		if p.curTokenIs(newline) || p.curTokenIs(eof) {
			fragment := &fragment{String: spanString}
			prefixFragment(fragment, "$", closing)
			return fragment
		}
	}

	children := []component{&fragment{String: strings.TrimSpace(spanString)}}

	p.nextToken()
	return &span{Properties: properties, Children: children}
}

func (p *parser) parseCodeBlock(properties []property, closing tokenType) component {
	p.nextToken()
	p.nextToken()

	var codeBlockString string
	for !(p.curTokenIs(caret) && p.peekTokenIs(caret)) {
		codeBlockString += p.currentTok.Literal
		p.nextToken()

		if p.curTokenIs(eof) {
			fragment := &fragment{String: "^^" + codeBlockString}
			return fragment
		}
	}
	p.nextToken()

	codeBlockString = strings.ReplaceAll(codeBlockString, "\\t", "    ")
	codeBlockString = strings.TrimPrefix(codeBlockString, "\\n")
	codeBlockString = strings.TrimSuffix(codeBlockString, "\\n")
	return &codeBlock{Properties: properties, Content: codeBlockString}
}

func (p *parser) parseComment() {
	for !p.curTokenIs(newline) {
		p.nextToken()
	}
}
