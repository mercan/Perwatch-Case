package services

import (
	"Perwatch-case/internal/models"
	"Perwatch-case/internal/repositories/database"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type FormService struct {
	FormRepository database.FormRepositoryInterface
}

type FormServiceInterface interface {
	CreateForm(form models.Form) error
	GetForm(userID, formID primitive.ObjectID) (models.Form, error)
	GetForms(userID primitive.ObjectID, page int) ([]models.Form, error)
	DeleteForm(userID, formID primitive.ObjectID) error
	UpdateFormName(userID, formID primitive.ObjectID, formName string) error

	CreateFormField(userID, formID primitive.ObjectID, request models.CreateFormFieldRequest) error
	GetFormFields(userID, formID primitive.ObjectID) (bson.A, error)
	GetFormField(userID, fieldID, formID primitive.ObjectID) (models.Field, error)
	DeleteFormField(userID, fieldID, formID primitive.ObjectID) error
}

func NewFormService() FormServiceInterface {
	return &FormService{
		FormRepository: database.NewFormRepository(),
	}
}

func (f *FormService) CreateForm(form models.Form) error {
	exists, err := f.FormRepository.CheckFormNameExists(form.UserID, form.Name)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("form name already exists")
	}

	if err := f.FormRepository.CreateForm(form); err != nil {
		return err
	}

	return nil
}

func (f *FormService) GetForm(userID, formID primitive.ObjectID) (models.Form, error) {
	return f.FormRepository.GetForm(userID, formID)
}

func (f *FormService) GetForms(userID primitive.ObjectID, page int) ([]models.Form, error) {
	return f.FormRepository.GetForms(userID, page)
}

func (f *FormService) DeleteForm(userID, formID primitive.ObjectID) error {
	return f.FormRepository.DeleteForm(userID, formID)
}

func (f *FormService) UpdateFormName(userID, formID primitive.ObjectID, formName string) error {
	exists, err := f.FormRepository.CheckFormNameExists(userID, formName)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("form name already exists")
	}

	return f.FormRepository.UpdateFormName(userID, formID, formName)
}

func (f *FormService) CreateFormField(userID, formID primitive.ObjectID, request models.CreateFormFieldRequest) error {
	formNameExists, err := f.FormRepository.CheckFormFieldNameExists(userID, formID, request.FieldName)
	if err != nil {
		return err
	}
	if formNameExists {
		return errors.New("form field already exists")
	}

	if request.Sort != 0 {
		formSortExists, err := f.FormRepository.CheckFormFieldSortExists(userID, formID, request.Sort)
		if err != nil {
			return err
		}

		if formSortExists {
			return errors.New("form field sort already exists")
		}
	}

	/*	field := models.Field{
			ID:            primitive.NewObjectID(),
			Name:          request.FieldName,
			Text:          request.Text,
			Number:        request.Number,
			NumberDecimal: request.NumberDecimal,
			Checkbox:      request.Checkbox,
			Combobox:      request.Combobox,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
	*/

	field := models.Field{
		ID:                   primitive.NewObjectID(),
		Name:                 request.FieldName,
		Type:                 request.FieldType,
		Values:               request.Values,
		MinValue:             request.MinValue,
		MaxValue:             request.MaxValue,
		MinValueDecimal:      request.MinValueDecimal,
		MaxValueDecimal:      request.MaxValueDecimal,
		MinLength:            request.MinLength,
		MaxLength:            request.MaxLength,
		DefaultString:        request.DefaultString,
		DefaultNumber:        request.DefaultNumber,
		DefaultNumberDecimal: request.DefaultNumberDecimal,
		Sort:                 request.Sort,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	if err := f.FormRepository.CreateFormField(userID, formID, field); err != nil {
		return err
	}

	return nil
}

func (f *FormService) GetFormFields(userID, formID primitive.ObjectID) (bson.A, error) {
	return f.FormRepository.GetFormFields(userID, formID)
}

func (f *FormService) GetFormField(userID, fieldID, formID primitive.ObjectID) (models.Field, error) {
	return f.FormRepository.GetFormField(userID, fieldID, formID)
}

func (f *FormService) DeleteFormField(userID, fieldID, formID primitive.ObjectID) error {
	return f.FormRepository.DeleteFormField(userID, fieldID, formID)
}
