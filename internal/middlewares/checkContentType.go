package middlewares

import (
	"Perwatch-case/internal/utils"
	"github.com/gofiber/fiber/v2"
)

// CheckContentType middleware checks if the request content type is application/json or application/json; charset=utf-8
func CheckContentType(ctx *fiber.Ctx) error {
	if ctx.Get("Content-Type") != fiber.MIMEApplicationJSON &&
		ctx.Get("Content-Type") != fiber.MIMEApplicationJSONCharsetUTF8 {
		return ctx.Status(fiber.StatusUnsupportedMediaType).JSON(utils.NewErrorResponse("Content type must be application/json or application/json; charset=utf-8"))
	}

	return ctx.Next()
}
