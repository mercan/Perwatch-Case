package controllers

import (
	"Perwatch-case/internal/models"
	"Perwatch-case/internal/services"
	"Perwatch-case/internal/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// StockController struct holds the StockService instance for interacting with form logic
type StockController struct {
	StockService services.StockServiceInterface
}

// NewStockController creates a new StockController instance with the injected StockService
func NewStockController() StockController {
	return StockController{
		StockService: services.NewStockService(),
	}
}

// CreateStock creates a new stock entry
func (st StockController) CreateStock(c *fiber.Ctx) error {
	// Parse form ID from request parameters
	var requestParams models.StockRequestParams
	if err := c.ParamsParser(&requestParams); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Invalid request"))
	}
	
	// Parse stock data from request body
	var request models.CreateStockRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Invalid request"))
	}

	// Get user ID from context
	userID := c.Locals("userID").(primitive.ObjectID)
	
	// Convert form ID string to a MongoDB ObjectID
	formID, err := primitive.ObjectIDFromHex(requestParams.FormID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Invalid form ID"))
	}

	// Validate if any fields are missing
	if len(request.Fields) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("All fields are required"))
	}

	// Check if all fields have names and values
	for i := 0; i < len(request.Fields); i++ {
		if request.Fields[i].Name == "" || request.Fields[i].Value == "" {
			return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("All fields are required"))
		}
	}

	// Call StockService to create the stock entry and handle any errors
	if err := st.StockService.CreateStock(userID, formID, request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse(err.Error()))
	}

	return c.Status(fiber.StatusCreated).JSON(utils.NewSuccessResponse(nil, "Stock created successfully"))
}

// GetStocks retrieves all stocks for a specific form
func (st StockController) GetStocks(c *fiber.Ctx) error {
	// Parse form ID from request parameters
	var requestParams models.StockRequestParams
	if err := c.ParamsParser(&requestParams); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Invalid request"))
	}

	// Get user ID from context
	userID := c.Locals("userID").(primitive.ObjectID)

	// Convert form ID string to a MongoDB ObjectID
	formID, err := primitive.ObjectIDFromHex(requestParams.FormID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Invalid form ID"))
	}

	// Call StockService to get all stocks for the form and handle any errors
	stocks, err := st.StockService.GetStocks(userID, formID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(utils.NewSuccessResponse(stocks, "Stocks retrieved successfully"))
}

// GetStock retrieves a specific stock entry
func (st StockController) GetStock(c *fiber.Ctx) error {
	// Parse form and stock ID from request parameters
	var requestParams models.StockRequestParams
	if err := c.ParamsParser(&requestParams); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Invalid request"))
	}

	// Get user ID from context
	userID := c.Locals("userID").(primitive.ObjectID)

	// Convert form ID string to a MongoDB ObjectID
	formID, err := primitive.ObjectIDFromHex(requestParams.FormID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Invalid form ID"))
	}
	
	// Convert stock ID string to a MongoDB ObjectID
	stockID, err := primitive.ObjectIDFromHex(requestParams.StockID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Invalid stock ID"))
	}

	// Call StockService to get the specific stock entry and handle any errors
	stock, err := st.StockService.GetStock(userID, formID, stockID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(utils.NewSuccessResponse(stock, "Stock retrieved successfully"))
}

// Parse form and stock IDs from request parameters
func (st StockController) DeleteStock(c *fiber.Ctx) error {
	// Parse form and stock ID from request parameters
	var requestParams models.StockRequestParams
	if err := c.ParamsParser(&requestParams); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Invalid request"))
	}

	// Get user ID from context
	userID := c.Locals("userID").(primitive.ObjectID)

	// Convert form ID string to a MongoDB ObjectID
	formID, err := primitive.ObjectIDFromHex(requestParams.FormID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Invalid form ID"))
	}

	// Convert stock ID string to a MongoDB ObjectID
	stockID, err := primitive.ObjectIDFromHex(requestParams.StockID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Invalid stock ID"))
	}

	// Call StockService to delete the stock entry and handle any errors
	if err := st.StockService.DeleteStock(userID, formID, stockID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(utils.NewSuccessResponse(nil, "Stock deleted successfully"))
}

// UpdateStockValue updates the value of a specific stock entry
func (st StockController) UpdateStockValue(c *fiber.Ctx) error {
	// Parse form and stock ID from request parameters
	var requestParams models.StockRequestParams
	if err := c.ParamsParser(&requestParams); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Invalid request"))
	}

	// Parse update request data from request body
	var request models.UpdateStockValueRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Invalid request"))
	}

	// Get user ID from context
	userID := c.Locals("userID").(primitive.ObjectID)
	
	// Convert form ID string to a MongoDB ObjectID
	formID, err := primitive.ObjectIDFromHex(requestParams.FormID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Invalid form ID"))
	}

	// Convert stock ID string to a MongoDB ObjectID
	stockID, err := primitive.ObjectIDFromHex(requestParams.StockID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Invalid stock ID"))
	}

	// Call StockService to update the stock value and handle any errors
	if err := st.StockService.UpdateStockValue(userID, formID, stockID, request.Value); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(utils.NewSuccessResponse(nil, "Stock value updated successfully"))
}
