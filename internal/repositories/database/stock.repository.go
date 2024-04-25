package database

import (
	"Perwatch-case/internal/config"
	"Perwatch-case/internal/helpers"
	"Perwatch-case/internal/models"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type StockRepository struct {
	Collection *mongo.Collection
}

type StockRepositoryInterface interface {
	CreateStock(userID, formID primitive.ObjectID, stock []models.Stock) error
	GetStocks(userID, formID primitive.ObjectID) ([]models.Stock, error)
	GetStock(userID, formID, stockID primitive.ObjectID) (models.Stock, error)
	DeleteStock(userID, formID, stockID primitive.ObjectID) error
	UpdateStockValue(userID, formID, stockID primitive.ObjectID, value interface{}) error
}

func NewStockRepository() StockRepositoryInterface {
	return &StockRepository{
		Collection: GetCollection(config.GetConfig().GetDatabaseConfig().Collections.Forms),
	}
}

func (f *StockRepository) CreateStock(userID, formID primitive.ObjectID, stock []models.Stock) error {
	ctx, cancel := helpers.ContextWithTimeout(10)
	defer cancel()

	filter := bson.D{{"_id", formID}, {"user_id", userID}, {"is_deleted", false}}
	update := bson.D{{"$push", bson.D{{"stocks", bson.D{{"$each", stock}}}}}}

	if _, err := f.Collection.UpdateOne(ctx, filter, update); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("form not found")
		}

		return err
	}

	return nil
}

func (f *StockRepository) GetStocks(userID, formID primitive.ObjectID) ([]models.Stock, error) {
	ctx, cancel := helpers.ContextWithTimeout(10)
	defer cancel()

	filter := bson.D{{"_id", formID}, {"user_id", userID}, {"is_deleted", false}}
	projection := bson.D{{"stocks", 1}}
	findOneOptions := options.FindOne().SetProjection(projection)

	var form models.Form
	if err := f.Collection.FindOne(ctx, filter, findOneOptions).Decode(&form); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("stocks not found")
		}

		return nil, err
	}

	var filteredStocks []models.Stock
	for _, stock := range form.Stock {
		if !stock.IsDeleted {
			filteredStocks = append(filteredStocks, stock)
		}
	}

	return filteredStocks, nil
}

func (f *StockRepository) GetStock(userID, formID, stockID primitive.ObjectID) (models.Stock, error) {
	ctx, cancel := helpers.ContextWithTimeout(10)
	defer cancel()

	filter := bson.D{
		{
			"$and",
			bson.A{
				bson.D{{"_id", formID}},
				bson.D{{"user_id", userID}},
				bson.D{{"stocks._id", stockID}},
				bson.D{{"is_deleted", false}},
			},
		},
	}
	projection := bson.D{{"_id", 0}, {"stocks.$", 1}}
	findOneOptions := options.FindOne().SetProjection(projection)

	var form models.Form
	if err := f.Collection.FindOne(ctx, filter, findOneOptions).Decode(&form); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Stock{}, errors.New("stock not found")
		}

		return models.Stock{}, err
	}

	stock := form.Stock[0]

	if stock.IsDeleted {
		return models.Stock{}, errors.New("stock not found")
	}

	return stock, nil
}

func (f *StockRepository) DeleteStock(userID, formID, stockID primitive.ObjectID) error {
	ctx, cancel := helpers.ContextWithTimeout(10)
	defer cancel()

	filter := bson.D{{"_id", formID}, {"user_id", userID}, {"stocks._id", stockID}, {"is_deleted", false}}
	update := bson.D{{"$set", bson.D{{"stocks.$.is_deleted", true}}}}

	result, err := f.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("stock not found")
		}
	}

	if result.ModifiedCount == 0 {
		return errors.New("stock not found")
	}

	return nil
}

func (f *StockRepository) UpdateStockValue(userID, formID, stockID primitive.ObjectID, value interface{}) error {
	ctx, cancel := helpers.ContextWithTimeout(10)
	defer cancel()

	filter := bson.D{{"_id", formID}, {"user_id", userID}, {"stocks._id", stockID}, {"is_deleted", false}, {"stocks.is_deleted", false}}
	update := bson.D{{"$set", bson.D{{"stocks.$.value", value}}}}

	result, err := f.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("stock not found")
		}
	}

	if result.ModifiedCount == 0 {
		return errors.New("stock not found")
	}

	return nil
}
