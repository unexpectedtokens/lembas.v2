package services

import (
	"context"

	dao "github.com/unexpectedtoken/recipes/repository"
	"github.com/unexpectedtoken/recipes/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

func (s *IngredientService) idsFromImplementedIngredients(implementedIngredients []types.IngredientInRecipe) []primitive.ObjectID {
	allIds := []primitive.ObjectID{}

	for _, ingr := range implementedIngredients {
		allIds = append(allIds, ingr.IngredientID)
	}

	return allIds
}

func (s *IngredientService) ingredientIDInRecipeIngredients(ID primitive.ObjectID, ingredients []types.IngredientInRecipe) bool {
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
			if impIngredient.IngredientID.Hex() == ingredient.ID.Hex() {
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
	ingredients, err := s.dao.GetList(ctx, dao.IngredientsInListQuery(idsInRecipe))

	if err != nil {
		return nil, err
	}

	return s.matchIngredientsToImplementation(ingredients, implementedIngredients), nil
}

func (s *IngredientService) GetIngredientsNotInList(ctx context.Context, ingredients []types.IngredientInRecipe) ([]types.Ingredient, error) {
	ids := []primitive.ObjectID{}

	for _, ingr := range ingredients {
		ids = append(ids, ingr.IngredientID)
	}

	return s.dao.GetList(ctx, dao.GetIngredientsNotInList(ids))
}
