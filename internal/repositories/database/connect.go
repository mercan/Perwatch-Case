package database

import (
	"Perwatch-case/internal/config"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
)

var client = Connect()

func Connect() *mongo.Client {
	databaseURI := config.GetConfig().GetDatabaseConfig().URI
	if databaseURI == "" {
		log.Panic("MongoDB URI not found")
		return nil
	}

	clientOptions := options.Client().ApplyURI(databaseURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Panicf("MongoDB connection error: %v", err)
	}

	if err := client.Ping(context.Background(), readpref.Primary()); err != nil {
		log.Panicf("MongoDB ping error: %v", err)
	}

	if err := createUserIndexes(client); err != nil {
		log.Panicf("MongoDB create user indexes error: %v", err)
	}

	log.Println("Connected to MongoDB")
	return client
}

func createUserIndexes(client *mongo.Client) error {
	databaseName := config.GetConfig().GetDatabaseConfig().Name
	usersCollectionName := config.GetConfig().GetDatabaseConfig().Collections.Users

	collection := client.Database(databaseName).Collection(usersCollectionName)
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{"username", 1}},
		Options: options.Index().SetUnique(true),
	}

	if _, err := collection.Indexes().CreateOne(context.Background(), indexModel); err != nil {
		return err
	}

	return nil
}

func GetCollection(collectionName string) *mongo.Collection {
	databaseName := config.GetConfig().GetDatabaseConfig().Name

	return client.Database(databaseName).Collection(collectionName)
}
