package gocalc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculatePostfix(t *testing.T) {
	var ir = &Interpreter{
		vars: map[string]float64{
			"x": 5,
			"y": 10.5,
		},
		funcs: map[string]*function{
			"foo": {
				params: []string{"a"},
				body:   []*Token{UnOp("-"), Var("a"), Op("*"), Num(2)},
			},
		},
	}
	type test struct {
		input  []*Token
		answer float64
	}

	tests := []test{
		{
			input: []*Token{
				Num(2), Num(2), Op("*"),
			},
			answer: 4,
		},
		{
			input: []*Token{
				Num(2), Num(5), Op("+"), Num(3), Op("*"),
			},
			answer: 21,
		},
		{
			input: []*Token{
				Num(2),
			},
			answer: 2,
		},
		{
			input: []*Token{
				Num(2), Num(3), Num(10), Op("-"), Op("-"),
			},
			answer: 9,
		},
		{
			input: []*Token{
				Num(10), Num(2), Op("/"), Num(2), Num(3), Op("+"), Op("*"), Num(25), Num(5), Op("/"), Op("*"),
			},
			answer: 125,
		},
		{
			input: []*Token{
				Var("y"), Num(0.5), Op("-"), Num(2), Op("/"), Num(2), Num(3), Op("+"), Op("*"), Num(25), Var("x"), Op("/"), Op("*"),
			},
			answer: 125,
		},
		{
			input: []*Token{
				Var("x"), UnOp("-"),
			},
			answer: -5,
		},
		{
			input: []*Token{
				Num(3), Num(5), Op("-"), UnOp("-"),
			},
			answer: 2,
		},
		{
			input: []*Token{
				Num(3), Num(5), Op("-"), UnOp("-"), UnOp("+"), Num(-2), Op("/"),
			},
			answer: -1,
		},
		{
			input: []*Token{
				Num(3), Func("foo"),
			},
			answer: -6,
		},
		{
			input: []*Token{
				Num(3), Num(2), Op("+"), Func("foo"),
			},
			answer: -10,
		},
		{
			input: []*Token{
				Num(3), Num(2), Op("+"), Num(5), Op("/"), Func("foo"), UnOp("-"),
			},
			answer: 2,
		},
	}

	ass := assert.New(t)
	for _, test := range tests {
		actualAnswer, err := ir.calculatePostfix(test.input)
		ass.NoError(err)
		ass.Equal(test.answer, actualAnswer)
	}
}
