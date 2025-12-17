package tinygo

type Node interface {
	TokenLiteral() string
	String() string
	First() Token
	Follow() []Token
}


type Cmd interface {
	Node
	CmdString() string
}

type Stmt interface {
	Cmd
	StmtString() string
}
type SimpleStmt interface {
	Stmt
	SimpleStmtString() string
}
type For interface {
	Stmt
	ForString() string
}
type Decl interface {
	Cmd
	DeclString() string
}

type Expr interface {
	Node
	ExprString() string
}

// Atom, Batom, Catom, Aatom은 치환 원칙에 의해 배제함
// 치환됨에도 살려두는 경우는 의미론적 차이에 입각할 경우인데,
// 개인적 생각으로 그렇게 의미론적인 차이가 크지 않다고 봤음

type Type interface {
	Node
	TypeString() string
}

type DeclType interface {
	Node
	DeclTypeString() string
}

type ArgTypes interface {
	Node
	ArgTypesString() string
}

type Args interface {
	Node
	ArgsString() string
}

type Params interface {
	Node
	ParamsString() string
}
