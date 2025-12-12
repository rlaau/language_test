package calclang

import (
	"bufio"
	"fmt"
	"os"
)

type Calculator struct {
	lexer  *Lexer
	parser *Parser
}

func NewCalculator() *Calculator {
	return &Calculator{
		lexer:  NewLexer(),
		parser: NewParser(),
	}
}

func (c *Calculator) Run() {
	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("입력: ")
		//엔터키 기준 블로킹
		if reader.Scan() {
			line := reader.Text()
			fmt.Println("입력받음:", line)
			c.PrintEval(line)
		}
		if err := reader.Err(); err != nil {
			fmt.Println("입력 에러:", err)
		}
	}
}

func (c *Calculator) PrintEval(s string) {
	lexed, err := c.lexer.Lex(s)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	c.parser.Parse(lexed)
	return
}
