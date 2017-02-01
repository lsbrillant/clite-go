package types

import (
	"fmt"

	. "github.com/mentalpumkins/clite-go/ast"
)

type DuplicateDeclerationError Variable

func (e DuplicateDeclerationError) Error() string {
	return fmt.Sprintf("Duplicate decleration %s", e)
}

func StaticCheck(prog Program) bool {
	//tm := Typing(prog)
	return true
}
