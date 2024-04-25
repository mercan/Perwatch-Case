package controllers

import (
	"Perwatch-case/internal/models"
	"Perwatch-case/internal/services"
	"Perwatch-case/internal/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserController struct {
	UserService services.UserServiceInterface
}

func NewUserController() UserController {
	return UserController{
		UserService: services.NewUserService(),
	}
}

func (u *UserController) Register(c *fiber.Ctx) error {
	var request models.UserRegisterRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Invalid request"))
	}

	if request.Firstname == "" || request.Lastname == "" || request.Username == "" || request.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("All fields are required"))
	}

	if len(request.Username) < 6 || len(request.Username) > 16 {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Username should be between 6 and 16 characters"))
	}

	if len(request.Password) < 6 || len(request.Password) > 16 {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Password should be between 6 and 16 characters"))
	}

	token, err := u.UserService.Register(request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse(err.Error()))
	}

	return c.Status(fiber.StatusCreated).JSON(utils.NewSuccessResponse(token, "User created successfully"))
}

func (u *UserController) Login(c *fiber.Ctx) error {
	var request models.UserLoginRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Invalid request"))
	}

	if len(request.Username) < 6 || len(request.Username) > 16 {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Username should be between 6 and 16 characters"))
	}

	if len(request.Password) < 6 || len(request.Password) > 16 {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Password should be between 6 and 16 characters"))
	}

	token, err := u.UserService.Login(request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(utils.NewSuccessResponse(token, "User logged in successfully"))
}

func (u *UserController) Account(c *fiber.Ctx) error {
	userID := c.Locals("userID").(primitive.ObjectID).Hex()
	firstname := c.Locals("firstname").(string)
	lastname := c.Locals("lastname").(string)
	username := c.Locals("username").(string)

	user := map[string]interface{}{
		"_id":       userID,
		"firstname": firstname,
		"lastname":  lastname,
		"username":  username,
	}

	return c.Status(fiber.StatusOK).JSON(utils.NewSuccessResponse(user, "User retrieved successfully"))
}
