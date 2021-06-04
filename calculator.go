package gocalc

// calculateExpression calculates expression in infix notation represented with string
func (ir *Interpreter) calculateExpression(tokens []*Token) (float64, error) {
	postfixTokens, err := ir.infixToPostfix(tokens)
	if err != nil {
		return 0, err
	}

	res, err := ir.calculatePostfix(postfixTokens)
	if err != nil {
		return 0, err
	}

	return res, nil
}
