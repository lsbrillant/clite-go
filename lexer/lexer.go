package lexer

import (
	"fmt"
	"os"

	"github.com/mentalpumkins/clite-go/token"
)

type ErrorHandler func(pos token.Position, msg string)

var DefaultErrorHandler ErrorHandler = func(pos token.Position, msg string) {
	fmt.Fprintf(os.Stderr, "error %s at %s\n", msg, pos)
}

type Lexer struct {
	src []byte
	err ErrorHandler

	ch         rune
	offset     int
	rdOffset   int
	lineCount  int
	lineOffset int

	ErrorCount int
}

func (l *Lexer) next() {
	if l.rdOffset < len(l.src) {
		l.offset = l.rdOffset
		if l.ch == '\n' {
			l.lineCount++
			l.lineOffset = l.offset
		}
		r, w := rune(l.src[l.rdOffset]), 1
		if r == 0 {
			l.error("illegal character NUL")
		}
		l.rdOffset += w
		l.ch = r
	} else {
		l.offset = len(l.src)
		if l.ch == '\n' {
			l.lineCount++
			l.lineOffset = l.offset
		}
		l.ch = -1 // eof
	}
}

func (l *Lexer) Lex() (pos token.Position, tok token.Token, lit string) {
scanAgain:
	l.skipWhitespace()
	pos = l.Pos()
	switch ch := l.ch; {
	case isLetter(ch):
		lit = l.scanIdentifier()
		if len(lit) > 1 {
			// keywords are longer than one letter - avoid lookup otherwise
			tok = token.Lookup(lit)
		} else {
			tok = token.IDENTIFIER
		}
	case isDigit(ch):
		tok, lit = l.scanNumber()
	default: //don't strictly know why I to do two switches
		l.next() //here but it does in the go one so here gos
		switch ch {
		case -1:
			tok = token.EOF
		// punctuation
		case ';':
			tok = token.SEMICOLON
		case ',':
			tok = token.COMMA
		case '\'':
			lit = string(l.ch)
			l.next() // get the trailing "'"
			l.next() // move to space after
			tok = token.CHARLITERAL
		// arithmetic operators
		case '+':
			lit = string(ch)
			tok = token.PLUS
		case '-':
			lit = string(ch)
			tok = token.MINUS
		case '*':
			lit = string(ch)
			tok = token.MULTIPLY
		case '/':
			if l.ch == '/' {
				// in a comment
				l.next()
				for l.ch != '\n' && l.ch >= 0 {
					l.next()
				}
				goto scanAgain
			} else {
				lit = string(ch)
				tok = token.DIVIDE
			}
		// seperators
		case '{':
			tok = token.LEFTBRACE
		case '}':
			tok = token.RIGHTBRACE
		case '(':
			tok = token.LEFTPAREN
		case ')':
			tok = token.RIGHTPAREN
		case '[':
			tok = token.LEFTBRACKET
		case ']':
			tok = token.RIGHTBRACKET

		case '&':
			l.match('&')
			tok = token.AND
			lit = tok.String()
		case '|':
			l.match('|')
			tok = token.OR
			lit = tok.String()

		case '=':
			tok, lit = l.switch2(token.ASSIGN, token.EQUALS)
		case '<':
			tok, lit = l.switch2(token.LESS, token.LESSEQUAL)
		case '>':
			tok, lit = l.switch2(token.GREATER, token.GREATEREQUAL)
		case '!':
			tok, lit = l.switch2(token.NOT, token.NOTEQUAL)
		default:
			l.error("Illegal character")
		}
	}
	return
}

func (l *Lexer) Init(src []byte) {
	l.src = src

	l.ch = ' '
	l.offset = 0
	l.rdOffset = 0
	l.lineCount = 1
	l.lineOffset = 0
	l.ErrorCount = 0

	l.err = DefaultErrorHandler
}

func (l *Lexer) Pos() token.Position {
	return token.Position{l.offset, l.lineCount, (l.offset - l.lineOffset)}
}

func (l *Lexer) error(msg string) {
	l.ErrorCount++
	l.err(l.Pos(), msg)
}

func (l *Lexer) AtEof() bool {
	return l.ch == -1
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.next()
	}
}

func (l *Lexer) scanNumber() (token.Token, string) {
	offs := l.offset
	tok := token.INTLITERAL
	for isDigit(l.ch) {
		l.next()
	}
	if l.ch == '.' {
		tok = token.FLOATLITERAL
		l.next()
		for isDigit(l.ch) {
			l.next()
		}
	}
	return tok, string(l.src[offs:l.offset])
}

func (l *Lexer) scanIdentifier() string {
	offs := l.offset
	for isLetter(l.ch) || isDigit(l.ch) {
		l.next()
	}
	return string(l.src[offs:l.offset])
}

func (l *Lexer) match(r rune) {
	if l.ch != r {
		l.error("couldn't match a character")
	} else {
		l.next()
	}
}

func (l *Lexer) switch2(tok0, tok1 token.Token) (token.Token, string) {
	// because of the double switch, ch is the char after the current next
	if l.ch == '=' {
		l.next()
		return tok1, tok0.String() + "="
	}
	return tok0, tok0.String()
}

func digitVal(ch rune) int {
	if isDigit(ch) {
		return int(ch - '0')
	} else {
		return 16
	}
}
func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func isLetter(ch rune) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z')
}
