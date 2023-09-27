package db

import (
	"context"
	"fmt"

	"github.com/saifuljnu/todo/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SetupMongoDB() (*mongo.Client, *mongo.Collection, error) {
	// Define MongoDB connection options.
	clientOptions := options.Client().ApplyURI(config.MongoDB_URI)

	// Initialize MongoDB client.
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}

	// Check the connection to MongoDB.
	err = client.Ping(context.Background(), nil)
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}

	fmt.Println("Connected to MongoDB!")

	// Access a specific database and collection.
	db := client.Database("todoDB")
	collection := db.Collection("todos")

	return client, collection, nil
}
