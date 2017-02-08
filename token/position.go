package token

import "fmt"

// Position type taken from go/token
type Position struct {
	Offset int // offset, starting at 0
	Line   int // line number, starting at 1
	Column int // column number, starting at 1 (byte count)
}

func (pos *Position) IsValid() bool { return pos.Line > 0 }

// String returns a string in one of two forms:
//
//	line:column         valid position without file name
//	-                   invalid position without file name
//
func (pos Position) String() (s string) {
	if pos.IsValid() {
		s = fmt.Sprintf("%d:%d", pos.Line, pos.Column)
	} else {
		s = "-"
	}
	return
}
