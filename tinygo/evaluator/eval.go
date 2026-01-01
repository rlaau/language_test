package evaluator

import "github.com/rlaaudgjs5638/langTest/tinygo/parser"


func (e *Evaluator) EvalDecl(decl parser.Decl) error {
	panic("not implemented")
}
func (e *Evaluator) EvalFuncDecl(funcDecl parser.FuncDecl) error {
	//1. FuncDecl의 경우엔,
	// callFrame의 env에 자신의 value값을 채울 떄에, Value의 closedEnv에 자기자신을 포함시킬 것
	// 전체적으로 리졸버와 동일하게 스코핑-동작하기
	// (ShortDecl, VarDecl같은 익명함수 대입 시엔, 함수의 closedEnv에 자기자신이 들어가지 못함)
	panic("not implemented")
}

func (e *Evaluator) EvalBlock(block parser.Block) error {
	panic("not implemented")
}
