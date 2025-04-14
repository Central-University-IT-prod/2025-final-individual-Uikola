package request

func validateDate(currentDate int, startDate, endDate int) bool {
	return (startDate >= currentDate) && (endDate >= currentDate)
}
