package gocalc

import "errors"

//opPriority = Supported operators with priority
var opPriority = map[string]int{
	"+": 1,
	"-": 1,
	"*": 2,
	"/": 2,
	"(": 3,
	")": 3,
}

// infixToPostfix converts infix notation to reverse polish notation
func (ir *Interpreter) infixToPostfix(input []*Token) ([]*Token, error) {
	output := make([]*Token, 0, len(input))
	stack := []*Token{}
	for _, tok := range input {
		switch tok.Type {
		case TokenNumber, TokenVariable:
			output = append(output, tok)
		case TokenOperator:
			if tok.Operator == "(" {
				stack = append(stack, tok)
				break
			}
			if tok.Operator == ")" {
				op := (*Token)(nil)
				for len(stack) > 0 {
					stack, op = stack[:len(stack)-1], stack[len(stack)-1]
					if op.Operator == "(" {
						break
					}
					output = append(output, op)
				}
				if op == nil || op.Operator != "(" {
					return nil, errors.New("parens not matching")
				}
				break
			}
			for len(stack) > 0 && stack[len(stack)-1].Operator != "(" &&
				opPriority[tok.Operator] <= opPriority[stack[len(stack)-1].Operator] {

				op := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				output = append(output, op)
			}
			stack = append(stack, tok)
		default:
			return nil, errors.New("unknown token type")
		}
	}

	for i := range stack {
		op := stack[len(stack)-1-i]
		if op.Operator == "(" {
			return nil, errors.New("parens not matching")
		}
		output = append(output, op)
	}

	return output, nil
}
