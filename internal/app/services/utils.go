package services

import (
	"edumatch/internal/config"
	"errors"
	"strconv"
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

func GenerateToken(userID uuid.UUID) (string, error) {
	//  exp time
	expTime, _ := strconv.Atoi(config.GetEnv("JWT_EXP_TIME", "24"))
	// Define the claims for the token
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24 * time.Duration(expTime)).Unix(),
	}

	// Generate the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := config.GetEnv("JWT_SECRET", "nothing")
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func GenerateRefreshToken(userID uuid.UUID) (string, error) {
	expTime, _ := strconv.Atoi(config.GetEnv("JWT_REFRESH_EXP_TIME", "30"))

	// Define the claims for the refresh token
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24 * time.Duration(expTime)).Unix(),
	}

	// Generate the refresh token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretToken := config.GetEnv("JWT_REFRESH_SECRET", "nothing")
	signedToken, err := token.SignedString([]byte(secretToken))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ValidateToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method used
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token")
		}
		secretKey := config.GetEnv("JWT_SECRET", "nothing")
		return []byte(secretKey), nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return uuid.Nil, errors.New("invalid token")
	}

	userID, ok := claims["user_id"].(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("invalid token")
	}

	return userID, nil
}

func ValidateRefreshToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method used
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token")
		}

		secretToken := config.GetEnv("JWT_REFRESH_SECRET", "nothing")

		return []byte(secretToken), nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return uuid.Nil, errors.New("invalid token")
	}

	userID, ok := claims["user_id"].(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("invalid token")
	}

	return userID, nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
