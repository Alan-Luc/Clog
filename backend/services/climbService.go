package services

import (
	"html"
	"strings"
	"time"

	"github.com/Alan-Luc/VertiLog/backend/database"
	"github.com/Alan-Luc/VertiLog/backend/models"
)

func CreateClimb(c *models.Climb) error {
	// create climb for requesting user
	if err := database.DB.Create(&c).Error; err != nil {
		return err
	}
	return nil
}

func PrepareClimb(c *models.Climb) error {
	session, err := FindOrCreateSessionByDate(c.UserID, &c.Date)
	if err != nil {
		return err
	}

	c.SessionID = session.ID
	c.RouteName = html.EscapeString(strings.TrimSpace(c.RouteName))
	c.Load = c.CalculateLoad()

	return nil
}

func FindClimbByID(userID, climbId int) (*models.Climb, error) {
	var climb models.Climb
	var err error

	err = climb.FindById(database.DB, userID, climbId)
	if err != nil {
		return nil, err
	}
	return &climb, nil
}

func FindClimbByDate(userID int, date time.Time) (*models.Climb, error) {
	var climb models.Climb
	var err error

	err = climb.FindByDate(database.DB, userID, date)
	if err != nil {
		return nil, err
	}
	return &climb, nil
}

func FindAllClimbsByUserID(userID, page, limit int) (*[]models.Climb, error) {
	var climb models.Climb
	var climbs []models.Climb
	var err error

	offset := (page - 1) * limit
	climbs, err = climb.FindAll(database.DB, userID, offset, limit)
	if err != nil {
		return nil, err
	}

	return &climbs, nil
}
