package database

import (
	"context"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func NewMongoDB(uri string) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// Ping the database to check if the connection is successful
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Extract database name frm URI  or use default
	dbName := "image-processing-service"
	if uri != "" {
		dbName = uri[strings.LastIndex(uri, "/")+1:]
	}

	return &MongoDB{
		Client:   client,
		Database: client.Database(dbName),
	}, nil
}

// Disconnect disconnects from the database
func (m *MongoDB) Disconnect(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	return m.Client.Disconnect(ctx)
}

// Collection returns a collection from the database
func (m *MongoDB) Collection(name string) *mongo.Collection {
	return m.Database.Collection(name)
}
