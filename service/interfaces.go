package services

import (
	"context"

	"github.com/unexpectedtoken/recipes/types"
)

type DAO[T types.RecipeDocumentModels] interface {
	GetList(ctx context.Context, query interface{}) ([]T, error)
	CreateDocument(ctx context.Context, document T) (insertedID string, err error)
	GetByID(ctx context.Context, id string) (T, error)
	GetListExclIDS(ctx context.Context, idList *[]string) ([]T, error)
	GetListInclIDS(ctx context.Context, idList *[]string) ([]T, error)
	DeleteByID(ctx context.Context, id string) error
}
