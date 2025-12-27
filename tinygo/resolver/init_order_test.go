package resolver

import "testing"

func initOrderNames(order InitOrder, hoist *HoistInfo) []string {
	names := make([]string, 0, len(order))
	for _, step := range order {
		sym := hoist.getById(step.VarId)
		if sym == nil {
			names = append(names, "<missing>")
			continue
		}
		names = append(names, sym.name)
	}
	return names
}

func TestInitOrder_TopologyAndZeroInit(t *testing.T) {
	input := "var a int = b; var b int = c; var c int; var d int = a;"
	pkg, table, hoist, err := resolveFromInput(t, input)
	if err != nil {
		t.Fatalf("unexpected resolve error: %v", err)
	}
	order, err := BuildInitOrder(pkg, table, hoist)
	if err != nil {
		t.Fatalf("unexpected init order error: %v", err)
	}
	names := initOrderNames(order, hoist)
	want := []string{"c", "b", "a", "d"}
	if len(names) != len(want) {
		t.Fatalf("init order length mismatch: got %v want %v", names, want)
	}
	for i := range want {
		if names[i] != want[i] {
			t.Fatalf("init order mismatch: got %v want %v", names, want)
		}
	}
}

func TestInitOrder_SkipFexp(t *testing.T) {
	input := "var f func() int = func() int { return 1; }; var a int = 1;"
	pkg, table, hoist, err := resolveFromInput(t, input)
	if err != nil {
		t.Fatalf("unexpected resolve error: %v", err)
	}
	order, err := BuildInitOrder(pkg, table, hoist)
	if err != nil {
		t.Fatalf("unexpected init order error: %v", err)
	}
	names := initOrderNames(order, hoist)
	want := []string{"f", "a"}
	if len(names) != len(want) {
		t.Fatalf("init order length mismatch: got %v want %v", names, want)
	}
	for i := range want {
		if names[i] != want[i] {
			t.Fatalf("init order mismatch: got %v want %v", names, want)
		}
	}
}

func TestInitOrder_VarFuncCycle(t *testing.T) {
	input := "var a int = f(); func f(){ return a; }"
	pkg, table, hoist, err := resolveFromInput(t, input)
	if err != nil {
		t.Fatalf("unexpected resolve error: %v", err)
	}
	_, ierr := BuildInitOrder(pkg, table, hoist)
	if ierr == nil {
		t.Fatalf("expected var-func cycle error but got nil")
	}
}

func TestInitOrder_VarFexpCycle(t *testing.T) {
	input := "var f func() int = func() int { return a; }; var a int = f();"
	pkg, table, hoist, err := resolveFromInput(t, input)
	if err != nil {
		t.Fatalf("unexpected resolve error: %v", err)
	}
	_, ierr := BuildInitOrder(pkg, table, hoist)
	if ierr == nil {
		t.Fatalf("expected var-fexp cycle error but got nil")
	}
}

func TestInitOrder_ZeroThenFexpThenExpr(t *testing.T) {
	input := "var z int; var f func() int = func() int { return 1; }; var a int = 1; var b int = a;"
	pkg, table, hoist, err := resolveFromInput(t, input)
	if err != nil {
		t.Fatalf("unexpected resolve error: %v", err)
	}
	order, err := BuildInitOrder(pkg, table, hoist)
	if err != nil {
		t.Fatalf("unexpected init order error: %v", err)
	}
	names := initOrderNames(order, hoist)
	want := []string{"z", "f", "a", "b"}
	if len(names) != len(want) {
		t.Fatalf("init order length mismatch: got %v want %v", names, want)
	}
	for i := range want {
		if names[i] != want[i] {
			t.Fatalf("init order mismatch: got %v want %v", names, want)
		}
	}
}

func TestInitOrder_VarVarCycle(t *testing.T) {
	input := "var a int = b; var b int = a;"
	pkg, table, hoist, err := resolveFromInput(t, input)
	if err != nil {
		t.Fatalf("unexpected resolve error: %v", err)
	}
	_, ierr := BuildInitOrder(pkg, table, hoist)
	if ierr == nil {
		t.Fatalf("expected var-var cycle error but got nil")
	}
}

func TestInitOrder_MultiDeclOrder(t *testing.T) {
	input := "var a, b int = 1, 2; var c int = a; var d int = b;"
	pkg, table, hoist, err := resolveFromInput(t, input)
	if err != nil {
		t.Fatalf("unexpected resolve error: %v", err)
	}
	order, err := BuildInitOrder(pkg, table, hoist)
	if err != nil {
		t.Fatalf("unexpected init order error: %v", err)
	}
	names := initOrderNames(order, hoist)
	want := []string{"a", "b", "c", "d"}
	if len(names) != len(want) {
		t.Fatalf("init order length mismatch: got %v want %v", names, want)
	}
	for i := range want {
		if names[i] != want[i] {
			t.Fatalf("init order mismatch: got %v want %v", names, want)
		}
	}
}

func TestInitOrder_MultiDeclWithZeroAndFexp(t *testing.T) {
	input := "var a, b int; var f func() int = func() int { return a; }; var c int = f();"
	pkg, table, hoist, err := resolveFromInput(t, input)
	if err != nil {
		t.Fatalf("unexpected resolve error: %v", err)
	}
	order, err := BuildInitOrder(pkg, table, hoist)
	if err != nil {
		t.Fatalf("unexpected init order error: %v", err)
	}
	names := initOrderNames(order, hoist)
	want := []string{"a", "b", "f", "c"}
	if len(names) != len(want) {
		t.Fatalf("init order length mismatch: got %v want %v", names, want)
	}
	for i := range want {
		if names[i] != want[i] {
			t.Fatalf("init order mismatch: got %v want %v", names, want)
		}
	}
}

func TestInitOrder_VarFexpCycle_MultiDecl(t *testing.T) {
	input := "var f func() int = func() int { return a; }; var a, b int = f(), 1;"
	pkg, table, hoist, err := resolveFromInput(t, input)
	if err != nil {
		t.Fatalf("unexpected resolve error: %v", err)
	}
	_, ierr := BuildInitOrder(pkg, table, hoist)
	if ierr == nil {
		t.Fatalf("expected var-fexp cycle error but got nil")
	}
}

func TestInitOrder_VarFuncCycle_Forward(t *testing.T) {
	input := "var a int = f(); func f(){ return a; }"
	pkg, table, hoist, err := resolveFromInput(t, input)
	if err != nil {
		t.Fatalf("unexpected resolve error: %v", err)
	}
	_, ierr := BuildInitOrder(pkg, table, hoist)
	if ierr == nil {
		t.Fatalf("expected var-func cycle error but got nil")
	}
}

func TestInitOrder_ComplexDependencies(t *testing.T) {
	input := "var z int; var f func() int = func() int { return z; }; var g func() int = func() int { return a; }; var a int = 1; var b int = f(); var c int = a; var d int = c + b; var e int = g();"
	pkg, table, hoist, err := resolveFromInput(t, input)
	if err != nil {
		t.Fatalf("unexpected resolve error: %v", err)
	}
	order, err := BuildInitOrder(pkg, table, hoist)
	if err != nil {
		t.Fatalf("unexpected init order error: %v", err)
	}
	names := initOrderNames(order, hoist)
	want := []string{"z", "f", "g", "a", "b", "c", "d", "e"}
	if len(names) != len(want) {
		t.Fatalf("init order length mismatch: got %v want %v", names, want)
	}
	for i := range want {
		if names[i] != want[i] {
			t.Fatalf("init order mismatch: got %v want %v", names, want)
		}
	}
}
