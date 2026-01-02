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
		return newErrorVal(nil)
	case parser.FuncionType:
		return newClosureVal(nil, nil, nil, parser.Block{}, nil)
	default:
		return nil
	}
}

// 이렇게 Value 인터페이스 기준으로 값 모델링을 해 두면
// 추후 Slice, Mape, Chan등의 런타임 구조체를 추가 시에
// 이들이 Object 인터페이스만 충족시키게 한 후에
// 별도 런타임 구조체로써 작동하게 할 수 있음
// ex: type Slice sturct{len:int, cap:int, array:[]T{}}이후
// ex: 이게 Value인터페이스 만족시키게 하면, 이 호스트의 표현을 값처럼 들고 다니기가 가능함.
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

// newErrorVal은 msg의 nil체크로 해당 에러가 ok 값인지 아닌지 검사함
func newErrorVal(msg *string) *ErrorValue {
	if msg == nil {
		return &ErrorValue{IsOk: true, ErrMsg: ""}
	}
	return &ErrorValue{IsOk: false, ErrMsg: *msg}
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
	ParentEnv   *EnvFrame // captured env
}

func newClosureVal(idOrNil *parser.Id, params []parser.Param, returnTypes []parser.Type, block parser.Block, parentEnv *EnvFrame) *ClosureValue {
	return &ClosureValue{
		IdOrNil:     idOrNil,
		Params:      params,
		ReturnTypes: returnTypes,
		Block:       block,
		ParentEnv:   parentEnv,
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
