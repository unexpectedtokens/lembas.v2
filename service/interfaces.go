package services

import (
	"context"

	"github.com/unexpectedtoken/recipes/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DAO[T types.RecipeDocumentModels] interface {
	GetList(ctx context.Context, query interface{}) ([]T, error)
	CreateDocument(ctx context.Context, document T) (insertedID primitive.ObjectID, err error)
	GetSingle(ctx context.Context, id primitive.ObjectID) (T, error)
}
