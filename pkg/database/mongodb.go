package database

import (
	"context"
	"fmt"
	"order_service/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoConnect() (*mongo.Collection, error) {

	mongoConf := config.GetConfig().Mongo

	// Connect to the database.
	clientOption := options.Client().ApplyURI(mongoConf.Url)
	client, err := mongo.Connect(context.Background(), clientOption)
	if err != nil {
		return nil, fmt.Errorf("MongoConnect: %v", err)
	}

	// Check the connection.
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("MongoConnect: %v", err)
	}

	// Create collection
	collection := client.Database(mongoConf.Database).Collection(mongoConf.Collection)
	if err != nil {
		return nil, fmt.Errorf("MongoConnect: %v", err)
	}

	fmt.Println("Connected to db")

	return collection, nil
}
