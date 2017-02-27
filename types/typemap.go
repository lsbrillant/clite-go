package types

import (
	"fmt"

	. "github.com/mentalpumkins/clite-go/ast"
)

type TypeMap map[Variable]Type

// Gives the typing of a Program.
// Creating a new TypeMap out of an ast.Program
// Typing Returns a new TypeMap of the typing of
// the suplied program.
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
				return nil, DuplicateDeclerationError(string(val))
			}
		}
	}
	return &tm, nil
}

func (tm *TypeMap) IsTypeCorrect(node Node) bool {
	switch n := node.(type) {
	case *Program:
	case *Block:
		//A Block is valid if all of its Statements are valid.
		valid := true
		for _, s := range n.Members {
			valid = valid && tm.IsTypeCorrect(s)
		}
		return valid
	case *Conditional:
		//A Conditional is valid if its test Expression is valid and has type bool, and
		//both its thenbranch and elsebranch Statements are valid.
		if valid := tm.IsTypeCorrect(n.Test); !valid {
			return false
		}
		if tm.typeOf(n.Test) != BOOL_TYPE {
			return false
		}
	case *Loop:
		//A Loop is valid if its test Expression is valid and has type bool, and its
		//Statement body is valid.
		if valid := tm.IsTypeCorrect(n.Test); !valid {
			return false
		}
		if tm.typeOf(n.Test) != BOOL_TYPE {
			return false
		}
		if valid := tm.IsTypeCorrect(n.Body); !valid {
			return false
		}
	case *Assignment:
		//An Assignment is valid !fall the following are true:
		//	(a) its target Variable is declared.
		if _, ok := (*tm)[n.Target]; !ok {
			return false
		}
		//	(b) Its source Expression is valid.
		if valid := tm.IsTypeCorrect(n.Source); !valid {
			return false
		}
		targetType := tm.typeOf(n.Target)
		sourceType := tm.typeOf(n.Source)
		switch targetType {
		case FLOAT_TYPE:
			//(c) If the type of its target Variable is float, then the type of its source
			//    Expression must he either float or int
			return sourceType == FLOAT_TYPE || sourceType == INT_TYPE
		case INT_TYPE:
			//(d) Otherwise, if the type of its target Variable is int, then the type of its
			//    source Expression must be either int or char.
			return sourceType == INT_TYPE || sourceType == CHAR_TYPE
		default:
			//(e) Otherwise, the type of its target Variable must be the same as the type of
			//    its source Expression.
			return targetType == sourceType
		}
	case *Binary:
		//A Binary is valid if all the following are true:
		//	(a) Its Expressions terml and term2 are valid.
		if !tm.IsTypeCorrect(n.Term1) || !tm.IsTypeCorrect(n.Term2) {
			return false
		}
		t1 := tm.typeOf(n.Term1)
		t2 := tm.typeOf(n.Term2)
		switch n.Op {
		case "-", "+", "*", "/":
			//(b) If its BinaryOp op is arithmetic ( +, - , *, /), then both its Expressions
			//    must be either int or float.
			if t1 != FLOAT_TYPE && t1 != INT_TYPE {
				return false
			}
			if t2 != FLOAT_TYPE && t2 != INT_TYPE {
				return false
			}
		case "==", "!=", "<", "<=", ">", ">=":
			//(c) if op is relational(==, !=, <. <=. >, >=),then both its Expressions must
			//    have the same type.
			return t1 == t2
		case "&&", "||":
			//(d) If op is boolean ( &&, || ), then both its Expressions must be bool.
			return t1 == BOOL_TYPE && t2 == BOOL_TYPE
		}
	case *Unary:
		//A Unary is valid if all the following are true:
		//	(a) Its Expression term is valid.
		if valid := tm.IsTypeCorrect(n.Term); !valid {
			return false
		}
		switch n.Op {
		case "!":
			// (b) If its UnaryOp op is !, then term must be bool.
			return tm.typeOf(n.Term) == BOOL_TYPE
		case "-":
			// (c) If op is -, then term must be int or float.
			t := tm.typeOf(n.Term)
			return t == FLOAT_TYPE || t == INT_TYPE
		case "float", "char":
			// (d) If op is the type conversion float() or char(), then term must be int.
			return tm.typeOf(n.Term) == INT_TYPE
		case "int":
			// (e) If op is the type conversion int(), then term must be float or char.
			t := tm.typeOf(n.Term)
			if t != FLOAT_TYPE && t != CHAR_TYPE {
				return false
			}
		}
	case Variable:
		//A Variable is valid if its id appears in the type map.
		_, ok := (*tm)[n]
		return ok
	case *Skip:
		// A Skip is always valid.
	case *VariableDecl:
		// do nothing for now.
		// if I add in decleration initilization
		// this will be the place to check.
	case Value:
		//A Value is valid.
	}
	return true
}

func (tm *TypeMap) typeOf(exp Expr) (t Type) {
	// Hooray for the go switch!
	switch e := exp.(type) {
	case Variable:
		//if the Expression is a Variable, then its result type is the type of that Variable.
		var ok bool
		t, ok = (*tm)[e]
		if !ok {
			// TODO actual error handleing
			//panic("Undifined Variable ref")
		}
	case *Binary:
		//If the Expression is a Binary, then:
		switch e.Op {
		// Arithmentic Ops
		case "-", "+", "*", "/":
			//(a) If the Operator is arithmetic ( +, -, *,or I) then its result type is the type of
			//    its operands. For example, the Expression x+ 1 requires x to be int (since
			//    1 is an int), so its result type is int.
			t1 := tm.typeOf(e.Term1)
			t2 := tm.typeOf(e.Term2)
			// Promote in all cases
			if t1 == FLOAT_TYPE || t2 == FLOAT_TYPE {
				t = FLOAT_TYPE
			} else {
				t = INT_TYPE
			}
		//	(b) If the Operator is relational ( <, <=, >, >=, ==, !=) or boolean
		//	    (&&,II). then its result type is bool.
		// Boolean Ops
		case "&&", "||":
			t = BOOL_TYPE
		// Relational Ops
		case "<=", ">=", ">", "<", "==":
			t = BOOL_TYPE
		}
	case *Unary:
		//if the Expression is a Unary, then:
		switch e.Op {
		case "!":
			//(a) If the Operator is ! then its result type is bool.
			t = BOOL_TYPE
		case "-":
			//(b) If the Operator is - then its result type is the type of its operand.
			t = tm.typeOf(e.Term)
		//(c) if the Operator is a type conversion, then the result type is given by the
		//    conversion.
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
		//if the Expression is a Value, then its result type is the type of that Value.
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

func (tm *TypeMap) String() string {
	s := "TypeMap {\n"
	for key, val := range *tm {
		s += fmt.Sprintf("  %s : %s,\n", key, val)
	}
	s += "}"
	return s
}
