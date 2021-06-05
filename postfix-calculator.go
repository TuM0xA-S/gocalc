package gocalc

import (
	"errors"
	"fmt"
)

// calculatePostfix calculates expression in postfix notation
func (ir *Interpreter) calculatePostfix(input []*Token) (float64, error) {
	if len(input) == 0 {
		return 0, errors.New("nothing to calculate")
	}

	stack := []float64{}
	for _, tok := range input {
		if tok.Type == TokenNumber {
			stack = append(stack, tok.Number)
			continue
		}

		if tok.Type == TokenVariable {
			val, ok := ir.vars[tok.Variable]
			if !ok {
				return 0, fmt.Errorf("unknow variable: %v", tok.Variable)
			}
			stack = append(stack, val)
			continue
		}
		if isUnary(tok) {
			if len(stack) < 1 {
				return 0, errors.New("not enough operands")
			}
			if tok.Operator == "u-" {
				stack[len(stack)-1] *= -1
			}
			continue
		}
		if tok.Type == TokenFunction {
			name := tok.Function
			fn, ok := ir.funcs[name]
			if !ok {
				return 0, errors.New("unknown function")
			}

			if len(stack) < len(fn.params) {
				return 0, errors.New("not enougn params to call function")
			}

			args := stack[len(stack)-len(fn.params):]
			stack = stack[:len(stack)-len(fn.params)]
			res, err := fn.call(args)
			if err != nil {
				return 0, err
			}
			stack = append(stack, res)
			continue
		}
		if tok.Type != TokenOperator {
			return 0, errors.New("unknown token type")
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
			return 0, errors.New("unknown operator")
		}
	}

	if len(stack) > 1 {
		return 0, errors.New("not enough operators")
	}

	return stack[0], nil
}
