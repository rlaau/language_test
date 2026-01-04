package evaluator

import (
	"fmt"

	"github.com/rlaaudgjs5638/langTest/tinygo/parser"
	"github.com/rlaaudgjs5638/langTest/tinygo/resolver"
)

func (e *Evaluator) Valuate(expr parser.Expr) ([]Value, *ControlSignal, error) {
	//Valueate는 다음과 같은 암묵적 입력을 지님
	//builtIns []Value, resolveTable *resolver.ResolveTable, envFrame *EnvFrame, callFrame *CallFrame
	//이런 인자 폭발 우려로 인해, Valuate는 evaluator의 메서드로 취급하여 암묵적 인자를 받기로 함.
	switch node := expr.(type) {
	case *parser.Binary:
		return e.ValuateBinary(node)
	case *parser.Unary:
		return e.ValuateUnary(node)
	case *parser.Primary:
		return e.ValuatePrimary(node)
	case *parser.Call:
		return e.ValuateCall(node)
	default:
		return nil, nil, fmt.Errorf("unknown expr node: %T", expr)
	}
}

func (e *Evaluator) ValuateUnary(u *parser.Unary) ([]Value, *ControlSignal, error) {
	values, ctrlSigOrNil, err := e.Valuate(u.Object)
	if err != nil || ctrlSigOrNil != nil {
		return nil, ctrlSigOrNil, err
	}
	v, err := expectSingle(values, "unary")
	if err != nil {
		return nil, nil, err
	}
	switch u.Op {
	//"-"연산은 int만 허용
	case parser.MinusUnary:
		intVal, ok := v.(*IntValue)
		if !ok {
			return nil, nil, fmt.Errorf("unary - expects int")
		}
		return []Value{newIntVal(-intVal.Value)}, nil, nil
	//"!"연산은 bool만 허용
	case parser.Not:
		boolVal, ok := v.(*BoolValue)
		if !ok {
			return nil, nil, fmt.Errorf("unary ! expects bool")
		}
		return []Value{newBoolVal(!boolVal.Value)}, nil, nil
	default:
		return nil, nil, fmt.Errorf("unknown unary op: %v", u.Op)
	}
}

func (e *Evaluator) ValuateBinary(b *parser.Binary) ([]Value, *ControlSignal, error) {
	switch b.Op {
	case parser.And, parser.Or:
		// And, Or최적화에 기반한 로직임
		// &&이면서 left가 거짓-> 무조건 거짓. || 이면서 left가 참->무조건 참
		// 그 외의 경우는, &&이면서 left가 참, ||이면서 left가 거짓인데, 이 경우의 진리값은 right와 동일
		leftValues, ctrlSigOrNil, err := e.Valuate(b.LeftExpr)
		if err != nil || ctrlSigOrNil != nil {
			return nil, ctrlSigOrNil, err
		}
		leftVal, err := expectSingle(leftValues, "binary")
		if err != nil {
			return nil, nil, err
		}
		leftBool, ok := leftVal.(*BoolValue)
		if !ok {
			return nil, nil, fmt.Errorf("logical op expects bool")
		}
		if b.Op == parser.And && !leftBool.Value {
			return []Value{newBoolVal(false)}, nil, nil
		}
		if b.Op == parser.Or && leftBool.Value {
			return []Value{newBoolVal(true)}, nil, nil
		}
		rightValues, ctrlSigOrNil, err := e.Valuate(b.RightExpr)
		if err != nil || ctrlSigOrNil != nil {
			return nil, ctrlSigOrNil, err
		}
		rightVal, err := expectSingle(rightValues, "binary")
		if err != nil {
			return nil, nil, err
		}
		rightBool, ok := rightVal.(*BoolValue)
		if !ok {
			return nil, nil, fmt.Errorf("logical op expects bool")
		}
		return []Value{newBoolVal(rightBool.Value)}, nil, nil
	}

	leftValues, ctrlSigOrNil, err := e.Valuate(b.LeftExpr)
	if err != nil || ctrlSigOrNil != nil {
		return nil, ctrlSigOrNil, err
	}
	leftVal, err := expectSingle(leftValues, "binary")
	if err != nil {
		return nil, nil, err
	}
	rightValues, ctrlSigOrNil, err := e.Valuate(b.RightExpr)
	if err != nil || ctrlSigOrNil != nil {
		return nil, ctrlSigOrNil, err
	}
	rightVal, err := expectSingle(rightValues, "binary")
	if err != nil {
		return nil, nil, err
	}

	switch b.Op {
	case parser.Plus, parser.MinusBinary, parser.Mul, parser.Div:
		leftInt, lok := leftVal.(*IntValue)
		rightInt, rok := rightVal.(*IntValue)
		// plus에 한해선
		// string+string연산을 지원함
		// 둘 다  int인 경우가 아닐 경우, string집합 위의 +연산인지 검증함
		if !lok || !rok {

			if b.Op == parser.Plus {
				leftStr, lsok := leftVal.(*StringValue)
				rightStr, rsok := rightVal.(*StringValue)
				if lsok && rsok {
					return []Value{newStringVal(leftStr.Value + rightStr.Value)}, nil, nil
				}
			}
			return nil, nil, fmt.Errorf("arithmetic op expects int")
		}
		switch b.Op {
		case parser.Plus:
			return []Value{newIntVal(leftInt.Value + rightInt.Value)}, nil, nil
		case parser.MinusBinary:
			return []Value{newIntVal(leftInt.Value - rightInt.Value)}, nil, nil
		case parser.Mul:
			return []Value{newIntVal(leftInt.Value * rightInt.Value)}, nil, nil
		case parser.Div:
			if rightInt.Value == 0 {
				return nil, nil, fmt.Errorf("division by zero")
			}
			return []Value{newIntVal(leftInt.Value / rightInt.Value)}, nil, nil
		}
	case parser.Equal, parser.NotEqual:
		eq, ok := equalValues(leftVal, rightVal)
		if !ok {
			return nil, nil, fmt.Errorf("equality op expects same comparable types")
		}
		if b.Op == parser.NotEqual {
			eq = !eq
		}
		return []Value{newBoolVal(eq)}, nil, nil
	case parser.GreaterThan, parser.GreaterOrEqual, parser.LessThan, parser.LessOrEqual:
		leftInt, lok := leftVal.(*IntValue)
		rightInt, rok := rightVal.(*IntValue)
		if !lok || !rok {
			return nil, nil, fmt.Errorf("comparison op expects int")
		}
		switch b.Op {
		case parser.GreaterThan:
			return []Value{newBoolVal(leftInt.Value > rightInt.Value)}, nil, nil
		case parser.GreaterOrEqual:
			return []Value{newBoolVal(leftInt.Value >= rightInt.Value)}, nil, nil
		case parser.LessThan:
			return []Value{newBoolVal(leftInt.Value < rightInt.Value)}, nil, nil
		case parser.LessOrEqual:
			return []Value{newBoolVal(leftInt.Value <= rightInt.Value)}, nil, nil
		}
	}
	return nil, nil, fmt.Errorf("unknown binary op: %v", b.Op)
}

func (e *Evaluator) ValuatePrimary(p *parser.Primary) ([]Value, *ControlSignal, error) {
	switch p.PrimaryKind {
	case parser.ExprPrimary:
		return e.Valuate(p.ExprOrNil)
	case parser.IdPrimary:
		val, err := e.valueForId(p.IdOrNil)
		if err != nil {
			return nil, nil, err
		}
		return []Value{val}, nil, nil
	case parser.ValuePrimary:
		val, err := e.ValuateValueForm(p.ValueOrNil)
		if err != nil {
			return nil, nil, err
		}
		return []Value{val}, nil, nil
	default:
		return nil, nil, fmt.Errorf("unknown primary kind: %v", p.PrimaryKind)
	}
}

func (e *Evaluator) ValuateValueForm(v *parser.ValueForm) (Value, error) {
	switch v.ValueKind {
	case parser.NumberValue:
		if v.NumberOrNil == nil {
			return nil, fmt.Errorf("number literal missing value")
		}
		return newIntVal(int64(*v.NumberOrNil)), nil
	case parser.BoolValue:
		if v.BoolOrNil == nil {
			return nil, fmt.Errorf("bool literal missing value")
		}
		return newBoolVal(*v.BoolOrNil), nil
	case parser.StrLitValue:
		if v.StrLitOrNil == nil {
			return nil, fmt.Errorf("string literal missing value")
		}
		return newStringVal(*v.StrLitOrNil), nil
	case parser.ErrValue:
		return newErrorVal(v.ErrOrNilIfOk), nil
	case parser.FexpValue:
		if v.FexpOrNil == nil {
			return nil, fmt.Errorf("func literal missing body")
		}
		fexp := v.FexpOrNil
		return newClosureVal(nil, fexp.ParamsOrNil, fexp.ReturnTypesOrNil, fexp.Block, e.CurrentEnv()), nil
	default:
		return nil, fmt.Errorf("unknown value kind: %v", v.ValueKind)
	}
}

func (e *Evaluator) ValuateCall(c *parser.Call) ([]Value, *ControlSignal, error) {

	//가장 처음 평가된 "표현"은 primary임.
	// 계속해서 평가를 리듀스 해 갈 예정
	// f()()-> g()->h
	appliedExpr, ctrlSigOrNil, err := e.Valuate(&c.PrimaryOrNil)
	if err != nil || ctrlSigOrNil != nil {
		return nil, ctrlSigOrNil, err
	}

	for _, argTuple := range c.ArgsList {
		if len(appliedExpr) != 1 {
			return nil, nil, fmt.Errorf("invalid call: the callee must evaluate to a single function")

		}
		callee := appliedExpr[0]
		// args: 한 번 호출에 필요한 arg 튜플
		args := make([]Value, 0, len(argTuple))
		for _, expr := range argTuple {
			// 인자는 현재 환경에서 평가
			values, ctrlSigOrNil, err := e.Valuate(expr)
			if err != nil || ctrlSigOrNil != nil {
				return nil, ctrlSigOrNil, err
			}
			argVal, err := expectSingle(values, "call arg")
			if err != nil {
				return nil, nil, err
			}
			args = append(args, argVal)
		}

		switch fn := callee.(type) {
		case *BuiltinFuncValue:
			values, ctrlSig, err := fn.Func.Impl(e, args)
			if err != nil || ctrlSig != nil {
				return nil, ctrlSig, err
			}
			appliedExpr = values
		case *ClosureValue:
			values, ctrlSig, err := e.callClosure(fn, args)
			if err != nil || ctrlSig != nil {
				return nil, ctrlSig, err
			}
			appliedExpr = values
		default:
			return nil, nil, fmt.Errorf("call target is not callable")
		}
	}

	return appliedExpr, nil, nil
}

func (e *Evaluator) callClosure(c *ClosureValue, args []Value) ([]Value, *ControlSignal, error) {
	if len(args) != len(c.Params) {
		return nil, nil, fmt.Errorf("arg count mismatch")
	}
	// 함수 호출 시엔, 기존의 EnvList에서 pop, push하지 않고,
	// 대신 새 콜스텍의 원소를 추가 후 그 위에서 pop,push를 함
	newCallFrame := &EnvFrame{Slots: make([]Value, e.maxSlotFromParams(c.Params)+1), ParentEnvFrame: c.ParentEnv}
	e.callStack.pushCallFrame(newCallFrame)
	defer e.callStack.popCallFrame()

	for i, param := range c.Params {
		ref, ok := e.resolveTable[param.Id.IdId]
		if !ok {
			return nil, nil, fmt.Errorf("missing resolve entry for param")
		}
		if ref.Kind != resolver.RefLocal {
			return nil, nil, fmt.Errorf("param resolved as non-local")
		}
		if ref.Slot < 0 {
			return nil, nil, fmt.Errorf("negative slot for param")
		}
		if ref.Slot >= len(newCallFrame.Slots) {
			newCallFrame.Slots = growSlots(newCallFrame.Slots, ref.Slot+1)
		}
		newCallFrame.Slots[ref.Slot] = args[i]
	}
	//리졸버와 스코핑 로직 일치:
	//param은 block과 같은 환경에서 존재
	ctrlSig, err := e.evalBlock(c.Block, true)
	if err != nil {
		return nil, nil, err
	}
	if ctrlSig == nil {
		return []Value{}, nil, nil
	}
	switch ctrlSig.Kind {
	case CtrlReturn:
		return ctrlSig.Values, nil, nil
	case CtrlPanic:
		return nil, ctrlSig, nil
	default:
		//return, panic외의 제어신호는 함수 바깥으로 전파되지 못함
		return nil, nil, fmt.Errorf("only \"return\" or \"panic\" control signals may propagate out of a function")

	}
}

// valueForId는 리졸브 테이블 기반으로 id를 평가함
// Env와 resolve scope간의 동치를 상정하고 작동함.
func (e *Evaluator) valueForId(id *parser.Id) (Value, error) {
	if id == nil {
		return nil, fmt.Errorf("nil id")
	}
	ref, ok := e.resolveTable[id.IdId]
	if !ok {
		return nil, fmt.Errorf("missing resolve entry for id: %s", id.String())
	}
	switch ref.Kind {
	case resolver.RefBuiltin:
		if ref.Slot < 0 || ref.Slot >= len(e.builtInSlots) {
			return nil, fmt.Errorf("builtin slot out of range")
		}
		// 리졸버가 리졸빙 과정에서
		// r.builtins-slot테이블에 맞춰서 레퍼런스에 알맞은 Slot을 주입했음을 가정함
		return e.builtInSlots[ref.Slot], nil
	case resolver.RefGlobal:
		env := e.globalEnvFrame
		if ref.Slot < 0 || ref.Slot >= len(env.Slots) {
			return nil, fmt.Errorf("env slot out of range")
		}
		return env.Slots[ref.Slot], nil
	case resolver.RefLocal:
		// 클로저가 제대로 동작한단 가정 하에 올바르게 동작함.
		env, err := envAtDistance(e.CurrentEnv(), ref.Distance)
		if err != nil {
			return nil, err
		}
		if ref.Slot < 0 || ref.Slot >= len(env.Slots) {
			return nil, fmt.Errorf("env slot out of range")
		}
		return env.Slots[ref.Slot], nil
	default:
		return nil, fmt.Errorf("unknown ref kind")
	}
}

func envAtDistance(start *EnvFrame, distance int) (*EnvFrame, error) {
	env := start
	for range distance {
		if env == nil || env.ParentEnvFrame == nil {
			return nil, fmt.Errorf("env frame distance out of range")
		}
		env = env.ParentEnvFrame
	}
	if env == nil {
		return nil, fmt.Errorf("nil env frame")
	}
	return env, nil
}

func expectSingle(values []Value, context string) (Value, error) {
	if len(values) != 1 {
		return nil, fmt.Errorf("%s expects single value", context)
	}
	return values[0], nil
}

// equalValues return equal, ok
func equalValues(left, right Value) (bool, bool) {
	switch lv := left.(type) {
	case *IntValue:
		rv, ok := right.(*IntValue)
		if !ok {
			return false, false
		}
		return lv.Value == rv.Value, true
	case *BoolValue:
		rv, ok := right.(*BoolValue)
		if !ok {
			return false, false
		}
		return lv.Value == rv.Value, true
	case *StringValue:
		rv, ok := right.(*StringValue)
		if !ok {
			return false, false
		}
		return lv.Value == rv.Value, true
	case *ErrorValue:
		rv, ok := right.(*ErrorValue)
		if !ok {
			return false, false
		}
		if lv.IsOk && rv.IsOk {
			return true, true
		}
		if lv.IsOk != rv.IsOk {
			return false, true
		}
		return lv.ErrMsg == rv.ErrMsg, true
	default:
		// 함수 값 간의 동등성 비교는 허용하지 않음
		return false, false
	}
}

func (e *Evaluator) maxSlotFromParams(params []parser.Param) int {
	max := -1
	for _, param := range params {
		ref, ok := e.resolveTable[param.Id.IdId]
		if !ok {
			continue
		}
		if ref.Slot > max {
			max = ref.Slot
		}
	}
	if max < 0 {
		return 0
	}
	return max
}

func growSlots(slots []Value, size int) []Value {
	if size <= len(slots) {
		return slots
	}
	newSlots := make([]Value, size)
	copy(newSlots, slots)
	return newSlots
}
