package evaluator

import (
	"github.com/rlaaudgjs5638/langTest/tinygo/parser"
)

func (e *Evaluator) Valuate(expr parser.Expr) ([]Value, error) {
	//Valueate는 다음과 같은 암묵적 입력을 지님
	//builtIns []Value, resolveTable *resolver.ResolveTable, envFrame *EnvFrame, callFrame *CallFrame
	//이런 인자 폭발 우려로 인해, Valuate는 evaluator의 메서드로 취급하여 암묵적 인자를 받기로 함.
	panic("not implemented")
}
