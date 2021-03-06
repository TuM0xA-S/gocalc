package gocalc

import (
	"fmt"
	"io"
	"strings"
	"unicode"
)

// Token types
const (
	TokenNumber = iota
	TokenOperator
	TokenFunction
	TokenVariable
	TokenDelimiter
	TokenMetaCommand
)

// Token (can be one of Token types)
type Token struct {
	Pos       int // position in input
	Type      int
	Operator  string
	Number    float64
	Variable  string
	Function  string
	Delimiter string
	Command   string
}

func (t *Token) String() string {
	switch t.Type {
	case TokenNumber:
		return fmt.Sprint(t.Number)

	case TokenMetaCommand:
		return ";" + t.Command

	case TokenFunction:
		return "@" + t.Function

	case TokenVariable:
		return t.Variable

	case TokenDelimiter:
		return t.Delimiter

	case TokenOperator:
		if isUnary(t) {
			return t.Operator[1:]
		}
		return t.Operator
	}

	return ""
}

// Op creates operator token
func Op(op string) *Token {
	return &Token{Type: TokenOperator, Operator: op}
}

// UnOp creates unary operator token
func UnOp(op string) *Token {
	return &Token{Type: TokenOperator, Operator: "u" + op}
}

// Num creates number token
func Num(num float64) *Token {
	return &Token{Type: TokenNumber, Number: num}
}

// Var creates variable token
func Var(name string) *Token {
	return &Token{Type: TokenVariable, Variable: name}
}

// Func creates function token
func Func(name string) *Token {
	return &Token{Type: TokenFunction, Function: name}
}

// Meta creates meta command
func Meta(name string) *Token {
	return &Token{Type: TokenMetaCommand, Command: name}
}

// Delim creates delimiter token
func Delim(delim string) *Token {
	return &Token{Type: TokenDelimiter, Delimiter: delim}
}

// EOF error
var EOF = io.EOF

// Tokenizer returns next token or error
type Tokenizer interface {
	NextToken() (*Token, error)
	Tokens() ([]*Token, error)
}

type tokenizer struct {
	data      string
	pos       int
	prevToken *Token
}

// ParseNumber parses float64
// and returns number and cnt - length of number in bytes
// if s does not have number as prefix returns cnt = 0
func ParseNumber(s string) (num float64, pos int) {
	if len(s) == 0 {
		return 0, 0
	}
	numLen := 0
	for pos < len(s) && s[pos] >= '0' && s[pos] <= '9' {
		numLen++
		digit := s[pos]
		num = num*10 + float64(digit-48)
		pos++
	}

	if numLen == 0 {
		return 0, 0
	}

	if pos == len(s) || s[pos] != '.' {
		return num, pos
	}
	pos++

	k := 0.1
	for pos < len(s) && s[pos] >= '0' && s[pos] <= '9' {
		digit := s[pos]
		num = num + k*float64(digit-48)
		k /= 10
		pos++
	}

	return num, pos
}

// ParseIdentifier parses indetifer
func ParseIdentifier(s string) (identifier string, pos int) {
	if len(s) == 0 {
		return "", 0
	}

	if !unicode.IsLetter(rune(s[pos])) {
		return "", 0
	}
	pos++

	for pos < len(s) && (unicode.IsLetter(rune(s[pos])) || unicode.IsDigit(rune(s[pos]))) {
		pos++
	}

	return s[:pos], pos
}

func (t *tokenizer) NextToken() (tok *Token, err error) {
	for t.pos < len(t.data) && t.data[t.pos] == ' ' {
		t.pos++
	}
	if t.pos >= len(t.data) {
		return nil, EOF
	}
	defer func(pos int) {
		if tok == nil {
			return
		}
		tok.Pos = pos
		t.prevToken = tok
	}(t.pos)
	num, cnt := ParseNumber(t.data[t.pos:])
	if cnt > 0 {
		t.pos += cnt
		return Num(num), nil
	}

	op := string(t.data[t.pos])
	if strings.Contains("+-", op) {
		t.pos++
		if t.prevToken != nil && (t.prevToken.Type == TokenVariable ||
			t.prevToken.Type == TokenNumber ||
			t.prevToken.Operator == ")") {

			return Op(op), nil
		}
		return UnOp(op), nil
	}
	if strings.Contains("/*()=", op) {
		t.pos++
		return Op(op), nil
	}
	if op == "," || op == ":" {
		t.pos++
		return Delim(op), nil
	}
	initial := t.pos
	isfunc := false
	ismeta := false
	if op == "@" {
		isfunc = true
		t.pos++
	}
	if op == ";" {
		ismeta = true
		t.pos++
	}
	identifier, cnt := ParseIdentifier(t.data[t.pos:])
	if identifier != "" {
		t.pos += cnt
		if isfunc {
			return Func(identifier), nil
		}
		if ismeta {
			return Meta(identifier), nil
		}
		return Var(identifier), nil
	}
	return nil, indexedError{initial, "bad token"}

}

func (t *tokenizer) Tokens() ([]*Token, error) {
	res := []*Token{}
	for {
		token, err := t.NextToken()
		if err == EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		res = append(res, token)
	}

	return res, nil
}

// NewStringTokenizer returns tokenizer for tokenize data string
func NewStringTokenizer(data string) Tokenizer {
	return &tokenizer{
		data: data,
	}
}

func buildExprFromTokens(tokens []*Token) string {
	buf := &strings.Builder{}
	for _, tok := range tokens {
		if tok.Type == TokenOperator && !isUnary(tok) && !strings.Contains("()", tok.Operator) {
			fmt.Fprintf(buf, " %s ", tok)
			continue
		}

		if tok.Type == TokenDelimiter {
			fmt.Fprintf(buf, "%s ", tok)
			continue
		}

		fmt.Fprint(buf, tok)

	}

	return buf.String()
}
