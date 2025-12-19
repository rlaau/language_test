package parser

func (p *Parser) parsePackage() *Package { panic("not implemented") }

func (p *Parser) parseDecl() Decl          { panic("") }
func (p *Parser) parseVarDecl() *ValDecl   { panic("") }
func (p *Parser) parseId() *Id             { panic("") }
func (p *Parser) parseFuncDecl() *FuncDecl { panic("") }
func (p *Parser) parseParam() *Param       { panic("") }
func (p *Parser) parseType() *Type         { panic("") }

func (p *Parser) parseStmt() Stmt                    { panic("") }
func (p *Parser) parseAssign() *Assign               { panic("") }
func (p *Parser) parseCallStmt() *CallStmt           { panic("") }
func (p *Parser) parseShortDecl() *ShortDecl         { panic("") }
func (p *Parser) parseReturn() *Return               { panic("") }
func (p *Parser) parseIf() *If                       { panic("") }
func (p *Parser) parseForBexp() *ForBexp             { panic("") }
func (p *Parser) parseForRangeAexp() *ForRangeAexp   { panic("") }
func (p *Parser) parseForWithAssign() *ForWithAssign { panic("") }
func (p *Parser) parseBlock() *Block                 { panic("") }

func (p *Parser) parseExpr() Expr { panic("") }

func (p *Parser) parseFexp() *Fexp { panic("") }

// Begin: Binary, Unary에 의한 추상 구문법으로 압축
func (p *Parser) parseLexp() Lexp   { panic("") }
func (p *Parser) parseBexp() Lexp   { panic("") }
func (p *Parser) parseAexp() Lexp   { panic("") }
func (p *Parser) parseTerm() Lexp   { panic("") }
func (p *Parser) parseFactor() Lexp { panic("") }

// End

func (p *Parser) parseAtom() *Atom           { panic("") }
func (p *Parser) parseCall() *Call           { panic("") }
func (p *Parser) parseValueForm() *ValueForm { panic("") }
func (p *Parser) parseUnary() *Unary         { panic("") }
func (p *Parser) parseBinary() *Binary       { panic("") }
