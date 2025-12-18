package tinygo

import "github.com/rlaaudgjs5638/langTest/tinygo/token"

type VarDecl struct {
	Id        token.Token
	DeclType  DeclType
	ExprOrNil Expr
}

type ShortDecl struct {
	Id   token.Token
	Expr Expr
}
