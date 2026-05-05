package utils

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var SECRET = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}


// Get expiration hours from env
func getExpirationHours() time.Duration {
	hoursStr := os.Getenv("JWT_EXPIRATION_HOURS")
	if hoursStr == "" {
		return 24 * time.Hour
	}

	hours, err := strconv.Atoi(hoursStr)
	if err != nil || hours <= 0 {
		return 24 * time.Hour
	}

	return time.Duration(hours) * time.Hour
}

func GenerateJWT(userID uint) (string, error) {
	expirationTime := time.Now().Add(getExpirationHours())


	claims := &Claims{
		UserID: userID,

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(SECRET)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}


func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return SECRET, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
