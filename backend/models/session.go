package models

import "time"

type Session struct {
	ID        int
	Date      time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
