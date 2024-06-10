package mdx

import (
	"bytes"
)

type lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
	prevToken    tokenType
}

func newLexer(input string) *lexer {
	l := &lexer{input: input}
	l.readChar()
	return l
}

func (l *lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *lexer) nextToken() token {
	var tok token

	switch l.ch {
	case '#':
		tok = newToken(hash, string(l.ch))
	case '{':
		tok = newToken(lsquirly, string(l.ch))
	case '}':
		tok = newToken(rsquirly, string(l.ch))
	case '(':
		tok = newToken(lparen, string(l.ch))
	case ')':
		tok = newToken(rparen, string(l.ch))
	case '[':
		tok = newToken(lbracket, string(l.ch))
	case ']':
		tok = newToken(rbracket, string(l.ch))
	case '!':
		tok = newToken(bang, string(l.ch))
	case '\t':
		tok = newToken(tab, "\\t")
	case '\n':
		tok = newToken(newline, "\\n")
	case '`':
		tok = newToken(backtick, string(l.ch))
	case '*':
		tok = newToken(asterisk, string(l.ch))
	case '<':
		tok = newToken(lt, string(l.ch))
	case '>':
		tok = newToken(gt, string(l.ch))
	case '.':
		tok = newToken(dot, string(l.ch))
	case '-':
		tok = newToken(dash, string(l.ch))
	case '_':
		tok = newToken(underscore, string(l.ch))
	case '/':
		tok = newToken(slash, string(l.ch))
	case '\\':
		tok = newToken(backslash, string(l.ch))
	case '@':
		tok = newToken(at, string(l.ch))
	case '=':
		tok = newToken(equals, string(l.ch))
	case '~':
		tok = newToken(tidle, string(l.ch))
	case '$':
		tok = newToken(dollar, string(l.ch))
	case '^':
		tok = newToken(caret, string(l.ch))
	case ' ':
		tok = newToken(space, string(l.ch))
	case 0:
		tok = newToken(eof, "")
	default:
		if isDigit(l.ch) {
			numberBuffer := bytes.Buffer{}
			numberBuffer.WriteByte(l.ch)
			for isDigit(l.peekChar()) {
				l.readChar()
				numberBuffer.WriteByte(l.ch)
			}

			if l.peekChar() == '.' {
				if l.prevToken == newline || l.prevToken == "" {
					tok = newToken(listelement, numberBuffer.String()+".")
					l.readChar()
					l.readChar()
				} else {
					l.readChar()
					tok = newToken(word, numberBuffer.String())
				}
			} else {
				l.readChar()
				wordToken := l.readWord()
				tok = newToken(word, numberBuffer.String()+wordToken)
			}
		} else {
			wordToken := l.readWord()
			tok = newToken(word, wordToken)
		}
		l.prevToken = tok.Type
		return tok
	}

	l.readChar()
	l.prevToken = tok.Type
	return tok
}

func (l *lexer) peekChar() byte {
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
	return ch == ']' || ch == ')' || ch == '>' || ch == '*' || ch == '`' || ch == '$' || ch == '^'
}

func (l *lexer) readWord() string {
	position := l.position
	for !isWhitespace(l.ch) && !isClosingPair(l.ch) && l.ch != '=' && l.ch != 0 {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch == '-'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
