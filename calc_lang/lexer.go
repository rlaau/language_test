package calclang

import (
	"errors"
	"unicode"
)

type Lexer struct {
	input      []byte
	currentPos int // 현재 바이트 위치
}

func NewLexer() *Lexer {
	return &Lexer{input: []byte{}, currentPos: 0}
}

func (lx *Lexer) Lex(input string) ([]Token, error) {
	lx.input = []byte(input)
	lx.currentPos = 0

	tokenList := []Token{}
	for {
		// 다음 토큰 읽기
		nextToken := lx.lookUpCurrentToken()
		if nextToken.Type == ILLEGAL {
			return tokenList, errors.New("invalid Token")
		}
		if nextToken.Type == EOF {
			return tokenList, nil
		}
		//올바른 토큰 시 토큰 리스트에 추가
		tokenList = append(tokenList, nextToken)
		// 토큰 소비 후 다음 포지션으로 이동
		nextPos := lx.reflectToken(nextToken)
		lx.currentPos = nextPos
	}
}

// 계산: 다음 토큰 이후의 인덱스 리턴
func (lx *Lexer) reflectToken(t Token) int {
	endPos := lx.currentPos + len(t.Lexed) - 1
	nextPos := endPos + 1
	return nextPos
}

// 읽기 연산
func (lx *Lexer) lookUpCurrentToken() Token {
	lx.skipWhitespace()

	start := lx.currentPos
	if start >= len(lx.input) {
		return Token{
			Type:  EOF,
			Pos:   start,
			Lexed: []byte{},
		}
	}
	// " 현재 지점" 부터 룩업
	ch := lx.getCurrentByte()

	// 긴 길이 문자열 처리 with lookUntil
	// true/false
	if isAlpha(ch) {
		bWord := string(lx.lookUntil(isAlpha))
		switch bWord {
		case "true":
			return Token{Type: TRUE, Pos: start, Lexed: []byte("true")}
		case "false":
			return Token{Type: FALSE, Pos: start, Lexed: []byte("false")}
		default:
			return Token{Type: ILLEGAL, Pos: start, Lexed: []byte{}}
		}
	}

	//가변 길이 문자열 처리 with lookUntil
	// int
	if isDigit(ch) {
		number := string(lx.lookUntil(isDigit))
		return Token{Type: INT, Pos: start, Lexed: []byte(number)}
	}
	//두개 길이 문자열
	// 두 칸으 여유공간 있을 시 비교
	if start+1 < len(lx.input) {
		twoChar := string(lx.input[start : start+2])
		lookUpTwoChar := func(t TokenType, s string) Token {
			return Token{Type: t, Pos: start, Lexed: []byte(s)}
		}
		switch twoChar {
		case "||":
			return lookUpTwoChar(OR, "||")
		case "&&":
			return lookUpTwoChar(AND, "&&")
		case ">=":
			return lookUpTwoChar(GTE, ">=")
		case "<=":
			return lookUpTwoChar(LTE, "<=")
		case "==":
			return lookUpTwoChar(EQ, "==")
		case "!=":
			return lookUpTwoChar(NEQ, "!=")
		}
	}
	oneChar := string(ch)
	lookUpOneChar := func(t TokenType, s string) Token {
		return Token{Type: t, Pos: start, Lexed: []byte(s)}
	}
	switch oneChar {
	case "(":
		return lookUpOneChar(LPAREN, "(")
	case ")":
		return lookUpOneChar(RPAREN, ")")
	case "!":
		return lookUpOneChar(NOT, "!")
	case "+":
		return lookUpOneChar(PLUS, "+")
	case "-":
		return lookUpOneChar(MINUS, "-")
	case "*":
		return lookUpOneChar(MUL, "*")
	case "/":
		return lookUpOneChar(DIV, "/")
	case "^":
		return lookUpOneChar(POW, "^")
	case "<":
		return lookUpOneChar(LT, "<")
	case ">":
		return lookUpOneChar(GT, ">")
	}

	return Token{Type: ILLEGAL, Pos: start, Lexed: []byte{}}
}

// lx의 인풋에서 공백 스킵
func (lx *Lexer) skipWhitespace() {
	for lx.currentPos < len(lx.input) {
		r := rune(lx.input[lx.currentPos])
		if unicode.IsSpace(r) {
			lx.currentPos++
			continue
		}
		break
	}
}

// 읽기 연산
func (lx *Lexer) getCurrentByte() byte {
	return lx.input[lx.currentPos]
}

// 읽기 연산
func (lx *Lexer) lookUntil(cond func(byte) bool) []byte {
	startIdx := lx.currentPos
	for i := range len(lx.input) - startIdx {
		if !cond(lx.input[startIdx+i]) {
			return lx.input[startIdx : startIdx+i]
		}
	}
	return lx.input[startIdx:len(lx.input)]
}

func isDigit(b byte) bool {
	return '0' <= b && b <= '9'
}

func isAlpha(b byte) bool {
	return ('a' <= b && b <= 'z') || ('A' <= b && b <= 'Z')
}
