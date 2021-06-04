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
	TokenVariable
)

// Token (can be one of Token types)
type Token struct {
	Type     int
	Operator string
	Number   float64
	Variable string
}

// Op creates operator token
func Op(op string) *Token {
	return &Token{Type: TokenOperator, Operator: op}
}

// Num creates number token
func Num(num float64) *Token {
	return &Token{Type: TokenNumber, Number: num}
}

// Var creates variable token
func Var(name string) *Token {
	return &Token{Type: TokenVariable, Variable: name}
}

// EOF error
var EOF = io.EOF

// Tokenizer returns next token or error
type Tokenizer interface {
	NextToken() (*Token, error)
	Tokens() ([]*Token, error)
}

type tokenizer struct {
	data string
	pos  int
}

// ParseNumber parses float64
// and returns number and cnt - length of number in bytes
// if s does not have number as prefix returns cnt = 0
func ParseNumber(s string) (num float64, pos int) {
	if len(s) == 0 {
		return 0, 0
	}
	sign := float64(1)
	if s[0] == '+' {
		pos++
	} else if s[0] == '-' {
		pos++
		sign = -1
	}
	numLen := 0
	for pos < len(s) && s[pos] >= '0' && s[pos] <= '9' {
		numLen++
		digit := s[pos]
		num = num*10 + sign*float64(digit-48)
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
		num = num + k*sign*float64(digit-48)
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

func (t *tokenizer) NextToken() (*Token, error) {
	for t.pos < len(t.data) && t.data[t.pos] == ' ' {
		t.pos++
	}
	if t.pos >= len(t.data) {
		return nil, EOF
	}
	num, cnt := ParseNumber(t.data[t.pos:])
	if cnt > 0 {
		t.pos += cnt
		return &Token{Type: TokenNumber, Number: num}, nil
	}

	op := string(t.data[t.pos])
	if strings.Contains("+-/*()=", op) {
		t.pos++
		return &Token{Type: TokenOperator, Operator: op}, nil
	}

	identifier, cnt := ParseIdentifier(t.data[t.pos:])
	if identifier != "" {
		t.pos += cnt
		return &Token{Type: TokenVariable, Variable: identifier}, nil
	}
	return nil, fmt.Errorf("bad token at %v", t.pos)

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
