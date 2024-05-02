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

// Define a list of valid form types (used for validation)
var validFormTypes = []string{"text", "number", "number_decimal", "checkbox", "combobox"}

// FormController struct holds the FormService instance for interacting with form logic
type FormController struct {
	FormService services.FormServiceInterface
}

// NewFormController creates a new FormController instance with the injected FormService
func NewFormController() FormController {
	return FormController{
		FormService: services.NewFormService(),
	}
}

// CreateForm handler creates a new form based on the request body data
func (f *FormController) CreateForm(c *fiber.Ctx) error {
	var request models.FormNameRequestBody

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.NewErrorResponse("Invalid request"))

	}

	// Validate form name presence and length
	if len(request.FormName) < 3 || len(request.FormName) > 50 || request.FormName == "" {
  		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.NewErrorResponse("Invalid form name. Name should be between 3 and 50 characters."))
	}

	// Create a new form object
	form := models.NewForm()
	form.UserID = c.Locals("userID").(primitive.ObjectID)
	form.Name = request.FormName

	// Call FormService to create the form and handle any errors
	if err := f.FormService.CreateForm(form); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(utils.NewSuccessResponse(nil, "Form created successfully"))
}

// GetForm handler retrieves a form based on the provided form ID from the request params
func (f *FormController) GetForm(c *fiber.Ctx) error {
	var request models.FormFieldRequestParams

	if err := c.ParamsParser(&request); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.NewErrorResponse("Invalid request"))
	}

	// Validate form ID presence
	if request.FormID == "" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.NewErrorResponse("Form ID is required"))
	}

	// Convert form ID string to a MongoDB ObjectID
	formID, err := primitive.ObjectIDFromHex(request.FormID)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.NewErrorResponse("Invalid form ID"))
	}

	// Get user ID from context
	userID := c.Locals("userID").(primitive.ObjectID)
	
	// Call FormService to retrieve the form and handle any errors
	form, err := f.FormService.GetForm(userID, formID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(utils.NewSuccessResponse(form, "Form retrieved successfully"))
}

// GetForms handler retrieves a list of forms for the current user based
func (f *FormController) GetForms(c *fiber.Ctx) error {
	var request models.GetFormsRequest

	if err := c.QueryParser(&request); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.NewErrorResponse("Invalid request"))
	}

	// Validate that the requested page number is greater than or equal to 1
	if request.Page < 1 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.NewErrorResponse("Invalid page number"))
	}
	
	// Get user ID from context
	userID := c.Locals("userID").(primitive.ObjectID)
	
	// Call FormService to retrieve the list of forms for the user and handle any errors
	formList, err := f.FormService.GetForms(userID, request.Page)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(utils.NewSuccessResponse(formList, "Forms retrieved successfully"))
}

// DeleteForm handler deletes a form based on the provided form ID from the request params
func (f *FormController) DeleteForm(c *fiber.Ctx) error {
	var request models.FormIDRequestBody

	if err := c.ParamsParser(&request); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.NewErrorResponse("Invalid request"))
	}

	// Validate that the form ID is present
	if request.FormID == "" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.NewErrorResponse("Form ID is required"))
	}

	// Convert the form ID string to a MongoDB ObjectID
	formID, err := primitive.ObjectIDFromHex(request.FormID)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.NewErrorResponse("Invalid form ID"))
	}

	// Get the user ID from the context
	userID := c.Locals("userID").(primitive.ObjectID)

	// Call FormService to delete the form and handle any errors
	if err := f.FormService.DeleteForm(userID, formID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Failed to delete form"))
	}

	return c.Status(fiber.StatusOK).JSON(utils.NewSuccessResponse(nil, "Form deleted successfully"))
}

func (f *FormController) UpdateFormName(c *fiber.Ctx) error {
	// Parse request parameters for form ID
	var requestParams models.FormFieldRequestParams
	if err := c.ParamsParser(&requestParams); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.NewErrorResponse("Invalid request"))
	}

	// Parse request body for new form name
	var requestBody models.FormNameRequestBody
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.NewErrorResponse("Invalid request"))
	}

	// Validate that both form ID and new form name are provided
	if requestParams.FormID == "" || requestBody.FormName == "" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.NewErrorResponse("Form ID and form name are required"))
	}

	// Convert form ID string to a MongoDB ObjectID
	formID, err := primitive.ObjectIDFromHex(requestParams.FormID)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.NewErrorResponse("Invalid form ID"))
	}

	// Get user ID from context
	userID := c.Locals("userID").(primitive.ObjectID)
	
	// Call FormService to update the form name and handle any errors
	if err := f.FormService.UpdateFormName(userID, formID, requestBody.FormName); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(utils.NewSuccessResponse(nil, "Form name updated successfully"))
}

// CreateFormField handler creates a new form field for a specific form
func (f *FormController) CreateFormField(c *fiber.Ctx) error {
	// Parse request parameters for form ID
	var requestParams models.FormFieldRequestParams
	if err := c.ParamsParser(&requestParams); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Invalid request"))
	}

	// Parse request body for form field details
	var request models.CreateFormFieldRequest
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

	// Validate that the field type is provided
	if request.FieldType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Field type is required"))
	}

	// Validate that the field type is a valid option from the list
	if !slices.Contains(validFormTypes, request.FieldType) {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Invalid field type"))
	}

	// Validate that the sort order is greater than 0
	if request.Sort < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Sort must be greater than 0"))
	}

	// Perform specific field-type validations based on the chosen type
	switch request.FieldType {
	case "text":
		if err := validateText(request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse(err.Error()))
		}
	case "number":
		if err := validateNumber(request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse(err.Error()))
		}
	case "number_decimal":
		if err := validateNumberDecimal(request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse(err.Error()))
		}
	case "combobox":
		if err := validateCombobox(request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse(err.Error()))
		}
	default:
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("invalid field type")
	}

	// Call FormService to create the form field and handle any errors
	if err := f.FormService.CreateFormField(userID, formID, request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse(err.Error()))
	}

	return c.Status(fiber.StatusCreated).JSON(utils.NewSuccessResponse(nil, "Form field created successfully"))
}

// GetFormFields handler retrieves all form fields for a specific form
func (f *FormController) GetFormFields(c *fiber.Ctx) error {
	// Parse request parameters for form ID
	var requestParams models.FormFieldRequestParams
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

	// Call FormService to retrieve the form fields for the specified form and handle any errors
	formFields, err := f.FormService.GetFormFields(userID, formID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(utils.NewSuccessResponse(formFields, "Form fields retrieved successfully"))
}


// GetFormField handler retrieves a specific form field for a form
func (f *FormController) GetFormField(c *fiber.Ctx) error {
	// Parse request parameters for form ID and field ID
	var requestParams models.FormFieldRequestParams
	if err := c.ParamsParser(&requestParams); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Invalid request"))
	}

	// Get user ID from context
	userID := c.Locals("userID").(primitive.ObjectID)

	// Convert form ID and field ID strings to MongoDB ObjectIDs
	formID, err := primitive.ObjectIDFromHex(requestParams.FormID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Invalid form ID"))
	}
	fieldID, err := primitive.ObjectIDFromHex(requestParams.FieldID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Invalid field ID"))
	}

	// Call FormService to retrieve the specific form field for the given IDs and user
	formField, err := f.FormService.GetFormField(userID, fieldID, formID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(utils.NewSuccessResponse(formField, "Form field retrieved successfully"))
}


// DeleteFormField handler deletes a specific form field from a form
func (f *FormController) DeleteFormField(c *fiber.Ctx) error {
	// Parse request parameters for form ID and field ID
	var requestParams models.FormFieldRequestParams
	if err := c.ParamsParser(&requestParams); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Invalid request"))
	}

	// Get user ID from context
	userID := c.Locals("userID").(primitive.ObjectID)

	// Convert form ID and field ID strings to MongoDB ObjectIDs
	formID, err := primitive.ObjectIDFromHex(requestParams.FormID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Invalid form ID"))
	}
	fieldID, err := primitive.ObjectIDFromHex(requestParams.FieldID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse("Invalid field ID"))
	}

	// Call FormService to delete the specified form field and handle any errors
	if err := f.FormService.DeleteFormField(userID, fieldID, formID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewErrorResponse(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(utils.NewSuccessResponse(nil, "Form field deleted successfully"))
}


// validateText validates the min and max length properties for a text field
func validateText(request models.CreateFormFieldRequest) error {
  // Ensure minimum and maximum lengths are non-negative
	if request.MinLength < 0 || request.MaxLength < 0 {
		return errors.New("min and max length cannot be negative")
	}

	// Validate min and max length only if both are set (not zero)
	if request.MinLength != 0 && request.MaxLength != 0 {
		// Check if minimum length is greater than maximum length
		if request.MinLength > request.MaxLength {
			return errors.New("min length cannot be greater than max length")
		}
		// Check if minimum and maximum lengths are the same
		if request.MinLength == request.MaxLength {
			return errors.New("min and max length cannot be the same")
		}
	}

	// Return nil if all validations pass
	return nil
}

// validateNumber validates the min, max, and default values for a number field
func validateNumber(request models.CreateFormFieldRequest) error {
	// Ensure minimum and maximum values are non-negative
	if request.MinValue < 0 || request.MaxValue < 0 {
		return errors.New("min and max value cannot be negative")
	}

	// Validate min and max values only if both are set (not zero)
	if request.MinValue != 0 && request.MaxValue != 0 {
		// Check if minimum value is greater than maximum value
		if request.MinValue > request.MaxValue {
			return errors.New("min value cannot be greater than max value")
		}
		// Check if minimum and maximum values are the same
		if request.MinValue == request.MaxValue {
			return errors.New("min and max value cannot be the same")
		}

		// Validate default value if provided (not zero)
		if request.DefaultNumber != 0 {
			// Check if default value is less than minimum value
			if request.DefaultNumber < request.MinValue {
				return errors.New("default value cannot be less than min value")
			}
			// Check if default value is greater than maximum value
			if request.DefaultNumber > request.MaxValue {
				return errors.New("default value cannot be greater than max value")
			}
		}
	}

	// Return nil if all validations pass
	return nil
}


// validateNumberDecimal validates the min, max, and default values for a decimal number field
func validateNumberDecimal(request models.CreateFormFieldRequest) error {
	// Ensure minimum and maximum values are non-negative
	if request.MinValueDecimal < 0 || request.MaxValueDecimal < 0 {
		return errors.New("min and max value cannot be negative")
	}

	// Validate min and max values only if both are set (not zero)
	if request.MinValueDecimal != 0 && request.MaxValueDecimal != 0 {
		// Check if minimum value is greater than maximum value
		if request.MinValueDecimal > request.MaxValueDecimal {
			return errors.New("min value cannot be greater than max value")
		}
		// Check if minimum and maximum values are the same
		if request.MinValueDecimal == request.MaxValueDecimal {
			return errors.New("min value and max value cannot be the same")
		}

		// Validate default value if provided (not zero)
		if request.DefaultNumberDecimal != 0 {
			// Check if default value is less than minimum value
			if request.DefaultNumberDecimal < request.MinValueDecimal {
				return errors.New("default value cannot be less than min value")
			}
			// Check if default value is greater than maximum value
			if request.DefaultNumberDecimal > request.MaxValueDecimal {
				return errors.New("default value cannot be greater than max value")
			}
		}
	}

	// Return nil if all validations pass
	return nil
}


// validateCombobox validates the options for a combobox field
func validateCombobox(request models.CreateFormFieldRequest) error {
	// Check if there are any options provided
	if len(request.Values) == 0 {
		return errors.New("values cannot be empty")
	}

	// Enforce a maximum of 10 options
	if len(request.Values) > 10 {
		return errors.New("values cannot be more than 10")
	}

	// Ensure no empty option values
	for _, value := range request.Values {
		if value == "" {
			return errors.New("values cannot be empty")
		}
	}

	// Validate for unique options (no duplicates)
	uniqueValues := make(map[string]bool)
	for _, value := range request.Values {
		if _, ok := uniqueValues[value]; ok {
			return errors.New("values must be unique")
		}

		uniqueValues[value] = true
	}

	// Validate default value if provided
	if request.DefaultString != "" {
		found := false // Flag to track if default value exists in options
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

	// Return nil if all validations pass
	return nil
}

