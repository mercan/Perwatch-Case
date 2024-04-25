package services

import (
	"Perwatch-case/internal/models"
	"Perwatch-case/internal/repositories/database"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"slices"
	"strconv"
	"time"
)

type StockService struct {
	FormRepository  database.FormRepositoryInterface
	StockRepository database.StockRepositoryInterface
}

type StockServiceInterface interface {
	CreateStock(userID, formID primitive.ObjectID, request models.CreateStockRequest) error
	GetStocks(userID, formID primitive.ObjectID) ([]models.Stock, error)
	GetStock(userID, formID, stockID primitive.ObjectID) (models.Stock, error)
	DeleteStock(userID, formID, stockID primitive.ObjectID) error
	UpdateStockValue(userID, formID, stockID primitive.ObjectID, value interface{}) error
}

func NewStockService() StockServiceInterface {
	return &StockService{
		FormRepository:  database.NewFormRepository(),
		StockRepository: database.NewStockRepository(),
	}
}

func (s StockService) CreateStock(userID, formID primitive.ObjectID, request models.CreateStockRequest) error {
	var stocks []models.Stock

	for i := 0; i < len(request.Fields); i++ {
		exists, err := s.FormRepository.CheckFormFieldNameExists(userID, formID, request.Fields[i].Name)
		if err != nil {
			return err
		}

		if !exists {
			return errors.New("form field does not exist")
		}
	}

	for i := 0; i < len(request.Fields); i++ {
		fields, err := s.FormRepository.FindByFieldName(userID, formID, request.Fields[i].Name)
		if err != nil {
			return err
		}

		for j := 0; j < len(fields); j++ {
			if fields[j].Name == request.Fields[i].Name {
				if fields[j].Type == "text" {
					if fields[j].MinLength != 0 {
						if len(request.Fields[i].Value.(string)) < fields[j].MinLength {
							return errors.New("value is too short")
						}
					}

					if fields[j].MaxLength != 0 {
						if len(request.Fields[i].Value.(string)) > fields[j].MaxLength {
							return errors.New("value is too long")
						}
					}
				}

				if fields[j].Type == "number" {
					value, err := strconv.Atoi(request.Fields[i].Value.(string))
					if err != nil {
						return errors.New("invalid number")
					}

					if fields[j].MinValue != 0 {
						if value < fields[j].MinValue {
							return errors.New("value is too small")
						}
					}

					if fields[j].MaxValue != 0 {
						if value > fields[j].MaxValue {
							return errors.New("value is too large")
						}
					}

					request.Fields[i].Value = value
				}

				if fields[j].Type == "number_decimal" {
					value, err := strconv.ParseFloat(strconv.Itoa(int(request.Fields[i].Value.(float64))), 64)
					if err != nil {
						return errors.New("invalid number")
					}

					if fields[j].MinValue != 0 {
						if value < fields[j].MinValueDecimal {
							return errors.New("value is too small")
						}
					}

					if fields[j].MaxValue != 0 {
						if value > fields[j].MaxValueDecimal {
							return errors.New("value is too large")
						}
					}

					request.Fields[i].Value = value
				}

				if fields[j].Type == "checkbox" {
					if request.Fields[i].Value != true && request.Fields[i].Value != false {
						return errors.New("invalid value")
					}
				}

				if fields[j].Type == "combobox" {
					value := request.Fields[i].Value.(string)
					values := fields[j].Values

					if !slices.Contains(values, value) {
						return errors.New("invalid value for combobox")
					}

				}
			}
		}

	}

	for i := 0; i < len(request.Fields); i++ {
		stocks = append(stocks, models.Stock{
			ID:        primitive.NewObjectID(),
			Name:      request.Fields[i].Name,
			Value:     request.Fields[i].Value,
			IsDeleted: false,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}

	if err := s.StockRepository.CreateStock(userID, formID, stocks); err != nil {
		return err
	}

	return nil
}

func (s StockService) GetStocks(userID, formID primitive.ObjectID) ([]models.Stock, error) {
	return s.StockRepository.GetStocks(userID, formID)
}

func (s StockService) GetStock(userID, formID, stockID primitive.ObjectID) (models.Stock, error) {
	return s.StockRepository.GetStock(userID, formID, stockID)
}

func (s StockService) DeleteStock(userID, formID, stockID primitive.ObjectID) error {
	return s.StockRepository.DeleteStock(userID, formID, stockID)
}

func (s StockService) UpdateStockValue(userID, formID, stockID primitive.ObjectID, value interface{}) error {
	return s.StockRepository.UpdateStockValue(userID, formID, stockID, value)
}
