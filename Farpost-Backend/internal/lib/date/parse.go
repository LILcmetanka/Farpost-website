package date

import (
	"fmt"
	"strings"
	"time"
)

func ParseQueryDate(date string) (string, error) {
	var result string
	if date == "" {
		result = "2019-12-30 23:59:59"
	} else {
		if strings.Contains(date, "T") {
			t, err := time.Parse(time.RFC3339, date)
			if err != nil {
				return "", err
			}
			result = t.Format("2006-01-02 15:04:05")
		} else if strings.Contains(date, "_") {
			result = strings.Replace(date, "_", " ", 1)
		} else {
			return "", fmt.Errorf("failed time parsing")
		}
	}

	return result, nil
}