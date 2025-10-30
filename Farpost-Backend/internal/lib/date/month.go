package date

import (
	"fmt"
	"time"
)

func GetAllDatesInMonth(startTime string) ([]string, error) {
	const op = "storage.sqlite.GetDatesMonth"

	startTimeParsed, err := time.Parse("2006-01", startTime)
	if err != nil {
		return nil, fmt.Errorf("%s: invalid time format: %w", op, err)
	}

	firstDayOfMonth := time.Date(startTimeParsed.Year(), startTimeParsed.Month(), 1, 0, 0, 0, 0, startTimeParsed.Location())

	firstDayOfNextMonth := firstDayOfMonth.AddDate(0, 1, 0)
	lastDayOfMonth := firstDayOfNextMonth.AddDate(0, 0, -1)

	daysInMonth := lastDayOfMonth.Day()

	dates := make([]string, 0, daysInMonth)

	for day := 1; day <= daysInMonth; day++ {
		currentDate := time.Date(firstDayOfMonth.Year(), firstDayOfMonth.Month(), day, 0, 0, 0, 0, firstDayOfMonth.Location())
		dates = append(dates, currentDate.Format("2006-01-02"))
	}

	return dates, nil
}