package calclang

import (
	"errors"
	"fmt"
)

type Parser struct {
	tokenList  []Token
	currentPos int
}

func NewParser() *Parser {
	return &Parser{}
}

type Symbol int

const (
	Expr Symbol = iota

	Bexp
	Bterm
	Bfact
	Batom

	Relop

	Aexp
	Term
	Power
	Factor
	Atom

	Int
)

func (p *Parser) Parse(ts []Token) {
	p.tokenList = ts
	p.currentPos = 0
	boolOrInt, err := p.parseExpr()
	if err != nil {
		fmt.Printf("파싱 중 에러 발생 %s\n", err.Error())
	}
	switch boolOrInt.(type) {
	case int:
		fmt.Printf("파싱 완료: 값= %d\n", boolOrInt)
	case bool:
		fmt.Printf("파싱 완료: 값= %v\n", boolOrInt)
	default:
		panic("논리 오류")
	}
	return
}

func (p *Parser) parseExpr() (any, error) {
	rb := p.rollback()
	b, err := p.parseBexp()
	if err == nil {
		return b, nil
	}
	rb()
	i, err := p.parseAexp()
	if err == nil {
		return i, nil
	}
	return nil, err
}

// 매치 검사 후 매치 시 한 칸 이동
func (p *Parser) match(t TokenType) (Token, bool) {
	if p.currentPos >= len(p.tokenList) {
		return Token{Type: EOF}, false
	}
	currentToken := p.tokenList[p.currentPos]
	if currentToken.Type == t {
		// 매치 시 한 칸 이동
		p.currentPos += 1
		return currentToken, true
	}
	return Token{}, false
}

func (p *Parser) matchOnly(t TokenType) bool {
	_, flag := p.match(t)
	return flag
}

// p의 당시 위치를 클로저로 메모 후, 호출 시 롤백 기능
func (p *Parser) rollback() func() {
	// p의 현재 위치를 클로저로 기억
	memo := p.currentPos
	return func() {
		p.currentPos = memo
	}
}

func (p *Parser) parseBexp() (bool, error) {
	v, err := p.parseBterm()
	if err != nil {
		return v, err
	}
	result := v
	for p.matchOnly(OR) {
		v, err := p.parseBterm()
		if err != nil {
			return v, err
		}
		result = result || v
	}
	return result, nil
}

func (p *Parser) parseBterm() (bool, error) {
	rb := p.rollback()
	v, err := p.parseBfact()
	if err != nil {
		rb()
		return v, err
	}
	result := v
	for p.matchOnly(AND) {
		v, err := p.parseBfact()
		if err != nil {
			rb()
			return v, err
		}
		result = result && v
	}

	return result, nil
}

func (p *Parser) parseBfact() (bool, error) {
	rb := p.rollback()
	isNot := p.matchOnly(NOT)
	v, err := p.parseBatom()
	if err == nil {
		if isNot {
			return !v, nil
		}
		return v, nil
	}
	rb()
	return v, err
}

func (p *Parser) parseBatom() (bool, error) {
	//작업 실패 시 롤백 하기 위해 작업 전 현재위치 기억
	rb := p.rollback()
	if p.matchOnly(TRUE) {
		return true, nil
	}
	if p.matchOnly(FALSE) {
		return false, nil
	}

	// 큰 작업 트랜잭션 선언
	parseAexpRelopAexp := func() (bool, error) {
		v1, err := p.parseAexp()
		if err != nil {
			return false, err
		}
		v2, err := p.parseRelop()
		if err != nil {
			return false, err
		}
		v3, err := p.parseAexp()
		if err != nil {
			return false, err
		}
		return computeCompare(v1, v2, v3)
	}
	v, err := parseAexpRelopAexp()
	// 정상적으로 파싱될 시에 값 리턴
	if err == nil {
		return v, nil
	}
	//에러 발생 시 상태 롤백
	rb()
	parseWrappedBexp := func() (bool, error) {
		if !p.matchOnly(LPAREN) {
			return v, LparenErr
		}
		v, err := p.parseBexp()
		if err != nil {
			return v, BexpERR
		}
		if !p.matchOnly(RPAREN) {
			return v, RparenErr
		}
		return v, nil
	}
	v, err = parseWrappedBexp()
	// 정상적 파싱 시 값 리턴
	if err == nil {
		return v, nil
	}
	//에러시 롤백
	rb()
	return false, err
}

func computeCompare(left int, binary TokenType, right int) (bool, error) {
	switch binary {
	case EQ:
		return left == right, nil
	case NEQ:
		return left != right, nil
	case GT:
		return left > right, nil
	case GTE:
		return left >= right, nil
	case LT:
		return left < right, nil
	case LTE:
		return left <= right, nil
	default:
		return false, errors.New("cant compare by binary")
	}

}
func (p *Parser) parseRelop() (TokenType, error) {
	for _, token := range []TokenType{EQ, NEQ, GT, GTE, LT, LTE} {
		_, isOk := p.match(token)
		if isOk {
			return token, nil
		}
	}
	return EQ, RelopErr
}
func (p *Parser) parseAexp() (int, error) {
	v, err := p.parseTerm()
	if err != nil {
		return v, err
	}
	for {
		if p.matchOnly(PLUS) {
			added, err := p.parseTerm()
			if err != nil {
				return added, err
			}
			v += added
			continue
		}
		if p.matchOnly(MINUS) {
			added, err := p.parseTerm()
			if err != nil {
				return added, err
			}
			v -= added
			continue
		}
		break
	}
	return v, nil
}

func (p *Parser) parseTerm() (int, error) {
	v, err := p.parsePower()
	if err != nil {
		return v, err
	}
	for {
		if p.matchOnly(MUL) {
			added, err := p.parsePower()
			if err != nil {
				return added, err
			}
			v *= added
			continue
		}
		if p.matchOnly(DIV) {
			added, err := p.parsePower()
			if err != nil {
				return added, err
			}
			v /= added
			continue
		}
		break
	}
	return v, nil
}

func (p *Parser) parsePower() (int, error) {
	powList := []int{}
	v, err := p.parseFactor()
	if err != nil {
		return v, err
	}
	powList = append(powList, v)

	for p.matchOnly(POW) {
		po, err := p.parseFactor()
		if err != nil {
			return po, err
		}
		powList = append(powList, po)

	}
	//우측 유도 연산
	result := rightMostPow(powList)
	return result, nil
}

func rightMostPow(a []int) int {
	if len(a) == 1 {
		return a[0]
	}
	return pow(a[0], rightMostPow(a[1:]))
}

func pow(a int, b int) int {
	total := 1
	for range b {
		total *= a
	}
	return total
}
func (p *Parser) parseFactor() (int, error) {
	_, isMinus := p.match(MINUS)
	v, err := p.parseAtom()
	if err != nil {
		return 0, err
	}
	if isMinus {
		return -1 * v, nil
	} else {
		return v, nil
	}

}
func (p *Parser) parseAtom() (int, error) {

	digitToken, isOk := p.match(INT)
	if isOk {
		return evalDigitString(string(digitToken.Lexed)), nil
	}

	rb := p.rollback()
	parseWrappedAexp := func() (int, error) {
		if _, isOk := p.match(LPAREN); !isOk {
			return 0, makeParsingError("wrapped aexp")
		}
		v, err := p.parseAexp()
		if err != nil {
			rb()
			return 0, err
		}
		if _, isOk := p.match(RPAREN); !isOk {
			rb()
			return 0, makeParsingError("wrapped aexp")
		}
		return v, nil
	}
	return parseWrappedAexp()
}

func evalDigitString(s string) int {
	digitString := make([]int, len(s))
	for i := range s {
		digitString[i] = int(s[i] - '0')
	}
	total := 0
	pow10 := 1
	for i := range len(digitString) {
		v := digitString[len(digitString)-i-1] * pow10
		total += v
		pow10 *= 10
	}
	return total
}

type ParseError error

func makeParsingError(s string) error {
	errMsg := fmt.Sprintf("%s parsing error", s)
	return errors.New(errMsg)
}

var BexpERR ParseError = makeParsingError("bexp")
var RelopErr ParseError = makeParsingError("relop")
var LparenErr ParseError = makeParsingError("lparen")
var RparenErr ParseError = makeParsingError("rparen")
