package auth

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Alan-Luc/VertiLog/backend/utils/apiErrors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

const TOKEN_TTL_HRS = "1"

func GenerateJWT(userId int) (string, error) {
	// ttlStr := os.Getenv("TOKEN_TTL_HRS")
	ttlStr := TOKEN_TTL_HRS
	if ttlStr == "" {
		return "", errors.New("TOKEN_TTL_HRS environment variable is missing")
	}

	ttl, err := strconv.Atoi(ttlStr)
	if err != nil {
		return "", errors.Wrap(err, "Failed to convert TOKEN_TTL_HRS to int")
	}

	// generate JWT claims
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["userId"] = userId
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(ttl)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// sign the JWT
	tokenString, err := token.SignedString([]byte(os.Getenv("API_SECRET")))
	if err != nil {
		return "", errors.Wrap(err, "Failed to sign JWT")
	}

	return tokenString, nil
}

func ValidateJWT(ctx *gin.Context) error {
	tokenString := ExtractJWT(ctx)
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	_, err := jwt.Parse(tokenString, ParseJWT)
	if err != nil {
		return errors.Wrap(err, "Failed to parse the JWT")
	}

	return nil
}

func ParseJWT(t *jwt.Token) (interface{}, error) {
	// validate the alg
	if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
		err := fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		return nil, errors.Wrap(err, "Error occurred when parsing JWT")
	}

	// return api secret
	return []byte(os.Getenv("API_SECRET")), nil
}

func ExtractJWT(ctx *gin.Context) string {
	// check request header
	bearerToken := ctx.Request.Header.Get("Authorization")
	splitToken := strings.Split(bearerToken, " ")
	if len(splitToken) == 2 {
		return splitToken[1]
	}

	return ""
}

func ExtractUserIdFromJWT(ctx *gin.Context) (int, error) {
	tokenString := ExtractJWT(ctx)
	t, err := jwt.Parse(tokenString, ParseJWT)
	if err != nil {
		return 0, errors.Wrap(err, "Error occurred when parsing JWT")
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	log.Println(claims["userId"])
	if ok && t.Valid {
		userId, ok := claims["userId"].(float64)
		if !ok {
			return 0, nil
		} else {
			return int(userId), nil
		}
	}

	return 0, nil
}

func JWTAuthMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := ValidateJWT(ctx)
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
