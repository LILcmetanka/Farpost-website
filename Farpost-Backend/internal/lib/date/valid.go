package date

import (
    "time"
)

func IsValidDate(dateStr string) bool {
    _, err := time.Parse("2006-01-02", dateStr)
    if err != nil {
        return false
    }
    return true
}