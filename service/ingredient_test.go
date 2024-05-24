package services

import (
	"testing"

	"github.com/unexpectedtoken/recipes/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestIngredientIDInRecipeIngredients(t *testing.T) {
	ingService := IngredientService{}
	id := primitive.NewObjectID().Hex()

	ingredientList := []types.IngredientInRecipe{
		{
			IngredientID: id,
			Amount:       20,
		},
	}

	if !ingService.ingredientIDInRecipeIngredients(id, ingredientList) {
		t.Errorf("Should have returned true since id is in ingredients list")
	}

	id = primitive.NewObjectID().Hex()

	ingredientList = []types.IngredientInRecipe{
		{
			Amount: 20,
		},
	}

	if ingService.ingredientIDInRecipeIngredients(id, ingredientList) {
		t.Errorf("Should have returned false since id is not in ingredients list")
	}
}
