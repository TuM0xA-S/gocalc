package gorpn

import (
	"fmt"
	"io"
	"strings"
)

// Token types
const (
	TokenNumber = iota
	TokenOperator
)

// Token (can be one of Token types)
type Token struct {
	Type     int
	Operator string
	Number   float64
}

// Op creates operator token
func Op(op string) *Token {
	return &Token{Type: TokenOperator, Operator: op}
}

// Num creates number token
func Num(num float64) *Token {
	return &Token{Type: TokenNumber, Number: num}
}

// EOF error
var EOF = io.EOF

// Tokenizer returns next token or error
type Tokenizer interface {
	NextToken() (*Token, error)
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
	if strings.Contains("+-/*()", op) {
		t.pos++
		return &Token{Type: TokenOperator, Operator: op}, nil
	}

	return nil, fmt.Errorf("bad token at %v", t.pos)

}

// NewStringTokenizer returns tokenizer for tokenize data string
func NewStringTokenizer(data string) Tokenizer {
	return &tokenizer{
		data: data,
	}
}
