package handlers

import (
	"html"
	"net/http"
	"strings"

	"github.com/Alan-Luc/VertiLog/backend/database"
	"github.com/Alan-Luc/VertiLog/backend/models"
	"github.com/Alan-Luc/VertiLog/backend/utils/auth"
	"github.com/Alan-Luc/VertiLog/backend/utils/gContext"
	"github.com/gin-gonic/gin"
)

func LogClimb(ctx *gin.Context) {
	var climb models.Climb
	var err error

	err = ctx.ShouldBindJSON(&climb)
	if gContext.HandleReqError(ctx, err, http.StatusBadRequest) {
		return
	}

	userId, err := auth.ExtractUserIdFromJWT(ctx)
	if gContext.HandleReqError(ctx, err, http.StatusUnauthorized) {
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

// helpers
func CreateClimb(c *models.Climb) error {
	// create climb for requesting user
	if err := database.DB.Create(&c).Error; err != nil {
		return err
	}
	return nil
}

func PrepareClimb(c *models.Climb) error {
	session, err := GetSessionByDate(c.UserID, &c.Date)
	if err != nil {
		return err
	}

	c.SessionID = session.ID
	c.RouteName = html.EscapeString(strings.TrimSpace(c.RouteName))
	c.Load = c.CalculateLoad()

	return nil
}
