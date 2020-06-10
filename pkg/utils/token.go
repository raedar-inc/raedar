package utils

import (
	"github.com/dgrijalva/jwt-go"
	"raedar/tools"
	"time"
)

func CreateResetPasswordToken(email string) (string, error) {
	ecdsaKey := tools.GetPrivEcdsaKey()

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
