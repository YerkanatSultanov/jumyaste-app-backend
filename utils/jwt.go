package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte("secretkey")

type Claims struct {
	UserID       int `json:"user_id"`
	RoleID       int `json:"role_id"`
	CompanyID    int `json:"company_id"`
	DepartmentID int `json:"department_id"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID, roleID, companyId, depId int) (string, error) {
	claims := Claims{
		UserID:       userID,
		RoleID:       roleID,
		CompanyID:    companyId,
		DepartmentID: depId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func GenerateRefreshToken(userID, roleId, companyId, depId int) (string, error) {
	claims := Claims{
		UserID:       userID,
		RoleID:       roleId,
		CompanyID:    companyId,
		DepartmentID: depId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey := jwtSecret
	refreshToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}

func ParseJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func ValidateRefreshToken(refreshToken string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %v", err)
	}

	if !token.Valid {
		return nil, errors.New("refresh token is not valid")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("cannot parse claims")
	}

	return claims, nil
}
