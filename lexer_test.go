package mdx

import (
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `# Heading
Paragraph section

1. List
2. Elements

{ .class=test }

[
	Div Section 1
]

(*!<>-_=~$^/@\)`

	input += "`"

	expectedTokens := []struct {
		expectedType    tokenType
		expectedLiteral string
	}{
		{hash, "#"},
		{space, " "},
		{word, "Heading"},
		{newline, "\\n"},
		{word, "Paragraph"},
		{space, " "},
		{word, "section"},
		{newline, "\\n"},
		{newline, "\\n"},
		{listelement, "1."},
		{space, " "},
		{word, "List"},
		{newline, "\\n"},
		{listelement, "2."},
		{space, " "},
		{word, "Elements"},
		{newline, "\\n"},
		{newline, "\\n"},
		{lsquirly, "{"},
		{space, " "},
		{dot, "."},
		{word, "class"},
		{equals, "="},
		{word, "test"},
		{space, " "},
		{rsquirly, "}"},
		{newline, "\\n"},
		{newline, "\\n"},
		{lbracket, "["},
		{newline, "\\n"},
		{tab, "\\t"},
		{word, "Div"},
		{space, " "},
		{word, "Section"},
		{space, " "},
		{word, "1"},
		{newline, "\\n"},
		{rbracket, "]"},
		{newline, "\\n"},
		{newline, "\\n"},
		{lparen, "("},
		{asterisk, "*"},
		{bang, "!"},
		{lt, "<"},
		{gt, ">"},
		{dash, "-"},
		{underscore, "_"},
		{equals, "="},
		{tidle, "~"},
		{dollar, "$"},
		{caret, "^"},
		{slash, "/"},
		{at, "@"},
		{backslash, "\\"},
		{rparen, ")"},
		{backtick, "`"},
	}

	l := newLexer(input)

	for _, token := range expectedTokens {
		actual := l.nextToken()

		if actual.Type != token.expectedType {
			t.Fatalf("Incorrect token type. Expected=%q, got=%q", token.expectedType, actual.Type)
		}

		if actual.Literal != token.expectedLiteral {
			t.Fatalf("Incorrect token literal. Expected=%q, got=%q", token.expectedLiteral, actual.Literal)
		}
	}
}
