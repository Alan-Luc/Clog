package handlers

import (
	"go/token"
	"html"
	"net/http"
	"strings"

	"github.com/Alan-Luc/VertiLog/backend/database"
	"github.com/Alan-Luc/VertiLog/backend/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(ctx *gin.Context) {
	var registerInput models.User
	var err error

	if err = ctx.ShouldBindJSON(&registerInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = PrepareUser(&registerInput)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = CreateUser(&registerInput)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "User registered successfully",
			"user_id": registerInput.ID,
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

func LoginUser(ctx *gin.Context) {
	var loginInput models.User
	var err error

	if err = ctx.ShouldBindJSON(&loginInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error ": err.Error()})
		return
	}

	jwt, err := VerifyUser(&loginInput.Username, &loginInput.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"Invalid credentials ": err.Error()})
		return
	} else {
		ctx.JSON(http.StatusOK, gin.H{"token": jwt})
	}

}

func VerifyUser(username, password *string) (string, error) {
	var err error
	var user models.User
	// check if user in db
	if err = database.DB.Model(&models.User{}).Find(&user, username).Error; err != nil {
		return "", err
	}

	// verify password
	var pw string = *password
	err = VerifyPassword(user.Password, pw)
	if err != nil {
		return "", err
	}

	jwt, err := auth.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return jwt, nil
}

func VerifyPassword(hashedPw, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPw), []byte(password))
}
