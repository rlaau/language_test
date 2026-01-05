package evaluator

import (
	"github.com/rlaaudgjs5638/langTest/tinygo/parser"
	"github.com/rlaaudgjs5638/langTest/tinygo/resolver"
)

func Evaluate(pkg parser.PackageAST, hoist *resolver.HoistInfo, initOrder resolver.InitOrder, table resolver.ResolveTable, builtins map[string]int) (*Evaluator, error) {
	// 1. NewEvaluator 생성
	// 2. 해당 Evaluator.EvalMainFunc() 실행
	// 3. 결과 리턴
	e, err := NewEvaluator(pkg, hoist, initOrder, table, builtins)
	if err != nil {
		return nil, err
	}
	if err := e.EvalMainFunc(); err != nil {
		return e, err
	}
	return e, nil
}
