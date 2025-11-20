/*
 * Copyright (c) 2025-11-20 shinoda4
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
	return ParseToken(tokenStr)
}

func ParseAndValidateRefresh(tokenStr string) (*Claims, error) {
	return ParseToken(tokenStr)
}

func ParseToken(tokenStr string) (*Claims, error) {
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
