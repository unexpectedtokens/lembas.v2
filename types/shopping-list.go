package types

type ShoppingListIngredient struct {
	Name          string
	MeasuringUnit string
	Amount        float32
}

type ShoppingList struct {
	Ingredients map[string]PopulatedIngredientInRecipe
}
