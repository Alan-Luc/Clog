package models

import "time"

type Session struct {
	ID        int       `gorm:"primaryKey"`
	Date      time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
	UserID    int       `gorm:"not null"`
}
