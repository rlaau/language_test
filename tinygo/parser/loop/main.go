package main

import (
	"fmt"

	"github.com/rlaaudgjs5638/langTest/tinygo/lexer"
	"github.com/rlaaudgjs5638/langTest/tinygo/parser"
)

func main() {
	lexer := lexer.NewLexer()
	lexer.Set(`(newError("ss")||print(x))&&k(t,p,q)(z,u) -2`)
	parser := parser.NewParser(lexer)
	parsed, err := parser.ParseExpr()
	if err != nil {
		fmt.Printf(err.Error())
	}
	fmt.Printf(parsed.String())

}
