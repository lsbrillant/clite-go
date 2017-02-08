package lexer

import (
	"github.com/mentalpumkins/clite-go/token"
	"testing"
)

type test struct {
	source  string
	tokType token.Token
	lit     string
}

var tests = [...]test{
	{"alphabet", token.IDENTIFIER, "alphabet"},
	{"=", token.ASSIGN, "="},
	{"==", token.EQUALS, "=="},
	{"||", token.OR, "||"},
	{"&&", token.AND, "&&"},
	{"12345", token.INTLITERAL, "12345"},
	{"3.14159", token.FLOATLITERAL, "3.14159"},
	{" \n \talphabet", token.IDENTIFIER, "alphabet"},
	{";", token.SEMICOLON, ""},
	{"'a'", token.CHARLITERAL, "a"},
	{"// some comments \n alpha", token.IDENTIFIER, "alpha"},
}

func TestLexer(ts *testing.T) {
	l := Lexer{}
	for _, t := range tests {
		l.Init([]byte(t.source))
		_, tok, lit := l.Lex()
		if tok != t.tokType || lit != t.lit {
			ts.Errorf("token %s with name %s", tok, lit)
		}
	}
}
