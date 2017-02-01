package types

import (
	. "github.com/mentalpumkins/clite-go/ast"
)

type TypeMap map[Variable]Type

// Gives the typing a Program.
// Creating a new TypeMap out of an ast.Program
func Typing(p *Program) (*TypeMap, error) {
	tm := TypeMap(make(map[Variable]Type))
	for _, decl := range p.DecPart {
		// type switch not needed right now will
		// be usefull for adding array declerations
		switch d := decl.(type) {
		case *VariableDecl:
			if val, duplicate := tm[d.Var]; !duplicate {
				tm[d.Var] = d.T
			} else {
				// You done screwed up.
				return nil, DuplicateDeclerationError(val)
			}
		}
	}
	return &tm, nil
}

func (tm *TypeMap) IsTypeCorrect(node Node) bool {
	switch n := node.(type) {
	case *Program:
	case *Block:
	case *Conditional:
		if tm.typeOf(n.Test) != BOOL_TYPE {
			return false
		}
	case *Loop:
		if tm.typeOf(n.Test) != BOOL_TYPE {
			return false
		}
	case *Assignment:
		if tm.typeOf(n.Source) != tm.typeOf(n.Target) {
			return false
		}
	case *Binary:
		// TODO - this
		//t1 := tm.typeOf(n.Term1)
		//t2 := tm.typeOf(n.Term2)
		switch n.Op {
		// Arithmentic Ops
		case "-", "+", "*", "/":
		}
	case *Unary:
		switch n.Op {
		case "-":
			t := tm.typeOf(n.Term)
			if t != FLOAT_TYPE && t != INT_TYPE {
				return false
			}
		default:
			// nothing for now
		}
	case Variable, *Skip:
		// do nothing
	case *VariableDecl:
		// do nothing for now.
		// if I add in decleration initilization
		// this will be the place to check.
	}
	return true
}

func (tm *TypeMap) typeOf(exp Expr) (t Type) {
	// Hooray for the go switch!
	switch e := exp.(type) {
	case Variable:
		var ok bool
		t, ok = (*tm)[e]
		if !ok {
			// TODO actual error handleing
			panic("Undifined Variable ref")
		}
	case *Binary:
		switch e.Op {
		// Arithmentic Ops
		case "-", "+", "*", "/":
			t1 := tm.typeOf(e.Term1)
			t2 := tm.typeOf(e.Term2)
			// Promote in all cases
			if t1 == FLOAT_TYPE || t2 == FLOAT_TYPE {
				t = FLOAT_TYPE
			} else {
				t = INT_TYPE
			}
		// Boolean Ops
		case "&&", "||":
			t = BOOL_TYPE
		// Relational Ops
		case "<=", ">=", ">", "<", "==":
			t = BOOL_TYPE
		}
	case *Unary:
		switch e.Op {
		case "-":
			t = tm.typeOf(e.Term)
		case "float":
			t = FLOAT_TYPE
		case "int":
			t = INT_TYPE
		case "char":
			t = CHAR_TYPE
		case "bool":
			t = BOOL_TYPE
		}
	case Value:
		switch e.(type) {
		case IntVal:
			t = INT_TYPE
		case BoolVal:
			t = BOOL_TYPE
		case FloatVal:
			t = FLOAT_TYPE
		case CharVal:
			t = CHAR_TYPE
		}
	}
	return
}
