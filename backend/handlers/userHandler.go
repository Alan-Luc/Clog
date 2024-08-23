package handlers

import (
	"html"
	"net/http"
	"strings"

	"github.com/Alan-Luc/VertiLog/backend/database"
	"github.com/Alan-Luc/VertiLog/backend/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(ctx *gin.Context) {
	var user models.User
	var err error

	if err = ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = PrepareUser(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = CreateUser(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "User registered successfully",
			"user_id": user.ID,
		})
	}
}

func CreateUser(u *models.User) error {
	if err := database.DB.Create(&u).Error; err != nil {
		return err
	}
	return nil
}

func PrepareUser(u *models.User) error {
	// hash password
	hashedPW, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPW)

	// trim username
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))

	return nil
}
