package gocalc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssignment(t *testing.T) {
	ir := NewInterpreter(false, 0)
	tests := []struct {
		input, answer string
	}{
		{"a = 2 + 2", "4"},
		{"b = a * 10", "40"},
		{"c = (a + b) / 2", "22"},
		{"b = a - c", "-18"},
	}
	ass := assert.New(t)
	for _, test := range tests {
		ass.Equal("", ir.ProcessInstruction(test.input))
		ass.Equal(test.answer, ir.ProcessInstruction(test.input[:1]))
	}
}
