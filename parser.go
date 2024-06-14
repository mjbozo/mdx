package mdx

import (
	"fmt"
	"strconv"
	"strings"
)

type parser struct {
	lex           *lexer
	previousToken token
	currentTok    token
	nextTok       token
}

func newParser(lex *lexer) *parser {
	parser := &parser{lex: lex}
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
	p.previousToken = p.currentTok
	p.currentTok = p.nextTok
	p.nextTok = p.lex.nextToken()
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

		if p.curTokenIs(delim) {
			break
		}

		p.nextToken()
	}

	return elements, nil
}

func (p *parser) parseComponent(properties []property, closing tokenType) component {
	previousToken := p.previousToken

	var element component

	switch p.currentTok.Type {
	case hash:
		element = p.parseHeader(properties, closing)
	case word:
		element = p.parseParagraph(properties, closing)
	case backtick:
		if p.peekTokenIs(backtick) {
			element = p.parseCodeDouble(properties)
		} else {
			element = p.parseCode(properties)
		}
	case asterisk:
		if p.peekTokenIs(asterisk) {
			element = p.parseStrong(properties, closing)
		} else {
			element = p.parseEm(properties, closing)
		}
	case gt:
		element, _ = p.parseBlockQuote(properties, closing, 0)
	case listelement:
		element = p.parseOrderedListElement(properties, closing)
	case dash:
		if p.peekTokenIs(space) {
			element = p.parseUnorderedList(properties, closing)
		} else if p.peekTokenIs(dash) {
			p.nextToken()
			if p.peekTokenIs(dash) {
				element = &horizontalRule{Properties: properties}
				p.nextToken()
			} else {
				element = p.parseFragment(closing)
				prefixFragment(element, "-")
			}
		} else {
			element = p.parseFragment(closing)
		}
	case bang:
		if p.peekTokenIs(lbracket) {
			element = p.parseImage(properties)
		} else {
			element = p.parseFragment(closing)
		}
	case underscore:
		if p.peekTokenIs(underscore) {
			p.nextToken()
			if p.peekTokenIs(underscore) {
				element = &horizontalRule{Properties: properties}
				p.nextToken()
			} else {
				element = p.parseFragment(closing)
				prefixFragment(element, "_")
			}
		} else {
			element = p.parseFragment(closing)
		}
	case lbracket:
		if p.peekTokenIs(space) || p.peekTokenIs(newline) {
			element = p.parseDiv(properties)
		} else {
			element = p.parseLink(properties)
		}
	case lt:
		element = p.parseShortLink(properties)
	case tidle:
		element = p.parseButton(properties)
	case at:
		element = p.parseNav(properties)
	case dollar:
		element = p.parseSpan(properties, closing)
	case caret:
		if p.peekTokenIs(caret) {
			element = p.parseCodeBlock(properties)
		} else {
			element = p.parseFragment(closing)
		}
	case newline:
		// might be easier to just not render line breaks?
		// element = &lineBreak{}
	case slash:
		if p.peekTokenIs(slash) {
			p.parseComment()
		} else {
			element = p.parseFragment(closing)
		}
	}

	// if block component, skip newlines
	if element != nil && isBlockElement(element) {
		for p.peekTokenIs(newline) {
			p.nextToken()
		}
	}

	// if line starts with an inline element and is followed by a paragraph, wrap the first inline element
	// in following paragraph
	if joinsParagraph(element) && !p.curTokenIs(newline) && p.peekTokenIs(word) && previousToken.Type == newline {
		pComponent := p.parseParagraph(nil, closing).(*paragraph)
		paragraphChildren := append([]component{element}, pComponent.Content...)
		pComponent.Content = paragraphChildren
		element = pComponent
	}

	return element
}

func joinsParagraph(comp component) bool {
	switch comp.(type) {
	case *code,
		*bold,
		*italic,
		*span,
		*button,
		*link:
		return true
	}
	return false
}

func isBlockElement(comp component) bool {
	switch comp.(type) {
	case *div,
		*codeBlock,
		*horizontalRule,
		*image,
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
				errorMessage := "Property formatted incorrectly. DOT must be followed by a WORD"
				return nil, &parseError{errorReason: errorMessage}, propsString
			}

			p.nextToken()
			propsString += p.currentTok.Literal
			key := p.currentTok.Literal

			if !p.peekTokenIs(equals) {
				errorMessage := "Property formatted incorrectly. KEY must be follwed by EQUALS"
				return nil, &parseError{errorReason: errorMessage}, propsString
			}

			p.nextToken()
			propsString += p.currentTok.Literal
			if !p.peekTokenIs(word) {
				errorMessage := "Property formatted incorrectly. EQUALS must be followed by VALUE"
				return nil, &parseError{errorReason: errorMessage}, propsString
			}

			p.nextToken()
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

func (p *parser) parseFragment(closing tokenType) *fragment {
	content := p.parseTextLine(closing)
	return &fragment{Value: content}
}

func (p *parser) parseTextLine(closing tokenType) string {
	var lineString string
	for !(p.curTokenIs(newline) || p.curTokenIs(closing)) {
		lineString += p.currentTok.Literal
		p.nextToken()
	}
	return lineString
}

// Appends a fragment containing fragmentValue to the lineElements slice after replacing '\\n' with spaces.
// Subsequently sets fragmentValue to an empty string.
func bankCurrentFragment(lineElements *[]component, fragmentValue *string) {
	if len(*fragmentValue) > 0 {
		*lineElements = append(*lineElements, &fragment{Value: strings.ReplaceAll(*fragmentValue, "\\n", " ")})
		*fragmentValue = ""
	}
}

func (p *parser) parseLine(closing tokenType) []component {
	lineElements := make([]component, 0)
	var lineString string

	for !(p.curTokenIs(newline) || p.curTokenIs(closing)) {
		if p.currentTok.IsElementToken() {
			bankCurrentFragment(&lineElements, &lineString)
			lineElements = append(lineElements, p.parseComponent(nil, closing))
		} else if p.curTokenIs(lsquirly) {
			properties, parseErr, propsText := p.parseProperties()
			if parseErr != nil {
				lineString += propsText
				p.nextToken()
			} else {
				for p.curTokenIs(tab) {
					p.nextToken()
				}
				bankCurrentFragment(&lineElements, &lineString)
				nextComponent := p.parseComponent(properties, closing)
				if nextComponent != nil {
					lineElements = append(lineElements, nextComponent)
				}
			}
		} else {
			if !(p.currentTok.Type == space && p.peekTokenIs(closing)) {
				lineString += p.currentTok.Literal
			}
			p.nextToken()
		}
	}

	bankCurrentFragment(&lineElements, &lineString)
	return lineElements
}

func (p *parser) parseBlockQuoteLine(closing tokenType) []component {
	lineElements := make([]component, 0)
	var lineString string

	for !(p.curTokenIs(newline) || p.curTokenIs(closing)) {
		if p.currentTok.IsElementToken() {
			bankCurrentFragment(&lineElements, &lineString)
			lineElements = append(lineElements, p.parseComponent(nil, closing))
		} else if p.curTokenIs(lsquirly) {
			properties, parseErr, propsText := p.parseProperties()
			if parseErr != nil {
				lineString += propsText
				p.nextToken()
			} else {
				for p.curTokenIs(tab) {
					p.nextToken()
				}
				bankCurrentFragment(&lineElements, &lineString)
				nextComponent := p.parseComponent(properties, closing)
				if nextComponent != nil {
					lineElements = append(lineElements, nextComponent)
				}
			}
		} else {
			if !(p.currentTok.Type == space && p.peekTokenIs(closing)) {
				lineString += p.currentTok.Literal
			}
			p.nextToken()
		}

		if p.curTokenIs(newline) && p.peekTokenIs(tab) {
			p.nextToken()
			for p.curTokenIs(tab) {
				p.nextToken()
			}
		}
	}

	bankCurrentFragment(&lineElements, &lineString)
	return lineElements
}

func (p *parser) parseLineDoubleClose(closing tokenType) []component {
	lineElements := make([]component, 0)
	var lineString string

	for !(p.curTokenIs(newline) || (p.curTokenIs(closing) && p.peekTokenIs(closing))) {
		if p.currentTok.IsElementToken() {
			bankCurrentFragment(&lineElements, &lineString)
			lineElements = append(lineElements, p.parseComponent(nil, closing))
		} else if p.curTokenIs(lsquirly) {
			properties, parseErr, propsText := p.parseProperties()
			if parseErr != nil {
				lineString += propsText
				p.nextToken()
			} else {
				for p.curTokenIs(tab) {
					p.nextToken()
				}
				bankCurrentFragment(&lineElements, &lineString)
				nextComponent := p.parseComponent(properties, closing)
				if nextComponent != nil {
					lineElements = append(lineElements, nextComponent)
				}
			}
		} else {
			lineString += p.currentTok.Literal
			p.nextToken()
		}

		if p.curTokenIs(newline) && p.peekTokenIs(tab) {
			p.nextToken()
			for p.curTokenIs(tab) {
				p.nextToken()
			}
		}
	}

	bankCurrentFragment(&lineElements, &lineString)
	return lineElements
}

func (p *parser) parseBlock(closing tokenType) []component {
	blockElements := make([]component, 0)
	var blockString string

	for !(p.curTokenIs(eof) || p.curTokenIs(closing) || p.isDoubleBreak() || p.isAfterNewline(closing) || p.isNextLineElement()) {
		if p.currentTok.IsElementToken() {
			bankCurrentFragment(&blockElements, &blockString)
			blockElements = append(blockElements, p.parseComponent(nil, closing))
		} else if p.curTokenIs(lsquirly) {
			properties, parseErr, propsText := p.parseProperties()
			if parseErr != nil {
				blockString += propsText
				p.nextToken()
			} else {
				for p.curTokenIs(tab) {
					p.nextToken()
				}
				bankCurrentFragment(&blockElements, &blockString)
				nextComponent := p.parseComponent(properties, closing)
				if nextComponent != nil {
					blockElements = append(blockElements, nextComponent)
				}
			}
		} else {
			if !(p.currentTok.Type == space && p.peekTokenIs(closing)) {
				blockString += p.currentTok.Literal
			}
			p.nextToken()
		}

		if p.curTokenIs(newline) && p.peekTokenIs(tab) {
			p.nextToken()
			for p.curTokenIs(tab) {
				p.nextToken()
			}
		}
	}

	bankCurrentFragment(&blockElements, &blockString)
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
		return p.parseFragment(closing)
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

func prefixFragment(comp component, prefix string) {
	switch c := (comp).(type) {
	case *fragment:
		c.Value = prefix + c.Value
	}
}

func (p *parser) parseCode(properties []property) component {
	p.nextToken()
	var codeString string

	for !p.curTokenIs(backtick) {
		if p.curTokenIs(newline) || p.curTokenIs(eof) {
			return &paragraph{Properties: properties, Content: []component{&fragment{Value: "`" + codeString}}}
		}
		codeString += p.currentTok.Literal
		p.nextToken()
	}

	p.nextToken()
	return &code{Properties: properties, Text: codeString}
}

func (p *parser) parseCodeDouble(properties []property) component {
	p.nextToken()
	p.nextToken()

	var codeText string
	for !(p.curTokenIs(eof) || (p.curTokenIs(backtick) && p.peekTokenIs(backtick))) {
		codeText += p.currentTok.Literal
		p.nextToken()
	}

	p.nextToken()
	p.nextToken()

	return &code{Properties: properties, Text: codeText}
}

func (p *parser) parseStrong(properties []property, closing tokenType) component {
	p.nextToken()
	if p.peekTokenIs(space) || p.peekTokenIs(newline) || p.peekTokenIs(eof) {
		content := p.parseTextLine(closing)
		return &paragraph{Properties: properties, Content: []component{&fragment{Value: "*" + content}}}
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
		return &paragraph{Properties: properties, Content: []component{&fragment{Value: content}}}
	}

	p.nextToken()
	content := p.parseLine(asterisk)

	p.nextToken()
	return &italic{Properties: properties, Content: content}
}

func (p *parser) parseBlockQuote(properties []property, closing tokenType, initialDepth int) (component, int) {
	content := make([]component, 0)
	depth := initialDepth

	for p.curTokenIs(gt) {
		depth += 1
		if depth > initialDepth+1 {
			nested, _ := p.parseBlockQuote(properties, closing, depth-1)
			content = append(content, nested)
			depth = initialDepth
		} else {
			p.nextToken()
		}

		for p.curTokenIs(space) {
			p.nextToken()
		}
	}

	for !(p.curTokenIs(newline) || p.curTokenIs(eof)) {
		lineContent := p.parseBlockQuoteLine(closing)
		content = append(content, lineContent...)
		p.nextToken()

		if p.curTokenIs(newline) || p.curTokenIs(eof) || p.curTokenIs(closing) {
			// double newline or eof or closing tag
			break
		}

		nextDepth := 0
		for p.curTokenIs(gt) {
			nextDepth += 1
			p.nextToken()
			for p.curTokenIs(space) {
				p.nextToken()
			}
		}

		if nextDepth < depth && p.curTokenIs(newline) {
			p.nextToken()
			return &blockQuote{Properties: properties, Content: content}, nextDepth
		}

		if nextDepth == depth && p.curTokenIs(newline) {
			content = append(content, &lineBreak{})
			p.nextToken()
			for p.curTokenIs(gt) {
				p.nextToken()
				for p.curTokenIs(space) {
					p.nextToken()
				}
			}
			continue
		}

		if nextDepth > depth {
			nested, d := p.parseBlockQuote(properties, closing, nextDepth)
			content = append(content, nested)
			if d < depth {
				return &blockQuote{Properties: properties, Content: content}, d
			}

			nextDepth = 0
			for p.curTokenIs(gt) {
				nextDepth += 1
				p.nextToken()
				for p.curTokenIs(space) {
					p.nextToken()
				}
			}
			if nextDepth < depth {
				if nextDepth != 0 && p.curTokenIs(newline) {
					p.nextToken()
				}
				return &blockQuote{Properties: properties, Content: content}, nextDepth
			}
		} else {
			content = append(content, &fragment{Value: " "})
		}
	}

	return &blockQuote{Properties: properties, Content: content}, 0
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
		element := listItem{Component: &paragraph{Content: []component{&fragment{Value: elementContent}}}}
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
		element := listItem{Component: &paragraph{Content: []component{&fragment{Value: elementContent}}}}
		listElements = append(listElements, element)
	}

	return &unorderedList{Properties: properties, ListItems: listElements}
}

func (p *parser) parseImage(properties []property) component {
	p.nextToken()
	p.nextToken()

	var altText string
	for !p.curTokenIs(rbracket) {
		altText += p.currentTok.Literal
		p.nextToken()

		if p.curTokenIs(newline) || p.curTokenIs(eof) {
			return &paragraph{Properties: properties, Content: []component{&fragment{Value: "![" + altText}}}
		}
	}

	if !p.peekTokenIs(lparen) {
		return &paragraph{Properties: properties, Content: []component{&fragment{Value: "![" + altText + "]"}}}
	}

	p.nextToken()
	p.nextToken()

	var urlString string
	for !p.curTokenIs(rparen) {
		urlString += p.currentTok.Literal
		p.nextToken()

		if p.curTokenIs(newline) || p.curTokenIs(eof) {
			fragment := &fragment{Value: "![" + altText + "](" + urlString}
			return &paragraph{Properties: properties, Content: []component{fragment}}
		}
	}

	return &image{Properties: properties, ImgUrl: urlString, AltText: altText}
}

func (p *parser) parseDiv(properties []property) component {
	p.nextToken()
	for p.curTokenIs(newline) {
		p.nextToken()
	}

	components, err := p.parse(rbracket)
	p.nextToken()
	if err != nil {
		panic(err.Error())
	}

	if p.peekTokenIs(newline) {
		p.nextToken()
	}

	return &div{Properties: properties, Children: components}
}

func (p *parser) parseLink(properties []property) component {
	p.nextToken()

	components, err := p.parse(rbracket)
	if err != nil {
		panic(err.Error())
	}

	if !p.peekTokenIs(lparen) {
		content := make([]component, 0)
		content = append(content, &fragment{Value: "["})
		content = append(content, components...)
		content = append(content, &fragment{Value: "]"})
		return &paragraph{Content: content}
	}

	p.nextToken()
	p.nextToken()

	var urlString string
	for !p.curTokenIs(rparen) {
		urlString += p.currentTok.Literal
		p.nextToken()

		if p.curTokenIs(newline) || p.curTokenIs(eof) {
			content := make([]component, 0)
			content = append(content, &fragment{Value: "["})
			content = append(content, components...)
			content = append(content, &fragment{Value: "](" + urlString})
			return &paragraph{Content: content}
		}
	}

	// if only child is a simple paragraph, replace with a fragment for cleaner output
	if len(components) == 1 {
		if p, ok := components[0].(*paragraph); ok {
			if len(p.Content) == 1 {
				if frag, ok := p.Content[0].(*fragment); ok {
					components = []component{frag}
				}
			}
		}
	}

	p.nextToken()
	return &link{Properties: properties, Url: urlString, Content: components}
}

func (p *parser) parseShortLink(properties []property) component {
	p.nextToken()

	var urlString string
	for !p.curTokenIs(gt) {
		urlString += p.currentTok.Literal
		p.nextToken()

		if p.curTokenIs(newline) || p.curTokenIs(eof) {
			return &paragraph{Properties: properties, Content: []component{&fragment{Value: "<" + urlString}}}
		}
	}

	return &link{Properties: properties, Url: urlString, Content: []component{&fragment{Value: urlString}}}
}

func (p *parser) parseButton(properties []property) component {
	p.nextToken()
	p.nextToken()

	components, err := p.parse(rbracket)
	if err != nil {
		panic(err.Error())
	}

	if !p.peekTokenIs(lparen) {
		content := []component{&fragment{Value: "~["}}
		content = append(content, components...)
		content = append(content, &fragment{Value: "]"})
		return &paragraph{Content: content}
	}

	p.nextToken()
	p.nextToken()

	var onClick string
	for !p.curTokenIs(rparen) {
		onClick += p.currentTok.Literal
		p.nextToken()

		if p.curTokenIs(newline) || p.curTokenIs(eof) {
			content := []component{&fragment{Value: "~["}}
			content = append(content, components...)
			content = append(content, &fragment{Value: "](" + onClick})
			return &paragraph{Content: content}
		}
	}

	p.nextToken()

	return &button{Properties: properties, OnClick: onClick, Content: components}
}

func (p *parser) parseNav(properties []property) component {
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
		p.nextToken()
		return &paragraph{Properties: properties, Content: []component{&fragment{Value: content}}}
	}

	p.nextToken()
	for p.curTokenIs(space) {
		p.nextToken()
	}
	content := p.parseLine(dollar)

	// remove trailing whitespace if last component is fragment
	if len(content) > 0 {
		if frag, ok := content[len(content)-1].(*fragment); ok {
			frag.Value = strings.TrimRight(frag.Value, " ")
		}
	}

	p.nextToken()
	return &span{Properties: properties, Content: content}
}

func (p *parser) parseCodeBlock(properties []property) component {
	p.nextToken()
	p.nextToken()

	var codeBlockString string
	for !(p.curTokenIs(caret) && p.peekTokenIs(caret)) {
		codeBlockString += p.currentTok.Literal
		p.nextToken()

		if p.curTokenIs(eof) {
			fragment := &fragment{Value: "^^" + codeBlockString}
			paragraph := &paragraph{Properties: properties, Content: []component{fragment}}
			return paragraph
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
