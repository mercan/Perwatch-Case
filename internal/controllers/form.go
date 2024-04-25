package controllers

import (
	"Perwatch-case/internal/models"
	"Perwatch-case/internal/services"
	"Perwatch-case/internal/utils"
	"errors"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"slices"
)

var validFormTypes = []string{"text", "number", "number_decimal", "checkbox", "combobox"}

type FormController struct {
	FormService services.FormServiceInterface
}

func NewFormController() FormController {
	return FormController{
		FormService: services.NewFormService(),
	}
}

func (f *FormController) CreateForm(c *fiber.Ctx) error {
	var request models.FormNameRequestBody

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.NewErrorResponse("Invalid request"))
	}

	if request.FormName == "" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.NewErrorResponse("Form name is required"))
	}

	if len(request.FormName) < 3 || len(request.FormName) > 50 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.NewErrorResponse("Form name should be between 3 and 50 characters"))
	}

	form := models.NewForm()
	form.UserID = c.Locals("userID").(primitive.ObjectID)
	form.Name = request.FormName

	if err := f.FormService.CreateForm(form); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(utils.NewSuccessResponse(nil, "Form created successfully"))
}

func (f *FormController) GetForm(c *fiber.Ctx) error {
	var request models.FormFieldRequestParams

	if err := c.ParamsParser(&request); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.NewErrorResponse("Invalid request"))
	}

	if request.FormID == "" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.NewErrorResponse("Form ID is required"))
	}

	formID, err := primitive.ObjectIDFromHex(request.FormID)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.NewErrorResponse("Invalid form ID"))
	}

	userID := c.Locals("userID").(primitive.ObjectID)
	form, err := f.FormService.GetForm(userID, formID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(utils.NewSuccessResponse(form, "Form retrieved successfully"))
}

func (f *FormController) GetForms(c *fiber.Ctx) error {
	var request models.GetFormsRequest

	if err := c.QueryParser(&request); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.NewErrorResponse("Invalid request"))
	}

	if request.Page < 1 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.NewErrorResponse("Invalid page number"))
	}

	userID := c.Locals("userID").(primitive.ObjectID)
	formList, err := f.FormService.GetForms(userID, request.Page)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(utils.NewSuccessResponse(formList, "Forms retrieved successfully"))
}

func (f *FormController) DeleteForm(c *fiber.Ctx) error {
	var request models.FormIDRequestBody

	if err := c.ParamsParser(&request); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.NewErrorResponse("Invalid request"))
	}

	if request.FormID == "" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.NewErrorResponse("Form ID is required"))
	}

	formID, err := primitive.ObjectIDFromHex(request.FormID)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.NewErrorResponse("Invalid form ID"))
	}

	userID := c.Locals("userID").(primitive.ObjectID)
	if err := f.FormService.DeleteForm(userID, formID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Failed to delete form"))
	}

	return c.Status(fiber.StatusOK).JSON(utils.NewSuccessResponse(nil, "Form deleted successfully"))
}

func (f *FormController) UpdateFormName(c *fiber.Ctx) error {
	var requestParams models.FormFieldRequestParams
	var requestBody models.FormNameRequestBody

	if err := c.ParamsParser(&requestParams); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.NewErrorResponse("Invalid request"))
	}

	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.NewErrorResponse("Invalid request"))
	}

	if requestParams.FormID == "" || requestBody.FormName == "" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.NewErrorResponse("Form ID and form name are required"))
	}

	formID, err := primitive.ObjectIDFromHex(requestParams.FormID)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.NewErrorResponse("Invalid form ID"))
	}

	userID := c.Locals("userID").(primitive.ObjectID)
	if err := f.FormService.UpdateFormName(userID, formID, requestBody.FormName); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(utils.NewSuccessResponse(nil, "Form name updated successfully"))
}

func (f *FormController) CreateFormField(c *fiber.Ctx) error {
	var requestParams models.FormFieldRequestParams
	var request models.CreateFormFieldRequest

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

	if request.FieldType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Field type is required"))
	}

	if !slices.Contains(validFormTypes, request.FieldType) {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Invalid field type"))
	}

	if request.Sort < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Sort must be greater than 0"))
	}

	if request.FieldType == "text" {
		if err := validateText(request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse(err.Error()))
		}
	}

	if request.FieldType == "number" {
		if err := validateNumber(request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse(err.Error()))
		}
	}

	if request.FieldType == "number_decimal" {
		if err := validateNumberDecimal(request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse(err.Error()))
		}
	}
	
	if request.FieldType == "combobox" {
		if err := validateCombobox(request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse(err.Error()))
		}
	}

	if err := f.FormService.CreateFormField(userID, formID, request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse(err.Error()))
	}

	return c.Status(fiber.StatusCreated).JSON(utils.NewSuccessResponse(nil, "Form field created successfully"))
}

func (f *FormController) GetFormFields(c *fiber.Ctx) error {
	var requestParams models.FormFieldRequestParams

	if err := c.ParamsParser(&requestParams); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Invalid request"))
	}

	userID := c.Locals("userID").(primitive.ObjectID)
	formID, err := primitive.ObjectIDFromHex(requestParams.FormID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Invalid form ID"))
	}

	formFields, err := f.FormService.GetFormFields(userID, formID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(utils.NewSuccessResponse(formFields, "Form fields retrieved successfully"))

}

func (f *FormController) GetFormField(c *fiber.Ctx) error {
	var requestParams models.FormFieldRequestParams

	if err := c.ParamsParser(&requestParams); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Invalid request"))
	}

	userID := c.Locals("userID").(primitive.ObjectID)
	formID, err := primitive.ObjectIDFromHex(requestParams.FormID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Invalid form ID"))
	}

	fieldID, err := primitive.ObjectIDFromHex(requestParams.FieldID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Invalid field ID"))
	}

	formField, err := f.FormService.GetFormField(userID, fieldID, formID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(utils.NewSuccessResponse(formField, "Form field retrieved successfully"))
}

func (f *FormController) DeleteFormField(c *fiber.Ctx) error {
	var requestParams models.FormFieldRequestParams

	if err := c.ParamsParser(&requestParams); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Invalid request"))
	}

	userID := c.Locals("userID").(primitive.ObjectID)
	formID, err := primitive.ObjectIDFromHex(requestParams.FormID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Invalid form ID"))
	}

	fieldID, err := primitive.ObjectIDFromHex(requestParams.FieldID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Invalid field ID"))
	}

	if err := f.FormService.DeleteFormField(userID, fieldID, formID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(utils.NewSuccessResponse(nil, "Form field deleted successfully"))
}

func validateText(request models.CreateFormFieldRequest) error {
	if request.MinLength < 0 || request.MaxLength < 0 {
		return errors.New("min and max length cannot be negative")
	}

	if request.MinLength != 0 && request.MaxLength != 0 {
		if request.MinLength > request.MaxLength {
			return errors.New("min length cannot be greater than max length")
		}

		if request.MinLength == request.MaxLength {
			return errors.New("min and max length cannot be the same")
		}
	}

	return nil
}

func validateNumber(request models.CreateFormFieldRequest) error {
	if request.MinValue < 0 || request.MaxValue < 0 {
		return errors.New("min and max value cannot be negative")
	}

	if request.MinValue != 0 && request.MaxValue != 0 {
		if request.MinValue > request.MaxValue {
			return errors.New("min value cannot be greater than max value")
		}

		if request.MinValue == request.MaxValue {
			return errors.New("min and max value cannot be the same")
		}

		if request.DefaultNumber != 0 {
			if request.DefaultNumber < request.MinValue {
				return errors.New("default value cannot be less than min value")
			}

			if request.DefaultNumber > request.MaxValue {
				return errors.New("default value cannot be greater than max value")
			}
		}
	}

	return nil
}

func validateNumberDecimal(request models.CreateFormFieldRequest) error {
	if request.MinValueDecimal < 0 || request.MaxValueDecimal < 0 {
		return errors.New("min and max value cannot be negative")
	}

	if request.MinValueDecimal != 0 && request.MaxValueDecimal != 0 {
		if request.MinValueDecimal > request.MaxValueDecimal {
			return errors.New("min value cannot be greater than max value")
		}

		if request.MinValueDecimal == request.MaxValueDecimal {
			return errors.New("min value and max value cannot be the same")
		}

		if request.DefaultNumberDecimal != 0 {
			if request.DefaultNumberDecimal < request.MinValueDecimal {
				return errors.New("default value cannot be less than min value")
			}

			if request.DefaultNumberDecimal > request.MaxValueDecimal {
				return errors.New("default value cannot be greater than max value")
			}
		}
	}

	return nil
}

func validateCombobox(request models.CreateFormFieldRequest) error {
	if len(request.Values) == 0 {
		return errors.New("values cannot be empty")
	}

	if len(request.Values) > 10 {
		return errors.New("values cannot be more than 10")
	}

	for _, value := range request.Values {
		if value == "" {
			return errors.New("values cannot be empty")
		}
	}

	uniqueValues := make(map[string]bool)
	for _, value := range request.Values {
		if _, ok := uniqueValues[value]; ok {
			return errors.New("values must be unique")
		}

		uniqueValues[value] = true
	}

	if request.DefaultString != "" {
		found := false // Flag to check if default value is in values
		for _, value := range request.Values {
			if request.DefaultString == value {
				found = true
				break
			}
		}

		if !found {
			return errors.New("default value must be in values")
		}
	}

	return nil
}
