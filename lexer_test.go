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
		expectedType    TokenType
		expectedLiteral string
	}{
		{HASH, "#"},
		{SPACE, " "},
		{WORD, "Heading"},
		{NEWLINE, "\\n"},
		{WORD, "Paragraph"},
		{SPACE, " "},
		{WORD, "section"},
		{NEWLINE, "\\n"},
		{NEWLINE, "\\n"},
		{LISTELEMENT, "1."},
		{SPACE, " "},
		{WORD, "List"},
		{NEWLINE, "\\n"},
		{LISTELEMENT, "2."},
		{SPACE, " "},
		{WORD, "Elements"},
		{NEWLINE, "\\n"},
		{NEWLINE, "\\n"},
		{LSQUIRLY, "{"},
		{SPACE, " "},
		{DOT, "."},
		{WORD, "class"},
		{EQUALS, "="},
		{WORD, "test"},
		{SPACE, " "},
		{RSQUIRLY, "}"},
		{NEWLINE, "\\n"},
		{NEWLINE, "\\n"},
		{LBRACKET, "["},
		{NEWLINE, "\\n"},
		{TAB, "\\t"},
		{WORD, "Div"},
		{SPACE, " "},
		{WORD, "Section"},
		{SPACE, " "},
		{WORD, "1"},
		{NEWLINE, "\\n"},
		{RBRACKET, "]"},
		{NEWLINE, "\\n"},
		{NEWLINE, "\\n"},
		{LPAREN, "("},
		{ASTERISK, "*"},
		{BANG, "!"},
		{LT, "<"},
		{GT, ">"},
		{DASH, "-"},
		{UNDERSCORE, "_"},
		{EQUALS, "="},
		{TIDLE, "~"},
		{DOLLAR, "$"},
		{CARET, "^"},
		{SLASH, "/"},
		{AT, "@"},
		{BACKSLASH, "\\"},
		{RPAREN, ")"},
		{BACKTICK, "`"},
	}

	l := NewLexer(input)

	for _, token := range expectedTokens {
		actual := l.NextToken()

		if actual.Type != token.expectedType {
			t.Fatalf("Incorrect token type. Expected=%q, got=%q", token.expectedType, actual.Type)
		}

		if actual.Literal != token.expectedLiteral {
			t.Fatalf("Incorrect token literal. Expected=%q, got=%q", token.expectedLiteral, actual.Literal)
		}
	}
}
