package models

import "time"

type Climb struct {
	ID                int
	UserID            int
	SessionID         int
	Date              time.Time
	RouteName         string
	VGrade            int
	Attempts          int
	CompletionPercent float32
	Load              float64
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (c *Climb) CalculateLoad() float64 {
	return 0
}
