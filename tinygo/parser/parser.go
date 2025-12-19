package parser

import (
	"github.com/rlaaudgjs5638/langTest/tinygo/token"

	"github.com/rlaaudgjs5638/langTest/tinygo/lexer"
)

type Parser struct {
	lexer      *lexer.Lexer
	errorStack Stack[string]
	tokenStack Stack[token.Token]
}
