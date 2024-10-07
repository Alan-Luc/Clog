package services

import (
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/Alan-Luc/VertiLog/backend/database"
	"github.com/Alan-Luc/VertiLog/backend/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func CreateClimb(c *models.Climb) error {
	// create climb for requesting user
	if err := database.DB.Create(&c).Error; err != nil {
		if err.Error() == gorm.ErrDuplicatedKey.Error() {
			return errors.Wrap(
				err,
				"Failed to create climb: a climb with the same ID already exists",
			)
		}
		return errors.Wrap(err, "Failed to create climb:")
	}
	return nil
}

func PrepareClimb(c *models.Climb) error {
	session, err := FindOrCreateSessionByDate(c.UserID, &c.Date)
	if err != nil {
		return errors.Wrap(
			err,
			"Error occurred finding or creating session",
		)
	}

	c.SessionID = session.ID
	c.RouteName = html.EscapeString(strings.TrimSpace(c.RouteName))
	c.Load = c.CalculateLoad()

	return nil
}

func FindClimbByID(userID, climbID int) (*models.Climb, error) {
	var climb models.Climb
	var err error

	err = climb.FindById(database.DB, userID, climbID)
	if err != nil {
		return nil, errors.Wrap(
			err,
			fmt.Sprintf(
				"Error occurred when finding climb with id %d",
				climbID,
			),
		)
	}
	return &climb, nil
}

func FindClimbByDate(userID int, date time.Time) (*models.Climb, error) {
	var climb models.Climb
	var err error

	err = climb.FindByDate(database.DB, userID, date)
	if err != nil {
		return nil, errors.Wrap(
			err,
			fmt.Sprintf(
				"Error occurred when finding climb on date %s",
				date,
			),
		)
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
		return nil, errors.WithMessage(
			err,
			fmt.Sprintf("Error occurred when finding climbs for user with id %d", userID),
		)
	}

	return &climbs, nil
}
