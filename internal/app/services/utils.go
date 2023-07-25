package services

import (
	custom_errors "edumatch/internal/app/errors"
	"edumatch/internal/app/models"
	"edumatch/internal/config"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), err
}

func GenerateToken(userID uuid.UUID, role models.Role, isRefreshToken bool) (string, error) {
	expTimeKey := "JWT_EXP_TIME"
	secretKey := "JWT_SECRET"

	if isRefreshToken {
		expTimeKey = "JWT_REFRESH_EXP_TIME"
		secretKey = "JWT_REFRESH_SECRET"
	}
	expTime, _ := strconv.Atoi(config.GetEnv(expTimeKey, "24"))

	// Define the claims for the token
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24 * time.Duration(expTime)).Unix(),
	}

	// Generate the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretToken := config.GetEnv(secretKey, "nothing")
	signedToken, err := token.SignedString([]byte(secretToken))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ValidateToken(tokenString string, isRefreshToken bool) (uuid.UUID, models.Role, error) {
	// Remove the "Bearer " prefix if it exists
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	var secretKey string
	if isRefreshToken {
		secretKey = config.GetEnv("JWT_REFRESH_SECRET", "nothing")
	} else {
		secretKey = config.GetEnv("JWT_SECRET", "nothing")
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method used
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, custom_errors.ErrInvalidToken
		}

		return []byte(secretKey), nil
	})

	if err != nil {
		return uuid.Nil, "", custom_errors.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return uuid.Nil, "", custom_errors.ErrInvalidToken
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return uuid.Nil, "", custom_errors.ErrInvalidToken
	}

	var userRole models.Role
	switch roleClaim := claims["role"].(type) {
	case string:
		userRole = models.Role(roleClaim)
	default:
		return uuid.Nil, "", custom_errors.ErrInvalidToken
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, "", custom_errors.ErrInvalidToken
	}

	return userID, userRole, nil
}

func CheckPassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}
