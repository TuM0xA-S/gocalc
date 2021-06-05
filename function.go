package gocalc

import "errors"

type function struct {
	params []string
	body   []*Token
}

func (f *function) call(args []float64) (float64, error) {
	if len(f.params) != len(args) {
		return 0, errors.New("wrong argument count")
	}

	vars := map[string]float64{}
	for i := range f.params {
		vars[f.params[i]] = args[i]
	}

	interpreter := &Interpreter{
		vars: vars,
	}

	return interpreter.calculateExpression(f.body)
}

func (ir *Interpreter) processFunctionDeclaration(tokens []*Token) error {
	if len(tokens) < 3 || tokens[0].Type != TokenFunction ||
		tokens[1].Operator != "=" ||
		tokens[2].Operator != "(" {

		return errors.New("not a function declaration")
	}
	function := &function{}
	pos := 3
	ok := false
	if pos < len(tokens) && tokens[pos].Operator == ")" {
		pos++
		ok = true
	}
	for pos+1 < len(tokens) {
		if tokens[pos].Type != TokenVariable {
			break
		}
		function.params = append(function.params, tokens[pos].Variable)
		if tokens[pos+1].Type == TokenDelimiter {
			pos += 2
			continue
		}
		if tokens[pos+1].Operator == ")" {
			pos += 2
			ok = true
			break
		}
		break
	}

	if !ok {
		return errors.New("bad parameter syntax")
	}

	if pos+1 >= len(tokens) || tokens[pos].Delimiter != ":" {
		return errors.New("bad body syntax")
	}
	pos++
	function.body = tokens[pos:]
	for pos < len(tokens) {
		if tokens[pos].Type == TokenFunction {
			return errors.New("subcalls not allowed")
		}
		pos++
	}

	ir.funcs[tokens[0].Function] = function

	return nil
}
