package ast

// A Visitor's Visit method is invoked for each node encountered by Walk.
// If the result visitor w is not nil, Walk visits each of the children
// of node with the visitor w, followed by a call of w.Visit(nil).
type Visitor interface {
	Visit(node Node) (w Visitor)
}

// Helper functions for common node lists. They may be empty.
func walkExprList(v Visitor, list []Expr) {
	for _, x := range list {
		Walk(v, x)
	}
}

func walkStmtList(v Visitor, list []Stmt) {
	for _, x := range list {
		Walk(v, x)
	}
}

func walkDeclList(v Visitor, list []Decl) {
	for _, x := range list {
		Walk(v, x)
	}
}

func Walk(v Visitor, node Node) {
	if v = v.Visit(node); v == nil {
		return
	}

	switch n := node.(type) {
	case *Program:
		walkDeclList(v, n.DecPart)
		walkStmtList(v, n.Body)
	case *Block:
		walkStmtList(v, n.Members)
	case *Conditional:
		Walk(v, n.Test)
		Walk(v, n.Body)
		if n.Else != nil {
			Walk(v, n.Else)
		}
	case *Loop:
		Walk(v, n.Test)
		Walk(v, n.Body)
	case *Assignment:
		Walk(v, n.Target)
		Walk(v, n.Source)
	case *Binary:
		Walk(v, n.Term1)
		Walk(v, n.Term2)
	case *Unary:
		Walk(v, n.Term)
	case Variable, *Skip:
		// do nothing
	case *VariableDecl:
		Walk(v, n.Var)
		Walk(v, n.T)
	}
}

type inspector func(Node) bool

func (f inspector) Visit(node Node) Visitor {
	if f(node) {
		return f
	}
	return nil
}

// Inspect traverses an AST in depth-first order: It starts by calling
// f(node); node must not be nil. If f returns true, Inspect invokes f
// recursively for each of the non-nil children of node, followed by a
// call of f(nil).
//
func Inspect(node Node, f func(Node) bool) {
	Walk(inspector(f), node)
}
