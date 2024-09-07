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
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(ctx *gin.Context) {
	var user models.User
	var err error

	err = ctx.ShouldBindJSON(&user)
	if gContext.HandleReqError(ctx, err, http.StatusBadRequest) {
		return
	}

	err = PrepareUser(&user)
	if gContext.HandleReqError(ctx, err, http.StatusInternalServerError) {
		return
	}

	err = CreateUser(&user)
	if gContext.HandleReqError(ctx, err, http.StatusInternalServerError) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User registered successfully",
		"user_id": user,
	})
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

	err = ctx.ShouldBindJSON(&loginInput)
	if gContext.HandleReqError(ctx, err, http.StatusBadRequest) {
		return
	}

	jwt, err := VerifyUser(&loginInput.Username, &loginInput.Password)
	if gContext.HandleReqError(ctx, err, http.StatusUnauthorized) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": jwt})
}

func VerifyUser(username, password *string) (string, error) {
	var user models.User
	var err error
	// check if user in db
	if err = database.DB.Model(&models.User{}).Where("username = ?", username).Find(&user).Error; err != nil {
		return "", err
	}

	// verify password
	var pw string = *password
	err = VerifyPassword(user.Password, pw)
	if err != nil {
		return "", err
	}

	jwt, err := auth.GenerateJWT(user.ID)
	if err != nil {
		return "", err
	}

	return jwt, nil
}

func VerifyPassword(hashedPw, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPw), []byte(password))
}

func GetCurrentUserId(ctx *gin.Context) (int, error) {
	userId, err := auth.ExtractUserIdFromJWT(ctx)
	if err != nil {
		return 0, err
	}

	return userId, nil
}

func GetCurrentUser(userId int) (models.User, error) {
	var user models.User
	var err error
	// check if user in db
	if err = database.DB.Model(&models.User{}).Where("id = ?", userId).Find(&user).Error; err != nil {
		return models.User{}, err
	}

	return user, nil
}
