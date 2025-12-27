package parser

import (
	"github.com/rlaaudgjs5638/langTest/tinygo/token"

	"github.com/rlaaudgjs5638/langTest/tinygo/lexer"
)

type TokenTape struct {
	lexer       *lexer.Lexer
	parser      *Parser
	tokenRecord []token.Token
	currentIdx  int
}

// NewTokenTape는 렉서에서 첫 토큰을 받으면서 초기화됨
func NewTokenTape(l *lexer.Lexer) *TokenTape {
	initToken := l.Next()
	return &TokenTape{
		lexer:       l,
		tokenRecord: []token.Token{initToken},
		currentIdx:  0,
	}
}
func (t *TokenTape) SetParser(p *Parser) {
	t.parser = p
}

// 현재 토큰 보기
func (t *TokenTape) CurrentToken() token.Token {
	return t.tokenRecord[t.currentIdx]
}

// 앞의 토큰 미리보기
func (t *TokenTape) Peek(n int) token.Token {
	if t.isNextNTokensExistOnRecord(n) {
		return t.tokenRecord[t.currentIdx+n]
	}
	for i := range n {
		// k번쨰 앞의 토큰 검사 = 1~k까지의 수
		// 그러므로 i= n+1
		if t.isNextNTokensExistOnRecord(i + 1) {
			continue
		}
		nextToken := t.lexer.Next()
		t.tokenRecord = append(t.tokenRecord, nextToken)
	}
	return t.tokenRecord[t.currentIdx+n]
}

// MoveToNextToken은 다음 토큰을 읽어들어서 자신의 레코드에 기록한다.
func (t *TokenTape) MoveToNextToken() {
	if t.isNextTokenExistOnRecord() {
		t.currentIdx++
		return
	}
	nextToken := t.lexer.Next()
	t.tokenRecord = append(t.tokenRecord, nextToken)
	t.currentIdx++
	return
}

func (t *TokenTape) isNextTokenExistOnRecord() bool {
	return t.currentIdx < len(t.tokenRecord)-1
}
func (t *TokenTape) isNextNTokensExistOnRecord(n int) bool {
	return t.currentIdx+n < len(t.tokenRecord)
}

// GetRollback은 호출 시, 지금 마킹한 시점의 위치로 복귀할 수 있는, 함수를 제공한다.
func (t *TokenTape) GetRollback() func() {
	memorizedPosition := t.currentIdx
	//id카운터도 테이프에 맞게 롤백됨.
	//id카운터가 정상작동했다면, 다음에 부여받을 id가 뭐였을지 기억
	currentIdId := t.parser.idIdCounter.ViewCurrentId()
	return func() {
		t.currentIdx = memorizedPosition
		t.parser.idIdCounter.SetCurrentId(currentIdId)
	}
}
