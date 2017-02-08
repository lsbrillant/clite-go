//Package clite-go/print provides a pretty printer for clite
package print

import (
	"fmt"
	"io"

	. "github.com/mentalpumkins/clite-go/ast"
)

func PrettyPrint(w io.Writer, prog *Program) {
	pp := PrettyPrinter{
		Target:      w,
		Indent:      "  ",
		Indentlevel: 0,
	}
	Walk(pp, prog)
}

type PrettyPrinter struct {
	Target      io.Writer
	Indent      string
	Indentlevel int
}

func (p PrettyPrinter) Visit(node Node) Visitor {
	//fmt.Print(node)
	switch n := node.(type) {
	case *Program:
		p.Printi("Program:\n")
	case *Block:
		p.Print("\n")
		p.Printi("Block:\n")
	case *Conditional:
		p.Print("Conditional: ")
	case *Loop:
		p.Printi("Loop: ")
	case *Assignment:
		p.Printi("Assignment: ")
	case *Binary:
		p.Print("Binary: %s ", n.Op)
	case *Unary:
		p.Print("Unary: %s ", n.Op)
	case Value:
		switch vt := n.(type) {
		case IntVal:
			p.Print("%d (%s) ", vt, n.GetType())
		case FloatVal:
			p.Print("%s (%s) ", vt, n.GetType())
		case CharVal:
			p.Print("%s (%s) ", vt, n.GetType())
		case BoolVal:
			p.Print("%s (%s) ", vt, n.GetType())
		}
	case *Skip:
		p.Print("\n")
		return nil
	case Variable, Type:
		p.Print("%s ", n)
	case *VariableDecl:
		p.Printi("Decl: %s %s\n", n.Var, n.T)
		return nil
	}
	// set indent to one more
	return PrettyPrinter{
		p.Target,
		p.Indent,
		p.Indentlevel + 1,
	}
}

// Prints without indentation
func (p PrettyPrinter) Print(s string, args ...interface{}) {
	fmt.Fprintf(p.Target, s, args...)
}

// Prints with indentation
func (p PrettyPrinter) Printi(s string, args ...interface{}) {
	padd := ""
	// padd out the indent
	for i := 0; i < p.Indentlevel; i++ {
		padd += p.Indent
	}
	// do the print
	p.Print(padd+s, args...)
}
