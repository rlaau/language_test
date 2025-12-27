package resolver

import (
	"testing"

	"github.com/rlaaudgjs5638/langTest/tinygo/lexer"
	"github.com/rlaaudgjs5638/langTest/tinygo/parser"
)

func resolveFromInput(t *testing.T, input string) (*parser.PackageAST, ResolveTable, error) {
	t.Helper()
	lx := lexer.NewLexer()
	lx.Set(input)
	ps := parser.NewParser(lx)
	pkg, err := ps.ParsePackage()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	rs := NewResolver()
	table, _, rerr := rs.ResolvePackage(pkg)
	return pkg, table, rerr
}

func TestResolveNoHoist_SuccessCases(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{
			name:  "global_var_then_func_uses",
			input: "var a int = 1; func f(){ a = 2; }",
		},
		{
			name:  "recursive_func_allowed",
			input: "func f(){ f(); }",
		},
		{
			name:  "short_decl_and_assign",
			input: "func f(){ a := 1; a = 2; }",
		},
		{
			name:  "multi_short_decl_and_assign",
			input: "func f(){ a, b := 1, 2; a, b = b, a; }",
		},
		{
			name:  "if_header_scope_visible_in_body",
			input: "func f(){ if a := 1; a == 1 { a = 2; } else { a = 3; } }",
		},
		{
			name:  "for_with_assign_header_scope",
			input: "func f(){ for i := 0; i < 2; i = i + 1; { print(i); } }",
		},
		{
			name:  "shadowing_allowed_in_inner_block",
			input: "func f(){ a := 1; if true { a := 2; } }",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := resolveFromInput(t, tc.input)
			if err != nil {
				t.Fatalf("unexpected resolve error: %v", err)
			}
		})
	}
}

func TestResolveNoHoist_FailureCases(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{

		{
			name:  "assign_undefined",
			input: "func f(){ a = 1; }",
		},
		{
			name:  "short_decl_requires_new",
			input: "func f(){ a := 1; a := 2; }",
		},
		{
			name:  "duplicate_var_decl_same_scope",
			input: "var a int = 1; var a int = 2;",
		},
		{
			name:  "builtin_shadowing_forbidden",
			input: "func f(){ var print int = 1; }",
		},
		{
			name:  "assign_to_builtin_forbidden",
			input: "func f(){ print = 1; }",
		},
		{
			name:  "short_decl_lhs_rhs_mismatch",
			input: "func f(){ a, b := 1; }",
		},
		{
			name:  "assign_lhs_rhs_mismatch",
			input: "func f(){ a, b = 1; }",
		},
		{
			name:  "var_decl_lhs_rhs_mismatch",
			input: "var a, b int = 1;",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := resolveFromInput(t, tc.input)
			if err == nil {
				t.Fatalf("expected resolve error but got nil")
			}
		})
	}
}

func TestResolveNoHoist_ResolvedKinds(t *testing.T) {
	input := "var a int = 1; func f(){ a = 2; print(a); }"
	pkg, table, err := resolveFromInput(t, input)
	if err != nil {
		t.Fatalf("unexpected resolve error: %v", err)
	}

	assign := pkg.DeclsOrNil[1].(*parser.FuncDecl).Block.StmtsOrNil[0].(*parser.Assign)
	assignId := assign.Ids[0]
	assignRef := table[assignId.IdId]
	if assignRef.Kind != RefGlobal {
		t.Fatalf("expected assign id to be global, got %v", assignRef.Kind)
	}

	callStmt := pkg.DeclsOrNil[1].(*parser.FuncDecl).Block.StmtsOrNil[1].(*parser.CallStmt)
	printId := *callStmt.Call.PrimaryOrNil.IdOrNil
	printRef := table[printId.IdId]
	if printRef.Kind != RefBuiltin {
		t.Fatalf("expected builtin print, got %v", printRef.Kind)
	}
}
