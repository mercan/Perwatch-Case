package helpers

import (
	"Perwatch-case/internal/config"
	"Perwatch-case/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func GenerateJWT(user *models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS512)
	now := time.Now().UTC()
	expiration := config.GetConfig().GetJWTConfig().Expiration
	expirationDuration, _ := time.ParseDuration(expiration + "h")

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = now.Add(expirationDuration).Unix()
	claims["iat"] = now.Unix()
	claims["id"] = user.ID
	claims["firstname"] = user.Firstname
	claims["lastname"] = user.Lastname
	claims["username"] = user.Username

	// Generate encoded token and send it as response.
	tokenString, err := token.SignedString([]byte(config.GetConfig().GetJWTConfig().Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
