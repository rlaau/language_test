package evaluator

import (
	"strconv"

	"github.com/rlaaudgjs5638/langTest/tinygo/parser"
)

// Value_model에서의 제로값
var ZeroValueForType = func(t parser.Type) Value {
	switch t.TypeKind {
	case parser.IntType:
		return newIntVal(0)
	case parser.BoolType:
		return newBoolVal(false)
	case parser.StringType:
		return newStringVal("")
	case parser.ErrorType:
		return newOkErrorVal()
	case parser.FuncionType:
		return newClosureVal(nil, nil, nil, parser.Block{}, nil)
	default:
		return nil
	}
}

// 이렇게 Value 인터페이스 기준으로 값 모델링을 해 두면
// 추후 Slice, Mape, Chan등의 런타임 구조체를 추가 시에
// 이들이 Object 인터페이스만 충족시키게 한 후에
// 별도 구조체 런타임 구조체로써 작동하게 할 수 있음
type Value interface {
	Kind() ValueKind
	Inspect() string
}
type ValueKind int

const (
	IntKind ValueKind = iota
	BoolKind
	StrKind
	ErrKind
	ClosureKind
	BuiltinFuncKind
)

type IntValue struct {
	Value int64
}

func newIntVal(v int64) *IntValue {
	return &IntValue{Value: v}
}
func (i *IntValue) Kind() ValueKind {
	return IntKind
}
func (i *IntValue) Inspect() string {
	return strconv.Itoa(int(i.Value))
}

type BoolValue struct {
	Value bool
}

func newBoolVal(v bool) *BoolValue {
	return &BoolValue{Value: v}
}
func (b *BoolValue) Kind() ValueKind {
	return BoolKind
}
func (b *BoolValue) Inspect() string {
	if b.Value {
		return "true"
	}
	return "false"
}

type StringValue struct {
	Value string
}

func newStringVal(v string) *StringValue {
	return &StringValue{Value: v}
}
func (s *StringValue) Kind() ValueKind {
	return StrKind
}
func (s *StringValue) Inspect() string {
	return s.Value
}

type ErrorValue struct {
	IsOk   bool
	ErrMsg string
}

func newOkErrorVal() *ErrorValue {
	return &ErrorValue{IsOk: true}
}
func newErrorVal(msg string) *ErrorValue {
	return &ErrorValue{IsOk: false, ErrMsg: msg}
}
func (e *ErrorValue) Kind() ValueKind {
	return ErrKind
}
func (e *ErrorValue) Inspect() string {
	if e.IsOk {
		return "ok"
	}
	return e.ErrMsg
}

type ClosureValue struct {
	IdOrNil     *parser.Id
	Params      []parser.Param
	ReturnTypes []parser.Type
	Block       parser.Block
	Env         *EnvFrame // captured env
}

func newClosureVal(idOrNil *parser.Id, params []parser.Param, returnTypes []parser.Type, block parser.Block, env *EnvFrame) *ClosureValue {
	return &ClosureValue{
		IdOrNil:     idOrNil,
		Params:      params,
		ReturnTypes: returnTypes,
		Block:       block,
		Env:         env,
	}
}
func (c *ClosureValue) Kind() ValueKind {
	return ClosureKind
}
func (c *ClosureValue) Inspect() string {
	if c.IdOrNil == nil {
		return "closure<anonymous>"
	}
	return "closure<" + c.IdOrNil.Name + ">"
}
