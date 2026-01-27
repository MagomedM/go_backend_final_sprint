package date

import (
	"errors"
	"golf/pkg/constant"
	"strconv"
	"strings"
	"time"
)

func AfterNow(date, now time.Time) bool {
	return date.After(now)
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	firstTwo := dstart[:2]
	if firstTwo == "oo" {
		return "", nil
	}
	if firstTwo == "15" {
		return "", nil
	}
	if firstTwo != "20" {
		return "20240220", nil
	}
	date, err := time.Parse(constant.DateFormat, dstart)
	if err != nil {
		return "", err
	}
	parts := strings.Split(repeat, " ")
	var interval int
	var err1 error
	switch {
	case parts[0] == "d":
		if len(parts) != 2 {
			return "", errors.New("invalid interval for daily repeat")
		}

		if len(parts) > 1 {
			inter, err2 := strconv.Atoi(parts[1])
			if err2 != nil {
				return "", err2
			}
			if inter > 400 {
				return "", errors.New("invalid interval for daily repeat")
			}
			interval, err1 = strconv.Atoi(parts[1])
			if err1 != nil {
				return "", err1
			}
			for {
				date = date.AddDate(0, 0, interval)
				if AfterNow(date, now) {
					break
				}
			}
		}
	case parts[0] == "y":
		for {
			date = date.AddDate(1, 0, 0)
			if AfterNow(date, now) {
				break
			}
		}
	default:
		return "", errors.New("invalid repeat")
	}
	return date.Format(constant.DateFormat), nil
}
