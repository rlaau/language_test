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
	//1. number, strlit, omit을 제외한 키워드
	//2. id
	//3. omit
	//4. number
	//5. strlit
	//6. operator
	//7. delimeter
	//8. special (EOF, ILLEGAL)

	//으로 "구성"된다. (서로소는 아님으로 "분할"까진 아님.)

	// 우선 EOF의 경우를 제거하자
	// ILLEGAL은 토크나이징 실패 케이스로써 처리하자
	if lx.isOvered() {
		return NewToken(EOF, lx.currentPosition)
	}

	// 그럼 이제 1,2,3,4,5,6,7 번의 경우만이 남는다.

	//롤백 등록
	rollBack := lx.rollback()

	// (1, 2 번 케이스의 토큰)라면 -> 현재 읽은 문자가 알파벳이다.
	// (현재 읽은 문자가 알파벳)라면-> (1,2번 케이스의 토큰)
	// 그러므로 (1, 2 번 케이스의 토큰)  ==  (현재 읽은 문자가 알파벳) 이다.
	// 참고: Union(1, 2번 케이스) = number, strlit아닌 키워드와 id이다
	if isAlpha(lx.currentByte()) {
		// 문자{문자|숫자} 전부 수집 후 다음 칸 이동
		candidate := lx.readWhile(isDigitOrAlpha)

		keyWord, isKeyWordExist := IsKeyWord(string(candidate))
		// 키워드인 경우 <-> 1번 케이스
		if isKeyWordExist {
			return NewToken(keyWord, lx.currentPosition)
		}
		// 식별자인 경우 <-> 2번 케이스
		idToken := NewToken(ID, lx.currentPosition)
		idToken.SetValue(string(candidate))
		return idToken
	}
	// 이제 3,4,5,6,7 케이스가 남는다.
	// 3번 케이스인 Omit 키워드 파싱한다
	// Omit은 연산자를 제외한 모든 토큰과 서로소이며, 연산자보다 길이가 길다.
	// 그러므로 이 시점에서 유일-타당하게 토크나이징 될 수 있다.
	isOmitStart := func(s string) (TokenKind, bool) {
		if s == "(" {
			return OMIT, true
		}
		return ILLLEGAL, false
	}
	isOmitEnd := func(s string) (TokenKind, bool) {
		if s == ")" {
			return OMIT, true
		}
		return ILLLEGAL, false
	}
	if _, isStartOfOmit := isOmitStart(string(lx.currentByte())); isStartOfOmit {
		_, _ = lx.readIf(string(lx.currentByte()), isOmitStart)
		if _, isEndOfOmit := lx.readIf(string(lx.currentByte()), isOmitEnd); isEndOfOmit {
			return NewToken(OMIT, lx.currentPosition)
		}
		rollBack()
	}

	// 이제 4,5,6,7 케이스가 남는다.
	// 4번 케이스. number처리.
	// 이 역시 필요충분 조건이며, 5,6,7 케이스와 서로소이다.
	if isDigit(lx.currentByte()) {
		candidate := lx.readWhile(isDigit)
		numberToken := NewToken(NUMBER, lx.currentPosition)
		numberToken.SetValue(string(candidate))
		return numberToken
	}

	// 이제 5,6,7 케이스가 남는다.
	// 5 케이스 검사
	// strlit은 "\"" 으로 시작하고, 이는 6,7 케이스의 첫 문자와 서로소이다.
	// 그러므로 첫 문자만 보고 유일한 케이스로 좁힐 수 있다.
	isStrEdge := func(s string) (TokenKind, bool) {
		if s == "\"" {
			return STRLIT, true
		}
		return ILLLEGAL, false
	}
	if _, isEdgeOfStrlit := lx.readIf(string(lx.currentByte()), isStrEdge); isEdgeOfStrlit {
		isInnerOfStrlit := func(b byte) bool { _, isStr := isStrEdge(string(b)); return !isStr }
		candidate := lx.readWhile(isInnerOfStrlit)
		if _, isStrEdge2 := lx.readIf(string(lx.currentByte()), isStrEdge); isStrEdge2 {
			strlitToken := NewToken(STRLIT, lx.currentPosition)
			// STRLT은 따옴표 제거한 "내부 문자" 만 토큰 값으로 넣는다.
			strlitToken.SetValue(string(candidate))
			return strlitToken
		}
		rollBack()
	}
	// 이제 6,7 케이스가 남는다.
	// 6 케이스 검사
	// 6의 연산자 집합에 대해, 모든 원소의 길이가 2 이하임을 상정한다
	// 6,7 케이스 역시 서로소이므로 순서에 의존하지 않지만, 길이가 더 긴 6케이스를 우선 검사한다.
	// 길이 2부터 검사
	if lx.isNextExist() {
		candidate := lx.input[lx.currentPosition : lx.currentPosition+2]
		if tokenKind, isOperator := lx.readIf(candidate, IsOperator); isOperator {
			return NewToken(tokenKind, lx.currentPosition)
		}
		// 한 칸 이동 후 토크나이징 시도했지만 실패 시, 다시 원위치로 롤백
		rollBack()
	}
	if tokenKind, isOperator := lx.readIf(string(lx.currentByte()), IsOperator); isOperator {
		return NewToken(tokenKind, lx.currentPosition)
	}

	// 이제 7케이스만이 남는다.
	// 7의 구분자 집합에 대해, 모든 원소의 길이가 1임을 상정한다.
	if tokenKind, isDelimeter := lx.readIf(string(lx.currentByte()), IsDelimeter); isDelimeter {
		return NewToken(tokenKind, lx.currentPosition)
	}

	// 이제 어느 케이스도 남지 않는다.
	//나머지 경우엔 올바른 토큰이 존재하지 않는다고 볼 수 있다.
	rollBack()
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
	return ILLLEGAL, false
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
