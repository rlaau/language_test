package evaluator

import "github.com/rlaaudgjs5638/langTest/tinygo/parser"

type ValueKind uint8

const (
	ValInt ValueKind = iota
	ValBool
	ValString
	ValError
	ValClosure
)

// 런타임용 값.
type Value struct {
	Kind ValueKind

	IntOrNil     *int
	BoolOrNil    *bool
	StrOrNil     *string
	ErrOrOkOrNil *string // ok시엔 err이 nil임

	ClosureOrNil *ClosureValue
}

type ClosureValue struct {
	// 함수가 익명함수일 수 있으므로 Id는 포인터로 추가
	IdOrNil      *parser.Id
	Params       []parser.Param
	ReturnTypes  []parser.Type
	Block        parser.Block
	caputuredEnv *EnvFrame
}
