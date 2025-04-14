package utils

import "github.com/shopspring/decimal"

func CalculateMinMLScore(
	mu, stdDev, usedImpressions, impressionsLimit, remainingDays, totalDays, z, k, d decimal.Decimal,
) decimal.Decimal {
	// 1. Базовый порог: μ + z * σ
	baseThreshold := mu.Add(z.Mul(stdDev))

	// 2. Порог с учётом лимита: μ + k * σ * (1 - usedImpressions/impressionsLimit)
	limitRatio := decimal.NewFromInt(1).Sub(usedImpressions.Div(impressionsLimit))
	limitThreshold := mu.Add(k.Mul(stdDev).Mul(limitRatio))

	// 3. Динамическое снижение по времени: -d * (1 - remainingDays/totalDays)
	timeRatio := decimal.NewFromInt(1).Sub(remainingDays.Div(totalDays))
	dynamicDecrease := d.Mul(timeRatio)

	limitThreshold = limitThreshold.Sub(dynamicDecrease)

	// 4. Итог: max(baseThreshold, limitThreshold)
	if baseThreshold.GreaterThan(limitThreshold) {
		return baseThreshold
	}
	return limitThreshold
}
