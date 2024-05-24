package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/unexpectedtoken/recipes/mealplan"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MealplanRepo struct {
	col *mongo.Collection
	MongoDAO[mealplan.MealplanEntryV2]
}

func NewMealplanRepo(db *mongo.Database) *MealplanRepo {
	col := db.Collection("mealplan_entries")
	return &MealplanRepo{
		col: col,
		MongoDAO: MongoDAO[mealplan.MealplanEntryV2]{
			col: col,
		},
	}
}

func (r *MealplanRepo) EntriesForDateWithType(ctx context.Context, fromDate time.Time, toDate time.Time, mealtype mealplan.MealType) (*mealplan.MealplanRecipeDataEntries, error) {
	filter := bson.M{
		"date": bson.M{
			"$gte": fromDate,
			"$lte": toDate,
		},
		"mealtype": mealtype,
	}

	docs, err := r.GetList(ctx, filter)

	if err != nil {
		return nil, fmt.Errorf("error getting mealplan entries: %w", err)
	}
	return (*mealplan.MealplanRecipeDataEntries)(&docs), nil
}

func (r *MealplanRepo) EntriesForDate(ctx context.Context, fromDate time.Time, toDate time.Time) (*mealplan.MealplanRecipeDataEntries, error) {
	filter := bson.M{
		"date": bson.M{
			"$gte": fromDate,
			"$lte": toDate,
		},
	}

	docs, err := r.GetList(ctx, filter)

	if err != nil {
		return nil, fmt.Errorf("error getting mealplan entries: %w", err)
	}
	return (*mealplan.MealplanRecipeDataEntries)(&docs), nil
}

func (r *MealplanRepo) RemoveRecipeEntries(ctx context.Context, recipeID string) error {
	_, err := r.col.DeleteMany(ctx, bson.M{
		"recipe_id": recipeID,
	})

	return err
}
