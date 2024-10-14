package models

import (
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
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
	Notes            string    `                                              json:"notes"`
}

func (c *Climb) CalculateLoad() float64 {
	return (c.FailedAttemptSum + float64(c.Tops)) * float64(c.VGrade)
}

func (c *Climb) FindAll(db *gorm.DB, userID, offset, limit int) ([]Climb, error) {
	var climbs []Climb

	err := db.
		Where("user_id = ?", userID).
		Offset(offset).
		Limit(limit).
		Find(&climbs).Error

	if err != nil {
		log.Printf("Database error: %v", err) // Log the error if query fails
		return nil, errors.Wrap(
			err,
			fmt.Sprintf("Error finding climbs for user with id %d", userID),
		)
	}

	return climbs, nil
}

func (c *Climb) FindByID(db *gorm.DB, userID, climbID int) error {
	err := db.
		Where("user_id = ? AND id = ?", userID, climbID).
		Take(c).Error

	if err != nil {
		log.Printf("Database error: %v", err) // Log the error if query fails
		return errors.Wrap(
			err,
			fmt.Sprintf("Error finding climb with id %d for user with id %d", climbID, userID),
		)
	}

	return nil
}

func (c *Climb) FindByDate(db *gorm.DB, userID int, climbDate time.Time) error {
	err := db.
		Where("user_id = ? AND date = ?", userID, climbDate).
		Take(c).Error

	if err != nil {
		log.Printf("Database error: %v", err) // Log the error if query fails
		return errors.Wrap(
			err,
			fmt.Sprintf("Error finding climb on date %s for user with id %d", climbDate, userID),
		)
	}

	return nil
}
