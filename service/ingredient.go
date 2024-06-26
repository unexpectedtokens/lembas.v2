package services

import (
	"context"

	"github.com/unexpectedtoken/recipes/types"
)

type IngredientDAO interface {
}

type IngredientService struct {
	dao DAO[types.Ingredient]
	genericCrudService[types.Ingredient]
}

func NewIngredientService(dao DAO[types.Ingredient]) *IngredientService {
	return &IngredientService{
		dao: dao,
		genericCrudService: genericCrudService[types.Ingredient]{
			dao: dao,
		},
	}
}

func (s *IngredientService) idsFromImplementedIngredients(implementedIngredients []types.IngredientInRecipe) []string {
	allIds := []string{}

	for _, ingr := range implementedIngredients {
		allIds = append(allIds, ingr.IngredientID)
	}

	return allIds
}

func (s *IngredientService) ingredientIDInRecipeIngredients(ID string, ingredients []types.IngredientInRecipe) bool {
	inList := false

	for _, ingredient := range ingredients {
		if ingredient.IngredientID == ID {
			inList = true
		}
	}

	return inList
}

func (s *IngredientService) matchIngredientsToImplementation(ingredients []types.Ingredient, implementedIngredients []types.IngredientInRecipe) *[]types.PopulatedIngredientInRecipe {
	popIngredients := []types.PopulatedIngredientInRecipe{}
	for _, impIngredient := range implementedIngredients {
		popIngredient := types.PopulatedIngredientInRecipe{}
		found := false
		for _, ingredient := range ingredients {
			if impIngredient.IngredientID == ingredient.ID {
				popIngredient.Ingredient = ingredient
				popIngredient.IngredientInRecipe = impIngredient
				found = true
			}
		}

		if found {
			popIngredients = append(popIngredients, popIngredient)
		}
	}
	return &popIngredients
}

func (s *IngredientService) PopulateIngredientList(ctx context.Context, implementedIngredients []types.IngredientInRecipe) (*[]types.PopulatedIngredientInRecipe, error) {
	idsInRecipe := s.idsFromImplementedIngredients(implementedIngredients)
	ingredients, err := s.dao.GetListInclIDS(ctx, &idsInRecipe)

	if err != nil {
		return nil, err
	}

	return s.matchIngredientsToImplementation(ingredients, implementedIngredients), nil
}

func (s *IngredientService) GetIngredientsNotInList(ctx context.Context, ingredients []types.IngredientInRecipe) ([]types.Ingredient, error) {
	ids := []string{}

	for _, ingr := range ingredients {
		ids = append(ids, ingr.IngredientID)
	}

	return s.dao.GetListExclIDS(ctx, &ids)
}
