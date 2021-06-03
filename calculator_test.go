package gorpn

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateExpression(t *testing.T) {
	ass := assert.New(t)

	tests := []struct {
		expr string
		ans  float64
	}{
		{"22", 22},
		{"2 + 3", 5},
		{"(2 + 2)*2", 8},
		{" 2  + 2 * 2", 6},
		{"( (2 - 4) *-1 * 3 * (2 + 3)) / 5", 6},
		{"55.66 - 55.66 + 22.4+ 2/2 * 110", 132.4},
		{"53 / 2 - 6.5 - 15.55", 4.45},
	}

	for _, test := range tests {
		actualAns, err := CalculateExpression(test.expr)
		ass.NoError(err)
		ass.True(math.Abs(actualAns-test.ans) < 2e-14, "%exp=v act=%v", test.ans, actualAns)
	}
}
