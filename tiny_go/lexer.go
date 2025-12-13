package tinygo

import (
	"fmt"
	"unicode"
)

type Lexer struct {
	input           string
	currentPosition int
}

func NewLexer() *Lexer {
	return &Lexer{input: "", currentPosition: 0}
}

func (lx *Lexer) Set(s string) {
	lx.input = s
}

// Next는 현재 위치에서의 토큰을 리턴한 후, 다음 위치로 렉서의 포지션을 옮긴다.
func (lx *Lexer) Next() Token {
	// 공백을 제거하면 문자와 맞닿게 된다.
	lx.skipWhitespace()
	// "문자"는 토크나이징 되거나, 토크나이징 되지 못한다.
	// 이떄 올바른 토큰이 존재한다 가정하면
	// 렉서가 파싱할 모든 토큰은
	//1. true, false
	//2. strlit의 '"' 값
	//3. number의 숫자 값
	//4. id의 문자열{문자열|숫자} 값
	//5. 키워드의 문자열 값
	//6. 구분자의 특수값
	//7. 연산자의 특수값
	//8. EOF
	//으로 "구성"된다. (서로소는 아님으로 "분할"까진 아님.)

	// 우선 EOF의 경우를 제거하자
	if lx.isOvered() {
		return NewToken(EOF, lx.currentPosition)
	}

	//롤백 등록
	rb := lx.rollback()

	// (1, 4, 5 번 케이스의 토큰)라면 -> 현재 읽은 문자가 알파벳이다.
	// (현재 읽은 문자가 알파벳)라면-> (1,4,5번 케이스의 토큰)
	// 그러므로 (1, 4, 5 번 케이스의 토큰) <-> (현재 읽은 문자가 알파벳)
	if isAlpha(lx.currentByte()) {

		// 문자{문자|숫자} 전부 수집 후 다음 칸 이동
		candidate := lx.readWhile(isDigitOrAlpha)

		keyWord, isKeyWordExist := IsKeyWord(string(candidate))
		// 키워드인 경우 <-> 1,5번
		if isKeyWordExist {
			return NewToken(keyWord, lx.currentPosition)
		}
		// 식별자인 경우
		idToken := NewToken(ID, lx.currentPosition)
		idToken.SetValue(string(candidate))
		return idToken
	}

	// 3번 케이스. number처리. 이 역시 필요충분 조건.
	if isDigit(lx.currentByte()) {
		candidate := lx.readWhile(isDigit)
		numberToken := NewToken(NUMBER, lx.currentPosition)
		numberToken.SetValue(string(candidate))
		return numberToken
	}

	// 2,6,7 케이스
	// 6의 구분자 집합과 2,7은 서로소 집합이므로, 문자열 길이가 1임에도 구분자 집합 먼저 검사
	// 6의 구분자 집합에 대해, 모든 원소의 길이가 1임을 상정함.
	if tokenKind, isDelimeter := lx.readIf(string(lx.currentByte()), IsDelimeter); isDelimeter {
		return NewToken(tokenKind, lx.currentPosition)
	}

	// 7 케이스 검사
	// 7의 연산자 집합에 대해, 모든 원소의 길이가 2 이하임을 상정함
	// 길이 2부터 검사
	if lx.isNextExist() {
		candidate := lx.input[lx.currentPosition : lx.currentPosition+2]
		if tokenKind, isOperator := lx.readIf(candidate, IsOperator); isOperator {
			return NewToken(tokenKind, lx.currentPosition)
		}
		// 한 칸 이동 후 토크나이징 시도했지만 실패 시, 다시 원위치로 롤백
		rb()
	}
	if tokenKind, isOperator := lx.readIf(string(lx.currentByte()), IsOperator); isOperator {
		return NewToken(tokenKind, lx.currentPosition)
	}

	// 2케이스 검사
	isStrCons := func(s string) (TokenKind, bool) {
		if s == StringSpec(STRCONS) {
			return STRCONS, true
		}
		return EOF, false
	}
	if tokenKind, isStrCons := lx.readIf(string(lx.currentByte()), isStrCons); isStrCons {
		return NewToken(tokenKind, lx.currentPosition)
	}

	//나머지 경우엔 올바른 토큰이 존재하지 않는다고 볼 수 있다.
	fmt.Printf("lexer position %d, lexer current string %s \n", lx.currentPosition, string(lx.currentByte()))
	panic("매칭되는 토큰이 존재하지 않음.")
}

// readIf는 받은 candidate가 isToken함수가 반환하는 키워드 집합에 속한 문자열일 시,
// lx가 그만큼 현재 칸을 이동하고 참을 반환하도록 한다.
// 상태 변경 함수다.
func (lx *Lexer) readIf(candidate string, isToken func(string) (TokenKind, bool)) (TokenKind, bool) {
	if tokenKind, isToken := isToken(candidate); isToken {
		lx.currentPosition += len(candidate)
		return tokenKind, true
	}
	return EOF, false
}

// readWhile은 렉서가 현재 위치한 문자가 cond조건을 만족하는 한 계속 문자를 읽고, 다음 칸으로 간다.
// EOF 케이스도 체크한다.
// 상태 변경 함수다.
func (lx *Lexer) readWhile(cond func(b byte) bool) []byte {
	candidate := []byte{}
	for !lx.isOvered() && cond(lx.currentByte()) {
		candidate = append(candidate, lx.currentByte())
		lx.currentPosition++
	}
	return candidate
}
func (lx *Lexer) currentByte() byte {
	return lx.input[lx.currentPosition]
}

// isOvered는 lx의 현재 위치가 자신의 input을 "이미 벗어난" 상태인 경우 true리턴
func (lx *Lexer) isOvered() bool {
	return lx.currentPosition >= len(lx.input)
}

// isNextExist는 lx의 현재 위치 앞에 "한 문자 이상 존재" 할 경우 true 리턴
func (lx *Lexer) isNextExist() bool {
	return lx.currentPosition < len(lx.input)-1
}

// rollback은 lx의 현재 위치를 저장되었던 위치로 롤백시킨다.
func (lx *Lexer) rollback() func() {
	memorizedPosition := lx.currentPosition
	return func() {
		lx.currentPosition = memorizedPosition
	}

}

// skipWhitespace는 lx의 인풋에서 공백을 스킵한다.
func (lx *Lexer) skipWhitespace() {
	for lx.currentPosition < len(lx.input) {
		r := rune(lx.input[lx.currentPosition])
		if unicode.IsSpace(r) {
			lx.currentPosition++
			continue
		}
		break
	}
}
