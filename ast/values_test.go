package ast

import (
	"testing"
)

var tests = [...]struct {
	v    Value
	goal interface{}
}{
	{IntVal(1337), 1337},
	{CharVal('q'), 'q'},
	{BoolVal(true), true},
	{FloatVal(3.14159), 3.14159},
}

func TestValues(t *testing.T) {
	var v Value
	for _, test := range tests {
		v = test.v
		if v.GetValue() != test.goal {
			t.Error("Value Doesnt make sense")
		}
	}
}
