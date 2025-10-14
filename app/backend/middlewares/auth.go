package middlewares

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JWTSecret = getenv("JWT_SECRET", "takopi-no-genzai")

func CreateJWT(userID int, username string) (string, error) {
	claims := jwt.MapClaims{
		"sub":      strconv.Itoa(userID),
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	log.Printf("JWT created for user %s (id=%d)", username, userID)
	return t.SignedString([]byte(JWTSecret))
}

func ParseJWT(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func GetUserIDFromCookie(r *http.Request) (int, error) {
	c, err := r.Cookie("token")
	if err != nil {
		return 0, errors.New("no token cookie")
	}
	claims, err := ParseJWT(c.Value)
	if err != nil {
		return 0, errors.New("invalid token")
	}
	sub, ok := claims["sub"].(string)
	if !ok {
		return 0, errors.New("invalid claims")
	}
	uid, _ := strconv.Atoi(sub)
	return uid, nil
}

func getenv(k, fallback string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return fallback
}