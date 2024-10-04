package handlers

import (
	"net/http"
	"strconv"

	"github.com/Alan-Luc/VertiLog/backend/models"
	"github.com/Alan-Luc/VertiLog/backend/services"
	"github.com/Alan-Luc/VertiLog/backend/utils/auth"
	"github.com/Alan-Luc/VertiLog/backend/utils/gContext"
	"github.com/Alan-Luc/VertiLog/backend/utils/params"
	"github.com/gin-gonic/gin"
)

func GetSessionByIDHandler(ctx *gin.Context) {
	var session *models.Session
	var sessionID int
	var userID int
	var page int
	var limit int
	var err error

	// pagination params
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

	sessionID, err = strconv.Atoi(ctx.Param("id"))
	if gContext.HandleReqError(ctx, err, http.StatusBadRequest) {
		return
	}

	userID, err = auth.ExtractUserIdFromJWT(ctx)
	if gContext.HandleReqError(ctx, err, http.StatusUnauthorized) {
		return
	}

	session, err = services.FindSessionByID(userID, sessionID, page, limit)
	if gContext.HandleReqError(ctx, err, http.StatusNotFound) {
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

func GetAllSessionsHandler(ctx *gin.Context) {
	var sessions *[]models.Session
	var userID int
	var page int
	var limit int
	var err error

	// pagination params
	pageParam := ctx.DefaultQuery("page", "1")
	limitParam := ctx.DefaultQuery("limit", "10")

	page, limit, err = params.ValidatePaginationParams(pageParam, limitParam)
	if gContext.HandleReqError(ctx, err, http.StatusBadRequest) {
		return
	}

	userID, err = auth.ExtractUserIdFromJWT(ctx)
	if gContext.HandleReqError(ctx, err, http.StatusUnauthorized) {
		return
	}

	sessions, err = services.FindAllSessionsByUserID(userID, page, limit)
	if gContext.HandleReqError(ctx, err, http.StatusNotFound) {
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

// func GetSessionSummariesByDateHandler(ctx *gin.Context) {
// 	var sessionSummaries *[]models.SessionSummary
// 	var userID int
// 	var startDate time.Time
// 	var endDate time.Time
// 	var err error
//
// }
