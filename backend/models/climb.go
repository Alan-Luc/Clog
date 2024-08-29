package models

import "time"

type Climb struct {
	ID                int       `gorm:"primaryKey"`
	UserID            int       `gorm:"not null"`
	SessionID         int       `gorm:"not null"`
	Date              time.Time `gorm:"not null"`
	Attempts          int       `gorm:"not null"`
	RouteName         string    `json:"route_name"`
	VGrade            int       `json:"v_grade"`
	CompletionPercent float32   `json:"completion_percent"`
	Load              float64
	CreatedAt         time.Time `gorm:"not null"`
	UpdatedAt         time.Time `gorm:"not null"`
}

func (c *Climb) CalculateLoad() float64 {
	return 0
}
