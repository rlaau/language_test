package evaluator

import (
	"testing"

	"github.com/rlaaudgjs5638/langTest/tinygo/lexer"
	"github.com/rlaaudgjs5638/langTest/tinygo/parser"
	"github.com/rlaaudgjs5638/langTest/tinygo/resolver"
)

func TestEvalMain_AssignGlobal(t *testing.T) {
	input := "var a int = 1; func main(){ a = 3; }"
	e, pkg := evalMainFromInput(t, input)
	aVal := getGlobalValue(t, e, pkg, "a")
	intVal, ok := aVal.(*IntValue)
	if !ok || intVal.Value != 3 {
		t.Fatalf("expected a=3, got %v", aVal.Inspect())
	}
}

func TestEvalMain_ForWithAssign_ControlFlow(t *testing.T) {
	input := "var sum int = 0; func main(){ for i := 0; i < 5; i = i + 1; { if i == 2 { continue; } if i == 4 { break; } sum = sum + i; } }"
	e, pkg := evalMainFromInput(t, input)
	sumVal := getGlobalValue(t, e, pkg, "sum")
	intVal, ok := sumVal.(*IntValue)
	if !ok || intVal.Value != 4 {
		t.Fatalf("expected sum=4, got %v", sumVal.Inspect())
	}
}
func TestEvalMain_For_Return_ControlFlow(t *testing.T) {
	input := "var sum int = 0; func add(){return 4;} func bb(){print(\"aa\");} func main(){ for i := 0; i < 5; i = i + 1; { if i == 2 { continue; } if i == 4 { break; } sum = sum + i; } }"
	e, pkg := evalMainFromInput(t, input)
	sumVal := getGlobalValue(t, e, pkg, "sum")
	intVal, ok := sumVal.(*IntValue)
	if !ok || intVal.Value != 4 {
		t.Fatalf("expected sum=4, got %v", sumVal.Inspect())
	}
}

func TestEvalMain_FunctionReturn(t *testing.T) {
	input := "func main(){ out = f(); } var out int = 0; func f() int { return 7; } "
	e, pkg := evalMainFromInput(t, input)
	outVal := getGlobalValue(t, e, pkg, "out")
	intVal, ok := outVal.(*IntValue)
	if !ok || intVal.Value != 7 {
		t.Fatalf("expected out=7, got %v", outVal.Inspect())
	}
}

func TestEvalMain_BlockShadowing(t *testing.T) {
	input := "var a int = 1; func main(){ { a := 2; } }"
	e, pkg := evalMainFromInput(t, input)
	aVal := getGlobalValue(t, e, pkg, "a")
	intVal, ok := aVal.(*IntValue)
	if !ok || intVal.Value != 1 {
		t.Fatalf("expected a=1, got %v", aVal.Inspect())
	}
}
func TestEvalMain_HigherOrderCounter(t *testing.T) {
	input := "var a int = 0; var b int = 0; var c int = 0; var d int = 0; func GetCounter() func() int { count := 0; func inc() int { count = count + 1; return count; } return inc; } func main(){ c1 := GetCounter(); c2 := GetCounter(); a = c1(); b = c1(); c = c2(); d = c2(); }"
	e, pkg := evalMainFromInput(t, input)
	aVal := getGlobalValue(t, e, pkg, "a").(*IntValue)
	bVal := getGlobalValue(t, e, pkg, "b").(*IntValue)
	cVal := getGlobalValue(t, e, pkg, "c").(*IntValue)
	dVal := getGlobalValue(t, e, pkg, "d").(*IntValue)
	if aVal.Value != 1 || bVal.Value != 2 || cVal.Value != 1 || dVal.Value != 2 {
		t.Fatalf("expected counters 1,2,1,2 got %d,%d,%d,%d", aVal.Value, bVal.Value, cVal.Value, dVal.Value)
	}
}

func TestEvalMain_ChainedCall(t *testing.T) {
	input := "var r int = 0; func make() func() int { return func() int { return 5; }; } func main(){ r = make()(); }"
	e, pkg := evalMainFromInput(t, input)
	rVal := getGlobalValue(t, e, pkg, "r").(*IntValue)
	if rVal.Value != 5 {
		t.Fatalf("expected r=5, got %v", rVal.Inspect())
	}
}

func TestEvalMain_ChainedCall_ErrorOnMultiReturn(t *testing.T) {
	input := "var r int = 0; func make() (func() int, int) { return func() int { return 1; }, 2; } func main(){ r = make()(); }"
	e, pkg := buildEvaluatorFromInput(t, input)
	err := e.EvalMainFunc()
	if err == nil {
		t.Fatalf("expected error")
	}
	if err.Error() != "invalid call: the callee must evaluate to a single function" {
		t.Fatalf("unexpected error: %v", err)
	}
	rVal := getGlobalValue(t, e, pkg, "r").(*IntValue)
	if rVal.Value != 0 {
		t.Fatalf("expected r=0, got %v", rVal.Inspect())
	}
}

func TestEvalMain_DivisionByZero_Panic(t *testing.T) {
	input := "func main(){ a := 1 / 0; }"
	_, err := evalMainExpectError(t, input)
	if err == nil {
		t.Fatalf("expected error")
	}
	if err.Error() != "division by zero" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEvalMain_MultiAssignAndReturn(t *testing.T) {
	input := "var a int = 0; var b int = 0; func pair() (int, int) { return 1, 2; } func main(){ a, b = pair(); a, b = b, a; }"
	e, pkg := evalMainFromInput(t, input)
	aVal := getGlobalValue(t, e, pkg, "a").(*IntValue)
	bVal := getGlobalValue(t, e, pkg, "b").(*IntValue)
	if aVal.Value != 2 || bVal.Value != 1 {
		t.Fatalf("expected a=2 b=1, got %d %d", aVal.Value, bVal.Value)
	}
}

func TestEvalMain_ApplyFunctionArg(t *testing.T) {
	input := "var r1 int = 0; var r2 int = 0; func apply(act func(int,int) int, a int, b int) int { return act(a,b); } func add(a int, b int) int { return a + b; } func mul(a int, b int) int { return a * b; } func main(){ r1 = apply(add, 2, 3); r2 = apply(mul, 2, 3); }"
	e, pkg := evalMainFromInput(t, input)
	r1Val := getGlobalValue(t, e, pkg, "r1").(*IntValue)
	r2Val := getGlobalValue(t, e, pkg, "r2").(*IntValue)
	if r1Val.Value != 5 || r2Val.Value != 6 {
		t.Fatalf("expected r1=5 r2=6, got %d %d", r1Val.Value, r2Val.Value)
	}
}

func TestEvalMain_PanicUnwind(t *testing.T) {
	input := "var after int = 0; func boom(){ panic(\"boom\"); } func main(){ after = 1; boom(); after = 2; }"
	e, pkg := buildEvaluatorFromInput(t, input)
	err := e.EvalMainFunc()
	if err == nil {
		t.Fatalf("expected error")
	}
	if err.Error() != "panic: boom" {
		t.Fatalf("unexpected error: %v", err)
	}
	afterVal := getGlobalValue(t, e, pkg, "after").(*IntValue)
	if afterVal.Value != 1 {
		t.Fatalf("expected after=1, got %v", afterVal.Inspect())
	}
}

func TestEvalMain_ReturnFromNestedBlock(t *testing.T) {
	input := "var out int = 0; func f() int { if true { { return 3; } } return 4; } func main(){ out = f(); }"
	e, pkg := evalMainFromInput(t, input)
	outVal := getGlobalValue(t, e, pkg, "out").(*IntValue)
	if outVal.Value != 3 {
		t.Fatalf("expected out=3, got %v", outVal.Inspect())
	}
}

func TestEvalMain_HoistingFunctionInit(t *testing.T) {
	input := "var a int = inc(); func inc() int { return 2; } func main(){ }"
	e, pkg := evalMainFromInput(t, input)
	aVal := getGlobalValue(t, e, pkg, "a").(*IntValue)
	if aVal.Value != 2 {
		t.Fatalf("expected a=2, got %v", aVal.Inspect())
	}
}

func TestEvalMain_BuiltinInterop(t *testing.T) {
	input := "var s string = \"\"; var l int = 0; func main(){ s = errString(newError(\"boom\")); l = len(\"hi\"); }"
	e, pkg := evalMainFromInput(t, input)
	sVal, ok := getGlobalValue(t, e, pkg, "s").(*StringValue)
	if !ok || sVal.Value != "boom" {
		t.Fatalf("expected s=boom, got %v", getGlobalValue(t, e, pkg, "s").Inspect())
	}
	lVal := getGlobalValue(t, e, pkg, "l").(*IntValue)
	if lVal.Value != 2 {
		t.Fatalf("expected l=2, got %v", lVal.Inspect())
	}
}

func evalMainFromInput(t *testing.T, input string) (*Evaluator, *parser.PackageAST) {
	t.Helper()
	e, pkg := buildEvaluatorFromInput(t, input)
	if err := e.EvalMainFunc(); err != nil {
		t.Fatalf("EvalMainFunc error: %v", err)
	}
	return e, pkg
}

func evalMainExpectError(t *testing.T, input string) (*Evaluator, error) {
	t.Helper()
	e, _ := buildEvaluatorFromInput(t, input)
	return e, e.EvalMainFunc()
}

func buildEvaluatorFromInput(t *testing.T, input string) (*Evaluator, *parser.PackageAST) {
	t.Helper()
	pkg := parsePackageForEval(t, input)
	table, hoist, order, builtins, err := resolver.Resolve(pkg)
	if err != nil {
		t.Fatalf("resolve error: %v", err)
	}
	e, err := NewEvaluator(*pkg, hoist, order, table, builtins)
	if err != nil {
		t.Fatalf("new evaluator error: %v", err)
	}
	return e, pkg
}

func parsePackageForEval(t *testing.T, input string) *parser.PackageAST {
	t.Helper()
	lx := lexer.NewLexer()
	lx.Set(input)
	ps := parser.NewParser(lx)
	pkg, err := ps.ParsePackage()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	return pkg
}

func getGlobalValue(t *testing.T, e *Evaluator, pkg *parser.PackageAST, name string) Value {
	t.Helper()
	id := findGlobalVarId(t, pkg, name)
	ref, ok := e.resolveTable[id.IdId]
	if !ok {
		t.Fatalf("missing resolve entry for %s", name)
	}
	if ref.Kind != resolver.RefGlobal {
		t.Fatalf("expected global ref for %s, got %v", name, ref.Kind)
	}
	if ref.Slot < 0 || ref.Slot >= len(e.globalEnvFrame.Slots) {
		t.Fatalf("global slot out of range for %s", name)
	}
	return e.globalEnvFrame.Slots[ref.Slot]
}

func findGlobalVarId(t *testing.T, pkg *parser.PackageAST, name string) parser.Id {
	t.Helper()
	for _, decl := range pkg.DeclsOrNil {
		vd, ok := decl.(*parser.VarDecl)
		if !ok {
			continue
		}
		for _, id := range vd.Ids {
			if id.Name == name {
				return id
			}
		}
	}
	t.Fatalf("global var not found: %s", name)
	return parser.Id{}
}
