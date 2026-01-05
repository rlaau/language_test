package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/rlaaudgjs5638/langTest/tinygo/token"
)

func (p *Parser) ParsePackage() (*PackageAST, error) {
	if !p.CheckProcessable() {
		return newPackage(nil), nil
	}
	decls := []Decl{}
	for {
		decl, err := p.parseDecl()
		if err != nil {
			errReasons := err.Error()
			if strings.Contains(errReasons, ErrNotProcesable.Error()) && p.CurrentToken().Kind == token.EOF {
				//더이상 파싱할 문자가 없어서 에러난 경우 정상 작동으로 처리
				break
			}
			return nil, NewParseError("Package", err)
		}
		decls = append(decls, decl)
	}

	return newPackage(decls), nil
}

func (p *Parser) parseDecl() (Decl, error) {
	if !p.CheckProcessable() {
		return nil, NewParseError("Decl", ErrNotProcesable)
	}
	if p.CurrentToken().Kind == token.VAR {
		varDecl, err := p.parseVarDecl()
		if err != nil {

			return nil, NewParseError("Decl", err)
		}
		return varDecl, nil
	}

	funcDecl, err := p.parseFuncDecl()
	if err != nil {
		return nil, NewParseError("Decl", err)
	}
	return funcDecl, nil

}
func (p *Parser) parseVarDecl() (*VarDecl, error) {
	if !p.CheckProcessable() {
		return nil, NewParseError("VarDecl", ErrNotProcesable)
	}
	if p.match(token.VAR) != nil {
		return nil, NewParseError("VarDecl", errors.New("VarDecl은 반드시 Var포함해야 함"))
	}

	ids, err := p.parseIdListLongerThan0()
	if err != nil {
		return nil, NewParseError("VarDecl", err)
	}
	typ, err := p.parseType()
	if err != nil {
		return nil, NewParseError("VarDecl", errors.New("VarDecl은 Type을 포함해야 합니다"))
	}

	if p.match(token.ASSIGN) != nil {
		//varDecl에서 명시적인 Expr을 할당하지 않은경우
		if p.match(token.SEMICOLON) != nil {
			return nil, NewParseError("VarDecl", ErrMissingSemicolon)
		}
		return newVarDecl(ids, *typ, []Expr{}), nil
	}
	exprs, err := p.parseExprListLongerThan0()
	if err != nil {
		return nil, NewParseError("VarDecl", err)
	}
	if p.match(token.SEMICOLON) != nil {
		return nil, NewParseError("VarDecl", ErrMissingSemicolon)
	}
	return newVarDecl(ids, *typ, exprs), nil

}
func (p *Parser) parseFuncDecl() (*FuncDecl, error) {
	if !p.CheckProcessable() {
		return nil, NewParseError("FuncDecl", ErrNotProcesable)
	}
	if p.match(token.FUNC) != nil {
		return nil, NewParseError("FuncDecl", errors.New("FuncDecl에서 Func키워드 누락"))
	}
	id, err := p.parseId()
	if err != nil {
		return nil, NewParseError("FuncDecl", err)
	}
	params, err := p.parseParams()
	if err != nil {
		return nil, NewParseError("FuncDecl", err)
	}
	rollBack := p.tape.GetRollback()
	returnTypesOrNil, err := p.parseReturnTypes()
	if err != nil {
		rollBack()
		returnTypesOrNil = []Type{}
	}
	block, err := p.parseBlock()
	if err != nil {
		return nil, NewParseError("FuncDecl", err)
	}

	return newFuncDecl(*id, params, returnTypesOrNil, *block), nil
}

func (p *Parser) parseStmt() (Stmt, error) {
	if !p.CheckProcessable() {
		return nil, NewParseError("Stmt", ErrNotProcesable)
	}

	rollBack := p.tape.GetRollback()
	switch p.CurrentToken().Kind {
	case token.ID:
		assign, err := p.parseAssign()
		if err == nil {
			return assign, nil
		}
		rollBack()
		shortDecl, err := p.parseShortDecl()
		if err == nil {
			return shortDecl, nil
		}
		rollBack()
		return p.parseCallStmt()
	case token.VAR:
		return p.parseVarDecl()
	case token.FUNC:
		return p.parseFuncDecl()
	case token.RETURN:
		return p.parseReturn()
	case token.BREAK:
		return p.parseBreak()
	case token.CONTINUE:
		return p.parseContinue()
	case token.IF:
		return p.parseIf()
	case token.FOR:
		forBexp, err := p.parseForBexp()
		if err == nil {
			return forBexp, nil
		}

		rollBack()
		return p.parseForWithAssign()
	case token.LBRACE:
		return p.parseBlock()
	default:
		return p.parseCallStmt()
	}
}

func (p *Parser) parseAssign() (*Assign, error) {
	if !p.CheckProcessable() {
		return nil, NewParseError("Assgin", ErrNotProcesable)
	}
	idList, err := p.parseIdListLongerThan0()
	if err != nil {
		return nil, NewParseError("Assign", err)
	}
	if p.match(token.ASSIGN) != nil {
		return nil, NewParseError("Assign", errors.New("\"=\"기호 부재"))
	}

	exprList, err := p.parseExprListLongerThan0()
	if err != nil {
		return nil, NewParseError("Assign", err)
	}
	if p.match(token.SEMICOLON) != nil {
		return nil, NewParseError("Assign", ErrMissingSemicolon)
	}
	return newAssign(idList, exprList), nil
}

func (p *Parser) parseCallStmt() (*CallStmt, error) {
	if !p.CheckProcessable() {
		return nil, NewParseError("CallStmt", ErrNotProcesable)
	}
	call, err := p.parseCall()
	if err != nil {
		return nil, NewParseError("CallStmt", err)
	}
	if p.match(token.SEMICOLON) != nil {
		return nil, NewParseError("CallStmt", ErrMissingSemicolon)
	}
	return newCallStmt(*call), nil
}
func (p *Parser) parseCall() (*Call, error) {
	if !p.CheckProcessable() {
		return nil, NewParseError("Call", ErrNotProcesable)
	}

	primary, err := p.parsePrimary()
	if err != nil {

		return nil, NewParseError("Call", err)
	}
	argsList := []Args{}
	for {
		rollBack2 := p.tape.GetRollback()
		args, err := p.parseArgs()
		if err != nil {
			rollBack2()
			break
		}
		argsList = append(argsList, *args)
	}

	if len(argsList) == 0 {
		return nil, NewParseError("Call", errors.New("Call은 Primary이후 하나 이상의 args가 와야 합나디."))
	}
	return newCall(*primary, argsList), nil
}

func (p *Parser) parseShortDecl() (*ShortDecl, error) {
	if !p.CheckProcessable() {
		return nil, NewParseError("ShortDecl", ErrNotProcesable)
	}
	idList, err := p.parseIdListLongerThan0()
	if err != nil {

		return nil, NewParseError("ShortDecl", err)
	}
	if p.match(token.DECLSIGN) != nil {
		return nil, NewParseError("ShortDecl", errors.New(":= 연산자 누락"))
	}
	exprList, err := p.parseExprListLongerThan0()
	if err != nil {
		return nil, NewParseError("ShortDecl", err)
	}
	if p.match(token.SEMICOLON) != nil {
		return nil, NewParseError("ShortDecl", ErrMissingSemicolon)
	}
	return newShortDecl(idList, exprList), nil
}

func (p *Parser) parseIdListLongerThan0() ([]Id, error) {
	if !p.CheckProcessable() {
		return nil, NewParseError("IdListLongerThan0", ErrNotProcesable)
	}
	ids := []Id{}
	idOnce, err := p.parseId()
	if err != nil {
		return nil, NewParseError("IdListLongerThan0", err)
	}
	ids = append(ids, *idOnce)
	for {
		if p.match(token.COMMA) != nil {
			break
		}
		id, err := p.parseId()
		if err != nil {
			return nil, NewParseError("IdListLongerThan0", err)
		}
		ids = append(ids, *id)
	}
	return ids, nil
}

func (p *Parser) parseExprListLongerThan0() ([]Expr, error) {
	if !p.CheckProcessable() {
		return nil, NewParseError("ExprListLongerThan0", ErrNotProcesable)
	}
	exprOnce, err := p.parseExpr()
	if err != nil {
		return nil, NewParseError("ExprListLongerThan0", err)
	}
	exprs := []Expr{}
	exprs = append(exprs, exprOnce)
	for {
		if p.match(token.COMMA) != nil {
			break
		}
		expr, err := p.parseExpr()
		if err != nil {
			return nil, NewParseError("ExprListLongerThan0", err)
		}
		exprs = append(exprs, expr)
	}
	return exprs, nil
}

func (p *Parser) parseReturn() (*Return, error) {
	if !p.CheckProcessable() {
		return nil, NewParseError("Return", ErrNotProcesable)
	}
	if p.match(token.RETURN) != nil {
		return nil, NewParseError("Return", errors.New("리터은 return 키워드로 시작해야 합니다."))
	}

	// 이번 단계에서 Expr로 파싱 가능한지 확인만 하기
	rollBack := p.tape.GetRollback()
	_, err := p.parseExpr()
	rollBack()
	if err != nil {
		//Expr로 파싱 실패 시, 리턴값 없는 것으로 취급
		if p.match(token.SEMICOLON) != nil {
			return nil, NewParseError("Return", ErrMissingSemicolon)
		}
		return newReturn([]Expr{}), nil
	}
	//여기서부턴 하나 이상의 리턴값 가진것이 담보됨
	exprs, err := p.parseExprListLongerThan0()
	if err != nil {
		return nil, NewParseError("Return", err)
	}
	if p.match(token.SEMICOLON) != nil {
		return nil, NewParseError("Return", ErrMissingSemicolon)
	}
	return newReturn(exprs), nil
}
func (p *Parser) parseBreak() (*Break, error) {
	if !p.CheckProcessable() {
		return nil, NewParseError("Break", ErrNotProcesable)
	}
	if p.match(token.BREAK) != nil {
		return nil, NewParseError("Break", errors.New("Break doesnt't match to \"break\""))
	}
	if p.match(token.SEMICOLON) != nil {
		return nil, NewParseError("Break", ErrMissingSemicolon)
	}
	return newBreak(), nil
}

func (p *Parser) parseContinue() (*Continue, error) {
	if !p.CheckProcessable() {
		return nil, NewParseError("Continue", ErrNotProcesable)
	}
	if p.match(token.CONTINUE) != nil {
		return nil, NewParseError("Continue", fmt.Errorf("Continue doesn't match to \"continue\""))
	}
	if p.match(token.SEMICOLON) != nil {
		return nil, NewParseError("Continue", ErrMissingSemicolon)
	}
	return newContinue(), nil

}
func (p *Parser) parseIf() (*If, error) {
	if !p.CheckProcessable() {
		return nil, NewParseError("If", ErrNotProcesable)
	}
	if p.match(token.IF) != nil {
		return nil, NewParseError("If", errors.New("If문은 if키워드로 시작해야 함"))
	}

	rollBack := p.tape.GetRollback()
	shortDeclOrNil, err := p.parseShortDecl()
	if err != nil {
		rollBack()
	}

	bexp, err := p.parseExpr()
	if err != nil {
		return nil, NewParseError("If", err)
	}

	thenBlock, err := p.parseBlock()
	if err != nil {
		return nil, NewParseError("If", err)
	}
	if p.match(token.ELSE) != nil {
		return newIf(shortDeclOrNil, bexp, *thenBlock, nil), nil
	}

	elseBlock, err := p.parseBlock()
	if err != nil {
		return nil, NewParseError("If", err)
	}
	return newIf(shortDeclOrNil, bexp, *thenBlock, elseBlock), nil
}

func (p *Parser) parseForBexp() (*ForBexp, error) {
	if !p.CheckProcessable() {
		return nil, NewParseError("ForBexp", ErrNotProcesable)
	}
	if p.match(token.FOR) != nil {
		return nil, NewParseError("ForBexp", errors.New("for키워드 누락"))
	}

	bexp, err := p.parseBexp()
	if err != nil {
		return nil, NewParseError("ForBexp", err)
	}
	block, err := p.parseBlock()
	if err != nil {
		return nil, NewParseError("ForBexp", err)
	}
	return newForBexp(bexp, *block), nil
}

func (p *Parser) parseForWithAssign() (*ForWithAssign, error) {
	if !p.CheckProcessable() {
		return nil, NewParseError("ForWithAssgin", ErrNotProcesable)
	}
	if p.match(token.FOR) != nil {
		return nil, NewParseError("ForWithAssign", errors.New("for키워드 누락"))
	}
	shortDecl, err := p.parseShortDecl()
	if err != nil {
		return nil, NewParseError("ForWithAssign", err)
	}
	bexp, err := p.parseBexp()
	if err != nil {
		return nil, NewParseError("ForWithAssign", err)
	}
	if p.match(token.SEMICOLON) != nil {
		return nil, NewParseError("ForWithAssign", ErrMissingSemicolon)
	}
	id, err := p.parseId()
	if err != nil {
		return nil, NewParseError("ForWithAssign", err)
	}
	if p.match(token.ASSIGN) != nil {
		return nil, NewParseError("ForWithAssign", errors.New(" \"=\" 기호 누락"))

	}

	expr, err := p.parseExpr()
	if err != nil {
		return nil, NewParseError("ForWithAssign", err)
	}
	if p.match(token.SEMICOLON) != nil {
		return nil, NewParseError("ForWithAssign", ErrMissingSemicolon)
	}
	assign := newAssign([]Id{*id}, []Expr{expr})
	block, err := p.parseBlock()
	if err != nil {
		return nil, NewParseError("ForWithAssign", err)
	}
	return newForWithAssign(*shortDecl, bexp, *assign, *block), nil
}

func (p *Parser) parseBlock() (*Block, error) {
	if !p.CheckProcessable() {
		return nil, NewParseError("Block", ErrNotProcesable)
	}
	if p.match(token.LBRACE) != nil {
		return nil, NewParseError("Block", errors.New("시작 위치에 \"{\" 기호가 존재하지 않음"))
	}
	stmts := []Stmt{}
	for {
		rollBack := p.tape.GetRollback()
		stms, err := p.parseStmt()
		if err != nil {
			rollBack()

			break
		}
		stmts = append(stmts, stms)
	}

	if p.match(token.RBRACE) != nil {
		return nil, NewParseError("Block", errors.New("맺음 위치에 \"}\"기호가 존재하지 않음."))
	}
	return newBlock(stmts), nil
}

// Begin: Binary, Unary에 의한 추상 구문법으로 압축

func (p *Parser) parseExpr() (Expr, error) {
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
			return nil, NewParseError("Bterm", err)
		}
		return newUnary(Not, bterm), nil
	}

	aexp, err := p.parseAexp()

	if err != nil {
		return nil, NewParseError("Btem", err)
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
		p.match(p.CurrentToken().Kind)
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
				return nil, NewParseError("Aexp", err)
			}
			binary = newBinary(Plus, binary, term)
			continue
		}
		if p.match(token.MINUS) == nil {
			term, err := p.parseTerm()
			if err != nil {
				return nil, NewParseError("Aexp", err)
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
	return newCall(*primary, argsOrZero), nil
}

// End
func (p *Parser) parsePrimary() (*Primary, error) {
	if !p.CheckProcessable() {
		return nil, NewParseError("Primary", ErrNotProcesable)
	}
	if p.match(token.LPAREN) == nil {
		expr, err := p.parseExpr()
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
		return nil, NewParseError("Args", ErrNotProcesable)
	}
	var args []Expr
	if p.match(token.OMIT) == nil {
		return newArgs(args), nil
	}
	if p.match(token.LPAREN) != nil {
		return nil, NewParseError("Args", errors.New("Args파싱 실패: Omit이 아닌 Args는 반드시 LPAREN과 EXPR이 필요함"))
	}
	// Omit이 아닌 한, 반드시 하나 이상의 expr은 있어야 함
	args, err := p.parseExprListLongerThan0()
	if err != nil {
		return nil, NewParseError("Args", err)
	}
	if p.match(token.RPAREN) != nil {
		return nil, NewParseError("Args", err)
	}
	return newArgs(args), nil
}

func (p *Parser) parseValueForm() (*ValueForm, error) {
	if !p.CheckProcessable() {
		return nil, NewParseError("ValueForm", ErrNotProcesable)
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
		return nil, NewParseError("Fexp", ErrNotProcesable)
	}
	if p.match(token.FUNC) != nil {
		return nil, NewParseError("Fexp", errors.New("Fexp에서 func토큰 발견 실패"))
	}

	params, err := p.parseParams()
	if err != nil {
		return nil, NewParseError("Fexp", err)
	}

	rollBack := p.tape.GetRollback()
	types, err := p.parseReturnTypes()
	if err != nil {
		types = []Type{}
		rollBack()
	}

	block, err := p.parseBlock()
	if err != nil {
		return nil, NewParseError("Fexp", err)
	}
	return newFexp(params, types, *block), nil
}

func (p *Parser) parseParams() ([]Param, error) {
	if !p.CheckProcessable() {
		return nil, NewParseError("Params", ErrNotProcesable)
	}
	var params []Param
	if p.match(token.OMIT) != nil {
		if p.match(token.LPAREN) != nil {

			return nil, NewParseError("Params", errors.New("Params에서 Omit이 아님에도, params도 없음"))
		}
		//Omit이 아니라면 최소 하나의 param필요
		paramOnce, err := p.parseParam()
		if err != nil {
			return nil, NewParseError("Params", err)
		}
		params = append(params, *paramOnce)

		for {
			if p.match(token.COMMA) != nil {
				break
			}
			param, err := p.parseParam()
			if err != nil {
				return nil, NewParseError("Params", err)
			}
			params = append(params, *param)
		}

		if p.match(token.RPAREN) != nil {
			return nil, NewParseError("Params", errors.New("Params파싱에서 \")\"를 기대했지만 나오지 않음"))
		}
	}
	return params, nil
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
		id := newId(p.tape.CurrentToken(), p.idIdCounter.GetNextID())
		p.match(token.ID)
		return id, nil
	}
	return nil, NewParseError("Id", errors.New("ID파싱 실패"))
}

func (p *Parser) parseReturnTypes() ([]Type, error) {
	if !p.CheckProcessable() {
		return nil, NewParseError("ReturnTypes", ErrNotProcesable)
	}
	rollBack := p.tape.GetRollback()
	if onlyType, err := p.parseType(); err == nil {
		//단일 타입으로 잘 파싱될 경우
		return []Type{*onlyType}, nil
	}
	rollBack()
	types := []Type{}
	if p.match(token.LPAREN) == nil {
		typeOnce, err := p.parseType()
		types = append(types, *typeOnce)
		if err != nil {
			return nil, NewParseError("ReturnTypes", errors.New("ReturnTypes. Omit이 아닌 경우, 괄호 안에는 리턴 타입 하나 이상 필수입니다"))
		}
		for {
			if p.match(token.COMMA) != nil {
				break
			}
			t, err := p.parseType()
			if err != nil {
				return nil, NewParseError("ReturnTypes", err)
			}
			types = append(types, *t)
		}

		if p.match(token.RPAREN) != nil {
			return nil, NewParseError("ReturnTypes", errors.New("리턴 타입의 닫는 괄호 부재"))
		}

		return types, nil
	}
	return nil, NewParseError("ReturnTypes", errors.New("ReturnTypes의 값이 존재하지 않음"))
}
func (p *Parser) parseType() (*Type, error) {
	if !p.CheckProcessable() {
		return nil, NewParseError("Type", ErrNotProcesable)
	}
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
	if !p.CheckProcessable() {
		return nil, NewParseError("FuncType", ErrNotProcesable)
	}
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
	rollback2 := p.tape.GetRollback()
	returnTypes, err = p.parseReturnTypes()
	if err != nil {
		returnTypes = []Type{}
		rollback2()
	}
	return newFuncType(argTypes, returnTypes), nil

}
