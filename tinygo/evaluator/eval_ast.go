package evaluator

import (
	"fmt"

	"github.com/rlaaudgjs5638/langTest/tinygo/parser"
	"github.com/rlaaudgjs5638/langTest/tinygo/resolver"
)

func (e *Evaluator) evalFuncDecl(funcDecl parser.FuncDecl) error {
	//1. FuncDecl의 경우엔,
	// callFrame의 env에 자신의 value값을 채울 떄에, Value의 closedEnv에 자기자신을 포함시킬 것
	// 전체적으로 리졸버와 동일하게 스코핑-동작하기
	// (ShortDecl, VarDecl같은 익명함수 대입 시엔, 함수의 closedEnv에 자기자신이 들어가지 못함)
	closure := newClosureVal(&funcDecl.Id, funcDecl.ParamsOrNil, funcDecl.ReturnTypesOrNil, funcDecl.Block, e.CurrentEnv())
	return e.setValueForId(funcDecl.Id, closure)
}
func (e *Evaluator) evalCallStmt(callStmt parser.CallStmt) (*ControlSignal, error) {
	// evalCallStmt는 직접적으로 ctrlSignal을 다루지 않음
	// ctrlSig가 return일때의 처리는 callValue를 Valueate할 때 일어남
	// 결국 CallStmt는 call된 "표현식의 값"을 받는 것임.
	// 함수 호출 표현식의 값을 만드는 과정에서 일어나는 return제어문에 대한 처리는
	// Valuate가 책임지고 맡아서 핸들링 후, eval에게는 "표현식 값"만 주는 것
	_, ctrlSigOrNil, err := e.Valuate(&callStmt.Call)
	if err != nil || ctrlSigOrNil != nil {
		return ctrlSigOrNil, err
	}
	return nil, nil
}
func (e *Evaluator) evalBlock(block parser.Block, reuseCurrentEnv bool) (*ControlSignal, error) {
	// 아무 제어도 없다면 *ControlSyntax==nil
	//reuseCurrent시엔 env pop,push없이 현재의 env재사용
	if !reuseCurrentEnv {
		e.pushEnvFrame(&EnvFrame{Slots: []Value{}})
		defer e.popEnvFrame()
	}
	for _, stmt := range block.StmtsOrNil {
		var ctrlSig *ControlSignal
		var err error
		switch node := stmt.(type) {
		case *parser.Assign:
			ctrlSig, err = e.evalAssign(node)
		case *parser.CallStmt:
			ctrlSig, err = e.evalCallStmt(*node)
		case *parser.ShortDecl:
			ctrlSig, err = e.evalShortDecl(node)
		case *parser.VarDecl:
			ctrlSig, err = e.evalVarDecl(node)
		case *parser.FuncDecl:
			err = e.evalFuncDecl(*node)
		case *parser.Return:
			ctrlSig, err = e.evalReturn(node)
		case *parser.Break:
			ctrlSig = newControlSignal(CtrlBreak, nil)
		case *parser.Continue:
			ctrlSig = newControlSignal(CtrlContinue, nil)
		case *parser.If:
			ctrlSig, err = e.EvalIf(*node)
		case *parser.ForBexp:
			ctrlSig, err = e.EvalForBexp(*node)
		case *parser.ForWithAssign:
			ctrlSig, err = e.EvalForWithAssign(*node)
		case *parser.Block:
			ctrlSig, err = e.evalBlock(*node, false)
		default:
			err = fmt.Errorf("unknown stmt node: %T", stmt)
		}
		if err != nil || ctrlSig != nil {
			return ctrlSig, err
		}
	}
	return nil, nil
}
func (e *Evaluator) EvalIf(ifNode parser.If) (*ControlSignal, error) {
	e.pushEnvFrame(&EnvFrame{Slots: []Value{}})
	defer e.popEnvFrame()

	if ifNode.ShortDeclOrNil != nil {
		ctrlSig, err := e.evalShortDecl(ifNode.ShortDeclOrNil)
		if err != nil || ctrlSig != nil {
			return ctrlSig, err
		}
	}
	cond, ctrlSig, err := e.evalBoolExpr(ifNode.Bexp)
	if err != nil || ctrlSig != nil {
		return ctrlSig, err
	}
	if cond {
		return e.evalBlock(ifNode.ThenBlock, false)
	}
	if ifNode.ElseOrNil != nil {
		return e.evalBlock(*ifNode.ElseOrNil, false)
	}
	return nil, nil
}

func (e *Evaluator) EvalForBexp(forNode parser.ForBexp) (*ControlSignal, error) {
	e.pushEnvFrame(&EnvFrame{Slots: []Value{}})
	defer e.popEnvFrame()

	for {
		cond, ctrlSig, err := e.evalBoolExpr(forNode.Bexp)
		if err != nil || ctrlSig != nil {
			return ctrlSig, err
		}
		if !cond {
			return nil, nil
		}
		ctrlSig, err = e.evalBlock(forNode.Block, false)
		if err != nil {
			return nil, err
		}
		if ctrlSig != nil {
			switch ctrlSig.Kind {
			case CtrlBreak:
				return nil, nil
			case CtrlContinue:
				continue
			default:
				return ctrlSig, nil
			}
		}
	}
}

func (e *Evaluator) EvalForWithAssign(node parser.ForWithAssign) (*ControlSignal, error) {
	e.pushEnvFrame(&EnvFrame{Slots: []Value{}})
	defer e.popEnvFrame()

	ctrlSig, err := e.evalShortDecl(&node.ShortDecl)
	if err != nil || ctrlSig != nil {
		return ctrlSig, err
	}

	for {
		cond, ctrlSig, err := e.evalBoolExpr(node.Bexp)
		if err != nil || ctrlSig != nil {
			return ctrlSig, err
		}
		if !cond {
			return nil, nil
		}
		ctrlSig, err = e.evalBlock(node.Block, false)
		if err != nil {
			return nil, err
		}
		if ctrlSig != nil {
			switch ctrlSig.Kind {
			case CtrlBreak:
				return nil, nil
			case CtrlContinue:
				ctrlSig, err = e.evalAssign(&node.Assign)
				if err != nil || ctrlSig != nil {
					return ctrlSig, err
				}
				continue
			default:
				return ctrlSig, nil
			}
		}
		ctrlSig, err = e.evalAssign(&node.Assign)
		if err != nil || ctrlSig != nil {
			return ctrlSig, err
		}
	}
}

func (e *Evaluator) evalAssign(assign *parser.Assign) (*ControlSignal, error) {

	values, ctrlSig, err := e.evalExprsAsSingles(assign.Exprs)
	if err != nil || ctrlSig != nil {
		return ctrlSig, err
	}
	for i, id := range assign.Ids {
		if err := e.setValueForId(id, values[i]); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (e *Evaluator) evalShortDecl(shortDecl *parser.ShortDecl) (*ControlSignal, error) {

	values, ctrlSig, err := e.evalExprsAsSingles(shortDecl.Exprs)
	if err != nil || ctrlSig != nil {
		return ctrlSig, err
	}
	for i, id := range shortDecl.Ids {
		if err := e.setValueForId(id, values[i]); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (e *Evaluator) evalVarDecl(node *parser.VarDecl) (*ControlSignal, error) {

	if len(node.ExprsOrNil) == 0 {
		zero := ZeroValueForType(node.Type)
		for _, id := range node.Ids {
			if err := e.setValueForId(id, zero); err != nil {
				return nil, err
			}
		}
		return nil, nil
	}
	values, ctrlSig, err := e.evalExprsAsSingles(node.ExprsOrNil)
	if err != nil || ctrlSig != nil {
		return ctrlSig, err
	}
	for i, id := range node.Ids {
		if err := e.setValueForId(id, values[i]); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (e *Evaluator) evalReturn(node *parser.Return) (*ControlSignal, error) {
	if len(node.ExprsOrNil) == 0 {
		return newControlSignal(CtrlReturn, []Value{}), nil
	}
	values, ctrlSig, err := e.evalExprsAsSingles(node.ExprsOrNil)
	if err != nil || ctrlSig != nil {
		return ctrlSig, err
	}
	return newControlSignal(CtrlReturn, values), nil
}

func (e *Evaluator) evalExprsAsSingles(exprs []parser.Expr) ([]Value, *ControlSignal, error) {
	values := make([]Value, 0, len(exprs))
	for _, expr := range exprs {
		vals, ctrlSig, err := e.Valuate(expr)
		if err != nil || ctrlSig != nil {
			return nil, ctrlSig, err
		}
		values = append(values, vals...)
	}
	return values, nil, nil
}

func (e *Evaluator) evalBoolExpr(expr parser.Expr) (bool, *ControlSignal, error) {
	values, ctrlSig, err := e.Valuate(expr)
	if err != nil || ctrlSig != nil {
		return false, ctrlSig, err
	}
	val, err := expectSingle(values, "condition")
	if err != nil {
		return false, nil, err
	}
	boolVal, ok := val.(*BoolValue)
	if !ok {
		return false, nil, fmt.Errorf("condition expects bool")
	}
	return boolVal.Value, nil, nil
}

func (e *Evaluator) setValueForId(id parser.Id, value Value) error {
	ref, ok := e.resolveTable[id.IdId]
	if !ok {
		return fmt.Errorf("missing resolve entry for id: %s", id.String())
	}
	switch ref.Kind {
	case resolver.RefBuiltin:
		return fmt.Errorf("cannot assign to builtin")
	case resolver.RefGlobal:
		env := e.globalEnvFrame
		if ref.Slot < 0 {
			return fmt.Errorf("negative slot for id")
		}
		if ref.Slot >= len(env.Slots) {
			env.Slots = growSlots(env.Slots, ref.Slot+1)
		}
		env.Slots[ref.Slot] = value
		return nil
	case resolver.RefLocal:
		env, err := envAtDistance(e.CurrentEnv(), ref.Distance)
		if err != nil {
			return err
		}
		if ref.Slot < 0 {
			return fmt.Errorf("negative slot for id")
		}
		if ref.Slot >= len(env.Slots) {
			env.Slots = growSlots(env.Slots, ref.Slot+1)
		}
		env.Slots[ref.Slot] = value
		return nil
	default:
		return fmt.Errorf("unknown ref kind")
	}
}
