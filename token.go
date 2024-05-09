package mdx

import (
	"fmt"
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	HASH        = "#"
	ASTERISK    = "*"
	BACKTICK    = "`"
	BANG        = "!"
	LT          = "<"
	GT          = ">"
	DASH        = "-"
	UNDERSCORE  = "_"
	LISTELEMENT = "x."

	LBRACKET = "["
	RBRACKET = "]"
	LPAREN   = "("
	RPAREN   = ")"
	LSQUIRLY = "{"
	RSQUIRLY = "}"

	EQUALS    = "="
	TIDLE     = "~"
	DOLLAR    = "$"
	CARET     = "^"
	DOT       = "."
	SLASH     = "/"
	AT        = "@"
	BACKSLASH = "\\"
	NEWLINE   = "\\n"
	TAB       = "\\t"
	SPACE     = "SPACE"

	WORD = "WORD"
)

func NewToken(t TokenType, literal string) Token {
	return Token{Type: t, Literal: literal}
}

func (t *Token) String() string {
	return fmt.Sprintf("Token [ Type=%s, Literal=%s ]\n", t.Type, t.Literal)
}
