package handlers

import (
	"html"
	"net/http"
	"strings"
	"time"

	"github.com/Alan-Luc/VertiLog/backend/database"
	"github.com/Alan-Luc/VertiLog/backend/models"
	"github.com/Alan-Luc/VertiLog/backend/utils/auth"
	"github.com/Alan-Luc/VertiLog/backend/utils/gContext"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func LogClimb(ctx *gin.Context) {
	var climb models.Climb
	var err error

	err = ctx.ShouldBindJSON(&climb)
	if gContext.HandleReqError(ctx, err, http.StatusBadRequest) {
		return
	}

	userId, err := auth.ExtractUserIdFromJWT(ctx)
	if gContext.HandleReqError(ctx, err, http.StatusInternalServerError) {
		return
	}
	climb.UserID = userId

	err = PrepareClimb(&climb)
	if gContext.HandleReqError(ctx, err, http.StatusInternalServerError) {
		return
	}

	err = CreateClimb(&climb)
	if gContext.HandleReqError(ctx, err, http.StatusInternalServerError) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Climb logged succesfully",
		"data":    climb,
	})
}

func CreateClimb(c *models.Climb) error {
	// create climb for requesting user
	if err := database.DB.Create(&c).Error; err != nil {
		return err
	}
	return nil
}

func PrepareClimb(c *models.Climb) error {
	session, err := GetCurrentSession(c.UserID, &c.Date)
	if err != nil {
		return err
	}

	c.SessionID = session.ID
	c.RouteName = html.EscapeString(strings.TrimSpace(c.RouteName))
	c.Load = c.CalculateLoad()

	return nil
}

func GetCurrentSession(userId int, date *time.Time) (*models.Session, error) {
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
		err := tx.
			Where("user_id = ? AND date = ?", userId, sessionDate).
			Take(&session).
			Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				// if there is no current session, create new session
				session = models.Session{
					UserID: userId,
					Date:   sessionDate,
				}
				// create session in transaction
				if err = tx.Create(&session).Error; err != nil {
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
