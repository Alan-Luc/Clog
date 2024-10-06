package handlers

import (
	"net/http"

	"github.com/Alan-Luc/VertiLog/backend/models"
	"github.com/Alan-Luc/VertiLog/backend/services"
	"github.com/Alan-Luc/VertiLog/backend/utils/gContext"
	"github.com/gin-gonic/gin"
)

func RegisterUserHandler(ctx *gin.Context) {
	var user models.User
	var err error

	err = ctx.ShouldBindJSON(&user)
	if gContext.HandleAPIError(
		ctx,
		"Invalid input. Please check the submitted data and try again.",
		err,
		http.StatusBadRequest,
	) {
		return
	}

	err = services.PrepareUser(&user)
	if gContext.HandleAPIError(
		ctx,
		"An error occurred while processing your request. Please try again later.",
		err,
		http.StatusInternalServerError,
	) {
		return
	}

	err = services.CreateUser(&user)
	if gContext.HandleAPIError(
		ctx,
		"We encountered an issue while registering your account. Please try again later.",
		err,
		http.StatusInternalServerError,
	) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User registered successfully",
		"user_id": user,
	})
}

func LoginUserHandler(ctx *gin.Context) {
	var loginInput models.User
	var err error

	err = ctx.ShouldBindJSON(&loginInput)
	if gContext.HandleAPIError(
		ctx,
		"Invalid input. Please check the submitted data and try again.",
		err,
		http.StatusBadRequest,
	) {
		return
	}

	jwt, err := services.VerifyUser(&loginInput.Username, &loginInput.Password)
	if gContext.HandleAPIError(
		ctx,
		"Invalid username or password. Please check your credentials and try again.",
		err,
		http.StatusUnauthorized,
	) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": jwt})
}
