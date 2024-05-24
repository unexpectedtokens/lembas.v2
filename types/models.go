package types

import "github.com/unexpectedtoken/recipes/mealplan"

type RecipeDocumentModels interface {
	Recipe | Ingredient | mealplan.MealplanEntry | mealplan.MealplanEntryV2
}
