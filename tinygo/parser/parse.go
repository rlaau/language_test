package parser

import (
	"fmt"

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
	if p.CheckProcessable() {
		return nil, ErrNotProcesable
	}

	if err := p.match(token.NOT); err == nil {
		lexp, err := p.parseLexp()
		if err != nil {
			return nil, fmt.Errorf("parseLexp failed: (%w)", err)
		}
		return &Unary{
			Op:     Not,
			Object: lexp,
		}, nil
	}

	firstBexp, err := p.parseBexp()
	if err != nil {
		return nil, fmt.Errorf("parseLexp: (%w)", err)
	}

	var bigBinary *Binary
	bigBinary.LeftExpr = nil
	bigBinary.RightExpr = nil
	nextLeftBinary := firstBexp
	buildBiggerBinary := func(op BinaryKind) error {
		newBexp, err := p.parseBexp()
		if err != nil {
			return err
		}
		bigBinary.LeftExpr = nextLeftBinary
		bigBinary.Op = op
		bigBinary.RightExpr = newBexp

		nextLeftBinary = bigBinary
		return nil
	}

	for {
		andErr := p.match(token.AND)
		if andErr == nil {
			if err := buildBiggerBinary(And); err != nil {
				return nil, err
			}
			continue
		}
		orErr := p.match(token.OR)
		if orErr == nil {
			if err := buildBiggerBinary(Or); err != nil {
				return nil, err
			}
			continue
		}
		break
	}
	if bigBinary.RightExpr == nil {
		return firstBexp, nil
	}
	return bigBinary, nil

}
func (p *Parser) parseBexp() (Lexp, error)   { panic("") }
func (p *Parser) parseAexp() (Lexp, error)   { panic("") }
func (p *Parser) parseTerm() (Lexp, error)   { panic("") }
func (p *Parser) parseFactor() (Lexp, error) { panic("") }
func (p *Parser) parseAtom() (*Atom, error)  { panic("") }

// End
func (p *Parser) parseCall() (*Call, error)           { panic("") }
func (p *Parser) parseValueForm() (*ValueForm, error) { panic("") }
func (p *Parser) parseUnary() (*Unary, error)         { panic("") }
func (p *Parser) parseBinary() (*Binary, error)       { panic("") }
