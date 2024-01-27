package services

import (
	"context"
	"time"

	"github.com/unexpectedtoken/recipes/types"
)

type MealplanRepo interface {
	DAO[types.MealplanEntry]
}

type MealplanService struct {
	repo MealplanRepo
	genericCrudService[types.MealplanEntry]
}

func NewMealplanService(repo MealplanRepo) *MealplanService {
	return &MealplanService{
		repo: repo,
		genericCrudService: genericCrudService[types.MealplanEntry]{
			dao: repo,
		},
	}
}

func (s *MealplanService) EntriesForDates(ctx context.Context, from time.Time, to time.Time) []types.MealplanEntry {
	return nil
}
