package parser

import (
	"errors"
	"strconv"

	"github.com/rlaaudgjs5638/langTest/tinygo/token"
)

func (p *Parser) ParsePackage() (*Package, error) {
	//! 최상단 패키지에선 받은 에러가 IsEof시에는 non-error로 처리하기
	if !p.CheckProcessable() {
		if IsEof(p.CurrentToken()) {
			return nil, nil
		}
		return nil, NewParseError("Expr", ErrNotProcesable)
	}

	panic("not implemented")
}

func (p *Parser) parseDecl() (Decl, error)          { panic("") }
func (p *Parser) parseVarDecl() (*ValDecl, error)   { panic("") }
func (p *Parser) parseFuncDecl() (*FuncDecl, error) { panic("") }

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

// Begin: Binary, Unary에 의한 추상 구문법으로 압축

func (p *Parser) ParseExpr() (Expr, error) {
	if !p.CheckProcessable() {
		return nil, NewParseError("Expr", ErrNotProcesable)
	}

	expr, err := p.parseBexp()
	if err != nil {
		return nil, NewParseError("Expr", err)
	}

	for {
		if p.match(token.OR) == nil {
			right, err := p.parseBexp()
			if err != nil {
				return nil, NewParseError("Expr", err)
			}
			expr = newBinary(Or, expr, right)
			continue
		}
		break
	}
	return expr, nil

}
func (p *Parser) parseBexp() (Expr, error) {
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

func (p *Parser) parseBterm() (Expr, error) {
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

func (p *Parser) parseAexp() (Expr, error) {
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
func (p *Parser) parseTerm() (Expr, error) {
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
func (p *Parser) parseFactor() (Expr, error) {
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

func (p *Parser) parseAtom() (Atom, error) {
	if !p.CheckProcessable() {
		return nil, NewParseError("Atom", ErrNotProcesable)
	}

	// 빌트인 콜인 경우를 가장 먼저 검사
	rollBack := p.tape.GetRollback()
	builtInCall, err := p.parseBuiltInCall()
	if err == nil {
		return builtInCall, nil
	}
	rollBack()

	primary, err := p.parsePrimary()
	if err != nil {
		return nil, NewParseError("Atom", err)
	}

	argsOrZero := []Args{}
	for {
		rollBack := p.tape.GetRollback()
		args, err := p.parseArgs()
		if err != nil {
			rollBack()
			break
		}
		argsOrZero = append(argsOrZero, *args)
	}

	// args==0 인 경우 primary로 리턴
	if len(argsOrZero) == 0 {
		return primary, nil
	}

	// args >=1 인 경우 call로 리턴
	return newCall(false, primary, NewErrorBuild, argsOrZero), nil
}

func (p *Parser) parseBuiltInCall() (*Call, error) {
	var builtInSet = map[token.TokenKind]BuiltInKind{
		token.NEWERROR:  NewErrorBuild,
		token.ERRSTRING: ErrStringBuild,
		token.SCAN:      ScanBuild,
		token.PRINT:     PrintBuild,
		token.PANIC:     PanicBuild,
		token.LEN:       LenBuild,
	}
	checkIsBuilInToken := func(t token.Token) (BuiltInKind, bool) {
		builtInKind, ok := builtInSet[t.Kind]
		return builtInKind, ok
	}
	if builtInKind, isBuiltIn := checkIsBuilInToken(p.CurrentToken()); isBuiltIn {
		p.match(p.CurrentToken().Kind)
		args, err := p.parseArgs()
		if err != nil {
			return nil, NewParseError("BuiltInCall", errors.New("빌트인 함수는 반드시 하나의 인자 세트 받아야 함"))
		}
		//* 빌트인 함수는 단 하나의 인자만을 소비함
		newBuiltInCall := newCall(true, nil, builtInKind, []Args{*args})
		//빌트인 함수에 하나 이상의 인자 세트를 넣었는지 검사함.
		rollBack := p.tape.GetRollback()
		_, err = p.parseArgs()
		if err == nil {
			return nil, NewParseError("BuiltInCall", errors.New("빌트인 함수는 단 하나의 인자 셋만을 받음. 연쇄호출 불가"))
		}
		rollBack()
		// 에러가 난 상홍이라면, 빌트인 함수에 단 하나의 인자 세트만이 들어간 상황이라는 뜻이므로 정상 검증 완료.
		return newBuiltInCall, nil
	}
	return nil, NewParseError("BuiltInCall", errors.New("빌트인 키워드에 해당하지 않음"))
}

// End
func (p *Parser) parsePrimary() (*Primary, error) {
	if !p.CheckProcessable() {
		return nil, NewParseError("Atom", ErrNotProcesable)
	}
	if p.match(token.LPAREN) == nil {
		expr, err := p.ParseExpr()
		if err != nil {
			return nil, NewParseError("Primary", err)
		}
		if p.match(token.RPAREN) != nil {
			return nil, NewParseError("Primary", p.match(token.RPAREN))
		}
		return newPrimary(ExprPrimary, expr, nil, nil), nil
	}

	rb := p.tape.GetRollback()
	id, err := p.parseId()
	if err == nil {
		return newPrimary(IdPrimary, nil, id, nil), nil
	}
	rb()

	valueForm, err := p.parseValueForm()
	if err != nil {
		return nil, NewParseError("Primary", err)
	}
	return newPrimary(ValuePrimary, nil, nil, valueForm), nil

}
func (p *Parser) parseArgs() (*Args, error) {
	if !p.CheckProcessable() {
		return nil, NewParseError("Atom", ErrNotProcesable)
	}
	var args []Expr
	if p.match(token.OMIT) == nil {
		return newArgs(args), nil
	}
	if p.match(token.LPAREN) != nil {
		return nil, NewParseError("Primary", errors.New("Args파싱 실패: Omit이 아닌 Args는 반드시 LPAREN과 EXPR이 필요함"))
	}
	// Omit이 아닌 한, 반드시 하나 이상의 expr은 있어야 함
	exprOnce, err := p.ParseExpr()
	args = append(args, exprOnce)
	if err != nil {
		return nil, NewParseError("Primary", err)
	}

	for {
		if p.match(token.COMMA) != nil {
			break
		}
		expr, err := p.ParseExpr()
		if err != nil {
			return nil, NewParseError("Primary", err)
		}
		args = append(args, expr)
	}
	if p.match(token.RPAREN) != nil {
		return nil, NewParseError("Primary", err)
	}
	return newArgs(args), nil
}

func (p *Parser) parseValueForm() (*ValueForm, error) {
	if !p.CheckProcessable() {
		return nil, NewParseError("Atom", ErrNotProcesable)
	}
	currentToken := p.tape.CurrentToken()
	if currentToken.Kind == token.FUNC {
		fexp, err := p.parseFexp()
		if err != nil {
			return nil, NewParseError("ValueForm", err)
		}
		return newValueForm(FexpValue, nil, nil, nil, nil, fexp), nil
	}
	switch currentToken.Kind {
	case token.NUMBER:
		num, err := strconv.Atoi(currentToken.Value)
		if err != nil {
			return nil, NewParseError("ValueForm", err)
		}
		p.match(token.NUMBER)
		return newValueForm(NumberValue, &num, nil, nil, nil, nil), nil
	case token.TRUE:
		t := true
		p.match(token.TRUE)
		return newValueForm(BoolValue, nil, &t, nil, nil, nil), nil
	case token.FALSE:
		f := false
		p.match(token.FALSE)
		return newValueForm(BoolValue, nil, &f, nil, nil, nil), nil
	case token.STRLIT:
		p.match(token.STRLIT)
		return newValueForm(StrLitValue, nil, nil, &currentToken.Value, nil, nil), nil
	case token.OK:
		// OK 값의 에러는 값은 nil, 타입은 error인 value로 처리함
		p.match(token.OK)
		return newValueForm(ErrValue, nil, nil, nil, nil, nil), nil
	default:
		return nil, NewParseError("ValueForm", errors.New("ValueForm 파싱에서 케이스 미스매치 발생"))
	}
}
func (p *Parser) parseFexp() (*Fexp, error) {
	if !p.CheckProcessable() {
		return nil, NewParseError("Atom", ErrNotProcesable)
	}
	if p.match(token.FUNC) != nil {
		return nil, NewParseError("Fexp", errors.New("Fexp에서 func토큰 발견 실패"))
	}

	var params []Param
	if p.match(token.OMIT) != nil {
		if p.match(token.LPAREN) != nil {

			return nil, NewParseError("Fexp", errors.New("Fexp에서 Omit이 아님에도, params도 없음"))
		}
		//Omit이 아니라면 최소 하나의 param필요
		paramOnce, err := p.parseParam()
		if err != nil {
			return nil, NewParseError("Fexp", err)
		}
		params = append(params, *paramOnce)

		for {
			if p.match(token.COMMA) != nil {
				break
			}
			param, err := p.parseParam()
			if err != nil {
				return nil, NewParseError("Fexp", err)
			}
			params = append(params, *param)
		}

		if p.match(token.RPAREN) != nil {
			return nil, NewParseError("Fexp", errors.New("Fexp파싱에서 \")\"를 기대했지만 나오지 않음"))
		}
	}

	rollBack := p.tape.GetRollback()
	onlyType, err := p.parseType()
	if err == nil {
		block, err := p.parseBlock()
		if err != nil {
			return nil, NewParseError("Fexp", err)
		}
		return newFexp(params, []Type{*onlyType}, *block), nil
	}
	rollBack()
	types := []Type{}
	if p.match(token.LPAREN) == nil {
		typeOnce, err := p.parseType()
		types = append(types, *typeOnce)
		if err != nil {
			return nil, NewParseError("Fexp", errors.New("Fexp. Omit이 아닌 경우, 괄호 안에는 리턴 타입 하나 이상 필수입니다"))
		}
		for {
			if p.match(token.COMMA) != nil {
				break
			}
			t, err := p.parseType()
			if err != nil {
				return nil, NewParseError("Fexp", err)
			}
			types = append(types, *t)
		}

		if p.match(token.RPAREN) != nil {
			return nil, NewParseError("Fexp", errors.New("리턴 타입의 닫는 괄호 부재"))
		}
	}

	block, err := p.parseBlock()
	if err != nil {
		return nil, NewParseError("Fexp", err)
	}
	return newFexp(params, types, *block), nil
}

func (p *Parser) parseParam() (*Param, error) {
	if !p.CheckProcessable() {
		return nil, NewParseError("Param", ErrNotProcesable)
	}
	id, err := p.parseId()
	if err != nil {
		return nil, NewParseError("Param", err)
	}
	t, err := p.parseType()
	if err != nil {
		return nil, NewParseError("Param", err)
	}

	return newParam(*id, *t), nil
}

func (p *Parser) parseId() (*Id, error) {
	if !p.CheckProcessable() {
		return nil, NewParseError("Id", ErrNotProcesable)
	}
	if p.tape.CurrentToken().Kind == token.ID {
		id := newId(p.tape.CurrentToken())
		p.match(token.ID)
		return id, nil
	}
	return nil, NewParseError("Id", errors.New("ID파싱 실패"))
}
func (p *Parser) parseType() (*Type, error) {
	currentToken := p.tape.CurrentToken()
	switch currentToken.Kind {
	case token.INT:
		p.match(currentToken.Kind)
		return newType(IntType, nil), nil
	case token.BOOL:
		p.match(currentToken.Kind)
		return newType(BoolType, nil), nil
	case token.STRING:
		p.match(currentToken.Kind)
		return newType(StringType, nil), nil
	case token.ERROR:
		p.match(currentToken.Kind)
		return newType(ErrorType, nil), nil
	}

	funcType, err := p.parseFuncType()
	if err != nil {
		return nil, NewParseError("Type", err)
	}
	return newType(FuncionType, funcType), nil
}

func (p *Parser) parseFuncType() (*FuncType, error) {
	if p.match(token.FUNC) != nil {
		return nil, NewParseError("FuncType", errors.New("FuncType파싱 중 Func키워드 미발견"))
	}

	argTypes := []Type{}
	if p.match(token.OMIT) != nil {
		if p.match(token.LPAREN) != nil {
			return nil, NewParseError("FuncType", NewParseError("FuncType", errors.New("FuncType은 func이후 omit or \"(\" 가 와야 함.)")))
		}
		tOnce, err := p.parseType()
		if err != nil {
			return nil, NewParseError("FuncType", errors.New("Omit이 아닌 arg타입은 반드시 하나 이상의 타입을 포함해야 함"))
		}
		argTypes = append(argTypes, *tOnce)

		for {
			if p.match(token.COMMA) != nil {
				break
			}
			t, err := p.parseType()
			if err != nil {
				return nil, NewParseError("FuncType", err)
			}
			argTypes = append(argTypes, *t)
		}
		if p.match(token.RPAREN) != nil {
			return nil, NewParseError("FuncType", errors.New("닫는 괄호가 없음"))
		}
	}
	rollBack := p.tape.GetRollback()
	onlyType, err := p.parseType()
	if err == nil {
		return newFuncType(argTypes, []Type{*onlyType}), nil
	}
	// 단일 타입인 경우가 아니라면 롤백
	rollBack()

	returnTypes := []Type{}

	if p.match(token.LPAREN) != nil {
		//좌괄호가 없는 경우, 리턴 타입이 없는 funcType(EBNF참고)
		return newFuncType(argTypes, returnTypes), nil
	}
	returnOnce, err := p.parseType()
	if err != nil {
		return nil, NewParseError("FuncType", err)
	}
	returnTypes = append(returnTypes, *returnOnce)

	for {
		if p.match(token.COMMA) != nil {
			break
		}
		rt, err := p.parseType()
		if err != nil {
			return nil, NewParseError("FuncType", err)
		}
		returnTypes = append(returnTypes, *rt)
	}
	if p.match(token.RPAREN) != nil {
		return nil, NewParseError("FuncType", errors.New("리턴 타입의 닫는괄호 누락"))
	}
	return newFuncType(argTypes, returnTypes), nil

}
