package shopping

import (
	"github.com/unexpectedtoken/recipes/mealplan"
	"github.com/unexpectedtoken/recipes/types"
)

type ShoppingList struct {
	DateRange *mealplan.DateRange

	Ingredients map[string]types.PopulatedIngredientInRecipe
}

func (s *ShoppingList) HasIngredient(ingredientID string) bool {
	present := false

	for _, ingredient := range s.Ingredients {
		if ingredient.IngredientID == ingredientID {
			present = true
			break
		}
	}

	return present
}
