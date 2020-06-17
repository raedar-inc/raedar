package utils

import (
	"github.com/dgrijalva/jwt-go"
	"raedar/tools"
	"time"
)

var ecdsaKey = tools.GetPrivEcdsaKey()

func CreateResetPasswordToken(email string) (string, error) {
	claims := jwt.MapClaims{}
	claims["email"] = email
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix() //Token expires after 1 hour
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	token, err := jwtToken.SignedString(ecdsaKey)
	if err != nil {
		return "", err
	}
	return token, nil
}

func DecodeToken(token string) (interface{}, error) {
	claims := jwt.MapClaims{}
	claims["email"] = ""
	decodedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return tools.GetPubEcdsaKey(), nil
	})
	if err != nil {
		return nil, err
	}

	return decodedToken.Claims, nil
}
