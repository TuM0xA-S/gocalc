package gocalc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculatePostfix(t *testing.T) {
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
	}

	ass := assert.New(t)
	for _, test := range tests {
		actualAnswer, err := CalculatePostfix(test.input)
		ass.NoError(err)
		ass.Equal(test.answer, actualAnswer)
	}
}
