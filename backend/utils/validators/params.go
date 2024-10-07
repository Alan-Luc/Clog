package validators

import (
	"strconv"
	"time"

	"github.com/pkg/errors"
)

func ValidatePaginationParams(page, limit string) (p, l int, e error) {
	var err error

	p, err = strconv.Atoi(page)
	if err != nil {
		return 0, 0, errors.Wrap(err, "Invalid page parameter!")
	}

	l, err = strconv.Atoi(limit)
	if err != nil {
		return 0, 0, errors.Wrap(err, "Invalid page parameter!")
	}

	if p < 1 {
		return 0, 0, errors.New("Page param must be greater than 1!")
	}

	if l < 1 || l > 10 {
		return 0, 0, errors.New("Limit param must be a number between 1 and 10!")
	}

	return p, l, nil
}

func ValidateAndParseDateParams(dateStr string) (time.Time, error) {
	var err error
	var parsedDate time.Time

	layout := "2006-01-02"

	parsedDate, err = time.Parse(layout, dateStr)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "Invalid date format, expected YYYY-MM-DD! %v")
	}

	truncatedDate := parsedDate.Truncate(24 * time.Hour)

	return truncatedDate, nil
}

func ValidateAndParseDateSpanParams(startDateStr, endDateStr string) (time.Time, time.Time, error) {
	var err error
	var parsedStartDate time.Time
	var parsedEndDate time.Time

	layout := "2006-01-02"

	parsedStartDate, err = time.Parse(layout, startDateStr)
	if err != nil {
		return time.Time{}, time.Time{}, errors.Wrap(
			err,
			"Invalid date format, expected YYYY-MM-DD! %v",
		)
	}

	parsedEndDate, err = time.Parse(layout, endDateStr)
	if err != nil {
		return time.Time{}, time.Time{}, errors.Wrap(
			err,
			"Invalid date format, expected YYYY-MM-DD! %v",
		)
	}

	startDate := parsedStartDate.Truncate(24 * time.Hour)
	endDate := parsedEndDate.Truncate(24 * time.Hour)

	if endDate.Before(startDate) || startDate.Equal(endDate) {
		return time.Time{}, time.Time{}, errors.New(
			"Start date must be at an earlier time than end date!",
		)
	}

	return startDate, endDate, nil
}
