package evaluator

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type BuiltinFunc struct {
	Name string
	Impl func(e *Evaluator, args []Value) ([]Value, *ControlSignal, error)
}

type BuiltinFuncValue struct {
	Func BuiltinFunc
}

func newBuiltinFuncVal(fn BuiltinFunc) *BuiltinFuncValue {
	return &BuiltinFuncValue{Func: fn}
}

func (b *BuiltinFuncValue) Kind() ValueKind {
	return BuiltinFuncKind
}
func (b *BuiltinFuncValue) Inspect() string {
	return "builtin<" + b.Func.Name + ">"
}

func builtinByName(name string) (BuiltinFunc, bool) {
	fn, ok := builtinRegistry[name]
	return fn, ok
}

var builtinRegistry = map[string]BuiltinFunc{
	"newError": {
		Name: "newError",
		Impl: func(e *Evaluator, args []Value) ([]Value, *ControlSignal, error) {
			if len(args) != 1 {
				return nil, nil, fmt.Errorf("newError expects 1 argument")
			}
			str, ok := args[0].(*StringValue)
			if !ok {
				return nil, nil, fmt.Errorf("newError expects string")
			}
			return []Value{newErrorVal(&str.Value)}, nil, nil
		},
	},
	"errString": {
		Name: "errString",
		Impl: func(e *Evaluator, args []Value) ([]Value, *ControlSignal, error) {
			if len(args) != 1 {
				return nil, nil, fmt.Errorf("errString expects 1 argument")
			}
			errVal, ok := args[0].(*ErrorValue)
			if !ok {
				return nil, nil, fmt.Errorf("errString expects error")
			}
			if errVal.IsOk {
				return []Value{newStringVal("")}, nil, nil
			}
			return []Value{newStringVal(errVal.ErrMsg)}, nil, nil
		},
	},
	"len": {
		Name: "len",
		Impl: func(e *Evaluator, args []Value) ([]Value, *ControlSignal, error) {
			if len(args) != 1 {
				return nil, nil, fmt.Errorf("len expects 1 argument")
			}
			str, ok := args[0].(*StringValue)
			if !ok {
				return nil, nil, fmt.Errorf("len expects string")
			}
			return []Value{newIntVal(int64(len(str.Value)))}, nil, nil
		},
	},
	"scan": {
		Name: "scan",
		Impl: func(e *Evaluator, args []Value) ([]Value, *ControlSignal, error) {
			if len(args) != 1 {
				return nil, nil, fmt.Errorf("scan expects 1 argument")
			}
			_, ok := args[0].(*StringValue)
			if !ok {
				return nil, nil, fmt.Errorf("scan expects string")
			}
			reader := bufio.NewReader(os.Stdin)
			line, err := reader.ReadString('\n')
			if err != nil && err != io.EOF {
				return nil, nil, err
			}
			line = strings.TrimRight(line, "\r\n")
			return []Value{newStringVal(line)}, nil, nil
		},
	},
	"print": {
		Name: "print",
		Impl: func(e *Evaluator, args []Value) ([]Value, *ControlSignal, error) {
			if len(args) != 1 {
				return nil, nil, fmt.Errorf("print expects 1 argument")
			}
			str, ok := args[0].(*StringValue)
			if !ok {
				return nil, nil, fmt.Errorf("print expects string")
			}
			fmt.Fprint(os.Stdout, str.Value)
			return []Value{}, nil, nil
		},
	},
	"panic": {
		Name: "panic",
		Impl: func(e *Evaluator, args []Value) ([]Value, *ControlSignal, error) {
			if len(args) != 1 {
				return nil, nil, fmt.Errorf("panic expects 1 argument")
			}
			_, ok := args[0].(*StringValue)
			if !ok {
				return nil, nil, fmt.Errorf("panic expects string")
			}
			return nil, newPanicSignal(args), nil
		},
	},
}
