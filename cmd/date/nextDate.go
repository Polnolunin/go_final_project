package date

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	Year              = "y"
	Day               = "d"
	DateFormat string = "20060102"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {

	repeatData := strings.Split(repeat, " ")

	if len(repeatData) == 0 {
		return "", fmt.Errorf("правило повторения не указано")
	}

	parsedDate, err := time.Parse(DateFormat, date)

	if err != nil {
		return "", fmt.Errorf("неверный формат даты: %v", err)
	}

	switch repeatData[0] {

	case Day:
		if len(repeatData) != 2 {
			return "", fmt.Errorf("неверное количество дней")
		}

		days, _ := strconv.Atoi(repeatData[1])

		if days > 400 || days < 0 {
			return "", fmt.Errorf("неверное количество дней")
		}
		for {
			parsedDate = parsedDate.AddDate(0, 0, days)
			if parsedDate.After(now) {
				break
			}
		}
	case Year:
		for {
			parsedDate = parsedDate.AddDate(1, 0, 0)
			if parsedDate.After(now) {
				break
			}
		}

	default:
		return "", fmt.Errorf("неверный формат повторения")
	}

	return parsedDate.Format(DateFormat), nil

}
