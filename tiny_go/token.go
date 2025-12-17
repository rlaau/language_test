package tinygo

import "fmt"

type Token struct {
	Kind  TokenKind
	Value string
	Pos   int
}

type TokenKind int

// 키워드 블록
const (
	// 표현식 키워드
	START_OF_KEYWORD TokenKind = iota
	TRUE
	FALSE
	NUMBER
	STRLIT
	FUNC

	// 타입 키워드
	BOOL
	INT
	STRING
	OMIT

	// 선언 키워드
	VAR

	//return 키워드
	RETURN

	// if 키워드
	IF
	ELSE

	// for 키워드
	FOR
	RANGE

	// let 키워드
	LET
	IN

	// 내장 함수 키워드
	SCAN
	PRINT
	END_OF_KEYWORD
)
const (
	//식별자
	START_OF_ID TokenKind = END_OF_KEYWORD + 1 + iota
	ID
	END_OF_ID
)
const (
	START_OF_DELIMETER TokenKind = END_OF_ID + 1 + iota
	// 구분자
	LBRACE
	RBRACE
	LBRACKET
	RBRACKET
	LPAREN
	RPAREN
	SEMICOLON
	COMMA
	END_OF_DELIMETER
)
const (
	START_OF_OPERATOR TokenKind = END_OF_DELIMETER + 1 + iota
	// 연산자 및 선언 연산자
	//대입
	ASSIGN
	// 짧은 선언 (선언 및 대입 연산)
	SHORTDECL
	//비교
	EQUAL
	NEQ
	LT
	LTE
	GT
	GTE
	// 논리
	AND
	OR
	NOT
	// 사칙
	PLUS
	MINUS
	MUL
	DIV
	END_OF_OPERATOR
)
const (
	START_OF_SPECIAL TokenKind = END_OF_OPERATOR + 1 + iota
	EOF
	ILLLEGAL
	END_OF_SPECIAL
)

// NewToken은 토큰 생성
func NewToken(t TokenKind, pos int) Token {
	return Token{
		Kind:  t,
		Value: StringSpec(t),
		Pos:   pos,
	}
}

// StringSpec은 tokenKind가 초기값으로 가져야 할 문자열 값 리턴.
func StringSpec(t TokenKind) string {
	switch t {
	case TRUE:
		return "true"
	case FALSE:
		return "false"
	case NUMBER:
		return ""
	case STRLIT:
		return ""
	case FUNC:
		return "func"

	case BOOL:
		return "bool"
	case INT:
		return "int"
	case STRING:
		return "string"
	case OMIT:
		return "()"

	case VAR:
		return "var"

	case RETURN:
		return "return"
	case IF:
		return "if"
	case ELSE:
		return "else"

	case FOR:
		return "for"
	case RANGE:
		return "range"

	case LET:
		return "let"
	case IN:
		return "in"

	case SCAN:
		return "scan"
	case PRINT:
		return "print"

	case ID:
		return ""

	case LBRACE:
		return "{"
	case RBRACE:
		return "}"
	case LBRACKET:
		return "["
	case RBRACKET:
		return "]"
	case LPAREN:
		return "("
	case RPAREN:
		return ")"
	case SEMICOLON:
		return ";"
	case COMMA:
		return ","

	case ASSIGN:
		return "="
	case SHORTDECL:
		return ":="
	case EQUAL:
		return "=="
	case NEQ:
		return "!="
	case LT:
		return "<"
	case LTE:
		return "<="
	case GT:
		return ">"
	case GTE:
		return ">="

	case AND:
		return "&&"
	case OR:
		return "||"
	case NOT:
		return "!"

	case PLUS:
		return "+"
	case MINUS:
		return "-"
	case MUL:
		return "*"
	case DIV:
		return "/"

	case EOF:
		//EOF는 "EOF"를 EOF로 토크나이징 하지는 않음.
		// 입력값 끝 시 나타날 뿐임.
		return "<<EOF>>"
	case ILLLEGAL:
		return "<<ILLEGAL>>"
	default:
		msg := fmt.Sprintf("StringSpec: 입력 %d에 대해 매치되는 토큰이 없습니다.", t)
		panic(msg)
	}
}

func IsKeyWord(s string) (TokenKind, bool) {
	return isStringInRange(s, START_OF_KEYWORD, END_OF_KEYWORD)
}
func IsDelimeter(s string) (TokenKind, bool) {
	return isStringInRange(s, START_OF_DELIMETER, END_OF_DELIMETER)
}

func IsOperator(s string) (TokenKind, bool) {
	return isStringInRange(s, START_OF_OPERATOR, END_OF_OPERATOR)
}

func isStringInRange(s string, start, end TokenKind) (TokenKind, bool) {
	if start > end {
		start, end = end, start
	}
	for offset := range end - start - 1 {
		startPoint := start + 1
		tokenKind := startPoint + offset
		if s == StringSpec(tokenKind) {
			return tokenKind, true
		}
	}
	return ILLLEGAL, false
}

func (t *Token) String() string {
	return t.Value
}
func (t *Token) SetValue(v string) {
	if t.Kind == ID || t.Kind == NUMBER || t.Kind == STRLIT {
		t.Value = v
		return
	}
	panic("ID, NUMBER, STR 외에는 SetValue 호출하면 안됩니다.")
}
