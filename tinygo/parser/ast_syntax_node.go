package parser

import (
	"fmt"
	"strconv"

	"github.com/rlaaudgjs5638/langTest/tinygo/token"
)

// Node
type Package struct {
	DeclsOrNil []Decl
}

func newPackage(declsOrNil []Decl) *Package {
	return &Package{
		DeclsOrNil: declsOrNil,
	}
}
func (p *Package) Print(depth int) []string {

	var pkgStrings []string
	pkgStart := LineWithDepth("Package Start -------", depth)
	pkgStrings = append(pkgStrings, pkgStart)

	for _, decl := range p.DeclsOrNil {
		stmts := decl.Print(depth)
		for _, stmt := range stmts {
			pkgStrings = append(pkgStrings, stmt)
		}
	}

	pkgStrings = append(pkgStrings, LineWithDepth("Package End -------", depth))
	return pkgStrings
}

func (p *Package) String() string {
	return JoinLines(p.Print(0))
}

// Decl
type VarDecl struct {
	Ids        []Id
	Type       Type
	ExprsOrNil []Expr
}

func newVarDecl(ids []Id, t Type, exprsOrNil []Expr) *VarDecl {
	return &VarDecl{
		Ids:        ids,
		Type:       t,
		ExprsOrNil: exprsOrNil,
	}
}

var _ Decl = (*VarDecl)(nil)

func (v *VarDecl) Print(depth int) []string {

	var lines []string
	start := "VarDecl("
	va := "var "
	typ := v.Type.String()
	ids := JoinWithSepG(v.Ids, ",")
	lines = append(lines, LineWithDepth(start, depth))
	lines = append(lines, LineWithDepth(va+ids, depth+1))
	lines = append(lines, LineWithDepth("type "+typ, depth+1))
	lines = append(lines, LineWithDepth("=", depth+1))

	if len(v.ExprsOrNil) == 0 {
		lines = append(lines, LineWithDepth("nothing assigned for var statement", depth+1))
	} else {
		for i, exp := range v.ExprsOrNil {
			lines = append(lines, exp.Print(depth+1)...)
			if i < len(v.ExprsOrNil)-1 {
				lines = append(lines, LineWithDepth(",", depth+1))
			}
		}
	}

	lines = append(lines, LineWithDepth(")", depth))
	return lines
}

func (v *VarDecl) String() string {
	return JoinLines(v.Print(0))
}
func (v *VarDecl) Decl() string {
	return v.String()
}

func (v *VarDecl) Stmt() string {
	return v.String()
}

// Decl
type FuncDecl struct {
	Id               Id
	ParamsOrNil      []Param
	ReturnTypesOrNil []Type
	Block            Block
}

func newFuncDecl(id Id, pOrNil []Param, rOrNil []Type, block Block) *FuncDecl {
	return &FuncDecl{
		Id:               id,
		ParamsOrNil:      pOrNil,
		ReturnTypesOrNil: rOrNil,
		Block:            block,
	}
}

var _ Decl = (*FuncDecl)(nil)

func (f *FuncDecl) Print(depth int) []string {
	start := "FuncDecl("
	lines := []string{}
	lines = append(lines, LineWithDepth(start, depth))
	lines = append(lines, LineWithDepth("ID:"+f.Id.String(), depth+1))
	paramStart := "Type: ["
	params := JoinWithSepG(f.ParamsOrNil, ",")
	paramEnd := "]"
	arrow := "=>"
	returnStart := "["
	returns := JoinWithSepG(f.ReturnTypesOrNil, ",")
	returnEnd := "]"
	funcTypeLine := LineWithDepth(paramStart+params+paramEnd+arrow+returnStart+returns+returnEnd, depth+1)
	lines = append(lines, funcTypeLine)

	blockLines := f.Block.Print(depth + 1)
	for _, bl := range blockLines {
		lines = append(lines, bl)
	}

	funcEndLine := LineWithDepth(")", depth)
	lines = append(lines, funcEndLine)

	return lines

}
func (f *FuncDecl) String() string {
	return JoinLines(f.Print(0))
}

func (f *FuncDecl) Decl() string {
	return f.String()
}

func (f *FuncDecl) Stmt() string {
	return f.String()
}

type Param struct {
	Id   Id
	Type Type
}

func newParam(id Id, t Type) *Param {
	return &Param{
		Id:   id,
		Type: t,
	}
}
func (p Param) String() string {
	return "<" + p.Id.String() + "," + p.Type.String() + ">"
}

type Id string

func newId(token token.Token) *Id {
	value := string(token.Value)
	id := Id(value)
	return &id
}
func (i Id) String() string {
	return string(i)
}

type Type struct {
	TypeKind      TypeKind
	FuncTypeOrNil *FuncType
}

func newType(kind TypeKind, funcTypeOrNil *FuncType) *Type {
	return &Type{
		TypeKind:      kind,
		FuncTypeOrNil: funcTypeOrNil,
	}
}
func (t Type) String() string {
	switch t.TypeKind {
	case IntType:
		return "int"
	case BoolType:
		return "bool"
	case StringType:
		return "string"
	case ErrorType:
		return "error"
	case FuncionType:
		return t.FuncTypeOrNil.String()
	default:
		panic("Type.String(): 스위치 미스매치")
	}
}

type TypeKind int

const (
	IntType TypeKind = iota
	BoolType
	StringType
	ErrorType
	FuncionType
)

type FuncType struct {
	ArgTypesOrNil    []Type
	ReturnTypesOrNil []Type
}

func newFuncType(argsTypes []Type, returnTypes []Type) *FuncType {
	return &FuncType{
		ArgTypesOrNil:    argsTypes,
		ReturnTypesOrNil: returnTypes,
	}
}
func (ft FuncType) String() string {
	start := "funcType"
	argStart := "["
	args := JoinWithSepG(ft.ArgTypesOrNil, ",")
	argEnd := "]"
	arrow := "->"
	returnStart := "["
	returns := JoinWithSepG(ft.ReturnTypesOrNil, ",")
	returnEnd := "]"

	strings := []string{start, argStart, args,
		argEnd, arrow, returnStart, returns, returnEnd}
	return JoinBuilder(strings)
}

// stmt
type Assign struct {
	Ids   []Id
	Exprs []Expr
}

func newAssign(ids []Id, exprs []Expr) *Assign {
	return &Assign{
		Ids:   ids,
		Exprs: exprs,
	}
}

var _ Stmt = (*Assign)(nil)

func (a *Assign) Print(depth int) []string {

	aStart := "Assign("

	lines := []string{}
	lines = append(lines, LineWithDepth(aStart, depth))
	idStart := "["
	ids := JoinWithSepG(a.Ids, ",")
	idEnd := "]"
	lines = append(lines, LineWithDepth(idStart+ids+idEnd, depth+1))
	lines = append(lines, LineWithDepth("=", depth+1))

	exprStart := "["
	exprs := JoinWithSepG(a.Exprs, ",")
	exprEnd := "]"

	lines = append(lines, LineWithDepth(exprStart+exprs+exprEnd, depth+1))

	aEnd := ")"
	lines = append(lines, LineWithDepth(aEnd, depth))
	return lines
}

func (a *Assign) String() string {
	return JoinLines(a.Print(0))
}
func (a *Assign) Stmt() string {
	return a.String()
}

// stmt
type CallStmt struct {
	//Call이 표현이 아닌 "Statement"로 쓰였음을 강조하기 위해서
	// 이렇게 따로 빼 둠
	Call Call
}

func newCallStmt(call Call) *CallStmt {
	return &CallStmt{Call: call}
}

var _ Stmt = (*CallStmt)(nil)

func (c *CallStmt) Print(depth int) []string {
	lines := []string{}
	lines = append(lines, LineWithDepth("CallStmt(", depth))
	calls := c.Call.Print(depth + 1)
	for _, cl := range calls {
		lines = append(lines, cl)
	}
	lines = append(lines, LineWithDepth("and semicolon )", depth))
	return lines
}

func (c *CallStmt) String() string {
	return JoinLines(c.Print(0))
}

func (c *CallStmt) Stmt() string {
	return c.String()
}

// stmt
type ShortDecl struct {
	Ids   []Id
	Exprs []Expr
}

func newShortDecl(ids []Id, exprs []Expr) *ShortDecl {
	return &ShortDecl{
		Ids:   ids,
		Exprs: exprs,
	}
}

var _ Stmt = (*ShortDecl)(nil)

func (s *ShortDecl) Print(depth int) []string {
	sS := "ShortDecl("
	lines := []string{}

	lines = append(lines, LineWithDepth(sS, depth))
	iS := "["
	i := JoinWithSepG(s.Ids, ",")
	iE := "]"
	lines = append(lines, LineWithDepth(iS+i+iE, depth+1))

	declSign := ":="
	lines = append(lines, LineWithDepth(declSign, depth+1))

	for _, exp := range s.Exprs {
		lines = append(lines, exp.Print(depth+1)...)
	}
	sE := ")"
	lines = append(lines, LineWithDepth(sE, depth))
	return lines
}
func (s *ShortDecl) String() string {
	return JoinLines(s.Print(0))
}
func (s *ShortDecl) Stmt() string {
	return s.String()
}

// stmt
type Return struct {
	ExprsOrNil []Expr
}

func newReturn(exprsOrNil []Expr) *Return {
	return &Return{
		ExprsOrNil: exprsOrNil,
	}
}

var _ Stmt = (*Return)(nil)

func (r *Return) Print(depth int) []string {
	lines := []string{}
	lines = append(lines, LineWithDepth("Return(", depth))
	for i, expr := range r.ExprsOrNil {
		lines = append(lines, expr.Print(depth+1)...)
		if i < len(r.ExprsOrNil)-1 {

			lines = append(lines, LineWithDepth(",", depth+1))
		}
	}
	lines = append(lines, LineWithDepth(")", depth))
	return lines
}

func (r *Return) String() string {
	return JoinLines(r.Print(0))
}
func (r *Return) Stmt() string {
	return r.String()
}

// stmt
type If struct {
	ShortDeclOrNil *ShortDecl
	Bexp           Expr
	ThenBlock      Block
	ElseOrNil      *Block
}

func newIf(shorDeclOrNil *ShortDecl, bexp Expr, thenBlock Block, elseOrNil *Block) *If {
	return &If{
		ShortDeclOrNil: shorDeclOrNil,
		Bexp:           bexp,
		ThenBlock:      thenBlock,
		ElseOrNil:      elseOrNil,
	}
}

var _ Stmt = (*If)(nil)

func (i *If) Print(depth int) []string {
	lines := []string{}
	iS := "IF("
	lines = append(lines, LineWithDepth(iS, depth))

	if i.ShortDeclOrNil != nil {
		lines = append(lines, i.ShortDeclOrNil.Print(depth+1)...)
	}

	lines = append(lines, i.Bexp.Print(depth+1)...)

	lines = append(lines, LineWithDepth("then", depth+1))

	lines = append(lines, i.ThenBlock.Print(depth+1)...)

	if i.ElseOrNil != nil {
		lines = append(lines, LineWithDepth("else", depth+1))
		lines = append(lines, i.ElseOrNil.Print(depth+1)...)
	}
	lines = append(lines, LineWithDepth(")", depth))
	return lines
}
func (i *If) String() string {
	return JoinLines(i.Print(0))
}
func (i *If) Stmt() string {
	return i.String()
}

// stmt
type ForBexp struct {
	Bexp  Expr
	Block Block
}

func newForBexp(bexp Expr, block Block) *ForBexp {
	return &ForBexp{
		Bexp:  bexp,
		Block: block,
	}
}

var _ Stmt = (*ForBexp)(nil)

func (f *ForBexp) Print(depth int) []string {
	fs := "For("
	lines := []string{}

	lines = append(lines, LineWithDepth(fs, depth))
	lines = append(lines, LineWithDepth("while", depth+1))
	lines = append(lines, f.Bexp.Print(depth+1)...)
	lines = append(lines, f.Bexp.Print(depth+1)...)
	lines = append(lines, LineWithDepth(")", depth))
	return lines
}
func (f *ForBexp) String() string {
	return JoinLines(f.Print(0))
}
func (f *ForBexp) Stmt() string {
	return f.String()
}

// stmt
type ForWithAssign struct {
	ShortDecl ShortDecl
	Bexp      Expr
	Assign    Assign
	Block     Block
}

func newForWithAssign(shortDecl ShortDecl, bexp Expr, assign Assign, block Block) *ForWithAssign {
	return &ForWithAssign{
		ShortDecl: shortDecl,
		Bexp:      bexp,
		Assign:    assign,
		Block:     block,
	}
}

var _ Stmt = (*ForWithAssign)(nil)

func (f *ForWithAssign) Print(depth int) []string {
	lines := []string{}
	lines = append(lines, LineWithDepth("For(", depth))

	lines = append(lines, f.ShortDecl.Print(depth+1)...)
	lines = append(lines, LineWithDepth("in", depth+1))
	lines = append(lines, f.Bexp.Print(depth+1)...)
	lines = append(lines, LineWithDepth("with", depth+1))
	lines = append(lines, f.Assign.Print(depth+1)...)
	lines = append(lines, f.Block.Print(depth+1)...)
	lines = append(lines, LineWithDepth(")", depth))
	return lines
}
func (f *ForWithAssign) String() string {
	return JoinLines(f.Print(0))
}
func (f *ForWithAssign) Stmt() string {
	return f.String()
}

// stmt
type ForRangeAexp struct {
	Aexp  Expr
	Block Block
}

func newForRangeAexp(aexp Expr, block Block) *ForRangeAexp {
	return &ForRangeAexp{
		Aexp:  aexp,
		Block: block,
	}
}

var _ Stmt = (*ForRangeAexp)(nil)

func (f *ForRangeAexp) Print(depth int) []string {
	lines := []string{}
	lines = append(lines, LineWithDepth("For(", depth))
	lines = append(lines, f.Aexp.Print(depth+1)...)
	lines = append(lines, LineWithDepth("times do", depth+1))
	lines = append(lines, f.Block.Print(depth+1)...)
	lines = append(lines, LineWithDepth(")", depth))
	return lines
}
func (f *ForRangeAexp) String() string {
	return JoinLines(f.Print(0))
}
func (f *ForRangeAexp) Stmt() string {
	return f.String()
}

// Stmt
type Block struct {
	StmtsOrNil []Stmt
}

func newBlock(stmtsOrNil []Stmt) *Block {
	return &Block{
		StmtsOrNil: stmtsOrNil,
	}
}

var _ Stmt = (*Block)(nil)

func (b *Block) Print(depth int) []string {
	lines := []string{}
	lines = append(lines, LineWithDepth("Block(", depth))
	for i, stmt := range b.StmtsOrNil {
		del := fmt.Sprintf("--- stmt %d ", i+1)
		lines = append(lines, LineWithDepth(del, depth+1))
		lines = append(lines, stmt.Print(depth+1)...)
	}
	lines = append(lines, LineWithDepth("--- block end", depth+1))
	lines = append(lines, LineWithDepth(")", depth))
	return lines
}
func (b *Block) String() string {
	return JoinLines(b.Print(0))
}
func (b *Block) Stmt() string {
	return b.String()
}

// Expr
type Unary struct {
	Op     UnaryKind
	Object Expr
}

func newUnary(op UnaryKind, expr Expr) *Unary {
	return &Unary{
		Op:     op,
		Object: expr,
	}
}

var _ Expr = (*Unary)(nil)

func (u *Unary) Print(depth int) []string {
	var op string
	switch u.Op {
	case MinusUnary:
		op = "-"
	case Not:
		op = "!"
	}

	lines := []string{}
	lines = append(lines, LineWithDepth("Unary(", depth))
	opString := fmt.Sprintf("op:%s", op)
	lines = append(lines, LineWithDepth(opString, depth+1))
	lines = append(lines, u.Object.Print(depth+1)...)
	lines = append(lines, LineWithDepth(")", depth))
	return lines
}

func (u *Unary) String() string {
	return JoinLines(u.Print(0))
}
func (u *Unary) Expr() string {
	return u.String()
}

type UnaryKind int

const (
	MinusUnary UnaryKind = UnaryKind(token.MINUS)
	Not        UnaryKind = UnaryKind(token.NOT)
)

// Expr
type Binary struct {
	Op        BinaryKind
	LeftExpr  Expr
	RightExpr Expr
}

func newBinary(op BinaryKind, left, right Expr) *Binary {
	return &Binary{
		Op:        op,
		LeftExpr:  left,
		RightExpr: right,
	}
}

var _ Expr = (*Binary)(nil)

func (b *Binary) Print(depth int) []string {
	var op string
	switch b.Op {
	case Plus:
		op = "+"
	case MinusBinary:
		op = "-"
	case Mul:
		op = "*"
	case Div:
		op = "/"
	case Equal:
		op = "=="
	case NotEqual:
		op = "!="
	case GreaterThan:
		op = "<"
	case GreaterOrEqual:
		op = "<="
	case LessThan:
		op = "<"
	case LessOrEqual:
		op = "<="
	case And:
		op = "&&"
	case Or:
		op = "||"
	}
	lines := []string{}
	lines = append(lines, LineWithDepth("Binary(", depth))
	lines = append(lines, b.LeftExpr.Print(depth+1)...)
	opString := "op:" + op
	lines = append(lines, LineWithDepth(opString, depth+1))
	lines = append(lines, b.RightExpr.Print(depth+1)...)
	lines = append(lines, LineWithDepth(")", depth))
	return lines
}
func (b *Binary) String() string {
	return JoinLines(b.Print(0))
}
func (b *Binary) Expr() string {
	return b.String()
}

// Expr

type Primary struct {
	PrimaryKind PrimaryKind

	ExprOrNil  Expr
	IdOrNil    *Id
	ValueOrNil *ValueForm
}

func newPrimary(primaryKind PrimaryKind, exprOrNil Expr, idOrNil *Id, valueOrNil *ValueForm) *Primary {
	return &Primary{
		PrimaryKind: primaryKind,
		ExprOrNil:   exprOrNil,
		IdOrNil:     idOrNil,
		ValueOrNil:  valueOrNil,
	}
}

var _ Atom = (*Primary)(nil)

func (p *Primary) Print(depth int) []string {
	lines := []string{}
	primaryStart := "Primary("
	lines = append(lines, LineWithDepth(primaryStart, depth))

	switch p.PrimaryKind {
	case ExprPrimary:
		lines = append(lines, p.ExprOrNil.Print(depth+1)...)
		lines = append(lines, LineWithDepth(")", depth))
		return lines
	case IdPrimary:
		idStr := "id " + p.IdOrNil.String()
		lines = append(lines, LineWithDepth(idStr, depth+1))
		lines = append(lines, LineWithDepth(")", depth))
		return lines

	case ValuePrimary:
		lines = append(lines, p.ValueOrNil.Print(depth+1)...)
		lines = append(lines, LineWithDepth(")", depth))
		return lines
	}
	panic("Primary.Print에서 스위치 케이스 미스매치")
}

func (p *Primary) String() string {
	return JoinLines(p.Print(0))
}

func (p *Primary) Expr() string {
	return p.String()
}
func (p *Primary) Atom() string {
	return p.String()
}

type PrimaryKind int

const (
	ExprPrimary PrimaryKind = iota
	IdPrimary
	ValuePrimary
)

type ValueForm struct {
	ValueKind ValueType

	NumberOrNil  *int
	BoolOrNil    *bool
	StrLitOrNil  *string
	ErrOrOkOrNil *string
	FexpOrNil    *Fexp
}

func newValueForm(valueKind ValueType, numberOrNil *int, boolOrNil *bool, strLitOrNil *string, errOrOkOrNil *string, fexoOrNil *Fexp) *ValueForm {
	return &ValueForm{
		ValueKind:    valueKind,
		NumberOrNil:  numberOrNil,
		BoolOrNil:    boolOrNil,
		StrLitOrNil:  strLitOrNil,
		ErrOrOkOrNil: errOrOkOrNil,
		FexpOrNil:    fexoOrNil,
	}
}
func (v *ValueForm) Print(depth int) []string {
	ss := func(s string) []string { return []string{LineWithDepth("valueForm<"+s+">", depth)} }
	switch v.ValueKind {
	case NumberValue:
		return ss("number: " + strconv.Itoa(*v.NumberOrNil))
	case BoolValue:
		if *v.BoolOrNil {
			return ss("bool: true")
		} else {
			return ss("bool: false")
		}
	case StrLitValue:
		return ss("strlit: " + *v.StrLitOrNil)
	case ErrValue:
		var errString string
		if v.ErrOrOkOrNil == nil {
			errString = "<OK>" //* ok값의 error는 nil포인터로 값을 담는다.
		} else {
			errString = *v.ErrOrOkOrNil
		}
		return ss("error: " + errString)
	case FexpValue:
		lines := []string{}
		lines = append(lines, LineWithDepth("valueForm<", depth))
		lines = append(lines, LineWithDepth("funcExpression: ", depth+1))
		lines = append(lines, v.FexpOrNil.Print(depth+1)...)
		lines = append(lines, LineWithDepth(">", depth))
		return lines
	default:
		panic("ValueForm.String() switch missmatch")
	}
}

func (v *ValueForm) String() string {
	return JoinLines(v.Print(0))
}

type ValueType int

const (
	NumberValue ValueType = iota
	BoolValue
	StrLitValue
	ErrValue
	FexpValue
)

type BinaryKind int

const (
	Plus        BinaryKind = BinaryKind(token.PLUS)
	MinusBinary            = BinaryKind(token.MINUS)
	Mul                    = BinaryKind(token.MUL)
	Div                    = BinaryKind(token.DIV)

	Equal    = BinaryKind(token.EQUAL)
	NotEqual = BinaryKind(token.NEQ)

	GreaterThan    = BinaryKind(token.GT)
	GreaterOrEqual = BinaryKind(token.GTE)
	LessThan       = BinaryKind(token.LT)
	LessOrEqual    = BinaryKind(token.LTE)

	And = BinaryKind(token.AND)
	Or  = BinaryKind(token.OR)
)

type Call struct {
	IsBuilinCall     bool
	PrimaryOrNil     *Primary
	BuiltInKindOrNil BuiltInKind
	ArgsList         []Args
}

var _ Atom = (*Call)(nil)

func newCall(isBuiltIn bool, primaryOrNil *Primary, builtInKindOrNil BuiltInKind, argsList []Args) *Call {
	return &Call{
		IsBuilinCall:     isBuiltIn,
		PrimaryOrNil:     primaryOrNil,
		BuiltInKindOrNil: builtInKindOrNil,
		ArgsList:         argsList,
	}
}
func (c *Call) Print(depth int) []string {
	lines := []string{}
	callStart := "call<"
	lines = append(lines, LineWithDepth(callStart, depth))
	builtInInfo := fmt.Sprintf("isBuiltIn: %v", c.IsBuilinCall)
	lines = append(lines, LineWithDepth(builtInInfo, depth+1))
	var ce []string
	toList := func(s string) []string { return []string{s} }
	if c.IsBuilinCall {
		switch c.BuiltInKindOrNil {
		case NewErrorBuild:
			ce = toList(LineWithDepth("id: newError", depth+1))
		case ErrStringBuild:
			ce = toList(LineWithDepth("id: errString", depth+1))
		case ScanBuild:
			ce = toList(LineWithDepth("id: scan", depth+1))
		case PrintBuild:
			ce = toList(LineWithDepth("id: print", depth+1))
		case PanicBuild:
			ce = toList(LineWithDepth("id: panic", depth+1))
		case LenBuild:
			ce = toList(LineWithDepth("id: len", depth+1))
		}
	} else {
		primary := c.PrimaryOrNil
		dummyBuiltIn := fmt.Sprintf("dummyBuiltInValue: %d", c.BuiltInKindOrNil)
		lines = append(lines, LineWithDepth(dummyBuiltIn, depth+1))
		switch primary.PrimaryKind {
		case IdPrimary:
			ce = toList(LineWithDepth("id: "+(*primary.IdOrNil).String(), depth+1))
		case ValuePrimary:
			lines = append(lines, LineWithDepth("valueForm:", depth+1))
			ce = primary.ValueOrNil.Print(depth + 1)
		case ExprPrimary:
			lines = append(lines, LineWithDepth("expr", depth+1))
			ce = primary.ExprOrNil.Print(depth + 1)
		}
	}

	lines = append(lines, ce...)
	for _, arg := range c.ArgsList {
		lines = append(lines, arg.Print(depth+1)...)
	}
	lines = append(lines, LineWithDepth(">", depth))

	return lines
}

func (c *Call) String() string {
	return JoinLines(c.Print(0))
}

func (c *Call) Expr() string {
	return c.String()
}

func (c *Call) Atom() string {
	return c.String()
}

type Args []Expr

func newArgs(exprs []Expr) *Args {
	a := Args(exprs)
	return &a
}

func (a Args) Print(depth int) []string {
	lines := []string{}
	lines = append(lines, LineWithDepth("args<", depth))
	for i, exp := range a {
		lines = append(lines, exp.Print(depth+1)...)
		if i < len(a)-1 {
			lines = append(lines, LineWithDepth(",", depth+1))
		}
	}
	lines = append(lines, LineWithDepth(">", depth))
	return lines
}
func (a Args) String() string {
	return JoinLines(a.Print(0))
}

type CallKind int

const (
	BuiltInCall CallKind = iota
	PrimaryIdCall
	PrimaryFexpCall
	PrimaryExprCall
)

type BuiltInKind int

const (
	NewErrorBuild BuiltInKind = iota
	ErrStringBuild
	ScanBuild
	PrintBuild
	PanicBuild
	LenBuild
)

type Fexp struct {
	ParamsOrNil      []Param
	ReturnTypesOrNil []Type
	Block            Block
}

func newFexp(paramOrNil []Param, returnOrNil []Type, body Block) *Fexp {
	return &Fexp{
		ParamsOrNil:      paramOrNil,
		ReturnTypesOrNil: returnOrNil,
		Block:            body,
	}
}
func (f *Fexp) Print(depth int) []string {
	lines := []string{}
	lines = append(lines, LineWithDepth("Fexp(", depth))
	pS := "["
	params := JoinWithSepG(f.ParamsOrNil, ",")
	pE := "]"

	lines = append(lines, LineWithDepth(pS+params+pE, depth+1))
	arrow := "=>"
	lines = append(lines, LineWithDepth(arrow, depth+1))
	rS := "["
	r := JoinWithSepG(f.ReturnTypesOrNil, ",")
	rE := "]"
	lines = append(lines, LineWithDepth(rS+r+rE, depth+1))

	lines = append(lines, f.Block.Print(depth+1)...)
	lines = append(lines, LineWithDepth(")", depth))
	return lines
}
