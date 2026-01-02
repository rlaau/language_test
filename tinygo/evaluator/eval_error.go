package evaluator

import "github.com/rlaaudgjs5638/langTest/tinygo/parser"

type EvalPanic struct {
	calledFuncId parser.Id
	occuredNode  parser.Node
	errMsg       string
	tailError    error
}

func (e *EvalPanic) Error() string {
	line1 := "runtime panic:" + e.errMsg
	line2 := e.calledFuncId.String()
	line3 := "beacuse->" + e.tailError.Error()
	lines := []string{line1, line2, line3}
	return parser.JoinLines(lines)
}
func NewEvalError(calledFuncId parser.Id, occuredNode parser.Node, errMsg string, tailError error) *EvalPanic {
	return &EvalPanic{
		calledFuncId: calledFuncId,
		occuredNode:  occuredNode,
		errMsg:       errMsg,
		tailError:    tailError,
	}
}
