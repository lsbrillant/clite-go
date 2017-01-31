package token

import "strconv"

type Token int

const (
	EOF Token = iota

	reserved_beg

	BOOL
	CHAR
	ELSE
	FALSE
	FLOAT
	IF
	INT
	MAIN
	TRUE
	WHILE

	reserved_end

	LEFTBRACE
	RIGHTBRACE
	LEFTBRACKET
	RIGHTBRACKET
	LEFTPAREN
	RIGHTPAREN

	SEMICOLON
	COMMA
	ASSIGN

	operator_beg

	EQUALS
	LESS
	LESSEQUAL
	GREATER
	GREATEREQUAL

	NOT
	NOTEQUAL
	PLUS
	MINUS
	MULTIPLY

	DIVIDE
	AND
	OR

	operator_end

	IDENTIFIER

	literal_beg

	INTLITERAL
	FLOATLITERAL
	CHARLITERAL

	literal_end
)

var tokens = [...]string{
	EOF: "<EOF>",

	BOOL:  "bool",
	CHAR:  "char",
	ELSE:  "else",
	FALSE: "false",
	FLOAT: "float",
	IF:    "if",
	INT:   "int",
	MAIN:  "main",
	TRUE:  "true",
	WHILE: "while",

	LEFTBRACE:    "{",
	RIGHTBRACE:   "}",
	LEFTBRACKET:  "[",
	RIGHTBRACKET: "]",
	LEFTPAREN:    "(",
	RIGHTPAREN:   ")",

	SEMICOLON: ";",
	COMMA:     ",",
	ASSIGN:    "=",

	EQUALS:       "==",
	LESS:         "<",
	LESSEQUAL:    "<=",
	GREATER:      ">",
	GREATEREQUAL: ">=",

	NOT:      "!",
	NOTEQUAL: "!=",
	PLUS:     "+",
	MINUS:    "-",
	MULTIPLY: "*",
	DIVIDE:   "/",
	AND:      "&&",
	OR:       "||",

	IDENTIFIER: "IDENT",

	INTLITERAL:   "INT",
	FLOATLITERAL: "FLOAT",
	CHARLITERAL:  "CHAR",
}

func (tok Token) String() string {
	s := ""
	if 0 <= tok && tok < Token(len(tokens)) {
		s = tokens[tok]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(tok)) + ")"
	}
	return s
}

var keywords map[string]Token

func init() {
	keywords = make(map[string]Token)
	for i := reserved_beg + 1; i < reserved_end; i++ {
		keywords[tokens[i]] = i
	}
}

// Lookup maps an identifier to its keyword token or IDENT (if not a keyword).
//
func Lookup(ident string) Token {
	if tok, is_keyword := keywords[ident]; is_keyword {
		return tok
	}
	return IDENTIFIER
}

// Predicates

// IsLiteral returns true for tokens corresponding to identifiers
// and basic type literals; it returns false otherwise.
//
func (tok Token) IsLiteral() bool { return literal_beg < tok && tok < literal_end }

// IsOperator returns true for tokens corresponding to operators and
// delimiters; it returns false otherwise.
//
func (tok Token) IsOperator() bool { return operator_beg < tok && tok < operator_end }

// IsKeyword returns true for tokens corresponding to keywords;
// it returns false otherwise.
//
func (tok Token) IsKeyword() bool { return reserved_beg < tok && tok < reserved_end }
