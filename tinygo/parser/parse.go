package parser

import (
	"errors"

	"github.com/rlaaudgjs5638/langTest/tinygo/token"
)

func (p *Parser) ParsePackage() (*Package, error) { panic("not implemented") }

func (p *Parser) parseDecl() (Decl, error)          { panic("") }
func (p *Parser) parseVarDecl() (*ValDecl, error)   { panic("") }
func (p *Parser) parseId() (*Id, error)             { panic("") }
func (p *Parser) parseFuncDecl() (*FuncDecl, error) { panic("") }
func (p *Parser) parseParam() (*Param, error)       { panic("") }
func (p *Parser) parseType() (*Type, error)         { panic("") }

func (p *Parser) parseStmt() (Stmt, error)                    { panic("") }
func (p *Parser) parseAssign() (*Assign, error)               { panic("") }
func (p *Parser) parseCallStmt() (*CallStmt, error)           { panic("") }
func (p *Parser) parseShortDecl() (*ShortDecl, error)         { panic("") }
func (p *Parser) parseReturn() (*Return, error)               { panic("") }
func (p *Parser) parseIf() (*If, error)                       { panic("") }
func (p *Parser) parseForBexp() (*ForBexp, error)             { panic("") }
func (p *Parser) parseForRangeAexp() (*ForRangeAexp, error)   { panic("") }
func (p *Parser) parseForWithAssign() (*ForWithAssign, error) { panic("") }
func (p *Parser) parseBlock() (*Block, error)                 { panic("") }

func (p *Parser) parseExpr() (Expr, error) { panic("") }

func (p *Parser) parseFexp() (*Fexp, error) { panic("") }

// Begin: Binary, Unary에 의한 추상 구문법으로 압축

func (p *Parser) parseLexp() (Lexp, error) {
	if !p.CheckProcessable() {
		return nil, NewParseError("Lexp", ErrNotProcesable)
	}

	lexp, err := p.parseBexp()
	if err != nil {
		return nil, NewParseError("Lexp", err)
	}

	for {
		if p.match(token.OR) == nil {
			right, err := p.parseBexp()
			if err != nil {
				return nil, NewParseError("Lexp", err)
			}
			lexp = newBinary(Or, lexp, right)
			continue
		}
		break
	}
	return lexp, nil

}
func (p *Parser) parseBexp() (Lexp, error) {
	if !p.CheckProcessable() {
		return nil, NewParseError("Bexp", ErrNotProcesable)
	}

	bexp, err := p.parseBterm()
	if err != nil {
		return nil, NewParseError("Bexp", err)
	}

	for {
		if p.match(token.AND) == nil {
			bterm, err := p.parseBterm()
			if err != nil {
				return nil, NewParseError("Bexp", err)
			}
			bexp = newBinary(And, bexp, bterm)
			continue
		}
		break
	}

	return bexp, nil

}

func (p *Parser) parseBterm() (Lexp, error) {
	if !p.CheckProcessable() {
		return nil, NewParseError("Bterm", ErrNotProcesable)
	}
	if p.match(token.NOT) == nil {
		bterm, err := p.parseBterm()
		if err != nil {
			return nil, NewParseError("Bexp", err)
		}
		return newUnary(Not, bterm), nil
	}

	aexp, err := p.parseAexp()
	if err != nil {
		return nil, NewParseError("Bexp", err)
	}

	matchRelop := func() (BinaryKind, error) {
		switch p.tape.CurrentToken().Kind {
		case token.EQUAL:
			return Equal, nil
		case token.NEQ:
			return NotEqual, nil
		case token.GT:
			return GreaterThan, nil
		case token.GTE:
			return GreaterOrEqual, nil
		case token.LT:
			return LessThan, nil
		case token.LTE:
			return LessOrEqual, nil
		default:
			return Equal, errors.New("isMatchRelop: RelOp 미스매치")
		}
	}

	if relop, err := matchRelop(); err == nil {
		secondAexp, err := p.parseAexp()
		if err != nil {
			return nil, NewParseError("Bterm", err)
		}
		return newBinary(relop, aexp, secondAexp), nil
	}
	return aexp, nil
}

func (p *Parser) parseAexp() (Lexp, error) {
	if !p.CheckProcessable() {
		return nil, NewParseError("Aexp", ErrNotProcesable)
	}
	binary, err := p.parseTerm()
	if err != nil {
		return nil, NewParseError("Aexp", err)
	}

	for {
		if p.match(token.PLUS) == nil {
			term, err := p.parseTerm()
			if err != nil {
				return nil, NewParseError("Bexp", err)
			}
			binary = newBinary(Plus, binary, term)
			continue
		}
		if p.match(token.MINUS) == nil {
			term, err := p.parseTerm()
			if err != nil {
				return nil, NewParseError("Bexp", nil)
			}
			binary = newBinary(MinusBinary, binary, term)
			continue
		}
		break
	}
	return binary, nil
}
func (p *Parser) parseTerm() (Lexp, error) {
	if !p.CheckProcessable() {
		return nil, NewParseError("Term", ErrNotProcesable)
	}

	binary, err := p.parseFactor()
	if err != nil {
		return nil, NewParseError("Term", err)
	}
	for {
		if p.match(token.MUL) == nil {
			factor, err := p.parseFactor()
			if err != nil {
				return nil, NewParseError("Term", err)
			}
			binary = newBinary(Mul, binary, factor)
			continue
		}
		if p.match(token.DIV) == nil {
			factor, err := p.parseFactor()
			if err != nil {
				return nil, NewParseError("Term", err)
			}
			binary = newBinary(Div, binary, factor)
			continue
		}
		break
	}

	return binary, nil
}
func (p *Parser) parseFactor() (Lexp, error) {
	if !p.CheckProcessable() {
		return nil, NewParseError("Factor", ErrNotProcesable)
	}
	isMinus := false
	if p.match(token.MINUS) == nil {
		isMinus = true
	}

	atom, err := p.parseAtom()
	if err != nil {
		return nil, NewParseError("Factor", err)
	}

	if isMinus {
		return newUnary(MinusUnary, atom), nil
	}
	return atom, nil
}

func (p *Parser) parseAtom() (Lexp, error) {
	if !p.CheckProcessable() {
		return nil, NewParseError("Atom", ErrNotProcesable)
	}

	//1, "("+Expr+")" case: 이 경우 first가 Call의 "("+Expr+")"+Args와 겹칩
	if p.match(token.LPAREN) == nil {
		_, err := p.parseExpr()
		if err != nil {
			return nil, NewParseError("Atom", err)
		}
		err = p.match(token.RPAREN)
		if err != nil {
			return nil, NewParseError("Atom", err)
		}
		//Expr 이후 Args가 존재하지 않는 경우를 먼저 체크
		peeked := p.tape.Peek(1)
		if !(peeked.Kind == token.OMIT) && !(peeked.Kind == token.LPAREN) {
			return nil, nil
		}
	}
	panic("잠시 대기")
}

// End
func (p *Parser) parseCall() (*Call, error)           { panic("") }
func (p *Parser) parseValueForm() (*ValueForm, error) { panic("") }
func (p *Parser) parseUnary() (*Unary, error)         { panic("") }
func (p *Parser) parseBinary() (*Binary, error)       { panic("") }
