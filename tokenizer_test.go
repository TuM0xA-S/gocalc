package gocalc

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenizer(t *testing.T) {
	ass := assert.New(t)
	type test struct {
		expr     string
		expected []*Token
	}
	tests := []test{
		{
			expr: "2 + 32.2",
			expected: []*Token{
				{
					Type:   TokenNumber,
					Number: 2,
				},
				{
					Type:     TokenOperator,
					Operator: "+",
				},
				{
					Type:   TokenNumber,
					Number: 32.2,
				},
			},
		},
		{
			expr: "-2 * -32.2",
			expected: []*Token{
				{
					Type:   TokenNumber,
					Number: -2,
				},
				{
					Type:     TokenOperator,
					Operator: "*",
				},
				{
					Type:   TokenNumber,
					Number: -32.2,
				},
			},
		},
		{
			expr: "-2*32.2",
			expected: []*Token{
				{
					Type:   TokenNumber,
					Number: -2,
				},
				{
					Type:     TokenOperator,
					Operator: "*",
				},
				{
					Type:   TokenNumber,
					Number: 32.2,
				},
			},
		},
		{
			expr: "  (1 + +2.2)/  0.3",
			expected: []*Token{
				{
					Type:     TokenOperator,
					Operator: "(",
				},
				{
					Type:   TokenNumber,
					Number: 1,
				},
				{
					Type:     TokenOperator,
					Operator: "+",
				},
				{
					Type:   TokenNumber,
					Number: 2.2,
				},
				{
					Type:     TokenOperator,
					Operator: ")",
				},
				{
					Type:     TokenOperator,
					Operator: "/",
				},
				{
					Type:   TokenNumber,
					Number: 0.3,
				},
			},
		},
		{
			expr: "  -92.2    */88-(   ( +1.2 ) -2.1    ",
			expected: []*Token{
				{
					Type:   TokenNumber,
					Number: -92.2,
				},
				{
					Type:     TokenOperator,
					Operator: "*",
				},
				{
					Type:     TokenOperator,
					Operator: "/",
				},
				{
					Type:   TokenNumber,
					Number: 88,
				},
				{
					Type:     TokenOperator,
					Operator: "-",
				},
				{
					Type:     TokenOperator,
					Operator: "(",
				},
				{
					Type:     TokenOperator,
					Operator: "(",
				},
				{
					Type:   TokenNumber,
					Number: 1.2,
				},
				{
					Type:     TokenOperator,
					Operator: ")",
				},
				{
					Type:   TokenNumber,
					Number: -2.1,
				},
			},
		},
		// test with new constructors
		{
			expr: "/1 + -1-(1)",
			expected: []*Token{
				Op("/"), Num(1), Op("+"), Num(-1), Op("-"), Op("("), Num(1), Op(")"),
			},
		},
	}

	for _, test := range tests {
		ok := true
		tokenizer := NewStringTokenizer(test.expr)

		actual := []*Token{}
		for {
			token, err := tokenizer.NextToken()
			if err == EOF {
				break
			}
			ass.NoError(err)
			actual = append(actual, token)
		}
		if !ass.Equal(len(test.expected), len(actual), "fail length: %v", test.expr) {
			continue
		}
		lastTokExp := (*Token)(nil)
		lastTokAct := (*Token)(nil)
		for i := range test.expected {
			ok = false
			tokExp := test.expected[i]
			tokAct := actual[i]
			lastTokExp = tokExp
			lastTokAct = tokAct

			if !ass.Equal(tokExp.Type, tokAct.Type, test.expr) {
				break
			}

			if !ass.True(math.Abs(tokExp.Number-tokAct.Number) < 2e-14, test.expr) {
				break
			}

			if !ass.Equal(tokExp.Operator, tokAct.Operator) {
				break
			}
			ok = true
		}

		if !ok {
			ass.Fail("bad token", "fail on tokens: exp=%#v act=%#v", lastTokExp, lastTokAct)
		}
	}
}
