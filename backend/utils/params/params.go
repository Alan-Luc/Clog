package params

import (
	"errors"
	"strconv"
)

func ValidPaginationParams(page, limit string) (p, l int, e error) {
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
