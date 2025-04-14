package utils_test

import (
	"api/internal/utils"
	"math"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestProbabilityOfClick(t *testing.T) {
	tests := []struct {
		name     string
		mlScore  decimal.Decimal
		mu       decimal.Decimal
		sigma    decimal.Decimal
		expected decimal.Decimal
	}{
		{
			name:     "Sigma is zero and mlScore equals mu",
			mlScore:  decimal.NewFromFloat(50),
			mu:       decimal.NewFromFloat(50),
			sigma:    decimal.Zero,
			expected: decimal.NewFromFloat(0.5),
		},
		{
			name:     "Sigma is zero and mlScore is greater than mu",
			mlScore:  decimal.NewFromFloat(60),
			mu:       decimal.NewFromFloat(50),
			sigma:    decimal.Zero,
			expected: decimal.NewFromFloat(1.0),
		},
		{
			name:     "Sigma is zero and mlScore is less than mu",
			mlScore:  decimal.NewFromFloat(40),
			mu:       decimal.NewFromFloat(50),
			sigma:    decimal.Zero,
			expected: decimal.NewFromFloat(0.0),
		},
		{
			name:     "Standard sigmoid calculation",
			mlScore:  decimal.NewFromFloat(55),
			mu:       decimal.NewFromFloat(50),
			sigma:    decimal.NewFromFloat(10),
			expected: decimal.NewFromFloat(1 / (1 + math.Exp(-0.5))), // Z = 0.5
		},
		{
			name:     "Negative Z value",
			mlScore:  decimal.NewFromFloat(45),
			mu:       decimal.NewFromFloat(50),
			sigma:    decimal.NewFromFloat(10),
			expected: decimal.NewFromFloat(1 / (1 + math.Exp(0.5))), // Z = -0.5
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.ProbabilityOfClick(tt.mlScore, tt.mu, tt.sigma)
			assert.InDelta(t, tt.expected.InexactFloat64(), result.InexactFloat64(), 0.0001, "expected %v, got %v", tt.expected, result)
		})
	}
}
