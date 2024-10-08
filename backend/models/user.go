package models

import (
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type User struct {
	ID        int       `gorm:"primaryKey"`
	Username  string    `gorm:"unique"     json:"username"   binding:"required"`
	Password  string    `                  json:"password"   binding:"required"`
	CreatedAt time.Time `                  json:"created_at"`
	UpdatedAt time.Time `                  json:"updated_at"`
}

type PasswordUpdateInput struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password"     binding:"required"`
}

func (u *User) FindByUsername(db *gorm.DB, username string) error {
	err := db.
		Where("username = ?", username).
		Find(u).Error

	if err != nil {
		log.Printf("Database error: %v", err) // Log the error if query fails
		return errors.Wrap(
			err,
			fmt.Sprintf("Failed to find user with username '%s'", username),
		)
	}

	return nil
}

func (u *User) FindByID(db *gorm.DB, userID int) error {
	err := db.
		Where("id = ?", userID).
		Find(u).Error

	if err != nil {
		log.Printf("Database error: %v", err) // Log the error if query fails
		return errors.Wrap(
			err,
			fmt.Sprintf("Failed to find user with id '%d'", userID),
		)
	}

	return nil
}

func (u *User) UpdatePassword(db *gorm.DB, userID int, newPW string) error {
	err := db.
		Model(u).
		Where("id = ?", userID).
		Update("password", newPW).Error

	if err != nil {
		log.Printf("Database error: %v", err) // Log the error if query fails
		return errors.Wrap(
			err,
			fmt.Sprintf("Failed to update password for user with id '%d'", userID),
		)
	}

	return nil
}
