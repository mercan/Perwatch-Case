package routes

import (
	"Perwatch-case/internal/controllers"
	"Perwatch-case/internal/middlewares"
	"github.com/gofiber/fiber/v2"
)

type FormRoutes struct {
	FormController  controllers.FormController
	StockController controllers.StockController
}

func NewFormRoutes() FormRoutes {
	return FormRoutes{
		FormController:  controllers.NewFormController(),
		StockController: controllers.NewStockController(),
	}
}

func (u FormRoutes) Setup(app *fiber.App) {
	// api/v1
	apiV1 := app.Group("/api/v1")

	// Routes
	apiV1.Post("/form/create", middlewares.CheckContentType, middlewares.IsAuthenticated, u.FormController.CreateForm)
	apiV1.Get("/form/:FormID", middlewares.IsAuthenticated, u.FormController.GetForm)
	apiV1.Get("/form", middlewares.IsAuthenticated, u.FormController.GetForms)
	apiV1.Delete("/form/:FormID", middlewares.IsAuthenticated, u.FormController.DeleteForm)
	apiV1.Put("/form/:FormID", middlewares.IsAuthenticated, u.FormController.UpdateFormName)

	// Fields
	apiV1.Post("/form/:FormID/field", middlewares.CheckContentType, middlewares.IsAuthenticated, u.FormController.CreateFormField)
	apiV1.Get("/form/:FormID/field", middlewares.IsAuthenticated, u.FormController.GetFormFields)
	apiV1.Get("/form/:FormID/field/:FieldID", middlewares.IsAuthenticated, u.FormController.GetFormField)
	apiV1.Delete("/form/:FormID/field/:FieldID", middlewares.IsAuthenticated, u.FormController.DeleteFormField)

	// Stocks
	apiV1.Post("/form/:FormID/stock", middlewares.IsAuthenticated, u.StockController.CreateStock)
	apiV1.Get("/form/:FormID/stock", middlewares.IsAuthenticated, u.StockController.GetStocks)
	apiV1.Get("/form/:FormID/stock/:StockID", middlewares.IsAuthenticated, u.StockController.GetStock)
	apiV1.Delete("/form/:FormID/stock/:StockID", middlewares.IsAuthenticated, u.StockController.DeleteStock)
	apiV1.Put("/form/:FormID/stock/:StockID", middlewares.IsAuthenticated, u.StockController.UpdateStockValue)

}
