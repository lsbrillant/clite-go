package ast

import (
//	"local/clite/token"
)

type Node interface {
	//	node()
}

// Expression = VariableRef | Value | Biniary | Unary
type Expr interface {
	Node
	exprNode()
}

// Statement = Skip | Block | Assignment | Conditional | Loop
type Stmt interface {
	Node
	stmtNode()
}

// Declaration = VariableDecl | ArrayDecl
type Decl interface {
	Node
	declNode()
}

// Other
type (
	// Program = Declarations Statements
	Program struct {
		DecPart []Decl
		Body    []Stmt
	}
)

func (n *Program) node() {}

// Expressions
type (
	Variable string
	Binary   struct {
		Op    string
		Term1 Expr
		Term2 Expr
	}
	Unary struct {
		Op   string
		Term Expr
	}
)

func (n Variable) node() {}
func (n *Binary) node()  {}
func (n *Unary) node()   {}

func (n Variable) exprNode() {}
func (n *Binary) exprNode()  {}
func (n *Unary) exprNode()   {}

// Statements
type (
	Conditional struct {
		Test Expr
		Body Stmt
		Else Stmt
	}
	Loop struct {
		Test Expr
		Body Stmt
	}
	Assignment struct {
		Target Variable
		Source Expr
	}
	Block struct {
		Members []Stmt
	}
	Skip struct{}
)

func (n *Conditional) node() {}
func (n *Loop) node()        {}
func (n *Assignment) node()  {}
func (n *Block) node()       {}
func (n *Skip) node()        {}

func (n *Conditional) stmtNode() {}
func (n *Loop) stmtNode()        {}
func (n *Assignment) stmtNode()  {}
func (n *Block) stmtNode()       {}
func (n *Skip) stmtNode()        {}

// Declerations
type (
	// Declaration = VariableDecl | ArrayDecl
	// VariableDecl = Variable Type
	VariableDecl struct {
		Var Variable
		T   Type
	}
)

func (n *VariableDecl) node()     {}
func (n *VariableDecl) declNode() {}
