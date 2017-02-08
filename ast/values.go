package ast

import (
	"fmt"
	"strconv"
)

/* Values */
type Value interface {
	Expr
	GetType() Type
	//IsUndef() bool
	GetValue() interface{}
}

/* Value = Inval | CharVal | FloatVal | BoolVal */
type (
	IntVal   int
	CharVal  rune
	FloatVal float64
	BoolVal  bool
)

type Type int

const (
	INT_TYPE Type = iota
	CHAR_TYPE
	FLOAT_TYPE
	BOOL_TYPE
)

var typeNameLiterals = [...]string{
	INT_TYPE:   "int",
	CHAR_TYPE:  "char",
	FLOAT_TYPE: "float",
	BOOL_TYPE:  "bool",
}

func (t Type) String() string {
	s := ""
	if 0 <= t && t < Type(len(typeNameLiterals)) {
		s = typeNameLiterals[t]
	}
	if s == "" {
		s = "type(" + strconv.Itoa(int(t)) + ")"
	}
	return s
}

func (n Type) node() {}

func (i IntVal) GetType() Type   { return INT_TYPE }
func (c CharVal) GetType() Type  { return CHAR_TYPE }
func (f FloatVal) GetType() Type { return FLOAT_TYPE }
func (b BoolVal) GetType() Type  { return BOOL_TYPE }

func (i IntVal) GetValue() interface{}   { return int(i) }
func (c CharVal) GetValue() interface{}  { return rune(c) }
func (f FloatVal) GetValue() interface{} { return float64(f) }
func (b BoolVal) GetValue() interface{}  { return bool(b) }

func (i IntVal) String() string   { return fmt.Sprintf("%d", i) }
func (c CharVal) String() string  { return fmt.Sprintf("%c", c) }
func (f FloatVal) String() string { return fmt.Sprintf("%f", f) }
func (b BoolVal) String() string  { return fmt.Sprintf("%s", b) }

// Values are nodes
func (n IntVal) node()   {}
func (n CharVal) node()  {}
func (n FloatVal) node() {}
func (n BoolVal) node()  {}

// Values are Expressions
func (n IntVal) exprNode()   {}
func (n CharVal) exprNode()  {}
func (n FloatVal) exprNode() {}
func (n BoolVal) exprNode()  {}
