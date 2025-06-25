package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Define the MongoDB DB name and collections name
const (
	databaseName     = "urldb"
	urlsCollection   = "urls" // collection is similar to table in SQL
	clicksCollection = "clicks"
)

// Global MongoDB client
var client *mongo.Client

// ConnectMongoDB helps to connect MongoDB
func ConnectMongoDB(uri, dbName string) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Failed to create MongoDB client: %v", err)
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the database
	// if err := client.Ping(ctx, nil); err != nil {
	// 	return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	// }
	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		log.Fatalf("Failed to ping MongoDB server: %v", err)
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	fmt.Println("Connected to MongoDB!")

	return client.Database(databaseName), nil
}
