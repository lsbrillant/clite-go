package types

import (
	"fmt"
	"os"

	. "github.com/mentalpumkins/clite-go/ast"
)

type DuplicateDeclerationError string

func (e DuplicateDeclerationError) Error() string {
	return fmt.Sprintf("Duplicate decleration %s", Variable(e))
}

/*
func StaticCheck(prog Program) bool {
	//tm := Typing(prog)
	return true
}
*/

type ErrorHandler func(string, ...interface{})

type TypeChecker struct {
	tm *TypeMap

	ErrCount int
	err      ErrorHandler
}

func (tc *TypeChecker) Init(prog *Program) error {
	var err error
	tc.tm, err = Typing(prog)
	if err != nil {
		return err
	}
	tc.err = func(s string, args ...interface{}) {
		fmt.Fprintf(os.Stderr, s, args...)
	}
	tc.ErrCount = 0
	return nil
}

func (tc *TypeChecker) Visit(node Node) Visitor {
	ans := tc.tm.IsTypeCorrect(node)
	if !ans {
		tc.error("Bad typeing for %s", node)
		return nil
	}
	// returns (copy?) itself
	return tc
}

func (tc *TypeChecker) error(msg string, args ...interface{}) {
	tc.ErrCount = tc.ErrCount + 1
	tc.err(msg, args...)
}

func (tc *TypeChecker) String() string {
	return fmt.Sprintf("TypeChecker Map: %s ErrCount: %d", tc.tm, tc.ErrCount)
}
