package evaluator

import (
	"fmt"

	"github.com/rlaaudgjs5638/langTest/tinygo/parser"
	"github.com/rlaaudgjs5638/langTest/tinygo/resolver"
)

type Evaluator struct {
	packageAST   parser.PackageAST
	resolveTable resolver.ResolveTable

	//Evaluator는 callFrame, envFrame 둘 다 들고 있음
	callStack      CallStack
	globalEnvFrame *EnvFrame

	// 빌트인 함수값들
	builtInSlots []Value
	//디버그 여부
	debug bool
}

// CallStack은 함수 호출마다 새로 생성되는 CallFrame을 스택 형태로 보관한다.
type CallStack struct {
	callFrames []CallFrame
}

// CallFrame은 함수 호출마다 생성되며,
// currentEnv는 함수 호출에서 발생하는 envFrame의 스택 중 top에 있는 envFrame이다.
type CallFrame struct {
	funcIdOrNil *parser.Id
	currentEnv  *EnvFrame
}

func (cf *CallFrame) String() string {
	idOrNil := cf.funcIdOrNil
	if idOrNil != nil {
		return fmt.Sprintf("<funcName: %s, IdId: %d \n currentEnv: %v>", idOrNil.Name, idOrNil.IdId, cf.currentEnv)
	}
	return fmt.Sprintf("anonymous func, IdId: %d\n currentEnv: %v", idOrNil.IdId, cf.currentEnv)
}
func (cs *CallStack) pushCallFrame(cf CallFrame) {
	cs.callFrames = append(cs.callFrames, cf)
}
func (cs *CallStack) popCallFrame() {
	cs.callFrames = cs.callFrames[0 : len(cs.callFrames)-1]
}

// peekMostCurrentEnv는 가장 최신 callFrame의 currentEnv를 리턴한다
func (cs *CallStack) peekMostCurrentEnv() *EnvFrame {
	return cs.callFrames[len(cs.callFrames)-1].currentEnv
}

// setMostCurrentEnv는 가장 최신 callFrame의 currentEnv의 값을 바꾼다.
func (cs *CallStack) setMostCurrentEnv(ef *EnvFrame) {
	cs.callFrames[len(cs.callFrames)-1].currentEnv = ef
}
func NewEvaluator(packageAst parser.PackageAST, hoistInfo *resolver.HoistInfo, initOrder resolver.InitOrder, resolveTable resolver.ResolveTable, builtins map[string]int) (*Evaluator, error) {

	hoistedFuncDecls := []*parser.FuncDecl{}
	hoistedVarTypeByIdId := map[parser.IdId]parser.Type{}
	maxGlobalSlot := -1

	//1.기초적인 검증. 최대 global SLot index 결정
	// 호이스팅 정보, AST, 리졸브 테이블 간 일관성 검증
	if hoistInfo != nil {
		for _, id := range hoistInfo.FuncIds() {
			decl := hoistInfo.GetFuncDeclById(id)
			if decl == nil {
				return nil, fmt.Errorf("missing hoisted func decl")
			}
			hoistedFuncDecls = append(hoistedFuncDecls, decl)
			ref, ok := resolveTable[id]
			if !ok {
				return nil, fmt.Errorf("missing resolve entry for hoisted func")
			}
			if ref.Kind == resolver.RefGlobal {
				if ref.Slot < 0 {
					return nil, fmt.Errorf("negative global slot for func")
				}
				if ref.Slot > maxGlobalSlot {
					maxGlobalSlot = ref.Slot
				}
			}
		}
		for _, idId := range hoistInfo.VarIds() {
			decl := hoistInfo.GetVarDeclById(idId)
			if decl == nil {
				return nil, fmt.Errorf("missing hoisted var decl")
			}
			hoistedVarTypeByIdId[idId] = decl.Type
			ref, ok := resolveTable[idId]
			if !ok {
				return nil, fmt.Errorf("missing resolve entry for hoisted var")
			}
			if ref.Kind == resolver.RefGlobal {
				if ref.Slot < 0 {
					return nil, fmt.Errorf("negative global slot for var")
				}
				if ref.Slot > maxGlobalSlot {
					maxGlobalSlot = ref.Slot
				}
			}
		}
	} else {
		return nil, fmt.Errorf("hoist info doesn't exist")
	}

	globalSlots := []Value{}
	if maxGlobalSlot >= 0 {
		globalSlots = make([]Value, maxGlobalSlot+1)
	}

	// 2. 분석 결과 바탕으로 고정된 크기의 global Slot생성
	globalEnv := &EnvFrame{Slots: globalSlots}

	//3. Evaluator생성
	e := &Evaluator{
		packageAST:   packageAst,
		resolveTable: resolveTable,
		callStack: CallStack{
			callFrames: []CallFrame{
				{
					currentEnv:  globalEnv,
					funcIdOrNil: nil,
				},
			},
		},
		globalEnvFrame: globalEnv,
		builtInSlots:   []Value{},
		debug:          false,
	}
	//4. 빌트인 레지스트리 생성. resolver가 제공한 builtins를 사용
	e.builtInSlots = make([]Value, maxBuiltinSlot(builtins)+1)
	for name, slot := range builtins {
		fn, ok := builtinByName(name)
		if !ok {
			return nil, fmt.Errorf("missing builtin implementation: %s", name)
		}
		if slot < 0 || slot >= len(e.builtInSlots) {
			return nil, fmt.Errorf("builtin slot out of range: %s", name)
		}
		e.builtInSlots[slot] = newBuiltinFuncVal(fn)
	}
	// 5. hoistInfo서 funcDecl 꺼낸 후 클로저로 바인딩
	// resolveTable에 정의된 slot에 맞게 해당 값 채우기
	// 최종적으론 호이스팅된 함수를 환경에 등록
	for _, fn := range hoistedFuncDecls {
		ref := resolveTable[fn.Id.IdId]
		if ref.Kind != resolver.RefGlobal {
			continue
		}
		if ref.Slot < 0 || ref.Slot >= len(globalEnv.Slots) {
			return nil, fmt.Errorf("global slot out of range for func")
		}
		closure := newClosureVal(&fn.Id, fn.ParamsOrNil, fn.ReturnTypesOrNil, fn.Block, globalEnv)

		globalEnv.Slots[ref.Slot] = closure
	}
	// 6. initOrder의 순서대로 varDecl꺼낸 후 평가
	// 6-1 zero init시 go의 zero값으로 채우기. ExprInit시엔 Valuate후 채우기
	// 6-2 resolveTable에 정의된 slot에 맞게 해당 값 채우기
	// 6-3 최종적으론 호이스팅된 벼수를 초기화 순서에 맞게 환경 등록
	for _, step := range initOrder {
		ref, ok := resolveTable[step.VarId]
		if !ok {
			return nil, fmt.Errorf("missing resolve entry for init var")
		}
		if ref.Kind != resolver.RefGlobal {
			continue
		}
		if ref.Slot < 0 || ref.Slot >= len(globalEnv.Slots) {
			return nil, fmt.Errorf("global slot out of range for init var")
		}

		if step.ZeroInit {
			typ, ok := hoistedVarTypeByIdId[step.VarId]
			if !ok {
				return nil, fmt.Errorf("missing var type for zero init")
			}
			globalEnv.Slots[ref.Slot] = ZeroValueForType(typ)
			continue
		}

		if step.ExprOrNil == nil {
			// 논리 오류 케이스임.
			// ZeroInit이 아니라면, ExprOrNil은 nil값이 아니여야 함.
			return nil, fmt.Errorf("missing init expr for var")
		}
		values, ctrlSigOrNil, err := e.Valuate(step.ExprOrNil)
		if err != nil {
			return nil, fmt.Errorf("init expr evaluation failed")
		}
		if ctrlSigOrNil != nil {
			if ctrlSigOrNil.Kind == CtrlPanic {
				return nil, fmt.Errorf("panic during hoisting: %s", ctrlSigOrNil.Values[0].Inspect())
			}
		}
		if len(values) != 1 {
			// 이 부분은 리졸버에 근거함.
			// 리졸버가 글로벌 레벨에서 호이스팅되는 다중 선언, 다중 할당은
			// 단일 선언, 단일 할당으로 분해 후에 initOrder로 제공하고 있기 때문임.
			return nil, fmt.Errorf("init expr must return exactly one value")
		}
		globalEnv.Slots[ref.Slot] = values[0]
	}
	// 7. 최종적으로
	// 8. 호이스팅에 맞춰 환경을 초기화한 Evaluator를 리턴
	return e, nil
}

// EnvFrame은 현재 스코프에서의 환경을 나타냄
type EnvFrame struct {
	Slots          []Value // SLot에 따른 Value
	ParentEnvFrame *EnvFrame
}

func (e *Evaluator) EvalMainFunc() error {
	// 1. main함수의 시그니처 검사하기
	// 2. env에서 main함수 꺼내서 호출
	// 3 . 평가 시엔 env의 스코핑&슬롯 규칙이 리졸버와 동일할 것. 그리고 ref가 빌트인일 시 빌트인 참조할 것.

	var mainDecl *parser.FuncDecl
	for _, decl := range e.packageAST.DeclsOrNil {
		fn, ok := decl.(*parser.FuncDecl)
		if !ok {
			continue
		}
		if fn.Id.Name == "main" {
			mainDecl = fn
			break
		}
	}
	if mainDecl == nil {
		return fmt.Errorf("missing main function")
	}
	if len(mainDecl.ParamsOrNil) != 0 || len(mainDecl.ReturnTypesOrNil) != 0 {
		return fmt.Errorf("main must have signature func()")
	}

	mainVal, err := e.valueForId(&mainDecl.Id)
	if err != nil {
		return err
	}
	mainClosure, ok := mainVal.(*ClosureValue)
	if !ok {
		return fmt.Errorf("main is not a function")
	}
	if len(mainClosure.Params) != 0 || len(mainClosure.ReturnTypes) != 0 {
		return fmt.Errorf("main must have signature func()")
	}
	_, ctrlSig, err := e.callClosure(mainClosure, []Value{})
	if err != nil {
		return err
	}
	if ctrlSig != nil {
		if ctrlSig.Kind == CtrlPanic {
			if len(ctrlSig.Values) > 0 {
				return fmt.Errorf("panic: %s", ctrlSig.Values[0].Inspect())
			}
			return fmt.Errorf("panic")
		}
		return fmt.Errorf("unexpected control signal: %v", ctrlSig.Kind)
	}
	if e.debug == true {
		for i, e := range e.callStack.callFrames {
			fmt.Printf("%dth CallFame:", i)
			fmt.Println(e.String())
		}
	}
	return nil
}

func (e *Evaluator) pushEnvFrame(ef *EnvFrame) {
	ef.ParentEnvFrame = e.CurrentEnv()
	e.callStack.setMostCurrentEnv(ef)
}
func (e *Evaluator) popEnvFrame() {
	e.callStack.setMostCurrentEnv(e.CurrentEnv().ParentEnvFrame)
}

func (e *Evaluator) CurrentEnv() *EnvFrame {
	return e.callStack.peekMostCurrentEnv()
}

func maxBuiltinSlot(slots map[string]int) int {
	max := -1
	for _, slot := range slots {
		if slot > max {
			max = slot
		}
	}
	return max
}
