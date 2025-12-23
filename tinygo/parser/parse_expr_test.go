package parser_test

import (
	"reflect"
	"testing"

	"github.com/rlaaudgjs5638/langTest/tinygo/lexer"
	"github.com/rlaaudgjs5638/langTest/tinygo/parser"
)

func parseExprFromInput(t *testing.T, input string) parser.Expr {
	t.Helper()
	lx := lexer.NewLexer()
	lx.Set(input)
	p := parser.NewParser(lx)
	expr, err := p.ParseExpr()
	if err != nil {
		t.Fatalf("ParseExpr failed: input=%q err=%v", input, err)
	}
	return expr
}

func intPtr(v int) *int       { return &v }
func boolPtr(v bool) *bool    { return &v }
func strPtr(v string) *string { return &v }

func numPrimary(n int) *parser.Primary {
	return &parser.Primary{
		PrimaryKind: parser.ValuePrimary,
		ValueOrNil: &parser.ValueForm{
			ValueKind:   parser.NumberValue,
			NumberOrNil: intPtr(n),
		},
	}
}

func boolPrimary(b bool) *parser.Primary {
	return &parser.Primary{
		PrimaryKind: parser.ValuePrimary,
		ValueOrNil: &parser.ValueForm{
			ValueKind: parser.BoolValue,
			BoolOrNil: boolPtr(b),
		},
	}
}

func strPrimary(s string) *parser.Primary {
	return &parser.Primary{
		PrimaryKind: parser.ValuePrimary,
		ValueOrNil: &parser.ValueForm{
			ValueKind:   parser.StrLitValue,
			StrLitOrNil: strPtr(s),
		},
	}
}

func idPrimary(s string) *parser.Primary {
	id := parser.Id(s)
	return &parser.Primary{
		PrimaryKind: parser.IdPrimary,
		IdOrNil:     &id,
	}
}

func exprPrimary(expr parser.Expr) *parser.Primary {
	return &parser.Primary{
		PrimaryKind: parser.ExprPrimary,
		ExprOrNil:   expr,
	}
}

func emptyArgs() parser.Args { return parser.Args(nil) }

func TestParseExpr_Table(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  parser.Expr
	}{
		{
			name:  "number",
			input: "42",
			want:  numPrimary(42),
		},
		{
			name:  "string_literal",
			input: "\"hi\"",
			want:  strPrimary("hi"),
		},
		{
			name:  "unary_minus",
			input: "-2",
			want:  &parser.Unary{Op: parser.MinusUnary, Object: numPrimary(2)},
		},
		{
			name:  "arithmetic_precedence",
			input: "1+2*3",
			want: &parser.Binary{
				Op:       parser.Plus,
				LeftExpr: numPrimary(1),
				RightExpr: &parser.Binary{
					Op:        parser.Mul,
					LeftExpr:  numPrimary(2),
					RightExpr: numPrimary(3),
				},
			},
		},
		{
			name:  "paren_grouping",
			input: "(1+2)*3",
			want: &parser.Binary{
				Op:        parser.Mul,
				LeftExpr:  exprPrimary(&parser.Binary{Op: parser.Plus, LeftExpr: numPrimary(1), RightExpr: numPrimary(2)}),
				RightExpr: numPrimary(3),
			},
		},
		{
			name:  "boolean_and_or_precedence",
			input: "true||false&&true",
			want: &parser.Binary{
				Op:       parser.Or,
				LeftExpr: boolPrimary(true),
				RightExpr: &parser.Binary{
					Op:        parser.And,
					LeftExpr:  boolPrimary(false),
					RightExpr: boolPrimary(true),
				},
			},
		},
		{
			name:  "builtin_call",
			input: "print(x)",
			want: &parser.Call{
				IsBuilinCall:     true,
				BuiltInKindOrNil: parser.PrintBuild,
				ArgsList: []parser.Args{
					{idPrimary("x")},
				},
			},
		},
		{
			name:  "call_chain_with_omit",
			input: "f()(x)",
			want: &parser.Call{
				IsBuilinCall:     false,
				PrimaryOrNil:     idPrimary("f"),
				BuiltInKindOrNil: parser.NewErrorBuild,
				ArgsList: []parser.Args{
					emptyArgs(),
					{idPrimary("x")},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseExprFromInput(t, tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("ast mismatch: input=%q\n-- got --\n%s\n-- want --\n%s", tt.input, got.String(), tt.want.String())
			}
			t.Logf("AST:\n%s", got.String())
		})
	}
}
