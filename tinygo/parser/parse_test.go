package parser

import (
	"testing"

	"github.com/rlaaudgjs5638/langTest/tinygo/lexer"
)

func parsePackageForTest(t *testing.T, input string) *Package {
	t.Helper()
	lx := lexer.NewLexer()
	lx.Set(input)
	ps := NewParser(lx)
	pkg, err := ps.ParsePackage()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	return pkg
}

func intPtr(v int) *int {
	return &v
}

func strPtr(s string) *string {
	return &s
}

func idPtr(s string) *Id {
	id := Id(s)
	return &id
}

func numPrimary(v int) *Primary {
	return newPrimary(ValuePrimary, nil, nil, newValueForm(NumberValue, intPtr(v), nil, nil, nil, nil))
}

func strPrimary(s string) *Primary {
	return newPrimary(ValuePrimary, nil, nil, newValueForm(StrLitValue, nil, nil, strPtr(s), nil, nil))
}

func boolPrimary(v bool) *Primary {
	return newPrimary(ValuePrimary, nil, nil, newValueForm(BoolValue, nil, &v, nil, nil, nil))
}

func okPrimary() *Primary {
	return newPrimary(ValuePrimary, nil, nil, newValueForm(ErrValue, nil, nil, nil, nil, nil))
}

func idPrimary(s string) *Primary {
	return newPrimary(IdPrimary, nil, idPtr(s), nil)
}

func TestParser_ParsePackage_Table(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  *Package
	}{
		{
			name:  "var_decl",
			input: "var a int = 10;",
			want: newPackage([]Decl{
				newVarDecl(
					[]Id{"a"},
					Type{TypeKind: IntType},
					[]Expr{numPrimary(10)},
				),
			}),
		},
		{
			name:  "func_add",
			input: "func add(a int, b int) int { return a + b; }",
			want: newPackage([]Decl{
				newFuncDecl(
					"add",
					[]Param{{Id: "a", Type: Type{TypeKind: IntType}}, {Id: "b", Type: Type{TypeKind: IntType}}},
					[]Type{{TypeKind: IntType}},
					Block{StmtsOrNil: []Stmt{
						newReturn([]Expr{
							newBinary(Plus, idPrimary("a"), idPrimary("b")),
						}),
					}},
				),
			}),
		},
		{
			name:  "for_range_print",
			input: "func testLoop() { for range 10 { print(\"hello\"); } }",
			want: newPackage([]Decl{
				newFuncDecl(
					"testLoop",
					[]Param{},
					[]Type{},
					Block{StmtsOrNil: []Stmt{
						newForRangeAexp(
							numPrimary(10),
							Block{StmtsOrNil: []Stmt{
								newCallStmt(*newCall(true, nil, PrintBuild, []Args{
									{strPrimary("hello")},
								})),
							}},
						),
					}},
				),
			}),
		},
		{
			name:  "short_decl",
			input: "func main() { a, b := 4, 2; }",
			want: newPackage([]Decl{
				newFuncDecl(
					"main",
					[]Param{},
					[]Type{},
					Block{StmtsOrNil: []Stmt{
						newShortDecl(
							[]Id{"a", "b"},
							[]Expr{numPrimary(4), numPrimary(2)},
						),
					}},
				),
			}),
		},
		{
			name:  "if_else_return",
			input: "func test(a int) int { if a == 0 { return 0; } else { return a; } }",
			want: newPackage([]Decl{
				newFuncDecl(
					"test",
					[]Param{{Id: "a", Type: Type{TypeKind: IntType}}},
					[]Type{{TypeKind: IntType}},
					Block{StmtsOrNil: []Stmt{
						newIf(
							nil,
							newBinary(Equal, idPrimary("a"), numPrimary(0)),
							Block{StmtsOrNil: []Stmt{
								newReturn([]Expr{numPrimary(0)}),
							}},
							&Block{StmtsOrNil: []Stmt{
								newReturn([]Expr{idPrimary("a")}),
							}},
						),
					}},
				),
			}),
		},
		{
			name:  "for_with_assign",
			input: "func loop() { for i := 0; i < 10; i = i + 1; { print(i); } }",
			want: newPackage([]Decl{
				newFuncDecl(
					"loop",
					[]Param{},
					[]Type{},
					Block{StmtsOrNil: []Stmt{
						newForWithAssign(
							*newShortDecl([]Id{"i"}, []Expr{numPrimary(0)}),
							newBinary(LessThan, idPrimary("i"), numPrimary(10)),
							*newAssign([]Id{"i"}, []Expr{newBinary(Plus, idPrimary("i"), numPrimary(1))}),
							Block{StmtsOrNil: []Stmt{
								newCallStmt(*newCall(true, nil, PrintBuild, []Args{
									{idPrimary("i")},
								})),
							}},
						),
					}},
				),
			}),
		},
		{
			name:  "built_in_new_error",
			input: "func fail() error { return newError(\"boom\"); }",
			want: newPackage([]Decl{
				newFuncDecl(
					"fail",
					[]Param{},
					[]Type{{TypeKind: ErrorType}},
					Block{StmtsOrNil: []Stmt{
						newReturn([]Expr{
							newCall(true, nil, NewErrorBuild, []Args{
								{strPrimary("boom")},
							}),
						}),
					}},
				),
			}),
		},
		{
			name:  "return_bool",
			input: "func okt() bool { return true; }",
			want: newPackage([]Decl{
				newFuncDecl(
					"okt",
					[]Param{},
					[]Type{{TypeKind: BoolType}},
					Block{StmtsOrNil: []Stmt{
						newReturn([]Expr{boolPrimary(true)}),
					}},
				),
			}),
		},
		{
			name:  "var_decl_no_assign",
			input: "var flag bool",
			want: newPackage([]Decl{
				newVarDecl(
					[]Id{"flag"},
					Type{TypeKind: BoolType},
					[]Expr{},
				),
			}),
		},
		{
			name:  "multi_return_ok",
			input: "func pair() (int, error) { return 1, ok; }",
			want: newPackage([]Decl{
				newFuncDecl(
					"pair",
					[]Param{},
					[]Type{{TypeKind: IntType}, {TypeKind: ErrorType}},
					Block{StmtsOrNil: []Stmt{
						newReturn([]Expr{
							numPrimary(1),
							okPrimary(),
						}),
					}},
				),
			}),
		},
		{
			name:  "builtin_len_call",
			input: "func count() int { return len(\"abc\"); }",
			want: newPackage([]Decl{
				newFuncDecl(
					"count",
					[]Param{},
					[]Type{{TypeKind: IntType}},
					Block{StmtsOrNil: []Stmt{
						newReturn([]Expr{
							newCall(true, nil, LenBuild, []Args{
								{strPrimary("abc")},
							}),
						}),
					}},
				),
			}),
		},
		{
			name:  "binary_precedence",
			input: "func calc() int { return 1 + 2 * 3; }",
			want: newPackage([]Decl{
				newFuncDecl(
					"calc",
					[]Param{},
					[]Type{{TypeKind: IntType}},
					Block{StmtsOrNil: []Stmt{
						newReturn([]Expr{
							newBinary(Plus, numPrimary(1), newBinary(Mul, numPrimary(2), numPrimary(3))),
						}),
					}},
				),
			}),
		},
		{
			name:  "call_stmt_id_call",
			input: "func run() { task(1, 2); }",
			want: newPackage([]Decl{
				newFuncDecl(
					"run",
					[]Param{},
					[]Type{},
					Block{StmtsOrNil: []Stmt{
						newCallStmt(*newCall(false, idPrimary("task"), NewErrorBuild, []Args{
							{numPrimary(1), numPrimary(2)},
						})),
					}},
				),
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parsePackageForTest(t, tt.input)
			gotStr := got.String()
			wantStr := tt.want.String()
			if gotStr != wantStr {
				t.Fatalf("ast mismatch:\n--- got ---\n%s\n--- want ---\n%s", gotStr, wantStr)
			}
		})
	}
}
