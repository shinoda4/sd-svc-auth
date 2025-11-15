package token

import (
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	secret       []byte
	expireHours  int
	refreshHours int
)

func init() {
	log.Println("JWT initialing...")
	secret = []byte(getenv("JWT_SECRET", "change_me"))
	expireHours = getenvInt("JWT_EXPIRE_HOURS", 1)
	refreshHours = getenvInt("JWT_REFRESH_HOURS", 72)
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getenvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return def
}

type Claims struct {
	TokenType string `json:"token_type"`
	UserID    string `json:"uid"`
	Email     string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID, email string) (string, time.Duration, error) {
	return generateToken(userID, email, time.Duration(expireHours)*time.Hour, "access")
}

func GenerateRefreshJWT(userID, email string) (string, time.Duration, error) {
	return generateToken(userID, email, time.Duration(refreshHours)*time.Hour, "refresh")
}

func generateToken(userID, email string, duration time.Duration, tokenType string) (string, time.Duration, error) {
	exp := time.Now().Add(duration)
	claims := &Claims{
		TokenType: tokenType,
		UserID:    userID,
		Email:     email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(secret)
	if err != nil {
		return "", 0, err
	}
	return ss, duration, nil
}
func ParseAndValidate(tokenStr string) (*Claims, error) {
	return parseToken(tokenStr)
}

func ParseAndValidateRefresh(tokenStr string) (*Claims, error) {
	return parseToken(tokenStr)
}

func parseToken(tokenStr string) (*Claims, error) {
	tok, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := tok.Claims.(*Claims); ok && tok.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
