package gorpn

// CalculateExpression calculates expression in infix notation represented with string
func CalculateExpression(expr string) (float64, error) {
	tokenizer := NewStringTokenizer(expr)
	infixTokens := []*Token{}
	for {
		token, err := tokenizer.NextToken()
		if err == EOF {
			break
		}
		if err != nil {
			return 0, err
		}
		infixTokens = append(infixTokens, token)
	}

	postfixTokens, err := InfixToPostfix(infixTokens)
	if err != nil {
		return 0, err
	}

	res, err := CalculatePostfix(postfixTokens)
	if err != nil {
		return 0, err
	}

	return res, nil
}
