package validator

import "strconv"

func ImpressionsLimit(impressionsLimitStr string, _ map[string]interface{}) bool {
	impressionsLimit, err := strconv.Atoi(impressionsLimitStr)
	if err != nil {
		return false
	}

	return impressionsLimit > 0
}

func ClicksLimit(clicksLimitStr string, _ map[string]interface{}) bool {
	clicksLimit, err := strconv.Atoi(clicksLimitStr)
	if err != nil {
		return false
	}

	return clicksLimit > 0
}

func CostPerImpression(costPerImpressionStr string, _ map[string]interface{}) bool {
	costPerImpression, err := strconv.ParseFloat(costPerImpressionStr, 64)
	if err != nil {
		return false
	}

	return costPerImpression > 0
}

func CostPerClick(ostPerClickStr string, _ map[string]interface{}) bool {
	costPerClick, err := strconv.ParseFloat(ostPerClickStr, 64)
	if err != nil {
		return false
	}

	return costPerClick > 0
}

func AdTitle(adTitle string, _ map[string]interface{}) bool {
	return len(adTitle) > 0
}

func AdText(adText string, _ map[string]interface{}) bool {
	return len(adText) > 0
}

func StartDate(startDateStr string, params map[string]interface{}) bool {
	startDate, err := strconv.Atoi(startDateStr)
	if err != nil {
		return false
	}

	currentDate, ok := params["currentDate"].(int)
	if !ok {
		return false
	}

	return startDate >= currentDate
}

func EndDate(endDateStr string, params map[string]interface{}) bool {
	endDate, err := strconv.Atoi(endDateStr)
	if err != nil {
		return false
	}

	startDateStr, ok := params["startDate"].(string)
	if !ok {
		return false
	}
	startDate, _ := strconv.Atoi(startDateStr)

	return endDate > startDate
}

func TargetingGender(gender string, _ map[string]interface{}) bool {
	return gender == "MALE" || gender == "FEMALE" || gender == "ALL"
}

func TargetingAgeFrom(ageFromStr string, _ map[string]interface{}) bool {
	ageFrom, err := strconv.Atoi(ageFromStr)
	if err != nil {
		return false
	}
	return ageFrom >= 0
}

func TargetingAgeTo(ageToStr string, params map[string]interface{}) bool {
	ageTo, err := strconv.Atoi(ageToStr)
	if err != nil {
		return false
	}

	ageFromStr, ok := params["ageFrom"].(string)
	if !ok {
		return false
	}
	ageFrom, _ := strconv.Atoi(ageFromStr)

	return ageTo >= ageFrom
}

func TargetingLocation(location string, _ map[string]interface{}) bool {
	return len(location) > 0
}

func ContextForGeneration(generationStr string, _ map[string]interface{}) bool {
	return len(generationStr) > 0
}
