package evaluator

type ControlSignal struct {
	Kind   ControlKind
	Values []Value
}
type ControlKind int

const (
	CtrlReturn ControlKind = iota
	CtrlBreak
	CtrlPanic
)
