package utils_test

import (
	"api/internal/utils"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestCalculateMinMLScore(t *testing.T) {
	tests := []struct {
		name             string
		mu               decimal.Decimal
		stdDev           decimal.Decimal
		usedImpressions  decimal.Decimal
		impressionsLimit decimal.Decimal
		remainingDays    decimal.Decimal
		totalDays        decimal.Decimal
		z                decimal.Decimal
		k                decimal.Decimal
		d                decimal.Decimal
		expected         decimal.Decimal
	}{
		{
			name:             "Base threshold higher than limit threshold",
			mu:               decimal.NewFromFloat(50),
			stdDev:           decimal.NewFromFloat(10),
			usedImpressions:  decimal.NewFromFloat(100),
			impressionsLimit: decimal.NewFromFloat(200),
			remainingDays:    decimal.NewFromFloat(5),
			totalDays:        decimal.NewFromFloat(10),
			z:                decimal.NewFromFloat(1.0),
			k:                decimal.NewFromFloat(0.5),
			d:                decimal.NewFromFloat(0.2),
			expected:         decimal.NewFromFloat(60), // μ + z * σ = 50 + 1 * 10
		},
		{
			name:             "No impressions used, maximum threshold",
			mu:               decimal.NewFromFloat(50),
			stdDev:           decimal.NewFromFloat(10),
			usedImpressions:  decimal.Zero,
			impressionsLimit: decimal.NewFromFloat(200),
			remainingDays:    decimal.NewFromFloat(10),
			totalDays:        decimal.NewFromFloat(10),
			z:                decimal.NewFromFloat(1.5),
			k:                decimal.NewFromFloat(0.7),
			d:                decimal.NewFromFloat(0.2),
			expected:         decimal.NewFromFloat(65), // μ + z * σ = 50 + 1.5 * 10
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.CalculateMinMLScore(
				tt.mu, tt.stdDev, tt.usedImpressions, tt.impressionsLimit, tt.remainingDays, tt.totalDays, tt.z, tt.k, tt.d,
			)
			assert.True(t, result.Equal(tt.expected), "expected %v, got %v", tt.expected, result)
		})
	}
}
