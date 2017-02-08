package types

import (
	"testing"

	. "github.com/mentalpumkins/clite-go/ast"
)

var staticCheckTestCases = [...]struct {
	program *Program
	numErr  int
}{
	{
		&Program{
			[]Decl{
				&VariableDecl{
					Variable("a"),
					INT_TYPE,
				},
				&VariableDecl{
					Variable("b"),
					FLOAT_TYPE,
				},
			},
			[]Stmt{
				&Assignment{
					Variable("a"),
					IntVal(1),
				},
				&Assignment{
					Variable("b"),
					FloatVal(2.0),
				},
				&Assignment{
					Variable("b"),
					&Binary{
						"*",
						Variable("b"),
						Variable("a"),
					},
				},
			},
		},
		0,
	},
	{
		&Program{
			[]Decl{
				&VariableDecl{
					Variable("a"),
					INT_TYPE,
				},
			},
			[]Stmt{
				&Assignment{
					Variable("a"),
					CharVal('a'),
				},
			},
		},
		1,
	},
}

func TestStaticTypeChecker(t *testing.T) {
	var tc *TypeChecker = new(TypeChecker)
	for i, test := range staticCheckTestCases {
		tc.Init(test.program)
		// so we can make sure that tc.err is called
		var ec int
		tc.err = func(s string, args ...interface{}) {
			ec += 1
			t.Logf(s, args...)
		}
		Walk(tc, test.program)
		if test.numErr != tc.ErrCount || ec != tc.ErrCount || ec != test.numErr {
			t.Errorf("error in test %d saw %d errors, expecting %d",
				i, tc.ErrCount, test.numErr)
		}
	}
}
