package resolver

import "testing"

func TestResolveHoist_SuccessForwardRefs(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{
			name:  "global_var_forward_ref",
			input: "var a int = b; var b int = 1;",
		},
		{
			name:  "global_func_forward_ref",
			input: "func g(){ f(); } func f(){ }",
		},
		{
			name:  "mutual_func_forward_ref",
			input: "func f(){ g(); } func g(){ f(); }",
		},
		{
			name:  "var_uses_func_forward_ref",
			input: "var a int = f(); func f(){ return 1; }",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, _, err := resolveFromInput(t, tc.input)
			if err != nil {
				t.Fatalf("unexpected resolve error: %v", err)
			}
		})
	}
}
