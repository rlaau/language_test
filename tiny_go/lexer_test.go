package tinygo

import (
	"testing"
)

//* 테스트코드는 AI가 작성함 *//

type expTok struct {
	kind  TokenKind
	value string // value 검증이 필요 없으면 ""로 둠
}

func assertTok(t *testing.T, got Token, want expTok) {
	t.Helper()
	if got.Kind != want.kind {
		t.Fatalf("kind mismatch: got=%v (%q) want=%v (%q) pos=%d",
			got.Kind, got.Value, want.kind, want.value, got.Pos)
	}
	if want.value != "" && got.Value != want.value {
		t.Fatalf("value mismatch: got=%q want=%q (kind=%v pos=%d)",
			got.Value, want.value, got.Kind, got.Pos)
	}
}

func lexAll(t *testing.T, input string) []Token {
	t.Helper()
	lx := NewLexer()
	lx.Set(input)

	out := []Token{}
	for {
		tok := lx.Next()
		out = append(out, tok)
		if tok.Kind == EOF {
			break
		}
	}
	return out
}

func TestLexer_Keywords_And_Identifiers(t *testing.T) {
	// EBNF에 필요한 키워드들(현재 TokenKind에 있는 것들만):
	// bool/int/string, if/then/else, for/range, let/in, scan/print, true/false
	toks := lexAll(t, "bool int string if then else for range let in scan print true false abc xyz123")

	want := []expTok{
		{BOOL, "bool"},
		{INT, "int"},
		{STRING, "string"},
		{IF, "if"},
		{THEN, "then"},
		{ELSE, "else"},
		{FOR, "for"},
		{RANGE, "range"},
		{LET, "let"},
		{IN, "in"},
		{SCAN, "scan"},
		{PRINT, "print"},
		{TRUE, "true"},
		{FALSE, "false"},
		{ID, "abc"},
		{ID, "xyz123"},
		{EOF, "EOF"},
	}

	if len(toks) != len(want) {
		t.Fatalf("token count mismatch: got=%d want=%d; got=%v", len(toks), len(want), toks)
	}
	for i := range want {
		assertTok(t, toks[i], want[i])
	}
}

func TestLexer_Numbers(t *testing.T) {
	toks := lexAll(t, "0 7 42 12345")

	want := []expTok{
		{NUMBER, "0"},
		{NUMBER, "7"},
		{NUMBER, "42"},
		{NUMBER, "12345"},
		{EOF, "EOF"},
	}

	if len(toks) != len(want) {
		t.Fatalf("token count mismatch: got=%d want=%d; got=%v", len(toks), len(want), toks)
	}
	for i := range want {
		assertTok(t, toks[i], want[i])
	}
}

func TestLexer_Delimiters(t *testing.T) {
	// 구분자: { } [ ] ( ) ; ,
	toks := lexAll(t, "{ } [ ] ( ) ; ,")

	want := []expTok{
		{LBRACE, "{"},
		{RBRACE, "}"},
		{LBRACKET, "["},
		{RBRACKET, "]"},
		{LPAREN, "("},
		{RPAREN, ")"},
		{SEMICOLON, ";"},
		{COMMA, ","},
		{EOF, "EOF"},
	}

	if len(toks) != len(want) {
		t.Fatalf("token count mismatch: got=%d want=%d; got=%v", len(toks), len(want), toks)
	}
	for i := range want {
		assertTok(t, toks[i], want[i])
	}
}

func TestLexer_Operators_OneChar(t *testing.T) {
	// 1글자 연산자: = < > ! + - * / ^
	toks := lexAll(t, "= < > ! + - * / ^")

	want := []expTok{
		{ASSIGN, "="},
		{LT, "<"},
		{GT, ">"},
		{NOT, "!"},
		{PLUS, "+"},
		{MINUS, "-"},
		{MUL, "*"},
		{DIV, "/"},
		{POW, "^"},
		{EOF, "EOF"},
	}

	if len(toks) != len(want) {
		t.Fatalf("token count mismatch: got=%d want=%d; got=%v", len(toks), len(want), toks)
	}
	for i := range want {
		assertTok(t, toks[i], want[i])
	}
}

func TestLexer_Operators_TwoChar(t *testing.T) {
	// 2글자 연산자: == != <= >= && ||
	toks := lexAll(t, "== != <= >= && ||")

	want := []expTok{
		{EQUAL, "=="},
		{NEQ, "!="},
		{LTE, "<="},
		{GTE, ">="},
		{AND, "&&"},
		{OR, "||"},
		{EOF, "EOF"},
	}

	if len(toks) != len(want) {
		t.Fatalf("token count mismatch: got=%d want=%d; got=%v", len(toks), len(want), toks)
	}
	for i := range want {
		assertTok(t, toks[i], want[i])
	}
}

func TestLexer_Rollback_When_TwoChar_Not_Operator(t *testing.T) {
	// 목적: 2글자 연산자 시도 → 실패 → rollback → 1글자 연산자로 다시 시도
	// "!x"에서 "!x"는 2글자 연산자가 아니므로, "!" + "x"로 나와야 함
	toks := lexAll(t, "!x")

	want := []expTok{
		{NOT, "!"},
		{ID, "x"},
		{EOF, "EOF"},
	}

	if len(toks) != len(want) {
		t.Fatalf("token count mismatch: got=%d want=%d; got=%v", len(toks), len(want), toks)
	}
	for i := range want {
		assertTok(t, toks[i], want[i])
	}
}

func TestLexer_StringConstructorQuote(t *testing.T) {
	// 현재 lexer는 STRLIT 전체가 아니라, 따옴표 자체(STRCONS)만 토큰화하는 흐름으로 보임
	toks := lexAll(t, "\"hello\"")

	// 기대: "  ID(hello)  "
	// (hello를 STRLIT으로 만들려면 lexer에 문자열 리터럴 읽기 로직이 추가돼야 함)
	want := []expTok{
		{STRCONS, "\""},
		{ID, "hello"},
		{STRCONS, "\""},
		{EOF, "EOF"},
	}

	if len(toks) != len(want) {
		t.Fatalf("token count mismatch: got=%d want=%d; got=%v", len(toks), len(want), toks)
	}
	for i := range want {
		assertTok(t, toks[i], want[i])
	}
}

func TestLexer_Whitespace_Is_Ignored(t *testing.T) {
	toks := lexAll(t, " \n\t  if   true  then { print 1; }  ")

	want := []expTok{
		{IF, "if"},
		{TRUE, "true"},
		{THEN, "then"},
		{LBRACE, "{"},
		{PRINT, "print"},
		{NUMBER, "1"},
		{SEMICOLON, ";"},
		{RBRACE, "}"},
		{EOF, "EOF"},
	}

	if len(toks) != len(want) {
		t.Fatalf("token count mismatch: got=%d want=%d; got=%v", len(toks), len(want), toks)
	}
	for i := range want {
		assertTok(t, toks[i], want[i])
	}
}

// 선택: unknown 토큰에서 panic이 나는지 확인 (현재 Next()는 매칭 실패 시 panic)
func TestLexer_UnknownToken_Panics(t *testing.T) {
	lx := NewLexer()
	lx.Set("@")

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic, but got none")
		}
	}()

	_ = lx.Next()
}
