package utils

import (
	"github.com/shopspring/decimal"
	"math"
)

func ProbabilityOfClick(mlScore, mu, sigma decimal.Decimal) decimal.Decimal {
	// если σ = 0
	if sigma.IsZero() {
		if mlScore.Equal(mu) {
			return decimal.NewFromFloat(0.5)
		} else if mlScore.GreaterThan(mu) {
			return decimal.NewFromFloat(1.0)
		}
		return decimal.NewFromFloat(0.0)
	}

	z := mlScore.Sub(mu).Div(sigma)
	zFloat, _ := z.Float64()
	result := 1 / (1 + math.Exp(-zFloat))
	return decimal.NewFromFloat(result)
}
