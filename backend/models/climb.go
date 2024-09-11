package models

import (
	"time"

	"gorm.io/gorm"
)

type Climb struct {
	ID               int       `gorm:"primaryKey"`
	UserID           int       `gorm:"not null"`
	SessionID        int       `gorm:"not null"`
	Date             time.Time `gorm:"not null"`
	Attempts         int       `gorm:"not null"            binding:"required"`
	RouteName        string    `gorm:"not null"            binding:"required" json:"route_name"`
	VGrade           int       `gorm:"not null"            binding:"required" json:"v_grade"`
	FailedAttemptSum float64   `gorm:"not null; default:0"                    json:"failed_attempt_sum"`
	Load             float64   `gorm:"not null"`
	CreatedAt        time.Time `gorm:"not null"`
	UpdatedAt        time.Time `gorm:"not null"`
	Tops             int       `gorm:"default:0"`
}

func (c *Climb) CalculateLoad() float64 {
	load := (c.FailedAttemptSum + float64(c.Tops)) * float64(c.VGrade)
	return load
}

func (c *Climb) FindById(db *gorm.DB, userId, climbId int) error {
	return db.Where("user_id = ? AND id = ?", userId, climbId).Take(c).Error
}

func (c *Climb) FindByDate(db *gorm.DB, userId int, climbDate time.Time) error {
	return db.Where("user_id = ? AND date = ?", userId, climbDate).Take(c).Error
}
