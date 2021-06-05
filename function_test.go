package gocalc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFunctionDeclaration(t *testing.T) {
	ir := NewInterpreter(false, 0)
	tests := []struct {
		input    string
		name     string
		function *function
	}{
		{
			input: "@a = (): 1",
			name:  "a",
			function: &function{
				params: nil,
				body: []*Token{
					Num(1),
				},
			},
		},
		{
			input: "@xyz = (a, b): 2 * a + b",
			name:  "xyz",
			function: &function{
				params: []string{"a", "b"},
				body: []*Token{
					Num(2), Op("*"), Var("a"), Op("+"), Var("b"),
				},
			},
		},
		{
			input: "@neg = (a): -a",
			name:  "neg",
			function: &function{
				params: []string{"a"},
				body: []*Token{
					UnOp("-"), Var("a"),
				},
			},
		},
	}
	ass := assert.New(t)
	for _, test := range tests {
		ass.Equal("", ir.ProcessInstruction(test.input))
		act := ir.funcs[test.name]
		ass.Equal(test.function, ir.funcs[test.name],
			"exp(params=%v; body=%v) act(params=%v; body=%v)",
			test.function.params, buildExprFromTokens(test.function.body), act.params, buildExprFromTokens(act.body))
	}
}

func TestFunctionCall(t *testing.T) {
	ass := assert.New(t)
	f := &function{
		params: []string{"a", "b"},
		body: []*Token{
			Num(2), Op("*"), Var("a"), Op("+"), Var("b"),
		},
	}

	res, err := f.call([]float64{3, 5})
	ass.NoError(err)
	ass.EqualValues(11, res)
}
