package ast

import "fmt"

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

type Type string

func (n Type) node() {}

func (i IntVal) GetType() Type   { return Type("int") }
func (c CharVal) GetType() Type  { return Type("char") }
func (f FloatVal) GetType() Type { return Type("float") }
func (b BoolVal) GetType() Type  { return Type("bool") }

func (i IntVal) GetValue() interface{}   { return int(i) }
func (c CharVal) GetValue() interface{}  { return rune(c) }
func (f FloatVal) GetValue() interface{} { return float64(f) }
func (b BoolVal) GetValue() interface{}  { return bool(b) }

func (i IntVal) String() string   { return fmt.Sprintf("%d", i) }
func (c CharVal) String() string  { return fmt.Sprintf("%s", c) }
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
