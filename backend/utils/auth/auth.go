package auth

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userId int) (string, error) {
	token_ttl, err := strconv.Atoi(os.Getenv("TOKEN_TTL_HRS"))
	if err != nil {
		return "", err
	}

	// generate JWT claims
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["userId"] = userId
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(token_ttl)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// sign the JWT
	tokenString, err := token.SignedString([]byte(os.Getenv("API_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateJWT(ctx *gin.Context) error {
	tokenString := ExtractJWT(ctx)
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	_, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// validate the alg
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}

		// return api secret
		return []byte(os.Getenv("API_SECRET")), nil
	})

	if err != nil {
		return err
	}

	return nil
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

func JWTAuthMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := ValidateJWT(ctx)
		if err != nil {
			ctx.String(http.StatusUnauthorized, "Unauthorized")
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
