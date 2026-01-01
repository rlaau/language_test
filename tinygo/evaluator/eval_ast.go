package evaluator

import "github.com/rlaaudgjs5638/langTest/tinygo/parser"

// Eval은 기본적으로 사이드 이펙트만 관리함
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
func (e *Evaluator) EvalCallStmt(callStmt parser.CallStmt) (*ControlSignal, error) {
	//기본적으로 call은
	//call내의 args를 모두 "Valuate"한 다음에
	// []Value형태로 함수에 전달하기
	//함수는 빌트인 함수를 구별해서 빌트인 레지스트리 or 환경에서 찾기
	panic("not implemented")
}
func (e *Evaluator) EvalBlock(block parser.Block) (*ControlSignal, error) {
	// 아무 제어도 없다면 *ControlSyntax==nil
	panic("not implemented")
}
func (e *Evaluator) EvalIf(ifNode parser.If) (*ControlSignal, error) {
	panic("not implemented")
}
