package mdx

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) NextToken() Token {
	var tok Token

	switch l.ch {
	case '#':
		tok = NewToken(HASH, string(l.ch))
	case '{':
		tok = NewToken(LSQUIRLY, string(l.ch))
	case '}':
		tok = NewToken(RSQUIRLY, string(l.ch))
	case '(':
		tok = NewToken(LPAREN, string(l.ch))
	case ')':
		tok = NewToken(RPAREN, string(l.ch))
	case '[':
		tok = NewToken(LBRACKET, string(l.ch))
	case ']':
		tok = NewToken(RBRACKET, string(l.ch))
	case '!':
		tok = NewToken(BANG, string(l.ch))
	case '\t':
		tok = NewToken(TAB, "\\t")
	case '\n':
		tok = NewToken(NEWLINE, "\\n")
	case '`':
		tok = NewToken(BACKTICK, string(l.ch))
	case '*':
		tok = NewToken(ASTERISK, string(l.ch))
	case '<':
		tok = NewToken(LT, string(l.ch))
	case '>':
		tok = NewToken(GT, string(l.ch))
	case '.':
		tok = NewToken(DOT, string(l.ch))
	case '-':
		tok = NewToken(DASH, string(l.ch))
	case '_':
		tok = NewToken(UNDERSCORE, string(l.ch))
	case '/':
		tok = NewToken(SLASH, string(l.ch))
	case '\\':
		tok = NewToken(BACKSLASH, string(l.ch))
	case '@':
		tok = NewToken(AT, string(l.ch))
	case '=':
		tok = NewToken(EQUALS, string(l.ch))
	case '~':
		tok = NewToken(TIDLE, string(l.ch))
	case '$':
		tok = NewToken(DOLLAR, string(l.ch))
	case '^':
		tok = NewToken(CARET, string(l.ch))
	case ' ':
		tok = NewToken(SPACE, string(l.ch))
	case 0:
		tok = NewToken(EOF, "")
	default:
		if isDigit(l.ch) && l.peekChar() == '.' {
			// TODO: This only works for single digits
			// maybe change this later in case some psycho starts the list as an insane number
			tok = NewToken(LISTELEMENT, string(l.ch)+".")
			l.readChar()
			l.readChar()
		} else {
			word := l.readWord()
			tok = NewToken(WORD, word)
		}
		return tok
	}

	l.readChar()
	return tok
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func isClosingPair(ch byte) bool {
	return ch == ']' || ch == ')' || ch == '>' || ch == '*' || ch == '`'
}

func (l *Lexer) readWord() string {
	position := l.position
	for !isWhitespace(l.ch) && !isClosingPair(l.ch) && l.ch != '=' {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch == '-'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
