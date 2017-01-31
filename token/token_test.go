package token

import (
	"testing"
)

func TestLookup(t *testing.T) {
	tok := Lookup("while")
	if tok != WHILE {
		t.Errorf("%s not WHILE", tok)
	}
	tok = Lookup("someVar")
	if tok != IDENTIFIER {
		t.Errorf("%s not IDENTIFIER", tok)
	}
}
func TestIsLiteral(t *testing.T) {
	if !Token(literal_beg + 1).IsLiteral() {
		t.Error("literal isnt literal and should be")
	}
	if Token(literal_beg - 1).IsLiteral() {
		t.Error("non-literal is literal when shouldnt be")
	}
}
