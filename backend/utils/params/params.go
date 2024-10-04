package params

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

func ValidatePaginationParams(page, limit string) (p, l int, e error) {
	var err error

	p, err = strconv.Atoi(page)
	if err != nil {
		return 0, 0, errors.New("Invalid page parameter!")
	}

	l, err = strconv.Atoi(limit)
	if err != nil {
		return 0, 0, errors.New("Invalid page parameter!")
	}

	if p < 1 {
		return 0, 0, errors.New("Page param must be greater than 1!")
	}

	if l < 1 || l > 10 {
		return 0, 0, errors.New("Limit param must be a number between 1 and 10!")
	}

	return p, l, nil
}

func ValidateDateParams(dateStr string) (time.Time, error) {
	var err error
	var parsedDate time.Time

	layout := "2006-01-02"

	parsedDate, err = time.Parse(layout, dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("Invalid date format, expected YYYY-MM-DD! %v", err)
	}

	truncatedDate := parsedDate.Truncate(24 * time.Hour)

	return truncatedDate, nil
}
