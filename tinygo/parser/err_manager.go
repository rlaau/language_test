package parser

type errManager struct {
	currentIdx int
	errorStack []string
}

func NewErrManager() *errManager {
	return &errManager{
		currentIdx: 0,
		errorStack: []string{},
	}
}
func (e *errManager) PushErrorMsg(eString string) {
	e.errorStack = append(e.errorStack, eString)
}

func (e *errManager) GetRollback() func() {
	memorizedPosition := e.currentIdx
	return func() {
		e.currentIdx = memorizedPosition
		e.errorStack = e.errorStack[0 : e.currentIdx+1]
	}
}
func (e *errManager) GetErrorStack() []string {
	return e.errorStack
}
