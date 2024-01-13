package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewDBConn(connString string) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connString))

	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)

	database := client.Database("recipe_base")

	if err != nil {
		return nil, err
	}

	return database, nil
}
