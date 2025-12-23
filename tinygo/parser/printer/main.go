package main

import (
	"fmt"

	"github.com/rlaaudgjs5638/langTest/tinygo/parser"
)

func main() {
	fmt.Println("=== Test 1: Simple Variable Declaration ===")
	testSimpleVarDecl()

	fmt.Println("\n=== Test 2: Function Declaration with Parameters ===")
	testFuncDeclWithParams()

	fmt.Println("\n=== Test 3: Complex Function with If Statement ===")
	testComplexFuncWithIf()

	fmt.Println("\n=== Test 4: For Loop with Range ===")
	testForLoop()

	fmt.Println("\n=== Test 5: Complete Package (divide example) ===")
	testCompletePackage()
}

// Test 1: var a int = 10
func testSimpleVarDecl() {
	num := 10
	pkg := &parser.Package{
		DeclsOrNil: []parser.Decl{
			&parser.VarDecl{
				Ids:  []parser.Id{"a"},
				Type: parser.Type{TypeKind: parser.IntType},
				ExprsOrNil: []parser.Expr{
					&parser.Primary{
						PrimaryKind: parser.ValuePrimary,
						ValueOrNil: &parser.ValueForm{
							ValueKind:   parser.NumberValue,
							NumberOrNil: &num,
						},
					},
				},
			},
		},
	}
	fmt.Println(pkg.String())
}

// Test 2: func add(a int, b int) int { return a + b }
func testFuncDeclWithParams() {
	pkg := &parser.Package{
		DeclsOrNil: []parser.Decl{
			&parser.FuncDecl{
				Id: "add",
				ParamsOrNil: []parser.Param{
					{Id: "a", Type: parser.Type{TypeKind: parser.IntType}},
					{Id: "b", Type: parser.Type{TypeKind: parser.IntType}},
				},
				ReturnTypesOrNil: []parser.Type{
					{TypeKind: parser.IntType},
				},
				Block: parser.Block{
					StmtsOrNil: []parser.Stmt{
						&parser.Return{
							ExprsOrNil: []parser.Expr{
								&parser.Binary{
									Op: parser.Plus,
									LeftExpr: &parser.Primary{
										PrimaryKind: parser.IdPrimary,
										IdOrNil:     ptrId("a"),
									},
									RightExpr: &parser.Primary{
										PrimaryKind: parser.IdPrimary,
										IdOrNil:     ptrId("b"),
									},
								},
							},
						},
					},
				},
			},
		},
	}
	fmt.Println(pkg.String())
}

// Test 3: func divide(a int, b int) (int, error) { if b == 0 { return 0, newError("error") } return a/b, nil }
func testComplexFuncWithIf() {
	zero := 0
	nilStr := "nil"
	errMsg := "can't divide by zero"

	pkg := &parser.Package{
		DeclsOrNil: []parser.Decl{
			&parser.FuncDecl{
				Id: "divide",
				ParamsOrNil: []parser.Param{
					{Id: "a", Type: parser.Type{TypeKind: parser.IntType}},
					{Id: "b", Type: parser.Type{TypeKind: parser.IntType}},
				},
				ReturnTypesOrNil: []parser.Type{
					{TypeKind: parser.IntType},
					{TypeKind: parser.ErrorType},
				},
				Block: parser.Block{
					StmtsOrNil: []parser.Stmt{
						&parser.If{
							Bexp: &parser.Binary{
								Op: parser.Equal,
								LeftExpr: &parser.Primary{
									PrimaryKind: parser.IdPrimary,
									IdOrNil:     ptrId("b"),
								},
								RightExpr: &parser.Primary{
									PrimaryKind: parser.ValuePrimary,
									ValueOrNil: &parser.ValueForm{
										ValueKind:   parser.NumberValue,
										NumberOrNil: &zero,
									},
								},
							},
							ThenBlock: parser.Block{
								StmtsOrNil: []parser.Stmt{
									&parser.Return{
										ExprsOrNil: []parser.Expr{
											&parser.Primary{
												PrimaryKind: parser.ValuePrimary,
												ValueOrNil: &parser.ValueForm{
													ValueKind:   parser.NumberValue,
													NumberOrNil: &zero,
												},
											},
											&parser.Call{
												IsBuilinCall:     true,
												BuiltInKindOrNil: parser.NewErrorBuild,
												ArgsList: []parser.Args{
													{
														&parser.Primary{
															PrimaryKind: parser.ValuePrimary,
															ValueOrNil: &parser.ValueForm{
																ValueKind:   parser.StrLitValue,
																StrLitOrNil: &errMsg,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						&parser.Return{
							ExprsOrNil: []parser.Expr{
								&parser.Binary{
									Op: parser.Div,
									LeftExpr: &parser.Primary{
										PrimaryKind: parser.IdPrimary,
										IdOrNil:     ptrId("a"),
									},
									RightExpr: &parser.Primary{
										PrimaryKind: parser.IdPrimary,
										IdOrNil:     ptrId("b"),
									},
								},
								&parser.Primary{
									PrimaryKind: parser.ValuePrimary,
									ValueOrNil: &parser.ValueForm{
										ValueKind:    parser.ErrValue,
										ErrOrOkOrNil: &nilStr,
									},
								},
							},
						},
					},
				},
			},
		},
	}
	fmt.Println(pkg.String())
}

// Test 4: for range 10 { print("hello") }
func testForLoop() {
	ten := 10
	hello := "hello"

	pkg := &parser.Package{
		DeclsOrNil: []parser.Decl{
			&parser.FuncDecl{
				Id:               "testLoop",
				ParamsOrNil:      []parser.Param{},
				ReturnTypesOrNil: []parser.Type{},
				Block: parser.Block{
					StmtsOrNil: []parser.Stmt{
						&parser.ForRangeAexp{
							Aexp: &parser.Primary{
								PrimaryKind: parser.ValuePrimary,
								ValueOrNil: &parser.ValueForm{
									ValueKind:   parser.NumberValue,
									NumberOrNil: &ten,
								},
							},
							Block: parser.Block{
								StmtsOrNil: []parser.Stmt{
									&parser.CallStmt{
										Call: parser.Call{
											IsBuilinCall:     true,
											BuiltInKindOrNil: parser.PrintBuild,
											ArgsList: []parser.Args{
												{
													&parser.Primary{
														PrimaryKind: parser.ValuePrimary,
														ValueOrNil: &parser.ValueForm{
															ValueKind:   parser.StrLitValue,
															StrLitOrNil: &hello,
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	fmt.Println(pkg.String())
}

// Test 5: Complete package with main and divide functions
func testCompletePackage() {
	four := 4
	two := 2
	zero := 0
	nilStr := "nil"
	errMsg := "can't divide by zero"

	pkg := &parser.Package{
		DeclsOrNil: []parser.Decl{
			// func main() { a, b := 4, 2; divided, err := divide(a,b); if err != nil { panic(err) } }
			&parser.FuncDecl{
				Id:               "main",
				ParamsOrNil:      []parser.Param{},
				ReturnTypesOrNil: []parser.Type{},
				Block: parser.Block{
					StmtsOrNil: []parser.Stmt{
						&parser.ShortDecl{
							Ids: []parser.Id{"a", "b"},
							Exprs: []parser.Expr{
								&parser.Primary{
									PrimaryKind: parser.ValuePrimary,
									ValueOrNil: &parser.ValueForm{
										ValueKind:   parser.NumberValue,
										NumberOrNil: &four,
									},
								},
								&parser.Primary{
									PrimaryKind: parser.ValuePrimary,
									ValueOrNil: &parser.ValueForm{
										ValueKind:   parser.NumberValue,
										NumberOrNil: &two,
									},
								},
							},
						},
						&parser.ShortDecl{
							Ids: []parser.Id{"divided", "err"},
							Exprs: []parser.Expr{
								&parser.Call{
									IsBuilinCall: false,
									PrimaryOrNil: &parser.Primary{
										PrimaryKind: parser.IdPrimary,
										IdOrNil:     ptrId("divide"),
									},
									ArgsList: []parser.Args{
										{
											&parser.Primary{PrimaryKind: parser.IdPrimary, IdOrNil: ptrId("a")},
											&parser.Primary{PrimaryKind: parser.IdPrimary, IdOrNil: ptrId("b")},
										},
									},
								},
							},
						},
						&parser.If{
							Bexp: &parser.Binary{
								Op: parser.NotEqual,
								LeftExpr: &parser.Primary{
									PrimaryKind: parser.IdPrimary,
									IdOrNil:     ptrId("err"),
								},
								RightExpr: &parser.Primary{
									PrimaryKind: parser.ValuePrimary,
									ValueOrNil: &parser.ValueForm{
										ValueKind:    parser.ErrValue,
										ErrOrOkOrNil: &nilStr,
									},
								},
							},
							ThenBlock: parser.Block{
								StmtsOrNil: []parser.Stmt{
									&parser.CallStmt{
										Call: parser.Call{
											IsBuilinCall:     true,
											BuiltInKindOrNil: parser.PanicBuild,
											ArgsList: []parser.Args{
												{
													&parser.Primary{PrimaryKind: parser.IdPrimary, IdOrNil: ptrId("err")},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			// func divide(a int, b int) (int, error) { ... }
			&parser.FuncDecl{
				Id: "divide",
				ParamsOrNil: []parser.Param{
					{Id: "a", Type: parser.Type{TypeKind: parser.IntType}},
					{Id: "b", Type: parser.Type{TypeKind: parser.IntType}},
				},
				ReturnTypesOrNil: []parser.Type{
					{TypeKind: parser.IntType},
					{TypeKind: parser.ErrorType},
				},
				Block: parser.Block{
					StmtsOrNil: []parser.Stmt{
						&parser.If{
							Bexp: &parser.Binary{
								Op: parser.Equal,
								LeftExpr: &parser.Primary{
									PrimaryKind: parser.IdPrimary,
									IdOrNil:     ptrId("b"),
								},
								RightExpr: &parser.Primary{
									PrimaryKind: parser.ValuePrimary,
									ValueOrNil: &parser.ValueForm{
										ValueKind:   parser.NumberValue,
										NumberOrNil: &zero,
									},
								},
							},
							ThenBlock: parser.Block{
								StmtsOrNil: []parser.Stmt{
									&parser.Return{
										ExprsOrNil: []parser.Expr{
											&parser.Primary{
												PrimaryKind: parser.ValuePrimary,
												ValueOrNil: &parser.ValueForm{
													ValueKind:   parser.NumberValue,
													NumberOrNil: &zero,
												},
											},
											&parser.Call{
												IsBuilinCall:     true,
												BuiltInKindOrNil: parser.NewErrorBuild,
												ArgsList: []parser.Args{
													{
														&parser.Primary{
															PrimaryKind: parser.ValuePrimary,
															ValueOrNil: &parser.ValueForm{
																ValueKind:   parser.StrLitValue,
																StrLitOrNil: &errMsg,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						&parser.Return{
							ExprsOrNil: []parser.Expr{
								&parser.Binary{
									Op: parser.Div,
									LeftExpr: &parser.Primary{
										PrimaryKind: parser.IdPrimary,
										IdOrNil:     ptrId("a"),
									},
									RightExpr: &parser.Primary{
										PrimaryKind: parser.IdPrimary,
										IdOrNil:     ptrId("b"),
									},
								},
								&parser.Primary{
									PrimaryKind: parser.ValuePrimary,
									ValueOrNil: &parser.ValueForm{
										ValueKind:    parser.ErrValue,
										ErrOrOkOrNil: &nilStr,
									},
								},
							},
						},
					},
				},
			},
		},
	}
	fmt.Println(pkg.String())
}

// Helper function to create pointer to Id
func ptrId(s string) *parser.Id {
	id := parser.Id(s)
	return &id
}
