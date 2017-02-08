package parser

import (
	"fmt"
	"os"
	"strconv"

	"github.com/mentalpumkins/clite-go/ast"
	"github.com/mentalpumkins/clite-go/lexer"
	"github.com/mentalpumkins/clite-go/token"
)

type Parser struct {
	tok token.Token
	lit string
	lex lexer.Lexer
	pos token.Position
}

func (p *Parser) Init(l lexer.Lexer) {
	p.lex = l
	p.nextTok()
}

func (p *Parser) Program() *ast.Program {
	boilerPlate := []token.Token{ // int main ()
		token.INT, token.MAIN, token.LEFTPAREN, token.RIGHTPAREN,
	}
	for _, t := range boilerPlate {
		p.match(t)
	}

	p.match(token.LEFTBRACE)

	decs := p.declarations()

	stmts := p.statements()

	p.match(token.RIGHTBRACE)

	return &ast.Program{decs, stmts}
}

func (p *Parser) declarations() []ast.Decl {
	var decls []ast.Decl = make([]ast.Decl, 5)
	for isType(p.tok) {
		t := p.sType()
		for p.tok != token.SEMICOLON {
			name := ast.Variable(p.identifier())
			dec := &ast.VariableDecl{name, t}
			decls = append(decls, ast.Decl(dec))
			if p.tok == token.COMMA {
				p.match(token.COMMA)
			}
		}
		p.match(token.SEMICOLON)
	}
	return decls
}

func (p *Parser) identifier() string {
	s := p.lit
	p.match(token.IDENTIFIER)
	return s
}

func (p *Parser) statements() []ast.Stmt {
	s := make([]ast.Stmt, 3)
	for p.tok != token.RIGHTBRACE {
		s = append(s, p.statement())
	}
	return s
}

func (p *Parser) block() *ast.Block {
	b := &ast.Block{}

	fmt.Println("block")
	p.match(token.LEFTBRACE)

	b.Members = p.statements()

	p.match(token.RIGHTBRACE)
	return b
}

func (p *Parser) statement() (s ast.Stmt) {
	switch p.tok {
	case token.LEFTBRACE:
		s = p.block()
	case token.IF:
		s = p.ifstmt()
	case token.WHILE:
		s = p.loop()
	case token.IDENTIFIER:
		s = p.assignment()
	case token.SEMICOLON:
		p.match(token.SEMICOLON)
		s = &ast.Skip{}
	default:
		p.error(fmt.Sprintf("expecting stmt found %s", p.tok))
	}
	return
}

func (p *Parser) assignment() *ast.Assignment {
	v := ast.Variable(p.identifier())
	p.match(token.ASSIGN)
	return &ast.Assignment{v, p.expression()}
}

func (p *Parser) ifstmt() (c *ast.Conditional) {
	p.match(token.IF)
	p.match(token.LEFTPAREN)

	e := p.expression()

	p.match(token.RIGHTPAREN)

	s := p.statement()

	if p.tok == token.ELSE {
		p.match(p.tok)
		s1 := p.statement()
		c = &ast.Conditional{e, s, s1}
	} else {
		c = &ast.Conditional{e, s, nil}
	}
	return
}

func (p *Parser) loop() *ast.Loop {
	p.match(token.WHILE)
	p.match(token.LEFTPAREN)

	e := p.expression()

	p.match(token.RIGHTPAREN)

	s := p.statement()

	return &ast.Loop{e, s}
}

func (p *Parser) sType() ast.Type {
	var t ast.Type
	switch p.tok {
	case token.INT:
		t = ast.INT_TYPE
	case token.FLOAT:
		t = ast.FLOAT_TYPE
	case token.CHAR:
		t = ast.CHAR_TYPE
	case token.BOOL:
		t = ast.BOOL_TYPE
	default:
		p.error("Expecting Type")
	}
	p.match(p.tok)
	return t
}
func (p *Parser) expression() ast.Expr {
	e := p.conjunction()
	for p.tok == token.OR {
		op := p.lit
		p.match(p.tok)
		term2 := p.conjunction()
		e = &ast.Binary{op, e, term2}
	}
	return e
}
func (p *Parser) conjunction() ast.Expr {
	e := p.equality()
	for p.tok == token.AND {
		op := p.lit
		p.match(p.tok)
		term2 := p.equality()
		e = &ast.Binary{op, e, term2}
	}
	return e
}
func (p *Parser) equality() ast.Expr {
	e := p.relation()
	if isRelOp(p.tok) {
		op := p.lit
		p.match(p.tok)
		term2 := p.relation()
		e = &ast.Binary{op, e, term2}
	}
	return e
}
func (p *Parser) relation() ast.Expr {
	e := p.addition()
	if isRelOp(p.tok) {
		op := p.lit
		p.match(p.tok)
		term2 := p.addition()
		e = &ast.Binary{op, e, term2}
	}
	return e
}
func (p *Parser) addition() ast.Expr {
	e := p.term()
	for isAddOp(p.tok) {
		op := p.lit
		p.match(p.tok)
		term2 := p.term()
		e = &ast.Binary{op, e, term2}
	}
	return e
}
func (p *Parser) term() ast.Expr {
	e := p.factor()
	for isMulOp(p.tok) {
		op := p.lit
		p.match(p.tok)
		term2 := p.factor()
		e = &ast.Binary{op, e, term2}
	}
	return e
}
func (p *Parser) factor() ast.Expr {
	if isUnaryOp(p.tok) {
		op := p.lit
		p.match(p.tok)
		term := p.primary()
		return &ast.Unary{op, term}
	} else {
		return p.primary()
	}
}
func (p *Parser) primary() ast.Expr {
	var e ast.Expr
	switch t := p.tok; {
	case t == token.IDENTIFIER:
		e = ast.Variable(p.lit)
		p.match(token.IDENTIFIER)
	case isLiteral(t):
		e = p.literal()
	case t == token.LEFTPAREN:
		p.match(token.LEFTPAREN)
		e = p.expression()
		p.match(token.RIGHTPAREN)
	case isType(t):
		op := p.lit
		p.match(p.tok)
		p.match(token.LEFTPAREN)
		term := p.expression()
		p.match(token.RIGHTPAREN)
		e = &ast.Unary{op, term}
	default:
		p.error("Expecting primary expression")
	}
	return e
}

func (p *Parser) literal() ast.Value {
	var t ast.Value
	switch p.tok {
	case token.INTLITERAL:
		v, _ := strconv.ParseInt(p.lit, 10, 32)
		t = ast.IntVal(v)
	case token.FLOATLITERAL:
		v, _ := strconv.ParseFloat(p.lit, 64)
		t = ast.FloatVal(v)
	case token.TRUE, token.FALSE:
		v, _ := strconv.ParseBool(p.lit)
		t = ast.BoolVal(v)
	case token.CHARLITERAL:
		t = ast.CharVal(p.lit[0])
	default:
		p.error("Expecting Literal")
	}
	p.nextTok()
	return t
}

func (p *Parser) nextTok() token.Token {
	p.pos, p.tok, p.lit = p.lex.Lex()
	return p.tok
}

func (p *Parser) match(t token.Token) {
	if p.tok != t {
		p.error(fmt.Sprintf("Expecting %s found %s\n", t, p.tok))
	} else {
		p.nextTok()
	}
}

func (p *Parser) error(msg string) {
	fmt.Fprintf(os.Stderr, "error %s at %s\n", msg, p.pos)
}

func isLiteral(t token.Token) bool {
	return t == token.INTLITERAL ||
		t == token.FLOATLITERAL ||
		t == token.CHARLITERAL ||
		t == token.TRUE ||
		t == token.FALSE
}

func isAddOp(t token.Token) bool {
	return t == token.PLUS || t == token.MINUS
}
func isMulOp(t token.Token) bool {
	return t == token.MULTIPLY || t == token.DIVIDE
}
func isUnaryOp(t token.Token) bool {
	return t == token.NOT || t == token.MINUS
}
func isRelOp(t token.Token) bool {
	return t == token.LESS ||
		t == token.LESSEQUAL ||
		t == token.GREATER ||
		t == token.GREATEREQUAL
}
func isType(t token.Token) bool {
	return t == token.INT ||
		t == token.BOOL ||
		t == token.CHAR ||
		t == token.FLOAT
}
