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
	if gContext.HandleReqErrorWithStatus(ctx, err, http.StatusBadRequest) {
		return
	}

	userId, err := auth.ExtractUserIdFromJWT(ctx)
	if gContext.HandleReqErrorWithStatus(ctx, err, http.StatusInternalServerError) {
		return
	}

	err = PrepareClimb(userId, &climb)
	if gContext.HandleReqErrorWithStatus(ctx, err, http.StatusInternalServerError) {
		return
	}

	err = CreateClimb(&climb)
	if gContext.HandleReqErrorWithStatus(ctx, err, http.StatusInternalServerError) {
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

func PrepareClimb(userId int, c *models.Climb) error {
	sessionId, err := GetCurrentSessionId(userId)
	if err != nil {
		return err
	}

	// get today's date with time set to 00:00:00
	today := time.Now().UTC().Truncate(24 * time.Hour)

	c.UserID = userId
	c.SessionID = sessionId
	c.RouteName = html.EscapeString(strings.TrimSpace(c.RouteName))
	c.Date = today

	return nil
}

func GetCurrentSessionId(userId int) (int, error) {
	var session models.Session
	var err error

	// get today's date with time set to 00:00:00
	today := time.Now().UTC().Truncate(24 * time.Hour)

	// check if today's session already exists for curr user
	err = database.DB.
		Model(&models.Session{}).
		Where("user_id = ? AND date = ?", userId, today).
		// truncate time
		Take(&session).
		Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// if there is no current session, create today's session
			session = models.Session{
				UserID: userId,
				Date:   today, // Truncate time portion
			}
			err = database.DB.Create(&session).Error
			if err != nil {
				return 0, err
			}
		} else {
			return 0, err
		}
	}

	return session.ID, nil
}
