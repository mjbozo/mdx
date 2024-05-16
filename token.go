package mdx

import (
	"fmt"
)

type tokenType string

type token struct {
	Type    tokenType
	Literal string
}

const (
	hash        = "#"
	asterisk    = "*"
	backtick    = "`"
	bang        = "!"
	lt          = "<"
	gt          = ">"
	dash        = "-"
	underscore  = "_"
	listelement = "x."

	lbracket = "["
	rbracket = "]"
	lparen   = "("
	rparen   = ")"
	lsquirly = "{"
	rsquirly = "}"

	equals    = "="
	tidle     = "~"
	dollar    = "$"
	caret     = "^"
	dot       = "."
	slash     = "/"
	at        = "@"
	backslash = "\\"
	newline   = "\\n"
	tab       = "\\t"

	space = "SPACE"
	word  = "WORD"
	eof   = "EOF"
)

func newToken(t tokenType, literal string) token {
	return token{Type: t, Literal: literal}
}

func (t *token) String() string {
	return fmt.Sprintf("Token [ Type=%s, Literal=%s ]\n", t.Type, t.Literal)
}

func (t *token) IsElementToken() bool {
	switch t.Type {
	case hash,
		asterisk,
		backtick,
		bang,
		gt,
		dash,
		listelement,
		lbracket,
		tidle,
		caret,
		dollar,
		at:
		return true
	}

	return false
}
