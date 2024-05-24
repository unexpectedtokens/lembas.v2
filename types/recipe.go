package types

type Recipe struct {
	ID                   string                         `schema:"-" bson:"_id,omitempty"`
	Title                string                         `schema:"title,required" bson:"title"`
	Ingredients          []IngredientInRecipe           `schema:"-" bson:"ingredients"`
	PopulatedIngredients *[]PopulatedIngredientInRecipe `bson:"-" schema:"-"`
	InGroceryList        bool                           `bson:"inGroceryList" schema:"-"`
}

func (r *Recipe) GetID() string {
	return r.ID
}
