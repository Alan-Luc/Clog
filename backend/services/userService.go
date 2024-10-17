package services

import (
	"fmt"
	"html"
	"strings"

	"github.com/Alan-Luc/VertiLog/backend/database"
	"github.com/Alan-Luc/VertiLog/backend/models"
	"github.com/Alan-Luc/VertiLog/backend/utils/auth"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateUser(u *models.User) error {
	if err := database.DB.Create(&u).Error; err != nil {
		if err.Error() == gorm.ErrDuplicatedKey.Error() {
			return errors.Wrap(
				err,
				"Failed to create user: a user with the same ID already exists",
			)
		}
		return errors.Wrap(err, "Failed to create user")
	}
	return nil
}

func PrepareUser(u *models.User) error {
	// hash password
	hashedPW, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.Wrap(
			err,
			"Error occurred when hashing password",
		)
	}
	u.Password = string(hashedPW)

	// trim username
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))

	return nil
}

func VerifyUser(username, password string) (string, string, error) {
	var user models.User
	var err error
	// check if user in db
	err = user.FindByUsername(database.DB, username)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to find user with username '%s'", username)
		return "", "", errors.WithMessage(err, errMsg)
	}

	// verify password
	var pw string = password
	err = VerifyPassword(user.Password, pw)
	if err != nil {
		return "", "", errors.Wrap(err, "Error occurred when verifying password")
	}

	// TODO make this its own function that checks cache for issued tokens
	accessToken, refreshToken, err := auth.GenerateJWT(user.ID)
	if err != nil {
		return "", "", errors.WithMessage(err, "Error occurred when generating JWT")
	}

	return accessToken, refreshToken, nil
}

func VerifyPassword(hashedPW, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPW), []byte(password))
}

func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func GetCurrentUserId(ctx *gin.Context) (int, error) {
	userID, err := auth.ExtractUserIdFromJWT(ctx)
	if err != nil {
		return 0, errors.WithMessage(err, "Error occurred when extracting user ID")
	}

	return userID, nil
}

func GetCurrentUser(userID int) (*models.User, error) {
	var user models.User
	var err error
	// check if user in db
	err = user.FindByID(database.DB, userID)
	if err != nil {
		return nil, errors.Wrap(
			err,
			"Error occurred when finding user",
		)
	}

	return &user, nil
}

func UpdateUserPassword(userID int, currPW, newPW string) error {
	var user *models.User
	var err error

	user, err = GetCurrentUser(userID)
	if err != nil || user == nil {
		return errors.Wrap(err, "Error occurred when fetching user for password update")
	}

	err = VerifyPassword(user.Password, currPW)
	if err != nil {
		return errors.Wrap(err, "Current password is incorrect")
	}

	hashedNewPW, err := HashPassword(newPW)
	if err != nil {
		return errors.Wrap(
			err,
			"Error occurred when hashing new password",
		)
	}

	err = user.UpdatePassword(database.DB, userID, string(hashedNewPW))
	if err != nil {
		return errors.Wrap(
			err,
			"Error occurred when updating password in database",
		)
	}

	return nil
}
