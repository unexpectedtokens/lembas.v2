package repository

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

func parseId(id string) (primitive.ObjectID, error) {
	objId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return objId, ErrInvalidId
	}

	return objId, nil
}

func parseIDList(ids *[]string) (*[]primitive.ObjectID, error) {
	objIds := []primitive.ObjectID{}

	var parseError error
	for _, id := range *ids {
		objID, err := primitive.ObjectIDFromHex(id)

		if err != nil {
			parseError = err
			break
		}

		objIds = append(objIds, objID)
	}

	if parseError != nil {
		return nil, parseError
	}

	return &objIds, nil
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

func (d *MongoDAO[T]) GetListExclIDS(ctx context.Context, list *[]string) ([]T, error) {
	objIds, err := parseIDList(list)

	if err != nil {
		return nil, err
	}

	return d.GetList(ctx, bson.M{
		"_id": bson.M{
			"$nin": objIds,
		},
	})
}

func (d *MongoDAO[T]) GetListInclIDS(ctx context.Context, list *[]string) ([]T, error) {
	objIds, err := parseIDList(list)

	if err != nil {
		return nil, err
	}

	return d.GetList(ctx, bson.M{
		"_id": bson.M{
			"$in": objIds,
		},
	})
}

func (d *MongoDAO[T]) CreateDocument(ctx context.Context, document T) (insertedID string, err error) {
	result, err := d.col.InsertOne(ctx, document)

	if err != nil {
		return insertedID, err
	}

	id, ok := result.InsertedID.(primitive.ObjectID)

	if ok {
		return id.Hex(), nil
	}

	return "", nil
}

func (d *MongoDAO[T]) GetByID(ctx context.Context, id string) (T, error) {
	var doc T

	objID, err := parseId(id)

	if err != nil {
		return doc, err
	}

	err = d.col.FindOne(ctx, bson.M{"_id": objID}).Decode(&doc)

	if err != nil {
		return doc, err
	}

	return doc, nil
}

func (d *MongoDAO[T]) DeleteByID(ctx context.Context, id string) error {
	objID, err := parseId(id)

	if err != nil {
		return err
	}
	_, err = d.col.DeleteOne(ctx, bson.M{"_id": objID})

	return err
}
