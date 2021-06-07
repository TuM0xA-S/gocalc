package gocalc

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/fiorix/go-readline"
)

// Interpreter interprets calculator commands
type Interpreter struct {
	vars      map[string]float64
	funcs     map[string]*function
	verbose   bool
	precision int
	prevLine  *string
}

type indexedError struct {
	index int
	msg   string
}

func newIndexedError(index int, msg string, args ...interface{}) indexedError {
	return indexedError{index, fmt.Sprintf(msg, args...)}
}

func (ie indexedError) Error() string {
	return fmt.Sprintf("at index %d: %s", ie.index, ie.msg)
}

// NewInterpreter from input to output
func NewInterpreter(verbose bool, precision int) *Interpreter {
	return &Interpreter{
		vars:      map[string]float64{},
		funcs:     map[string]*function{},
		verbose:   verbose,
		precision: precision,
	}
}

// Start interpreting
func (ir *Interpreter) Start(input io.Reader, output io.Writer) error {
	prompt := ir.printPrompt()
	for {
		line := readline.Readline(&prompt)
		switch {
		case line == nil:
			return nil
		case *line != "":
			if ir.prevLine == nil || *ir.prevLine != *line {
				readline.AddHistory(*line)
			}
			if res := ir.ProcessInstruction(*line); res != "" {
				fmt.Fprintln(output, res)
			}
			ir.prevLine = line
		}
	}
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

// ProcessMetaCommand processes meta command
func (ir *Interpreter) ProcessMetaCommand(token *Token) (string, error) {
	if token.Type != TokenMetaCommand {
		return "", fmt.Errorf("not a meta command")
	}
	switch token.Command {
	case "mem":
		buf := &strings.Builder{}
		fmt.Fprintln(buf, "memory:")
		for k, v := range ir.vars {
			fmt.Fprintf(buf, "%s\t= %.*f\n", k, ir.precision, v)
		}
		for k, v := range ir.funcs {
			fmt.Fprintf(buf, "@%s\t= %s\n", k, v)
		}
		return buf.String(), nil
	default:
		return "", newIndexedError(token.Pos, "unknown meta command ;%s", token.Command)
	}

}

// ProcessInstruction processes instruction
func (ir *Interpreter) ProcessInstruction(input string) string {
	tokenizer := NewStringTokenizer(input)
	tokens, err := tokenizer.Tokens()
	if err != nil {
		return ir.printError(err)
	}
	if len(tokens) == 0 {
		return ""
	}
	if tokens[0].Type == TokenMetaCommand {
		res, err := ir.ProcessMetaCommand(tokens[0])
		if err != nil {
			return ir.printError(err)
		}
		return res
	}
	if len(tokens) >= 2 && tokens[1].Operator == "=" {
		var err error
		switch tokens[0].Type {
		case TokenVariable:
			err = ir.processAssignment(tokens)
		case TokenFunction:
			err = ir.processFunctionDeclaration(tokens)
		default:
			err = errors.New("invalid assignment")
		}
		if err != nil {
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
