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
			&parser.ValDecl{
				Ids:  []parser.Id{"a"},
				Type: parser.Type{TypeKind: parser.IntType},
				ExprsOrNil: []parser.Expr{
					&parser.Atom{
						AtomKind: parser.ValueAtom,
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
									LeftExpr: &parser.Atom{
										AtomKind: parser.IdAtom,
										IdOrNil:  ptrId("a"),
									},
									RightExpr: &parser.Atom{
										AtomKind: parser.IdAtom,
										IdOrNil:  ptrId("b"),
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
								LeftExpr: &parser.Atom{
									AtomKind: parser.IdAtom,
									IdOrNil:  ptrId("b"),
								},
								RightExpr: &parser.Atom{
									AtomKind: parser.ValueAtom,
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
											&parser.Atom{
												AtomKind: parser.ValueAtom,
												ValueOrNil: &parser.ValueForm{
													ValueKind:   parser.NumberValue,
													NumberOrNil: &zero,
												},
											},
											&parser.Atom{
												AtomKind: parser.CallAtom,
												CallOrNil: &parser.Call{
													CallKind:    parser.BuiltInCall,
													BuiltInKind: parser.NewErrorBuild,
													ArgsList: []parser.Args{
														{
															&parser.Atom{
																AtomKind: parser.ValueAtom,
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
						},
						&parser.Return{
							ExprsOrNil: []parser.Expr{
								&parser.Binary{
									Op: parser.Div,
									LeftExpr: &parser.Atom{
										AtomKind: parser.IdAtom,
										IdOrNil:  ptrId("a"),
									},
									RightExpr: &parser.Atom{
										AtomKind: parser.IdAtom,
										IdOrNil:  ptrId("b"),
									},
								},
								&parser.Atom{
									AtomKind: parser.ValueAtom,
									ValueOrNil: &parser.ValueForm{
										ValueKind: parser.NilValue,
										NilOrNil:  &nilStr,
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
				Id:          "testLoop",
				ParamsOrNil: []parser.Param{},
				ReturnTypesOrNil: []parser.Type{},
				Block: parser.Block{
					StmtsOrNil: []parser.Stmt{
						&parser.ForRangeAexp{
							Aexp: &parser.Atom{
								AtomKind: parser.ValueAtom,
								ValueOrNil: &parser.ValueForm{
									ValueKind:   parser.NumberValue,
									NumberOrNil: &ten,
								},
							},
							Block: parser.Block{
								StmtsOrNil: []parser.Stmt{
									&parser.CallStmt{
										Call: parser.Call{
											CallKind:    parser.BuiltInCall,
											BuiltInKind: parser.PrintKBuild,
											ArgsList: []parser.Args{
												{
													&parser.Atom{
														AtomKind: parser.ValueAtom,
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
								&parser.Atom{
									AtomKind: parser.ValueAtom,
									ValueOrNil: &parser.ValueForm{
										ValueKind:   parser.NumberValue,
										NumberOrNil: &four,
									},
								},
								&parser.Atom{
									AtomKind: parser.ValueAtom,
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
								&parser.Atom{
									AtomKind: parser.CallAtom,
									CallOrNil: &parser.Call{
										CallKind: parser.IdCall,
										IdOrNil:  ptrId("divide"),
										ArgsList: []parser.Args{
											{
												&parser.Atom{AtomKind: parser.IdAtom, IdOrNil: ptrId("a")},
												&parser.Atom{AtomKind: parser.IdAtom, IdOrNil: ptrId("b")},
											},
										},
									},
								},
							},
						},
						&parser.If{
							Bexp: &parser.Binary{
								Op: parser.NotEqual,
								LeftExpr: &parser.Atom{
									AtomKind: parser.IdAtom,
									IdOrNil:  ptrId("err"),
								},
								RightExpr: &parser.Atom{
									AtomKind: parser.ValueAtom,
									ValueOrNil: &parser.ValueForm{
										ValueKind: parser.NilValue,
										NilOrNil:  &nilStr,
									},
								},
							},
							ThenBlock: parser.Block{
								StmtsOrNil: []parser.Stmt{
									&parser.CallStmt{
										Call: parser.Call{
											CallKind:    parser.BuiltInCall,
											BuiltInKind: parser.PanicBuild,
											ArgsList: []parser.Args{
												{
													&parser.Atom{AtomKind: parser.IdAtom, IdOrNil: ptrId("err")},
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
								LeftExpr: &parser.Atom{
									AtomKind: parser.IdAtom,
									IdOrNil:  ptrId("b"),
								},
								RightExpr: &parser.Atom{
									AtomKind: parser.ValueAtom,
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
											&parser.Atom{
												AtomKind: parser.ValueAtom,
												ValueOrNil: &parser.ValueForm{
													ValueKind:   parser.NumberValue,
													NumberOrNil: &zero,
												},
											},
											&parser.Atom{
												AtomKind: parser.CallAtom,
												CallOrNil: &parser.Call{
													CallKind:    parser.BuiltInCall,
													BuiltInKind: parser.NewErrorBuild,
													ArgsList: []parser.Args{
														{
															&parser.Atom{
																AtomKind: parser.ValueAtom,
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
						},
						&parser.Return{
							ExprsOrNil: []parser.Expr{
								&parser.Binary{
									Op: parser.Div,
									LeftExpr: &parser.Atom{
										AtomKind: parser.IdAtom,
										IdOrNil:  ptrId("a"),
									},
									RightExpr: &parser.Atom{
										AtomKind: parser.IdAtom,
										IdOrNil:  ptrId("b"),
									},
								},
								&parser.Atom{
									AtomKind: parser.ValueAtom,
									ValueOrNil: &parser.ValueForm{
										ValueKind: parser.NilValue,
										NilOrNil:  &nilStr,
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
