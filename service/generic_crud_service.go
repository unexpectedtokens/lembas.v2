package services

import (
	"context"

	"github.com/unexpectedtoken/recipes/types"
	"go.mongodb.org/mongo-driver/bson"
)

type genericCrudService[T types.RecipeDocumentModels] struct {
	dao DAO[T]
}

func (s *genericCrudService[T]) GetList(ctx context.Context) ([]T, error) {
	return s.dao.GetList(ctx, bson.M{})
}

func (s *genericCrudService[T]) Create(ctx context.Context, doc T) (string, error) {
	return s.dao.CreateDocument(ctx, doc)
}

func (s *genericCrudService[T]) GetByID(ctx context.Context, id string) (T, error) {
	return s.dao.GetByID(ctx, id)
}

func (s *genericCrudService[T]) DeleteByID(ctx context.Context, id string) error {
	return s.dao.DeleteByID(ctx, id)
}
