package services

import (
	"context"

	"github.com/unexpectedtoken/recipes/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type genericCrudService[T types.RecipeDocumentModels] struct {
	dao DAO[T]
}

func (s *genericCrudService[T]) GetList(ctx context.Context) ([]T, error) {
	return s.dao.GetList(ctx, bson.M{})
}

func (s *genericCrudService[T]) Create(ctx context.Context, doc T) (primitive.ObjectID, error) {
	return s.dao.CreateDocument(ctx, doc)
}

func (s *genericCrudService[T]) GetByID(ctx context.Context, id primitive.ObjectID) (T, error) {
	return s.dao.GetSingle(ctx, id)
}
