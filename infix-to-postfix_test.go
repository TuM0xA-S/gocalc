package gorpn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInfixToPostfix(t *testing.T) {
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
			},
		},
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
				Op("("), Num(2), Op(")"), Op("*"), Num(3), Op("*"), Num(4),
			},
			output: []*Token{
				Num(2), Num(3), Op("*"), Num(4), Op("*"),
			},
		},
	}
	for _, test := range tests {
		actualOutput, err := InfixToPostfix(test.input)
		ass.NoError(err)
		ass.Equal(test.output, actualOutput)
	}
}

func TestInfixToPostfixErrors(t *testing.T) {
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
		_, err := InfixToPostfix(test.input)
		ass.Error(err)
	}

}
