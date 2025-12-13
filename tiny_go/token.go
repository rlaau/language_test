package tinygo

type Token struct {
	Kind  TokenKind
	Value string
	Pos   int
}

type TokenKind int

// ! Append-only
const (
	//값
	TRUE TokenKind = iota
	FALSE
	NUMBER
	STRLIT

	// 타입
	BOOL
	INT
	STRING
	//식별자
	ID

	// if 키워드
	IF
	THEN
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

	// 구분자
	EOF
	LBRACE
	RBRACE
	LBRACKET
	RBRACKET
	LPAREN
	RPAREN
	SEMICOLON
	COMMA

	// 연산자
	//대입
	ASSIGN
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
	POW

	// 생성자
	STRCONS
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

	case BOOL:
		return "bool"
	case INT:
		return "int"
	case STRING:
		return "string"

	case ID:
		return ""

	case IF:
		return "if"
	case THEN:
		return "then"
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

	case EOF:
		//EOF는 "EOF"를 EOF로 토크나이징 하지는 않음.
		// 입력값 끝 시 나타날 뿐임.
		return "EOF"
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
	case POW:
		return "^"

	case STRCONS:
		return "\""
	default:
		panic("매치되는 토큰이 없습니다.")
	}
}

// IsKeyword는 받은 문자열이 키워드인지 판단해서 맞다면 true, 아니면 false리탄힌디.
// false 인 경우 빋은 문자열이 식별자(ID)인 것을 알 수 있다.
func IsKeyWord(s string) (TokenKind, bool) {
	// enum 정의 상에서, EOF이후의 수는 전부 구분자, 연산자임
	// EOF 이후의 tokenKind는 키워드가 아니므로 더이상 비교하지 않음
	// EOF는 키워드 취급이 아님. 입력값 끝을 나타내는 표식일 뿐임.
	for tokenKind := range EOF {
		if s == StringSpec(tokenKind) {
			return tokenKind, true
		}
	}
	return EOF, false
}

func IsDelimeter(s string) (TokenKind, bool) {
	for offset := range ASSIGN - LBRACE {
		tokenKind := LBRACE + offset
		if s == StringSpec(tokenKind) {
			return tokenKind, true
		}
	}
	return EOF, false
}

func IsOperator(s string) (TokenKind, bool) {
	for offset := range STRCONS - ASSIGN {
		tokenKind := ASSIGN + offset
		if s == StringSpec(tokenKind) {
			return tokenKind, true
		}
	}
	return EOF, false
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
