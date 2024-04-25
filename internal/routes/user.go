package routes

import (
	"Perwatch-case/internal/controllers"
	"Perwatch-case/internal/middlewares"
	"github.com/gofiber/fiber/v2"
)

type UserRoutes struct {
	userController controllers.UserController
}

func NewUserRoutes() UserRoutes {
	return UserRoutes{
		userController: controllers.NewUserController(),
	}
}

func (u UserRoutes) Setup(app *fiber.App) {
	// api/v1
	apiV1 := app.Group("/api/v1")

	// Routes
	apiV1.Post("/register", middlewares.CheckContentType, u.userController.Register)
	apiV1.Post("/login", middlewares.CheckContentType, u.userController.Login)

	apiV1.Get("/account", middlewares.IsAuthenticated, u.userController.Account)

}
