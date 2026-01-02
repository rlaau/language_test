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
	// 런타임에서 호스트가 발생시킨 에러는 ControlSignal이 아닌 error리턴으로 처리됨
	//CtrlPanic은 사용자가 panic()문으로 발생시킨 패닉에 대한 것임
	CtrlPanic
)

func newControlSignal(kind ControlKind, values []Value) *ControlSignal {
	return &ControlSignal{Kind: kind, Values: values}
}

func newPanicSignal(values []Value) *ControlSignal {
	return newControlSignal(CtrlPanic, values)
}
