package calclang

import "fmt"

type Token struct {
	Type  TokenType
	Pos   int // byte offset in input
	Lexed []byte
}

type TokenType int

const (
	// Special
	EOF TokenType = iota
	ILLEGAL

	// Literals / keywords
	INT
	TRUE
	FALSE

	// Delimiters
	LPAREN // (
	RPAREN // )

	// Operators
	OR  // ||
	AND // &&
	NOT // !

	PLUS  // +
	MINUS // -
	MUL   // *
	DIV   // /
	POW   // ^

	// Relational operators
	EQ  // ==
	NEQ // !=
	LT  // <
	LTE // <=
	GT  // >
	GTE // >=
)

func (t TokenType) String() string {
	switch t {
	case EOF:
		return "EOF"
	case ILLEGAL:
		return "ILLEGAL"

	case INT:
		return "INT(n)"

	case TRUE:
		return "TRUE(true)"
	case FALSE:

		return "FALSE(false)"
	case OR:
		return "OR(||)"
	case AND:
		return "AND(&&)"

	case EQ:
		return "EQ(==)"
	case NEQ:
		return "NEQ(!=)"

	case LPAREN:
		return "LPAREN('(')"
	case RPAREN:
		return "RPAREN(')')"
	case NOT:
		return "NOT(!)"
	case PLUS:
		return "PLUS(+)"
	case MINUS:
		return "MINUS(-)"
	case MUL:
		return "MUL(*)"
	case DIV:
		return "DIV(/)"
	case POW:
		return "POW(^)"

	case LT:
		return "LT(<)"
	case LTE:
		return "LTE(<=)"
	case GT:
		return "GT(>)"
	case GTE:
		return "GTE(>=)"
	default:
		return fmt.Sprintf("TokenType(%d)", int(t))
	}
}
