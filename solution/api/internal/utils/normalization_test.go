package utils_test

import (
	"math"
	"testing"

	"api/internal/utils"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestSigmoidNormalization(t *testing.T) {
	tests := []struct {
		name     string
		x        decimal.Decimal
		mu       decimal.Decimal
		sigma    decimal.Decimal
		expected decimal.Decimal
	}{
		{
			name:     "Sigma is zero and x equals mu",
			x:        decimal.NewFromFloat(50),
			mu:       decimal.NewFromFloat(50),
			sigma:    decimal.Zero,
			expected: decimal.NewFromFloat(0.5),
		},
		{
			name:     "Sigma is zero and x is greater than mu",
			x:        decimal.NewFromFloat(60),
			mu:       decimal.NewFromFloat(50),
			sigma:    decimal.Zero,
			expected: decimal.NewFromFloat(1.0),
		},
		{
			name:     "Sigma is zero and x is less than mu",
			x:        decimal.NewFromFloat(40),
			mu:       decimal.NewFromFloat(50),
			sigma:    decimal.Zero,
			expected: decimal.NewFromFloat(0.0),
		},
		{
			name:     "Standard sigmoid calculation",
			x:        decimal.NewFromFloat(55),
			mu:       decimal.NewFromFloat(50),
			sigma:    decimal.NewFromFloat(10),
			expected: decimal.NewFromFloat(1 / (1 + math.Exp(-0.5))), // Z = 0.5
		},
		{
			name:     "Negative Z value",
			x:        decimal.NewFromFloat(45),
			mu:       decimal.NewFromFloat(50),
			sigma:    decimal.NewFromFloat(10),
			expected: decimal.NewFromFloat(1 / (1 + math.Exp(0.5))), // Z = -0.5
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.SigmoidNormalization(tt.x, tt.mu, tt.sigma)
			assert.InDelta(t, tt.expected.InexactFloat64(), result.InexactFloat64(), 0.0001, "expected %v, got %v", tt.expected, result)
		})
	}
}
