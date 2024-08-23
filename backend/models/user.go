package models

import "time"

type User struct {
	ID        int       `gorm:"primaryKey"`
	Username  string    `gorm:"unique"     json:"username"   binding:"required"`
	Password  string    `                  json:"password"   binding:"required"`
	CreatedAt time.Time `                  json:"created_at"`
	UpdatedAt time.Time `                  json:"updated_at"`
}
