package nextdate

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const FormatTime = "20060102"

func NextDate(now time.Time, date string, repeat string) (string, error) {

	startDate, err := time.Parse(FormatTime, date)
	if err != nil {
		return "", fmt.Errorf("incorrect date: %w", err)
	}

	if repeat == "" {
		return "", fmt.Errorf("empty repeat rule")
	}

	switch repeat[0] {
	case 'y':
		return annually(now, startDate)
	case 'd':
		return daily(now, startDate, repeat)
	default:
		return "", fmt.Errorf("wrong format: %s", repeat)
	}
}

// Ежегодно
func annually(now time.Time, startDate time.Time) (string, error) {

	nextDate := startDate.AddDate(1, 0, 0)
	for !nextDate.After(now) {
		nextDate = nextDate.AddDate(1, 0, 0)
	}

	return nextDate.Format(FormatTime), nil
}

// Ежедневно
func daily(now, startDate time.Time, repeat string) (string, error) {
	newSlice := strings.Split(repeat, " ")
	if len(newSlice) != 2 {
		return "", fmt.Errorf("wrong format: %s", repeat)
	}
	days, err := strconv.Atoi(newSlice[1])
	if err != nil {
		return "", fmt.Errorf("wrong format: %s", repeat)
	} else if days < 1 || days > 400 {
		return "", fmt.Errorf("maximum interval exceeded: %d", days)
	}
	nextDate := startDate.AddDate(0, 0, days)
	for !nextDate.After(now) {
		nextDate = nextDate.AddDate(0, 0, days)
	}
	return nextDate.Format(FormatTime), nil
}
