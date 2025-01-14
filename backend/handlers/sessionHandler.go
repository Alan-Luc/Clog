package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Alan-Luc/clog/backend/models"
	"github.com/Alan-Luc/clog/backend/services"
	"github.com/Alan-Luc/clog/backend/utils/apiErrors"
	"github.com/Alan-Luc/clog/backend/utils/auth"
	"github.com/Alan-Luc/clog/backend/utils/validators"
	"github.com/gin-gonic/gin"
)

func SessionLogHandler(ctx *gin.Context) {
	var session models.Session
	var userID int
	var err error

	err = ctx.ShouldBindJSON(&session)
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
	session.UserID = userID
}

func SessionGetByIDHandler(ctx *gin.Context) {
	var session *models.Session
	var sessionID int
	var userID int
	var page int
	var limit int
	var err error

	// pagination validators
	pageParam := ctx.DefaultQuery("page", "1")
	limitParam := ctx.DefaultQuery("limit", "10")

	// pagination for climbs within a session
	page, limit, err = validators.ValidatePaginationParams(pageParam, limitParam)
	if apiErrors.HandleAPIError(
		ctx,
		"Invalid pagination parameters. Please provide valid numeric values for page and limit.",
		err,
		http.StatusBadRequest,
	) {
		return
	}

	sessionID, err = strconv.Atoi(ctx.Param("id"))
	if apiErrors.HandleAPIError(
		ctx,
		"Invalid session ID. Please ensure the session ID is a valid number.",
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

	session, err = services.FindSessionByID(userID, sessionID, page, limit)
	if apiErrors.HandleAPIError(
		ctx,
		"Session not found. Please check the session ID and try again.",
		err,
		http.StatusNotFound,
	) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": session,
		"metadata": map[string]int{
			"page":  page,
			"count": len(session.Climbs),
		},
	})
}

func SessionGetAllHandler(ctx *gin.Context) {
	var sessions *[]models.Session
	var userID int
	var page int
	var limit int
	var err error

	// pagination validators
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

	sessions, err = services.FindAllSessionsByUserID(userID, page, limit)
	if apiErrors.HandleAPIError(
		ctx,
		"No sessions found. Please try again later.",
		err,
		http.StatusNotFound,
	) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": sessions,
		"metadata": map[string]int{
			"page":  page,
			"count": len(*sessions),
		},
	})
}

func SessionGetSummariesByDateHandler(ctx *gin.Context) {
	var sessionSummaries *[]models.SessionSummary
	var userID int
	var startDate time.Time
	var endDate time.Time
	var err error

	startDateStr := ctx.Query("startDate")
	endDateStr := ctx.Query("endDate")

	startDate, endDate, err = validators.ValidateAndParseDateSpanParams(startDateStr, endDateStr)
	if apiErrors.HandleAPIError(
		ctx,
		"Invalid date range. Please ensure that the start date and end date are valid.",
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

	sessionSummaries, err = services.FindSessionsSummariesByDate(userID, startDate, endDate)
	if apiErrors.HandleAPIError(
		ctx,
		"No session summaries found for the given date range.",
		err,
		http.StatusNotFound,
	) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": sessionSummaries,
		"metadata": map[string]any{
			"count":     len(*sessionSummaries),
			"startDate": startDate,
			"endDate":   endDate,
		},
	})
}
