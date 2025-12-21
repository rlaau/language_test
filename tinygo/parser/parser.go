package parser

import (
	"errors"
	"fmt"

	"github.com/rlaaudgjs5638/langTest/tinygo/lexer"
	"github.com/rlaaudgjs5638/langTest/tinygo/token"
)

type Parser struct {
	errorManager *errManager
	tape         *TokenTape
}

// ! EOF 토큰은 추후 고려
func NewParser(l *lexer.Lexer) *Parser {
	tokenManager := NewTokenManager(l)
	errManager := NewErrManager()
	return &Parser{
		tape:         tokenManager,
		errorManager: errManager,
	}
}

func (p *Parser) match(t token.TokenKind) error {
	if p.tape.CurrentToken().Kind == t {
		p.tape.MoveToNextToken()
		return nil
	}
	matchErrorMsg := fmt.Sprintf("p.match: 토큰 미스매치. 받은 토큰 %v != 테이프 토큰%v", t, p.tape.CurrentToken())
	return errors.New(matchErrorMsg)
}

// ErrNotProcesable 은 파서가 파싱 불가 상태 마주 시 리턴할 에러이다.
var ErrNotProcesable error = errors.New("파서가 현재 위치에서는 더 이상 파싱을 진행할 수 없습니다. (EOF or ILLEGAL)")

// CheckProcessable은 파싱 준비에 앞서, 파서가 더 진행 가능한지를 확인한다.
func (p *Parser) CheckProcessable() bool {
	return IsProcesable(p.tape.CurrentToken())
}

func IsProcesable(t token.Token) bool {
	return !IsEof(t) && !IsIllegal(t)
}
func IsEof(t token.Token) bool {
	return t.Kind == token.EOF
}
func IsIllegal(t token.Token) bool {
	return t.Kind == token.ILLLEGAL
}
