package database

import (
	"Perwatch-case/internal/config"
	"Perwatch-case/internal/helpers"
	"Perwatch-case/internal/models"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository struct {
	Collection *mongo.Collection
}

type UserRepositoryInterface interface {
	Create(user *models.User) error
	FindByUsername(username string) (*models.User, error)
	CheckUsernameExists(username string) (bool, error)
	CheckFirstnameAndLastnameExists(firstname, lastname string) (bool, error)
}

func NewUserRepository() UserRepositoryInterface {
	userCollection := config.GetConfig().GetDatabaseConfig().Collections.Users

	return &UserRepository{
		Collection: GetCollection(userCollection),
	}
}

func (u *UserRepository) Create(user *models.User) error {
	ctx, cancel := helpers.ContextWithTimeout(10)
	defer cancel()

	_, err := u.Collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserRepository) FindByUsername(username string) (*models.User, error) {
	var user *models.User

	ctx, cancel := helpers.ContextWithTimeout(10)
	defer cancel()

	filter := bson.D{{"username", username}}
	if err := u.Collection.FindOne(ctx, filter).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return user, errors.New("user not found")
		}

		return nil, err
	}

	return user, nil
}

func (u *UserRepository) CheckUsernameExists(username string) (bool, error) {
	ctx, cancel := helpers.ContextWithTimeout(10)
	defer cancel()

	filter := bson.D{{"username", username}}
	project := bson.D{{"username", 1}}
	setProjection := options.FindOne().SetProjection(project)

	var result bson.M
	err := u.Collection.FindOne(ctx, filter, setProjection).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (u *UserRepository) CheckFirstnameAndLastnameExists(firstname, lastname string) (bool, error) {
	ctx, cancel := helpers.ContextWithTimeout(10)
	defer cancel()

	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{"firstname", firstname}},
				bson.D{{"lastname", lastname}},
			},
		},
	}

	var result bson.M
	err := u.Collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
