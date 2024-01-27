package repository

import (
	"github.com/unexpectedtoken/recipes/types"
	"go.mongodb.org/mongo-driver/mongo"
)

type MealplanRepo struct {
	col *mongo.Collection
	MongoDAO[types.MealplanEntry]
}

func NewMealplanRepo(db *mongo.Database) *MealplanRepo {
	col := db.Collection("mealplan_entries")
	return &MealplanRepo{
		col: col,
		MongoDAO: MongoDAO[types.MealplanEntry]{
			col: col,
		},
	}
}
