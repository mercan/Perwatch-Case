package main

import (
	"Perwatch-case/internal/config"
	"Perwatch-case/internal/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New(fiber.Config{
		AppName:       "Perwatch-Case",
		CaseSensitive: true,
	})

	// Routes
	routes.NewUserRoutes().Setup(app)
	routes.NewFormRoutes().Setup(app)

	port := config.GetConfig().GetServer().Port
	if err := app.Listen(":" + port); err != nil {
		panic(err)
	}
}
