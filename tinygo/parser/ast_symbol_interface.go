package parser

type Node interface {
	String() string
	Print(depth int) []string
	//TODO First() []token.TokenKind
}

type Decl interface {
	Node
	Decl() string // == String()보다 더 구조화
}
type Stmt interface {
	Node
	Stmt() string
}

type Expr interface {
	Node
	Expr() string
}

type Lexp interface {
	Expr
	Lexp() string //ring으로 감싸기
}
