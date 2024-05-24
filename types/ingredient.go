package types

type IngredientID string

type Ingredient struct {
	ID            string `bson:"_id,omitempty" schema:"-"`
	Name          string `bson:"name" schema:"name"`
	MeasuringUnit string `bson:"measuring_unit" schema:"measuring-unit"`
}

func (i *Ingredient) GetID() string {
	return i.ID
}

type IngredientInRecipe struct {
	IngredientID string  `bson:"ingredient_id"`
	Amount       float32 `bson:"amount"`
}

type PopulatedIngredientInRecipe struct {
	Ingredient
	IngredientInRecipe
}
