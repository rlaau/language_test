package evaluator

import "github.com/rlaaudgjs5638/langTest/tinygo/parser"

type EvalPanic struct {
	callFrame   CallFrame
	occuredNode parser.Node
	errMsg      string
	tailError   error
}

func (e *EvalPanic) Error() string {
	line1 := "runtime panic:" + e.errMsg
	line2 := e.callFrame.String()
	line3 := "beacuse->" + e.tailError.Error()
	lines := []string{line1, line2, line3}
	return parser.JoinLines(lines)
}
func NewEvalError(callFrame CallFrame, occuredNode parser.Node, errMsg string, tailError error) *EvalPanic {
	return &EvalPanic{
		callFrame:   callFrame,
		occuredNode: occuredNode,
		errMsg:      errMsg,
		tailError:   tailError,
	}
}
