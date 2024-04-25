package controllers

import (
	"Perwatch-case/internal/models"
	"Perwatch-case/internal/services"
	"Perwatch-case/internal/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StockController struct {
	StockService services.StockServiceInterface
}

func NewStockController() StockController {
	return StockController{
		StockService: services.NewStockService(),
	}
}

func (st StockController) CreateStock(c *fiber.Ctx) error {
	var requestParams models.StockRequestParams
	var request models.CreateStockRequest

	if err := c.ParamsParser(&requestParams); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Invalid request"))
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Invalid request"))
	}

	userID := c.Locals("userID").(primitive.ObjectID)
	formID, err := primitive.ObjectIDFromHex(requestParams.FormID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Invalid form ID"))
	}

	if len(request.Fields) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("All fields are required"))
	}

	for i := 0; i < len(request.Fields); i++ {
		if request.Fields[i].Name == "" || request.Fields[i].Value == "" {
			return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("All fields are required"))
		}
	}

	if err := st.StockService.CreateStock(userID, formID, request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse(err.Error()))
	}

	return c.Status(fiber.StatusCreated).JSON(utils.NewSuccessResponse(nil, "Stock created successfully"))
}

func (st StockController) GetStocks(c *fiber.Ctx) error {
	var requestParams models.StockRequestParams

	if err := c.ParamsParser(&requestParams); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Invalid request"))
	}

	userID := c.Locals("userID").(primitive.ObjectID)
	formID, err := primitive.ObjectIDFromHex(requestParams.FormID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Invalid form ID"))
	}

	stocks, err := st.StockService.GetStocks(userID, formID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(utils.NewSuccessResponse(stocks, "Stocks retrieved successfully"))
}

func (st StockController) GetStock(c *fiber.Ctx) error {
	var requestParams models.StockRequestParams

	if err := c.ParamsParser(&requestParams); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Invalid request"))
	}

	userID := c.Locals("userID").(primitive.ObjectID)
	formID, err := primitive.ObjectIDFromHex(requestParams.FormID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Invalid form ID"))
	}

	stockID, err := primitive.ObjectIDFromHex(requestParams.StockID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Invalid stock ID"))
	}

	stock, err := st.StockService.GetStock(userID, formID, stockID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(utils.NewSuccessResponse(stock, "Stock retrieved successfully"))
}

func (st StockController) DeleteStock(c *fiber.Ctx) error {
	var requestParams models.StockRequestParams

	if err := c.ParamsParser(&requestParams); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Invalid request"))
	}

	userID := c.Locals("userID").(primitive.ObjectID)
	formID, err := primitive.ObjectIDFromHex(requestParams.FormID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Invalid form ID"))
	}

	stockID, err := primitive.ObjectIDFromHex(requestParams.StockID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Invalid stock ID"))
	}

	if err := st.StockService.DeleteStock(userID, formID, stockID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(utils.NewSuccessResponse(nil, "Stock deleted successfully"))
}

func (st StockController) UpdateStockValue(c *fiber.Ctx) error {
	var requestParams models.StockRequestParams
	var request models.UpdateStockValueRequest

	if err := c.ParamsParser(&requestParams); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Invalid request"))
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Invalid request"))
	}

	userID := c.Locals("userID").(primitive.ObjectID)
	formID, err := primitive.ObjectIDFromHex(requestParams.FormID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Invalid form ID"))
	}

	stockID, err := primitive.ObjectIDFromHex(requestParams.StockID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Invalid stock ID"))
	}

	if err := st.StockService.UpdateStockValue(userID, formID, stockID, request.Value); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(utils.NewSuccessResponse(nil, "Stock value updated successfully"))
}
