package services

import (
	"time"

	"github.com/Alan-Luc/VertiLog/backend/database"
	"github.com/Alan-Luc/VertiLog/backend/models"
	"gorm.io/gorm"
)

func FindAllSessionsByUserID(userID, page, limit int) (*[]models.Session, error) {
	var sessions *[]models.Session
	var err error

	// TODO: cache page limit

	offset := (page - 1) * limit
	sessions, err = models.FindAllSessions(database.DB, userID, offset, limit)
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

func FindSessionByID(userID, sessionID, page, limit int) (*models.Session, error) {
	var session models.Session
	var err error

	// TODO: cache page limit

	offset := (page - 1) * limit
	err = session.FindById(database.DB, userID, sessionID, offset, limit)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func FindOrCreateSessionByDate(userID int, date *time.Time) (*models.Session, error) {
	var session models.Session
	var err error
	var sessionDate time.Time

	// get today's date with time set to 00:00:00
	// if there is no current session, create new session
	today := time.Now().UTC().Truncate(24 * time.Hour)
	if date == nil {
		sessionDate = today
	} else {
		sessionDate = *date
	}

	// check if today's session already exists for curr user
	// transaction to avoid race conditions
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		err = session.FindByDate(database.DB, userID, sessionDate)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				// if there is no current session, create new session
				session = models.Session{
					UserID: userID,
					Date:   sessionDate,
				}
				// create session in transaction
				err = tx.Create(&session).Error
				if err != nil {
					return err // rollback transaction if error
				}
			} else {
				return err
			}
		}
		// if no errors in transaction
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &session, nil
}

func FindSessionsSummariesByDate(
	userID int,
	startDate, endDate time.Time,
) (*[]models.SessionSummary, error) {
	var sessionSummaries *[]models.SessionSummary
	var err error

	sessionSummaries, err = models.FindSessionSummaries(database.DB, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	return sessionSummaries, nil
}
