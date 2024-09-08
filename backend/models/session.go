package models

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	ID        int       `gorm:"primaryKey"`
	Date      time.Time `gorm:"not null;index"       binding:"required"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
	UserID    int       `gorm:"not null;index"`
	Climbs    []Climb   `gorm:"foreignKey:SessionID"` // not a column, is a relationship
}

func (s *Session) FindById(db *gorm.DB, userId, sessionId int) error {
	return db.Preload("Climbs").Where("user_id = ? AND id = ?", userId, sessionId).Take(s).Error
}

func (s *Session) FindByDate(db *gorm.DB, userId int, sessionDate time.Time) error {
	return db.Preload("Climbs").Where("user_id = ? AND date = ?", userId, sessionDate).Take(s).Error
}
