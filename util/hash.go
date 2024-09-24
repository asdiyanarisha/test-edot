package util

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"test-edot/src/models"
	"time"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func GenerateJWT(user models.User) (string, error) {
	secretKey := []byte(GetEnv("JWT_SECRET_KEY", ""))

	claims := jwt.MapClaims{
		"authorized": true,
		"userId":     strconv.Itoa(user.Id),
		"role":       user.Role,
		"exp":        time.Now().Add(time.Hour * 1).Unix(), // Token berlaku selama 1 jam
	}

	// Membuat token dengan algoritma signing HMAC SHA256 dan klaim yang sudah diset
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Menandatangani token dengan kunci rahasia
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	secretKey := []byte(GetEnv("JWT_SECRET_KEY", ""))

	// Parsing token dan validasi
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Memastikan algoritma yang digunakan adalah HMAC SHA256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Memeriksa apakah token valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}
