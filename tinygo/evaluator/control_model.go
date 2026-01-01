package evaluator

type ControlSignal struct {
	Kind   ControlKind
	Values []Value
}
type ControlKind int

const (
	CtrlReturn ControlKind = iota
	CtrlBreak
	CtrlContinue
	CtrlPanic
)

func newControlSignal(kind ControlKind, values []Value) *ControlSignal {
	return &ControlSignal{Kind: kind, Values: values}
}

func newPanicSignal(values []Value) *ControlSignal {
	return newControlSignal(CtrlPanic, values)
}
