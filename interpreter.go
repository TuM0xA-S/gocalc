package gocalc

import (
	"bufio"
	"fmt"
	"io"
)

// Interpreter interprets calculator commands
type Interpreter struct {
	scn       bufio.Scanner
	vars      map[string]float64
	verbose   bool
	precision int
}

// NewInterpreter from input to output
func NewInterpreter(verbose bool, precision int) *Interpreter {
	return &Interpreter{
		vars:      map[string]float64{},
		verbose:   verbose,
		precision: precision,
	}
}

// Start interpreting
func (ir *Interpreter) Start(input io.Reader, output io.Writer) error {
	scn := bufio.NewScanner(input)
	for {
		fmt.Fprint(output, ir.printPrompt())
		if !scn.Scan() {
			break
		}
		if res := ir.ProcessInstruction(scn.Text()); res != "" {
			fmt.Fprintln(output, res)
		}
	}

	return scn.Err()
}

func (ir *Interpreter) printError(err error) string {
	if ir.verbose {
		return fmt.Sprintf("error: %v", err)
	}
	return fmt.Sprintf("%v", err)
}

func (ir *Interpreter) printPrompt() string {
	return "eval> "
}

func (ir *Interpreter) printResult(res float64) string {
	if ir.verbose {
		return fmt.Sprintf("= %.*f", ir.precision, res)
	}
	return fmt.Sprintf("%.*f", ir.precision, res)
}

// ProcessInstruction processes instruction
func (ir *Interpreter) ProcessInstruction(input string) string {
	tokenizer := NewStringTokenizer(input)
	tokens, err := tokenizer.Tokens()
	if len(tokens) >= 2 && tokens[1].Operator == "=" {
		if err := ir.processAssignment(tokens); err != nil {
			return ir.printError(err)
		}
		return ""
	}
	res, err := ir.calculateExpression(tokens)
	if err != nil {
		return ir.printError(err)
	}
	return ir.printResult(res)
}
