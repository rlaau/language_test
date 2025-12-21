package lexer

import (
	"testing"

	"github.com/rlaaudgjs5638/langTest/tinygo/token"
)

//* 테스트코드는 AI가 작성함 *//

type expTok struct {
	kind  token.TokenKind
	value string // value 검증이 필요 없으면 ""로 둠
}

func assertTok(t *testing.T, got token.Token, want expTok) {
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

func lexAll(t *testing.T, input string) []token.Token {
	t.Helper()
	lx := NewLexer()
	lx.Set(input)

	out := []token.Token{}
	for {
		tok := lx.Next()
		out = append(out, tok)
		if tok.Kind == token.EOF {
			break
		}
	}
	return out
}

func TestLexer_Keywords_And_Identifiers(t *testing.T) {
	// EBNF에 필요한 키워드들(현재 TokenKind에 있는 것들만):
	// bool/int/string, if/else, for/range, let/in, scan/print, true/false, func/return
	toks := lexAll(t, "ok var bool int string if else for range let in scan print true false abc xyz123 func return len()")

	want := []expTok{
		{token.OK, "ok"},
		{token.VAR, "var"},
		{token.BOOL, "bool"},
		{token.INT, "int"},
		{token.STRING, "string"},
		{token.IF, "if"},
		{token.ELSE, "else"},
		{token.FOR, "for"},
		{token.RANGE, "range"},
		{token.LET, "let"},
		{token.IN, "in"},
		{token.SCAN, "scan"},
		{token.PRINT, "print"},
		{token.TRUE, "true"},
		{token.FALSE, "false"},
		{token.ID, "abc"},
		{token.ID, "xyz123"},
		{token.FUNC, "func"},
		{token.RETURN, "return"},
		{token.LEN, "len"},
		{token.OMIT, "()"},
		{token.EOF, "<<EOF>>"},
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
		{token.NUMBER, "0"},
		{token.NUMBER, "7"},
		{token.NUMBER, "42"},
		{token.NUMBER, "12345"},
		{token.EOF, "<<EOF>>"},
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
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.LBRACKET, "["},
		{token.RBRACKET, "]"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.COMMA, ","},
		{token.EOF, "<<EOF>>"},
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
	toks := lexAll(t, "= < > ! + - * /")

	want := []expTok{
		{token.ASSIGN, "="},
		{token.LT, "<"},
		{token.GT, ">"},
		{token.NOT, "!"},
		{token.PLUS, "+"},
		{token.MINUS, "-"},
		{token.MUL, "*"},
		{token.DIV, "/"},
		{token.EOF, "<<EOF>>"},
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
	toks := lexAll(t, "== != <= >= && || :=")

	want := []expTok{
		{token.EQUAL, "=="},
		{token.NEQ, "!="},
		{token.LTE, "<="},
		{token.GTE, ">="},
		{token.AND, "&&"},
		{token.OR, "||"},
		{token.DECLSIGN, ":="},
		{token.EOF, "<<EOF>>"},
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
		{token.NOT, "!"},
		{token.ID, "x"},
		{token.EOF, "<<EOF>>"},
	}

	if len(toks) != len(want) {
		t.Fatalf("token count mismatch: got=%d want=%d; got=%v", len(toks), len(want), toks)
	}
	for i := range want {
		assertTok(t, toks[i], want[i])
	}
}

func TestLexer_StringLiteral_STRLIT(t *testing.T) {
	// Goal: strlit이 Satom -> strlit 로 정의되어 있으므로
	// lexer가 "..." 전체를 STRLIT 하나로 토큰화하는지 확인한다.
	//
	// 이 테스트가 통과하려면 Lexer에 문자열 리터럴 읽기 로직이 있어야 함:
	// - 시작 따옴표(")를 만나면
	// - 다음 따옴표(")가 나올 때까지 내용을 읽고
	// - STRLIT 토큰(값은 hello 같은 내부 내용 또는 "hello" 전체 중 설계에 맞는 것)으로 반환

	toks := lexAll(t, "\"hello\" \"a b\" \"\"")

	want := []expTok{
		// value는 구현 정책에 따라 다를 수 있음:
		// 1) 내부만 저장: hello / a b / ""
		// 2) 따옴표 포함 저장: "hello" / "a b" / "\"\""
		//
		// 여기서는 "내부만 저장" 정책을 가정한다.
		{token.STRLIT, "hello"},
		{token.STRLIT, "a b"},
		{token.STRLIT, ""},
		{token.EOF, "<<EOF>>"},
	}

	if len(toks) != len(want) {
		t.Fatalf("token count mismatch: got=%d want=%d; got=%v", len(toks), len(want), toks)
	}
	for i := range want {
		assertTok(t, toks[i], want[i])
	}
}

func TestLexer_Whitespace_Is_Ignored(t *testing.T) {
	toks := lexAll(t, " var my_x_ error = newError(\"err\") \n\t  if   true   { print errString(my_x_); }  ")

	want := []expTok{
		{token.VAR, "var"},
		{token.ID, "my_x_"},
		{token.ERROR, "error"},
		{token.ASSIGN, "="},
		{token.NEWERROR, "newError"},
		{token.LPAREN, "("},
		{token.STRLIT, "err"},
		{token.RPAREN, ")"},
		{token.IF, "if"},
		{token.TRUE, "true"},
		{token.LBRACE, "{"},
		{token.PRINT, "print"},

		{token.ERRSTRING, "errString"},
		{token.LPAREN, "("},
		{token.ID, "my_x_"},
		{token.RPAREN, ")"},

		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.EOF, "<<EOF>>"},
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
