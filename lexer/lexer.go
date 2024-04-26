package lexer

import (
	"github.com/matt-bourke/mdx/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func New(input string) *Lexer {
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

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	switch l.ch {
	case '#':
		tok = token.New(token.HASH, string(l.ch))
	case '{':
		tok = token.New(token.LSQUIRLY, string(l.ch))
	case '}':
		tok = token.New(token.RSQUIRLY, string(l.ch))
	case '(':
		tok = token.New(token.LPAREN, string(l.ch))
	case ')':
		tok = token.New(token.RPAREN, string(l.ch))
	case '[':
		tok = token.New(token.LBRACKET, string(l.ch))
	case ']':
		tok = token.New(token.RBRACKET, string(l.ch))
	case '!':
		tok = token.New(token.BANG, string(l.ch))
	case '\t':
		tok = token.New(token.TAB, "\\t")
	case '\n':
		tok = token.New(token.NEWLINE, "\\n")
	case '`':
		tok = token.New(token.BACKTICK, string(l.ch))
	case '*':
		tok = token.New(token.ASTERISK, string(l.ch))
	case '<':
		tok = token.New(token.LT, string(l.ch))
	case '>':
		tok = token.New(token.GT, string(l.ch))
	case '.':
		tok = token.New(token.DOT, string(l.ch))
	case '-':
		tok = token.New(token.DASH, string(l.ch))
	case '_':
		tok = token.New(token.UNDERSCORE, string(l.ch))
	case '/':
		tok = token.New(token.SLASH, string(l.ch))
	case '\\':
		tok = token.New(token.BACKSLASH, string(l.ch))
	case '@':
		tok = token.New(token.AT, string(l.ch))
	case '=':
		tok = token.New(token.EQUALS, string(l.ch))
	case '~':
		tok = token.New(token.TIDLE, string(l.ch))
	case '$':
		tok = token.New(token.DOLLAR, string(l.ch))
	case '^':
		tok = token.New(token.CARET, string(l.ch))
	case ' ':
		tok = token.New(token.SPACE, string(l.ch))
	case 0:
		tok = token.New(token.EOF, "")
	default:
		if isDigit(l.ch) && l.peekChar() == '.' {
			// TODO: This only works for single digits
			// maybe change this later in case some psycho starts the list as an insane number
			tok = token.New(token.LISTELEMENT, string(l.ch)+".")
			l.readChar()
			l.readChar()
		} else {
			word := l.readWord()
			tok = token.New(token.WORD, word)
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
