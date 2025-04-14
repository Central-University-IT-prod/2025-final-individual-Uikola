package utils

import (
	"math"

	"github.com/shopspring/decimal"
)

func SigmoidNormalization(x, mu, sigma decimal.Decimal) decimal.Decimal {
	// если σ = 0
	if sigma.IsZero() {
		if x.Equal(mu) {
			return decimal.NewFromFloat(0.5)
		} else if x.GreaterThan(mu) {
			return decimal.NewFromFloat(1.0)
		}
		return decimal.NewFromFloat(0.0)
	}

	// Z = (X - μ) / σ
	z := x.Sub(mu).Div(sigma)

	// Применяем сигмоиду: 1 / (1 + e^(-Z))
	zFloat, _ := z.Float64()
	result := 1 / (1 + math.Exp(-zFloat)) // Сигмоида

	return decimal.NewFromFloat(result)
}
