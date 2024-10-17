package auth

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Alan-Luc/VertiLog/backend/utils/apiErrors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

const (
	ACCESS_TOKEN_TTL  = "1h"
	REFRESH_TOKEN_TTL = "168h" // 1 week
)

type CustomClaims struct {
	UserID int `json:"userID"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID int) (access, refresh string, err error) {
	apiSecret := os.Getenv("API_SECRET")
	if apiSecret == "" {
		return "", "", errors.New("Missing env var API_SECRET")
	}

	refreshSecret := os.Getenv("REFRESH_SECRET")
	if refreshSecret == "" {
		return "", "", errors.New("Missing env var REFRESH_SECRET")
	}

	accessTTL, err := time.ParseDuration(ACCESS_TOKEN_TTL)
	if err != nil {
		return "", "", errors.Wrap(err, "Failed to parse ACCESS_TOKEN_TTL")
	}

	refreshTTL, err := time.ParseDuration(REFRESH_TOKEN_TTL)
	if err != nil {
		return "", "", errors.Wrap(err, "Failed to parse REFRESH_TOKEN_TTL")
	}

	// generate access token claims
	accessTokenClaims := CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTTL)),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)

	// generate refresh token claims
	refreshTokenClaims := CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTTL)),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)

	// sign the access token
	accessTokenStr, err := accessToken.SignedString([]byte(apiSecret))
	if err != nil {
		return "", "", errors.Wrap(err, "Failed to sign access token")
	}

	// sign the refresh token
	refreshTokenStr, err := refreshToken.SignedString([]byte(refreshSecret))
	if err != nil {
		return "", "", errors.Wrap(err, "Failed to sign refresh token")
	}

	return accessTokenStr, refreshTokenStr, nil
}

func ValidateJWT(ctx *gin.Context, isRefreshToken bool) error {
	var secret string
	var tokenString string
	var err error

	if isRefreshToken {
		secret = os.Getenv("REFRESH_SECRET")
		tokenString, err = ExtractRefreshToken(ctx)
	} else {
		secret = os.Getenv("API_SECRET")
		tokenString, err = ExtractAccessToken(ctx)
	}

	if err != nil {
		return errors.WithMessage(err, "Error occurred when extracting token")
	}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&CustomClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return ParseJWTSecret(t, secret)
		},
	)
	if err != nil {
		return errors.Wrap(err, "Failed to parse the JWT")
	}

	claims, ok := token.Claims.(*CustomClaims)
	if ok && token.Valid {
		ctx.Set("userID", claims.UserID)
	}

	return nil
}

func ParseJWTSecret(t *jwt.Token, secret string) (interface{}, error) {
	if secret == "" {
		return nil, errors.New("JWT secret is missing")
	}

	// validate the alg
	if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
		err := fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		return nil, errors.Wrap(err, "Error occurred when parsing JWT")
	}

	// return api secret
	return []byte(secret), nil
}

func ExtractAccessToken(ctx *gin.Context) (string, error) {
	// check request header
	bearerToken := ctx.Request.Header.Get("Authorization")
	splitToken := strings.Split(bearerToken, " ")
	if len(splitToken) == 2 {
		return splitToken[1], nil
	}

	return "", errors.New("Access token is missing or was not provided")
}

func ExtractRefreshToken(ctx *gin.Context) (string, error) {
	refreshToken, err := ctx.Cookie("refreshToken")
	if err != nil {
		return "", errors.Wrap(err, "Refresh token is missing or was not provided")
	}

	return refreshToken, nil
}

func ExtractUserIdFromJWT(ctx *gin.Context) (int, error) {
	tokenString, err := ExtractAccessToken(ctx)
	if err != nil {
		return 0, errors.WithMessage(err, "Error occurred when extracting access token")
	}

	secret := os.Getenv("API_SECRET")

	token, err := jwt.ParseWithClaims(
		tokenString,
		&CustomClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return ParseJWTSecret(t, secret)
		},
	)
	if err != nil {
		return 0, errors.Wrap(err, "Failed to parse the JWT")
	}

	claims, ok := token.Claims.(*CustomClaims)
	if ok && token.Valid {
		return claims.UserID, nil
	}

	return 0, errors.New("Failed to extract userID")
}

func JWTAuthMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := ValidateJWT(ctx, false)
		if apiErrors.HandleAPIError(
			ctx,
			"Authentication is required to access this resource.",
			errors.WithMessage(err, "Error occurred when validating JWT"),
			http.StatusUnauthorized,
		) {
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
