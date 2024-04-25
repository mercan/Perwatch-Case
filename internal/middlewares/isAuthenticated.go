package middlewares

import (
	"Perwatch-case/internal/config"
	"Perwatch-case/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
)

func extractToken(Bearer string) (string, error) {
	splitToken := strings.Split(Bearer, "Bearer ")
	if len(splitToken) != 2 {
		return "", fiber.ErrUnauthorized
	}

	if splitToken[1] != "" {
		return splitToken[1], nil
	}

	return "", fiber.ErrUnauthorized
}

// IsAuthenticated middleware checks if the request has an Authorization header
func IsAuthenticated(ctx *fiber.Ctx) error {
	Bearer := ctx.Get("Authorization")
	if Bearer == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(utils.NewErrorResponse("Unauthorized"))
	}

	token, err := extractToken(Bearer)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(utils.NewErrorResponse("Unauthorized"))
	}

	// Parse and validate JWT token
	claims := jwt.MapClaims{}
	jwtToken, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.ErrUnauthorized
		}

		return []byte(config.GetConfig().GetJWTConfig().Secret), nil
	})

	if err != nil || !jwtToken.Valid {
		return ctx.Status(fiber.StatusUnauthorized).JSON(utils.NewErrorResponse("Unauthorized"))
	}

	userId, err := primitive.ObjectIDFromHex(claims["id"].(string))
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(utils.NewErrorResponse("Unauthorized"))
	}

	// Set user context with extracted data
	ctx.Locals("userID", userId)
	ctx.Locals("firstname", claims["firstname"])
	ctx.Locals("lastname", claims["lastname"])
	ctx.Locals("username", claims["username"])
	ctx.Locals("token", token)

	return ctx.Next()
}
