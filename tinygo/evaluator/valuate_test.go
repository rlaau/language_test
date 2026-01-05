package evaluator

import (
	"testing"

	"github.com/rlaaudgjs5638/langTest/tinygo/parser"
	"github.com/rlaaudgjs5638/langTest/tinygo/resolver"
)

func TestValuate_ComplexExpression(t *testing.T) {
	e := newTestEvaluator(resolver.ResolveTable{}, &EnvFrame{Slots: []Value{}})
	expr := bin(
		parser.MinusBinary,
		bin(parser.Plus, intExpr(1), bin(parser.Mul, intExpr(2), intExpr(3))),
		bin(parser.Div, intExpr(4), intExpr(2)),
	)
	values, ctrl, err := e.Valuate(expr)
	if err != nil || ctrl != nil {
		t.Fatalf("unexpected error: %v ctrl=%v", err, ctrl)
	}
	if len(values) != 1 {
		t.Fatalf("expected single value, got %d", len(values))
	}
	got, ok := values[0].(*IntValue)
	if !ok || got.Value != 5 {
		t.Fatalf("expected 5, got %v", values[0].Inspect())
	}
}

func TestValuate_TypeOperations(t *testing.T) {
	tests := []struct {
		name    string
		expr    parser.Expr
		wantVal Value
		wantErr string
	}{
		{
			name:    "string_plus_string",
			expr:    bin(parser.Plus, strExpr("a"), strExpr("b")),
			wantVal: newStringVal("ab"),
		},
		{
			name:    "string_minus_string",
			expr:    bin(parser.MinusBinary, strExpr("a"), strExpr("b")),
			wantErr: "arithmetic op expects int",
		},
		{
			name:    "bool_plus_int",
			expr:    bin(parser.Plus, boolExpr(true), intExpr(1)),
			wantErr: "arithmetic op expects int",
		},
		{
			name:    "bool_and_bool",
			expr:    bin(parser.And, boolExpr(true), boolExpr(false)),
			wantVal: newBoolVal(false),
		},
		{
			name:    "equality_type_mismatch",
			expr:    bin(parser.Equal, intExpr(1), boolExpr(true)),
			wantErr: "equality op expects same comparable types",
		},
	}

	e := newTestEvaluator(resolver.ResolveTable{}, &EnvFrame{Slots: []Value{}})
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			values, ctrl, err := e.Valuate(tc.expr)
			if tc.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error %q", tc.wantErr)
				}
				if err.Error() != tc.wantErr {
					t.Fatalf("expected error %q, got %q", tc.wantErr, err.Error())
				}
				if ctrl != nil {
					t.Fatalf("unexpected control signal: %v", ctrl.Kind)
				}
				return
			}
			if err != nil || ctrl != nil {
				t.Fatalf("unexpected error: %v ctrl=%v", err, ctrl)
			}
			if len(values) != 1 {
				t.Fatalf("expected single value, got %d", len(values))
			}
			if values[0].Inspect() != tc.wantVal.Inspect() || values[0].Kind() != tc.wantVal.Kind() {
				t.Fatalf("expected %s, got %s", tc.wantVal.Inspect(), values[0].Inspect())
			}
		})
	}
}

func TestValuate_ErrorMessages(t *testing.T) {
	tests := []struct {
		name    string
		expr    parser.Expr
		wantErr string
	}{
		{
			name:    "unary_minus_on_string",
			expr:    &parser.Unary{Op: parser.MinusUnary, Object: strExpr("x")},
			wantErr: "unary - expects int",
		},
		{
			name:    "logical_and_on_ints",
			expr:    bin(parser.And, intExpr(1), intExpr(2)),
			wantErr: "logical op expects bool",
		},
	}

	e := newTestEvaluator(resolver.ResolveTable{}, &EnvFrame{Slots: []Value{}})
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, ctrl, err := e.Valuate(tc.expr)
			if err == nil {
				t.Fatalf("expected error %q", tc.wantErr)
			}
			if err.Error() != tc.wantErr {
				t.Fatalf("expected error %q, got %q", tc.wantErr, err.Error())
			}
			if ctrl != nil {
				t.Fatalf("unexpected control signal: %v", ctrl.Kind)
			}
		})
	}
}

func TestValuate_DivisionByZero(t *testing.T) {
	e := newTestEvaluator(resolver.ResolveTable{}, &EnvFrame{Slots: []Value{}})
	expr := bin(parser.Div, intExpr(1), intExpr(0))
	_, ctrl, err := e.Valuate(expr)
	if err == nil {
		t.Fatalf("expected error: division by zero")
	}

	if ctrl != nil && ctrl.Kind == CtrlPanic {
		t.Fatalf("expected panic by error. NOT BY CONTROL SIGNAL")
	}
}

func TestValuate_FunctionValue(t *testing.T) {
	env := &EnvFrame{Slots: []Value{}}
	e := newTestEvaluator(resolver.ResolveTable{}, env)
	expr := funcExpr(nil, nil, parser.Block{})
	values, ctrl, err := e.Valuate(expr)
	if err != nil || ctrl != nil {
		t.Fatalf("unexpected error: %v ctrl=%v", err, ctrl)
	}
	if len(values) != 1 {
		t.Fatalf("expected single value, got %d", len(values))
	}
	closure, ok := values[0].(*ClosureValue)
	if !ok {
		t.Fatalf("expected closure value")
	}
	if closure.ParentEnv != env {
		t.Fatalf("expected closure to capture current env")
	}
}

func TestValuate_Builtins(t *testing.T) {
	rt := resolver.ResolveTable{}
	env := &EnvFrame{Slots: []Value{}}
	e := newTestEvaluator(rt, env)

	addBuiltinRef(rt, 10, "newError")
	addBuiltinRef(rt, 11, "errString")
	addBuiltinRef(rt, 12, "len")
	addBuiltinRef(rt, 13, "panic")

	newErrCall := &parser.Call{
		PrimaryOrNil: *idPrimary("newError", 10),
		ArgsList:     []parser.Args{{strExpr("boom")}},
	}
	values, ctrl, err := e.Valuate(newErrCall)
	if err != nil || ctrl != nil {
		t.Fatalf("unexpected error: %v ctrl=%v", err, ctrl)
	}
	errVal, ok := values[0].(*ErrorValue)
	if !ok || errVal.IsOk || errVal.ErrMsg != "boom" {
		t.Fatalf("unexpected newError result: %v", values[0].Inspect())
	}

	errStringCall := &parser.Call{
		PrimaryOrNil: *idPrimary("errString", 11),
		ArgsList:     []parser.Args{{errExpr(nil)}},
	}
	values, ctrl, err = e.Valuate(errStringCall)
	if err != nil || ctrl != nil {
		t.Fatalf("unexpected error: %v ctrl=%v", err, ctrl)
	}
	strVal, ok := values[0].(*StringValue)
	if !ok || strVal.Value != "" {
		t.Fatalf("unexpected errString result: %v", values[0].Inspect())
	}

	lenCall := &parser.Call{
		PrimaryOrNil: *idPrimary("len", 12),
		ArgsList:     []parser.Args{{strExpr("hi")}},
	}
	values, ctrl, err = e.Valuate(lenCall)
	if err != nil || ctrl != nil {
		t.Fatalf("unexpected error: %v ctrl=%v", err, ctrl)
	}
	intVal, ok := values[0].(*IntValue)
	if !ok || intVal.Value != 2 {
		t.Fatalf("unexpected len result: %v", values[0].Inspect())
	}

	panicCall := &parser.Call{
		PrimaryOrNil: *idPrimary("panic", 13),
		ArgsList:     []parser.Args{{strExpr("panic")}},
	}
	_, ctrl, err = e.Valuate(panicCall)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ctrl == nil || ctrl.Kind != CtrlPanic {
		t.Fatalf("expected CtrlPanic, got %v", ctrl)
	}
}

func newTestEvaluator(rt resolver.ResolveTable, env *EnvFrame) *Evaluator {
	if env == nil {
		env = &EnvFrame{Slots: []Value{}}
	}
	e := &Evaluator{
		resolveTable: rt,
		callStack: CallStack{callFrames: []CallFrame{
			{
				currentEnv: env,
			},
		}},
		globalEnvFrame: env,
		builtInSlots:   []Value{},
	}
	builtins := map[string]int{}
	for i, name := range resolver.Builtins {
		builtins[name] = i
	}
	e.builtInSlots = make([]Value, maxBuiltinSlot(builtins)+1)
	for name, slot := range builtins {
		fn, ok := builtinByName(name)
		if !ok {
			continue
		}
		e.builtInSlots[slot] = newBuiltinFuncVal(fn)
	}
	return e
}

func addBuiltinRef(rt resolver.ResolveTable, idId int64, name string) {
	slot := -1
	for i, n := range resolver.Builtins {
		if n == name {
			slot = i
			break
		}
	}
	if slot < 0 {
		return
	}
	rt[parser.IdId(idId)] = resolver.ResolvedRef{
		Kind:        resolver.RefBuiltin,
		Distance:    0,
		Slot:        slot,
		Name:        name,
		RefIdNodeId: parser.IdId(-1),
	}
}

func intExpr(v int) parser.Expr {
	val := v
	return &parser.Primary{
		PrimaryKind: parser.ValuePrimary,
		ValueOrNil: &parser.ValueForm{
			ValueKind:   parser.NumberValue,
			NumberOrNil: &val,
		},
	}
}

func boolExpr(v bool) parser.Expr {
	val := v
	return &parser.Primary{
		PrimaryKind: parser.ValuePrimary,
		ValueOrNil: &parser.ValueForm{
			ValueKind: parser.BoolValue,
			BoolOrNil: &val,
		},
	}
}

func strExpr(s string) parser.Expr {
	val := s
	return &parser.Primary{
		PrimaryKind: parser.ValuePrimary,
		ValueOrNil: &parser.ValueForm{
			ValueKind:   parser.StrLitValue,
			StrLitOrNil: &val,
		},
	}
}

func errExpr(msg *string) parser.Expr {
	return &parser.Primary{
		PrimaryKind: parser.ValuePrimary,
		ValueOrNil: &parser.ValueForm{
			ValueKind:    parser.ErrValue,
			ErrOrNilIfOk: msg,
		},
	}
}

func idPrimary(name string, idId int64) *parser.Primary {
	return &parser.Primary{
		PrimaryKind: parser.IdPrimary,
		IdOrNil:     &parser.Id{Name: name, IdId: parser.IdId(idId)},
	}
}

func bin(op parser.BinaryKind, left, right parser.Expr) parser.Expr {
	return &parser.Binary{Op: op, LeftExpr: left, RightExpr: right}
}

func funcExpr(params []parser.Param, returns []parser.Type, block parser.Block) parser.Expr {
	return &parser.Primary{
		PrimaryKind: parser.ValuePrimary,
		ValueOrNil: &parser.ValueForm{
			ValueKind: parser.FexpValue,
			FexpOrNil: &parser.Fexp{
				ParamsOrNil:      params,
				ReturnTypesOrNil: returns,
				Block:            block,
			},
		},
	}
}
