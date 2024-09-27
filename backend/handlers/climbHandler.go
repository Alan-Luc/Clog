package handlers

import (
	"net/http"
	"strconv"

	"github.com/Alan-Luc/VertiLog/backend/models"
	"github.com/Alan-Luc/VertiLog/backend/services"
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

	err = services.PrepareClimb(&climb)
	if gContext.HandleReqError(ctx, err, http.StatusInternalServerError) {
		return
	}

	err = services.CreateClimb(&climb)
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

	climb, err = services.FindClimbByID(userID, climbID)
	if gContext.HandleReqError(ctx, err, http.StatusNotFound) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": climb,
	})
}

func GetAllClimbsHandler(ctx *gin.Context) {
	var climbs *[]models.Climb
	var userID int
	var page int
	var limit int
	var err error

	pageParam := ctx.DefaultQuery("page", "1")
	limitParam := ctx.DefaultQuery("limit", "10")

	page, err = strconv.Atoi(pageParam)
	if gContext.HandleReqError(ctx, err, http.StatusBadRequest) {
		return
	}

	limit, err = strconv.Atoi(limitParam)
	if gContext.HandleReqError(ctx, err, http.StatusBadRequest) {
		return
	}

	userID, err = auth.ExtractUserIdFromJWT(ctx)
	if gContext.HandleReqError(ctx, err, http.StatusUnauthorized) {
		return
	}

	climbs, err = services.FindAllClimbsByUserID(userID, page, limit)
	if gContext.HandleReqError(ctx, err, http.StatusNotFound) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": climbs,
		"metadata": map[string]int{
			"page":  page,
			"count": len(*climbs),
		},
	})
}
