package jwt

import (
	"time"

	"github.com/RobertoSuarez/api_metal/models"
	"github.com/dgrijalva/jwt-go"
)

var myKey = []byte("super_secret")

func GenerateJWT(user models.User) (string, error) {

	payload := jwt.MapClaims{
		"ID":     user.ID.String(),
		"correo": user.Correo,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenStr, err := token.SignedString(myKey)
	if err != nil {
		return tokenStr, err
	}

	return tokenStr, nil
}
