package handlers

import (
	"net/http"
	"strconv"

	"github.com/Alan-Luc/clog/backend/models"
	"github.com/Alan-Luc/clog/backend/services"
	"github.com/Alan-Luc/clog/backend/utils/apiErrors"
	"github.com/Alan-Luc/clog/backend/utils/auth"
	"github.com/Alan-Luc/clog/backend/utils/validators"
	"github.com/gin-gonic/gin"
)

func ClimbLogHandler(ctx *gin.Context) {
	var climb models.Climb
	var userID int
	var err error

	err = ctx.ShouldBindJSON(&climb)
	if apiErrors.HandleAPIError(
		ctx,
		"Invalid input. Please check the submitted data and try again.",
		err,
		http.StatusBadRequest,
	) {
		return
	}

	userID, err = auth.ExtractUserIdFromJWT(ctx)
	if apiErrors.HandleAPIError(
		ctx,
		"Authorization token is invalid or missing. Please log in and try again.",
		err,
		http.StatusUnauthorized,
	) {
		return
	}
	climb.UserID = userID

	err = services.PrepareClimb(&climb)
	if apiErrors.HandleAPIError(
		ctx,
		"An error occurred while processing your request. Please try again later.",
		err,
		http.StatusInternalServerError,
	) {
		return
	}

	err = services.CreateClimb(&climb)
	if apiErrors.HandleAPIError(
		ctx,
		"An error occurred while logging the climb. Please try again later.",
		err,
		http.StatusInternalServerError,
	) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Climb logged succesfully",
		"data":    climb,
	})
}

func ClimbGetByIDHandler(ctx *gin.Context) {
	var climb *models.Climb
	var climbID int
	var userID int
	var err error

	climbID, err = strconv.Atoi(ctx.Param("id"))
	if apiErrors.HandleAPIError(
		ctx,
		"Invalid climb ID. Please ensure the climb ID is a valid number.",
		err,
		http.StatusBadRequest,
	) {
		return
	}

	userID, err = auth.ExtractUserIdFromJWT(ctx)
	if apiErrors.HandleAPIError(
		ctx,
		"Authorization token is invalid or missing. Please log in and try again.",
		err,
		http.StatusUnauthorized,
	) {
		return
	}

	climb, err = services.FindClimbByID(userID, climbID)
	if apiErrors.HandleAPIError(
		ctx,
		"Climb not found. Please check the climb ID and try again.",
		err,
		http.StatusNotFound,
	) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": climb,
	})
}

func ClimbGetAllHandler(ctx *gin.Context) {
	var climbs *[]models.Climb
	var userID int
	var page int
	var limit int
	var err error

	pageParam := ctx.DefaultQuery("page", "1")
	limitParam := ctx.DefaultQuery("limit", "10")

	page, limit, err = validators.ValidatePaginationParams(pageParam, limitParam)
	if apiErrors.HandleAPIError(
		ctx,
		"Invalid pagination parameters. Please provide valid numeric values for page and limit.",
		err,
		http.StatusBadRequest,
	) {
		return
	}

	userID, err = auth.ExtractUserIdFromJWT(ctx)
	if apiErrors.HandleAPIError(
		ctx,
		"Authorization token is invalid or missing. Please log in and try again.",
		err,
		http.StatusUnauthorized,
	) {
		return
	}

	climbs, err = services.FindAllClimbsByUserID(userID, page, limit)
	if apiErrors.HandleAPIError(
		ctx,
		"No climbs found. Please try again later.",
		err,
		http.StatusNotFound,
	) {
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
