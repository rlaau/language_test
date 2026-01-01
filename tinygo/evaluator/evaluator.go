package evaluator

import (
	"github.com/rlaaudgjs5638/langTest/tinygo/parser"
	"github.com/rlaaudgjs5638/langTest/tinygo/resolver"
)

type Evaluator struct {
	resolveTable resolver.ResolveTable

	//Evaluator는 callFrame, envFrame 둘 다 들고 있음
	currentEnvFrame  *EnvFrame
	currentCallFrame *CallFrame

	// 빌트인 함수값들
	builtIns []Value
}

type CallFrame struct {
	//* q2. Evaluator가 "값"으로 취급하는 것들을 Value로 정의해두려 하는데
	//* 이 ValueType을 parser.ValueForm을 이용해서 정의해도 되나?
	funcId          parser.Id
	Args            []Value
	ReturnValues    []Value
	LocalEnv        *EnvFrame
	ParentCallFrame *CallFrame
}
type EnvFrame struct {
	SlotAndValueList []Value
	ParentEnvFrame   *EnvFrame
}

type Value struct {
	// int, bool, string...의 경우엔 valueForm의 값 사용
	valueForm parser.ValueForm
	// 값이 fexp일 시엔 closedEnv까지 같이 활용
	closedEnv *EnvFrame
}

// type EvalError
func EvalPackage(packageAst parser.PackageAST, hoistInfo *resolver.HoistInfo, initOrder resolver.InitOrder, resolveTable resolver.ResolveTable) error {
	// 1. hoistInfo서 funcDecl 꺼낸 후 평가 => resolveTable에 정의된 slot에 맞게 해당 값 채우기
	// 2. initOrder의 순서대로 varDecl꺼낸 후 평가 => resolveTable에 정의된 slot에 맞게 해당 값 채우기
	// 3. 1,2번 과정을 통해 envFrame완성함
	//4. ast에서 funcDecl, 그 중에서 이름이 "main"인 funcDecl에 진입.
	// 5. main을 위한 CallFrame생성. Args, ReturnValues다 []value{}로 비우고, parentCallFrame=nil
	// 6. resolveTable, 만들어진 envFrame, 만들어진 callFrame바탕으로 newEvaluator생성
	// 7. evaluator에 빌트인 함수 등록. 빌트인도 리졸버와 순서 맞추기. (빌트인 함수는 리졸버의 빌트인 로드와, 내 언어 스펙 참고)
	// 7. 해당 상태로 main함수의 "block"평가 시작
	// 8 . 평가 시엔 env의 스코핑&슬롯 규칙이 리졸버와 동일할 것!
	// 9 . block성공적으로 평가되면 return nil
	return nil
}
func (e *Evaluator) EvalFuncDecl(decl parser.FuncDecl, callFrame *CallFrame) error {
	//1. FuncDecl의 경우엔,
	// callFrame의 env에 자신의 value값을 채울 떄에, Value의 closedEnv에 자기자신을 포함시킬 것
	// 전체적으로 리졸버와 동일하게 스코핑-동작하기
	// (ShortDecl, VarDecl같은 익명함수 대입 시엔, 함수의 closedEnv에 자기자신이 들어가지 못함)
	panic("not implemented")
}
