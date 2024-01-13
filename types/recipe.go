package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Recipe struct {
	ID                   primitive.ObjectID             `schema:"-" bson:"_id,omitempty"`
	Title                string                         `schema:"title,required" bson:"title"`
	Ingredients          []IngredientInRecipe           `schema:"-" bson:"ingredients"`
	PopulatedIngredients *[]PopulatedIngredientInRecipe `bson:"skip" schema:"skip"`
	InGroceryList        bool                           `bson:"inGroceryList" schema:"-"`
}
