package utils

func (now *Now) GetCurrentWeeks() (int, int) {
	var startOfMonth = now.BeginningOfMonth()
	var _, startWeek = startOfMonth.ISOWeek()
	var _, endWeek = now.ISOWeek()
	return startWeek, endWeek
}

func (now *Now) GetCurrentDay() (int, int) {
	var startOfWeek = now.BeginningOfWeek()
	var starDay = startOfWeek.Day()
	var endDay = now.Day()
	return starDay, endDay
}

func (now *Now) GetCurrentMonth() (int, int) {
	var startOfYear = now.BeginningOfYear()
	var startMonth = startOfYear.Month()
	var endMonth = now.Month()
	return int(startMonth), int(endMonth)
}
