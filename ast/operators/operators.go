package operators

type Operator string

/* Boolean */
type BooleanOp Operator

const (
	AND BooleanOp = "&&"
	OR            = "||"
)

/* Relational */
type RelationalOp Operator

const (
	LT RelationalOp = "<"
	LE              = "<="
	EQ              = "=="
	NE              = "!="
	GT              = ">"
	GE              = ">="
)

/* Arithmetic */
type ArithmeticOp Operator

const (
	PLUS  ArithmeticOp = "+"
	MINUS              = "-"
	TIMES              = "*"
	DIV                = "/"
)

/* Unary */
const (
	NOT Operator = "!"
	NEG          = "-"
)

/* Cast */
const (
	INT   Operator = "int"
	FLOAT          = "float"
	CHAR           = "char"
)
