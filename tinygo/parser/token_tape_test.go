package parser

import (
	"testing"

	"github.com/rlaaudgjs5638/langTest/tinygo/lexer"
	"github.com/rlaaudgjs5638/langTest/tinygo/token"
)

// 헬퍼: 렉서에 입력을 주입하고 TokenTape 생성
func newTokenTape(t *testing.T, input string) *TokenTape {
	t.Helper()
	lx := lexer.NewLexer()
	lx.Set(input)
	return NewTokenTape(lx)
}

// 헬퍼: 토큰 검증
func assertToken(t *testing.T, got token.Token, wantKind token.TokenKind, wantValue string) {
	t.Helper()
	if got.Kind != wantKind {
		t.Errorf("token kind mismatch: got=%v want=%v (value: %q)", got.Kind, wantKind, got.Value)
	}
	if wantValue != "" && got.Value != wantValue {
		t.Errorf("token value mismatch: got=%q want=%q (kind: %v)", got.Value, wantValue, got.Kind)
	}
}

// Test 1: 기본 CurrentToken 동작
func TestTokenTape_CurrentToken(t *testing.T) {
	tape := newTokenTape(t, "var x = 10")

	// 초기 토큰은 "var"
	tok := tape.CurrentToken()
	assertToken(t, tok, token.VAR, "var")
}

// Test 2: MoveToNextToken 기본 동작
func TestTokenTape_MoveToNextToken_Basic(t *testing.T) {
	tape := newTokenTape(t, "var x = 10")

	// var -> x -> = -> 10 -> EOF
	assertToken(t, tape.CurrentToken(), token.VAR, "var")

	tape.MoveToNextToken()
	assertToken(t, tape.CurrentToken(), token.ID, "x")

	tape.MoveToNextToken()
	assertToken(t, tape.CurrentToken(), token.ASSIGN, "=")

	tape.MoveToNextToken()
	assertToken(t, tape.CurrentToken(), token.NUMBER, "10")

	tape.MoveToNextToken()
	assertToken(t, tape.CurrentToken(), token.EOF, "<<EOF>>")
}

// Test 3: Peek(1) - 다음 토큰 미리보기
func TestTokenTape_Peek_NextToken(t *testing.T) {
	tape := newTokenTape(t, "if true { }")

	// 현재: if
	assertToken(t, tape.CurrentToken(), token.IF, "if")

	// Peek(1): true (현재는 여전히 if)
	peeked := tape.Peek(1)
	assertToken(t, peeked, token.TRUE, "true")
	assertToken(t, tape.CurrentToken(), token.IF, "if") // 현재 토큰 변화 없음
}

// Test 4: Peek(n) - 여러 칸 앞 미리보기
func TestTokenTape_Peek_Multiple(t *testing.T) {
	tape := newTokenTape(t, "for x in range")

	// 현재: for
	assertToken(t, tape.CurrentToken(), token.FOR, "for")

	// Peek(2): in
	peeked2 := tape.Peek(2)
	assertToken(t, peeked2, token.IN, "in")

	// Peek(3): range
	peeked3 := tape.Peek(3)
	assertToken(t, peeked3, token.RANGE, "range")

	// 현재 토큰은 여전히 for
	assertToken(t, tape.CurrentToken(), token.FOR, "for")
}

// Test 5: Peek 후 MoveToNextToken - 레코드 재사용
func TestTokenTape_Peek_Then_Move(t *testing.T) {
	tape := newTokenTape(t, "let a = 5")

	// 현재: let
	assertToken(t, tape.CurrentToken(), token.LET, "let")

	// Peek(1)로 미리 읽기
	tape.Peek(1)

	// MoveToNextToken으로 이동 - 이미 레코드에 있으므로 렉서 호출 없음
	tape.MoveToNextToken()
	assertToken(t, tape.CurrentToken(), token.ID, "a")

	tape.MoveToNextToken()
	assertToken(t, tape.CurrentToken(), token.ASSIGN, "=")
}

// Test 6: GetRollback 기본 동작
func TestTokenTape_Rollback_Basic(t *testing.T) {
	tape := newTokenTape(t, "func add ( x")

	// 현재: func
	assertToken(t, tape.CurrentToken(), token.FUNC, "func")

	// 롤백 포인트 생성
	rollback := tape.GetRollback()

	// 토큰 3개 이동: func -> add -> ( -> x
	tape.MoveToNextToken()
	tape.MoveToNextToken()
	tape.MoveToNextToken()
	assertToken(t, tape.CurrentToken(), token.ID, "x")

	// 롤백 실행
	rollback()

	// 원래 위치(func)로 복귀
	assertToken(t, tape.CurrentToken(), token.FUNC, "func")
}

// Test 7: 중첩 롤백
func TestTokenTape_Rollback_Nested(t *testing.T) {
	tape := newTokenTape(t, "a b c d e")

	// 현재: a
	assertToken(t, tape.CurrentToken(), token.ID, "a")

	// 첫 번째 롤백 포인트 (a)
	rollback1 := tape.GetRollback()

	tape.MoveToNextToken() // -> b
	tape.MoveToNextToken() // -> c
	assertToken(t, tape.CurrentToken(), token.ID, "c")

	// 두 번째 롤백 포인트 (c)
	rollback2 := tape.GetRollback()

	tape.MoveToNextToken() // -> d
	tape.MoveToNextToken() // -> e
	assertToken(t, tape.CurrentToken(), token.ID, "e")

	// 두 번째 롤백 실행 (c로 복귀)
	rollback2()
	assertToken(t, tape.CurrentToken(), token.ID, "c")

	// 첫 번째 롤백 실행 (a로 복귀)
	rollback1()
	assertToken(t, tape.CurrentToken(), token.ID, "a")
}

// Test 8: Peek(0) - 현재 토큰
func TestTokenTape_Peek_Zero(t *testing.T) {
	tape := newTokenTape(t, "print scan")

	assertToken(t, tape.CurrentToken(), token.PRINT, "print")

	// Peek(0)는 현재 토큰과 동일
	peeked := tape.Peek(0)
	assertToken(t, peeked, token.PRINT, "print")
}

// Test 9: EOF까지 Peek
func TestTokenTape_Peek_ToEOF(t *testing.T) {
	tape := newTokenTape(t, "x y")

	// 현재: x
	assertToken(t, tape.CurrentToken(), token.ID, "x")

	// Peek(1): y
	assertToken(t, tape.Peek(1), token.ID, "y")

	// Peek(2): EOF
	assertToken(t, tape.Peek(2), token.EOF, "<<EOF>>")
	assertToken(t, tape.Peek(50), token.EOF, "<<EOF>>")
}

// Test 10: EOF 이후 MoveToNextToken (경계 테스트)
func TestTokenTape_MoveToNextToken_AfterEOF(t *testing.T) {
	tape := newTokenTape(t, "ok")

	assertToken(t, tape.CurrentToken(), token.OK, "ok")

	tape.MoveToNextToken()
	assertToken(t, tape.CurrentToken(), token.EOF, "<<EOF>>")

	// EOF 이후에도 MoveToNextToken 호출 가능 (패닉 없어야 함)
	tape.MoveToNextToken()

	// tokenRecord 길이 확인
	if len(tape.tokenRecord) < 2 {
		t.Errorf("expected tokenRecord to grow, got length=%d", len(tape.tokenRecord))
	}
}

// Test 11: 빈 입력
func TestTokenTape_EmptyInput(t *testing.T) {
	tape := newTokenTape(t, "")

	// 빈 입력은 즉시 EOF
	assertToken(t, tape.CurrentToken(), token.EOF, "<<EOF>>")
}

// Test 12: Peek 범위 확장 (레코드에 없는 토큰 여러 개)
func TestTokenTape_Peek_ExpandRecord(t *testing.T) {
	tape := newTokenTape(t, "a b c d e f")

	// 현재: a
	assertToken(t, tape.CurrentToken(), token.ID, "a")

	// Peek(5)로 f까지 미리 읽기
	peeked := tape.Peek(5)
	assertToken(t, peeked, token.ID, "f")

	// tokenRecord에 a~f + EOF까지 저장되어야 함
	if len(tape.tokenRecord) < 6 {
		t.Errorf("expected tokenRecord length >= 6, got=%d", len(tape.tokenRecord))
	}

	// 현재 토큰은 여전히 a
	assertToken(t, tape.CurrentToken(), token.ID, "a")
}

// Test 13: Rollback 후 Peek
func TestTokenTape_Rollback_Then_Peek(t *testing.T) {
	tape := newTokenTape(t, "x := 10")

	// 현재: x
	assertToken(t, tape.CurrentToken(), token.ID, "x")
	rollback := tape.GetRollback()

	tape.MoveToNextToken() // :=
	tape.MoveToNextToken() // 10
	assertToken(t, tape.CurrentToken(), token.NUMBER, "10")

	// 롤백
	rollback()
	assertToken(t, tape.CurrentToken(), token.ID, "x")

	// 롤백 후 Peek(1)
	peeked := tape.Peek(1)
	assertToken(t, peeked, token.DECLSIGN, ":=")
}

// Test 14: Off-by-one 검증 - isNextNTokensExistOnRecord
func TestTokenTape_OffByOne_Peek(t *testing.T) {
	tape := newTokenTape(t, "true false")

	// 현재: true (currentIdx=0, tokenRecord=[true])
	assertToken(t, tape.CurrentToken(), token.TRUE, "true")

	// Peek(1) 호출 - false를 읽어와야 함
	peeked := tape.Peek(1)
	assertToken(t, peeked, token.FALSE, "false")

	// tokenRecord=[true, false], currentIdx=0
	if len(tape.tokenRecord) != 2 {
		t.Errorf("expected tokenRecord length=2, got=%d", len(tape.tokenRecord))
	}
	if tape.currentIdx != 0 {
		t.Errorf("expected currentIdx=0, got=%d", tape.currentIdx)
	}
}

// Test 15: Off-by-one 검증 - isNextTokenExistOnRecord
func TestTokenTape_OffByOne_MoveToNextToken(t *testing.T) {
	tape := newTokenTape(t, "int bool")

	// 현재: int (currentIdx=0, tokenRecord=[int])
	assertToken(t, tape.CurrentToken(), token.INT, "int")

	// Peek(1)로 bool 미리 읽기
	tape.Peek(1)
	// tokenRecord=[int, bool], currentIdx=0

	// MoveToNextToken - 이미 레코드에 있으므로 렉서 호출 없이 currentIdx++
	tape.MoveToNextToken()
	assertToken(t, tape.CurrentToken(), token.BOOL, "bool")

	if tape.currentIdx != 1 {
		t.Errorf("expected currentIdx=1, got=%d", tape.currentIdx)
	}
}

// Test 16: 복잡한 시나리오 - Peek, Move, Rollback 혼합
func TestTokenTape_Complex_Scenario(t *testing.T) {
	tape := newTokenTape(t, "if x > 5 { return true }")

	// 현재: if
	assertToken(t, tape.CurrentToken(), token.IF, "if")
	rollback1 := tape.GetRollback()

	// Peek(3): 5
	assertToken(t, tape.Peek(3), token.NUMBER, "5")

	tape.MoveToNextToken() // x
	assertToken(t, tape.CurrentToken(), token.ID, "x")
	rollback2 := tape.GetRollback()

	tape.MoveToNextToken() // >
	tape.MoveToNextToken() // 5
	tape.MoveToNextToken() // {
	assertToken(t, tape.CurrentToken(), token.LBRACE, "{")

	// rollback2 실행 (x로 복귀)
	rollback2()
	assertToken(t, tape.CurrentToken(), token.ID, "x")

	// 다시 진행
	tape.MoveToNextToken() // >
	tape.MoveToNextToken() // 5
	assertToken(t, tape.CurrentToken(), token.NUMBER, "5")

	// rollback1 실행 (if로 복귀)
	rollback1()
	assertToken(t, tape.CurrentToken(), token.IF, "if")
}
