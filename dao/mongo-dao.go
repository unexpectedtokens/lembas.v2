package dao

import (
	"context"
	"fmt"

	"github.com/unexpectedtoken/recipes/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDAO[T types.RecipeDocumentModels] struct {
	col *mongo.Collection
}

func NewMongoDAO[T types.RecipeDocumentModels](db *mongo.Database, collection string) *MongoDAO[T] {
	return &MongoDAO[T]{col: db.Collection(collection)}
}

func (d *MongoDAO[T]) GetList(ctx context.Context, query interface{}) ([]T, error) {
	slic := []T{}
	cursor, err := d.col.Find(ctx, query)

	if err != nil {
		if err == mongo.ErrNilDocument {
			return slic, nil
		}

		return nil, fmt.Errorf("error getting documents in collection %s: %w", d.col.Name(), err)
	}

	err = cursor.All(ctx, &slic)

	if err != nil {
		return nil, fmt.Errorf("error getting cursor values: %w", err)
	}

	return slic, nil
}

func (d *MongoDAO[T]) CreateDocument(ctx context.Context, document T) (insertedID primitive.ObjectID, err error) {
	result, err := d.col.InsertOne(ctx, document)

	if err != nil {
		return insertedID, err
	}

	id, ok := result.InsertedID.(primitive.ObjectID)

	if ok {
		return id, nil
	}

	return insertedID, nil
}

func (d *MongoDAO[T]) GetSingle(ctx context.Context, id primitive.ObjectID) (T, error) {
	var doc T
	err := d.col.FindOne(ctx, bson.M{"_id": id}).Decode(&doc)

	if err != nil {
		return doc, err
	}

	return doc, nil
}
