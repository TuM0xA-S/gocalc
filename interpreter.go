package gocalc

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/fiorix/go-readline"
)

// Interpreter interprets calculator commands
type Interpreter struct {
	vars        map[string]float64
	funcs       map[string]*function
	interactive bool
	precision   int
	prevLine    *string
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
		vars:        map[string]float64{},
		funcs:       map[string]*function{},
		interactive: verbose,
		precision:   precision,
	}
}

func (ir *Interpreter) completer(input, line string, start, end int) []string {
	if len(input) == 0 {
		return []string{"", "NOTHING TO COMPLETE"}
	}

	names := []string{}
	for k := range ir.funcs {
		k = "@" + k
		names = append(names, k)
	}
	for k := range ir.vars {
		names = append(names, k)
	}

	res := []string{}
	fullmatch := false
	for _, name := range names {
		if name == input {
			fullmatch = true
		}
		if strings.HasPrefix(name, input) {
			res = append(res, name)
		}
	}
	if len(res) == 0 {
		return []string{"", "NO MATCHES"}
	}
	if fullmatch && len(res) > 1 {
		res = append(res, "")
	}
	sort.Strings(res)
	return res
}

// Start interpreting
func (ir *Interpreter) Start(input io.Reader, output io.Writer) error {
	if ir.interactive {
		prompt := ir.printPrompt()
		readline.SetCompletionFunction(ir.completer)
		readline.SetCompleterDelims(" =\t,:()")

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
	} else {
		scn := bufio.NewScanner(input)
		for scn.Scan() {
			if res := ir.ProcessInstruction(scn.Text()); res != "" {
				fmt.Fprintln(output, res)
			}
		}
		return scn.Err()
	}
}

func (ir *Interpreter) printError(err error) string {
	return fmt.Sprintf("error: %v", err)
}

func (ir *Interpreter) printPrompt() string {
	return "eval> "
}

func (ir *Interpreter) printResult(res float64) string {
	if ir.interactive {
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
