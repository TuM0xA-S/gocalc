package gocalc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInfixToPostfix(t *testing.T) {
	var ir = &Interpreter{}
	type test struct {
		input  []*Token
		output []*Token
	}

	ass := assert.New(t)
	tests := []test{
		{
			input: []*Token{
				Num(1), Op("+"), Num(2),
			},
			output: []*Token{
				Num(1), Num(2), Op("+"),
			},
		},
		{
			input: []*Token{
				Num(1), Op("+"), Num(2), Op("*"), Num(3),
			},
			output: []*Token{
				Num(1), Num(2), Num(3), Op("*"), Op("+"),
			}},
		{
			input: []*Token{
				Op("("), Num(1), Op("+"), Num(2), Op(")"), Op("*"), Num(3),
			},
			output: []*Token{
				Num(1), Num(2), Op("+"), Num(3), Op("*"),
			},
		},
		{
			input: []*Token{
				Num(2), Op("-"), Num(10), Op("-"), Num(5),
			},
			output: []*Token{
				Num(2), Num(10), Op("-"), Num(5), Op("-"),
			},
		},
		{
			input: []*Token{
				Num(2), Op("-"), Op("("), Num(10), Op("-"), Num(5), Op(")"),
			},
			output: []*Token{
				Num(2), Num(10), Num(5), Op("-"), Op("-"),
			},
		},
		{
			input: []*Token{
				Var("x"), Op("-"), Op("("), Num(10), Op("-"), Var("y"), Op(")"),
			},
			output: []*Token{
				Var("x"), Num(10), Var("y"), Op("-"), Op("-"),
			},
		},
		{
			input: []*Token{
				Op("("), Num(2), Op("+"), Num(3), Op(")"),
				Op("*"),
				Op("("), Num(10), Op("-"), Num(5), Op(")"),
			},
			output: []*Token{
				Num(2), Num(3), Op("+"), Num(10), Num(5), Op("-"), Op("*"),
			},
		},
		{
			input: []*Token{
				Op("("),
				Num(5), Op("/"), Op("("), Num(2), Op("+"), Num(3), Op(")"),
				Op(")"), Op("*"), Op("("),
				Op("("), Num(10), Op("-"), Num(5), Op(")"), Op("*"), Num(10),
				Op(")"),
			},
			output: []*Token{
				Num(5), Num(2), Num(3), Op("+"), Op("/"),
				Num(10), Num(5), Op("-"), Num(10), Op("*"),
				Op("*"),
			},
		},
		{
			input: []*Token{
				Op("("), Num(2), Op(")"), Op("*"), Var("abc"), Op("*"), Num(4),
			},
			output: []*Token{
				Num(2), Var("abc"), Op("*"), Num(4), Op("*"),
			},
		},
		{
			input: []*Token{
				UnOp("-"), Num(2),
			},
			output: []*Token{
				Num(2), UnOp("-"),
			},
		},
		{
			input: []*Token{
				UnOp("-"), UnOp("-"), Num(2),
			},
			output: []*Token{
				Num(2), UnOp("-"), UnOp("-"),
			},
		},
		{
			input: []*Token{
				Num(2), Op("-"), UnOp("-"), Num(2),
			},
			output: []*Token{
				Num(2), Num(2), UnOp("-"), Op("-"),
			},
		},
		{
			input: []*Token{
				UnOp("-"), Op("("), UnOp("+"), Num(2), Op(")"),
			},
			output: []*Token{
				Num(2), UnOp("+"), UnOp("-"),
			},
		},
		{
			input: []*Token{
				UnOp("-"), Op("("), Num(3), Op("*"), Num(2), Op(")"),
			},
			output: []*Token{
				Num(3), Num(2), Op("*"), UnOp("-"),
			},
		},
		{
			input: []*Token{
				UnOp("-"), Op("("), UnOp("-"), Num(3), Op("*"), UnOp("-"), Num(2), Op(")"),
			},
			output: []*Token{
				Num(3), UnOp("-"), Num(2), UnOp("-"), Op("*"), UnOp("-"),
			},
		},
		{
			input: []*Token{
				UnOp("+"), UnOp("+"), Op("("), UnOp("+"), UnOp("-"), Num(1), Op(")"),
			},
			output: []*Token{
				Num(1), UnOp("-"), UnOp("+"), UnOp("+"), UnOp("+"),
			},
		},
		{
			input: []*Token{
				Func("a"), Op("("), Num(2), Op(")"),
			},
			output: []*Token{
				Num(2), Func("a"),
			},
		},
		{
			input: []*Token{
				Func("abc"), Op("("), Num(2), Op("+"), Num(2), Delim(","), UnOp("-"), Num(3), Op(")"),
			},
			output: []*Token{
				Num(2), Num(2), Op("+"), Num(3), UnOp("-"), Func("abc"),
			},
		},
		{
			input: []*Token{
				Num(10), Op("*"), Func("abc"), Op("("), Num(2), Op("+"), Num(2), Delim(","),
				UnOp("-"), Num(3), Op("*"), Num(8), Op(")"), Op("+"), Num(22),
			},
			output: []*Token{
				Num(10), Num(2), Num(2), Op("+"), Num(3), UnOp("-"), Num(8), Op("*"), Func("abc"), Op("*"), Num(22), Op("+"),
			},
		},
	}
	for _, test := range tests {
		actualOutput, err := ir.infixToPostfix(test.input)
		expr := buildExprFromTokens(test.input)
		ass.NoError(err, "%#v", expr)
		act := buildExprFromTokens(actualOutput)
		exp := buildExprFromTokens(test.output)
		ass.Equal(test.output, actualOutput, "exp=%s act=%s", exp, act)
	}
}

func TestInfixToPostfixErrors(t *testing.T) {
	var ir = &Interpreter{}
	type test struct {
		input []*Token
	}

	ass := assert.New(t)
	tests := []test{
		{
			input: []*Token{
				Op("("), Num(1), Op("*"), Num(2),
			},
		},
		{
			input: []*Token{
				Num(1), Op("*"), Num(2), Op(")"),
			},
		},
		{
			input: []*Token{
				Num(1), Op("-"), Op("("), Num(2),
			},
		},
		{
			input: []*Token{
				Num(1), Op("-"), Op(")"), Op("("), Num(2), Op(")"),
			},
		},
	}

	for _, test := range tests {
		_, err := ir.infixToPostfix(test.input)
		ass.Error(err)
	}

}
