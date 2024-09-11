package handlers

import (
	"html"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Alan-Luc/VertiLog/backend/database"
	"github.com/Alan-Luc/VertiLog/backend/models"
	"github.com/Alan-Luc/VertiLog/backend/utils/auth"
	"github.com/Alan-Luc/VertiLog/backend/utils/gContext"
	"github.com/gin-gonic/gin"
)

func LogClimbHandler(ctx *gin.Context) {
	var climb models.Climb
	var err error

	err = ctx.ShouldBindJSON(&climb)
	if gContext.HandleReqError(ctx, err, http.StatusBadRequest) {
		return
	}

	userID, err := auth.ExtractUserIdFromJWT(ctx)
	if gContext.HandleReqError(ctx, err, http.StatusUnauthorized) {
		return
	}
	climb.UserID = userID

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

func GetClimbByIDHandler(ctx *gin.Context) {
	var climb *models.Climb
	var climbID int
	var userID int
	var err error

	climbID, err = strconv.Atoi(ctx.Param("id"))
	if gContext.HandleReqError(ctx, err, http.StatusBadRequest) {
		return
	}

	userID, err = auth.ExtractUserIdFromJWT(ctx)
	if gContext.HandleReqError(ctx, err, http.StatusUnauthorized) {
		return
	}

	climb, err = FindClimbByID(userID, climbID)
	if gContext.HandleReqError(ctx, err, http.StatusNotFound) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": climb,
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
