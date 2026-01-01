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
	builtInSlots []Value
	//디버그 여부
	debug bool
}

type CallFrame struct {
	funcIdOrNil     *parser.Id
	Args            []Value
	ReturnValues    []Value
	ParentCallFrame *CallFrame
}

func (cf *CallFrame) String() string {
	if cf.funcIdOrNil == nil {
		return "callFrame" + "anonymous func"
	}
	return "callFrame:" + cf.funcIdOrNil.String()
}

type EnvFrame struct {
	Slots          []Value // SLot에 따른 Value
	ParentEnvFrame *EnvFrame
}

func EvalPackage(packageAst parser.PackageAST, hoistInfo *resolver.HoistInfo, initOrder resolver.InitOrder, resolveTable resolver.ResolveTable) error {

	// 1. resolveTable 및 globalEnv 바탕 고정 크기 env생성=> evaluator생성
	// 2. Evaluator의 Builtins 로드. 리졸버의 Builtins참고.리졸버와 순서 맞추기. (빌트인 함수는 리졸버의 빌트인 로드와, 내 언어 스펙 참고)
	// 3. hoistInfo서 funcDecl 꺼낸 후 클로저로 바인딩 => resolveTable에 정의된 slot에 맞게 해당 값 채우기
	// 4. initOrder의 순서대로 varDecl꺼낸 후 평가
	// 4-1 zero init시 go의 zero값으로 채우기. ExprInit시엔 Valuate후 채우기
	// 4-2 resolveTable에 정의된 slot에 맞게 해당 값 채우기
	// 5. main함수의 시그니처 검사하기
	// 6. env에서 main함수 꺼내서 호출
	// 7 . 평가 시엔 env의 스코핑&슬롯 규칙이 리졸버와 동일할 것. 그리고 ref가 빌트인일 시 빌트인 참조할 것.
	// 8 . block성공적으로 평가되면 return nil
	panic("not implemented")
}

func (e *Evaluator) pushCallFrame(f *CallFrame) {
	f.ParentCallFrame = e.currentCallFrame
	e.currentCallFrame = f
}
func (e *Evaluator) popCallFrame() { e.currentCallFrame = e.currentCallFrame.ParentCallFrame }

func (e *Evaluator) pushEnvFrame(ef *EnvFrame) {
	ef.ParentEnvFrame = e.currentEnvFrame
	e.currentEnvFrame = ef
}
func (e *Evaluator) popEnvFrame() { e.currentEnvFrame = e.currentEnvFrame.ParentEnvFrame }
