package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Ingredient struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" schema:"-"`
	Name          string             `bson:"name" schema:"name"`
	MeasuringUnit string             `bson:"measuring_unit" schema:"measuring-unit"`
}

type IngredientInRecipe struct {
	IngredientID primitive.ObjectID `bson:"ingredient_id"`
	Amount       float32            `bson:"amount"`
}

type PopulatedIngredientInRecipe struct {
	Ingredient
	IngredientInRecipe
}
