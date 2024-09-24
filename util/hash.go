package util

import (
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
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

	// Membuat klaim (claims) untuk token JWT
	claims := jwt.MapClaims{
		"authorized": true,
		"userId":     user.Id,
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
