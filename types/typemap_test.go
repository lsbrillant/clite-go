package types

import (
	"testing"

	"github.com/mentalpumkins/clite-go/ast"
	"github.com/mentalpumkins/clite-go/ast/operators"
)

var testCases = [...]struct {
	tm           TypeMap
	toType       ast.Expr
	expectedType ast.Type
}{
	{
		TypeMap(map[ast.Variable]ast.Type{
			ast.Variable("a"): ast.INT_TYPE,
			ast.Variable("b"): ast.INT_TYPE,
		}),
		ast.Variable("a"),
		ast.INT_TYPE,
	},
	{
		TypeMap(map[ast.Variable]ast.Type{
			ast.Variable("a"): ast.INT_TYPE,
			ast.Variable("b"): ast.FLOAT_TYPE,
		}),
		&ast.Binary{operators.Operator("-"), ast.Variable("a"), ast.Variable("b")},
		ast.FLOAT_TYPE,
	},
}

func TestTypeOf(t *testing.T) {
	for i, test := range testCases {
		if compType := test.tm.typeOf(test.toType); compType != test.expectedType {
			t.Errorf("error in test cast %d : expected type : %s saw type %s", i, test.expectedType, compType)
		}
	}
}
