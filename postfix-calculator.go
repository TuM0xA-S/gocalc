package gorpn

import "errors"

// CalculatePostfix calculates expression in postfix notation
func CalculatePostfix(input []*Token) (float64, error) {
	if len(input) == 0 {
		return 0, errors.New("nothing to calculate")
	}

	stack := []float64{}
	for _, tok := range input {
		if tok.Type == TokenNumber {
			stack = append(stack, tok.Number)
			continue
		}

		if tok.Type != TokenOperator {
			panic("unknown token type")
		}

		if len(stack) < 2 {
			return 0, errors.New("not enough operands")
		}
		b := stack[len(stack)-1]
		a := stack[len(stack)-2]
		stack = stack[:len(stack)-2]
		switch tok.Operator {
		case "+":
			stack = append(stack, a+b)
		case "-":
			stack = append(stack, a-b)
		case "*":
			stack = append(stack, a*b)
		case "/":
			stack = append(stack, a/b)
		default:
			panic("unknown operator")
		}
	}

	if len(stack) > 1 {
		return 0, errors.New("not enough operators")
	}

	return stack[0], nil
}
