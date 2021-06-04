package gocalc

import "errors"

func (ir *Interpreter) processAssignment(tokens []*Token) error {
	if !(len(tokens) >= 2 && tokens[1].Operator == "=") {
		return errors.New("invalid assignment")
	}
	if tokens[0].Type != TokenVariable {
		return errors.New("invalid assignment: no variable on left side")
	}
	varname := tokens[0].Variable
	expr := tokens[2:]
	varval, err := ir.calculateExpression(expr)
	if err != nil {
		return err
	}
	ir.vars[varname] = varval
	return nil
}
