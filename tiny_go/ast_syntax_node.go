package tinygo

type VarDecl struct {
	Id        Token
	DeclType  DeclType
	ExprOrNil Expr
}

type ShortDecl struct {
	Id   Token
	Expr Expr
}
