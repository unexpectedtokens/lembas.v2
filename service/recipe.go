package services

import (
	"context"
	"fmt"

	"github.com/unexpectedtoken/recipes/types"
)

type RecipeDAO interface {
	DAO[types.Recipe]
	AddIngredientToRecipe(ctx context.Context, ID string, ingredient types.IngredientInRecipe) error
	UpdateRecipeIngredient(ctx context.Context, recipeID string, ingredient types.IngredientInRecipe) error
	DeleteIngredientFromRecipe(ctx context.Context, ID string, ingredientID string) error
}

type RecipeService struct {
	dao RecipeDAO
	genericCrudService[types.Recipe]
	ingredientService *IngredientService
}

func NewRecipeService(dao RecipeDAO, ingredientsService *IngredientService) *RecipeService {

	return &RecipeService{
		dao: dao,
		genericCrudService: genericCrudService[types.Recipe]{
			dao: dao,
		},
		ingredientService: ingredientsService,
	}
}

func (s *RecipeService) GetRecipeByID(ctx context.Context, id string, populateIngredients bool) (types.Recipe, error) {
	recipe, err := s.dao.GetByID(ctx, id)

	if err != nil {
		return recipe, err
	}

	if populateIngredients {
		recipe.PopulatedIngredients, err = s.ingredientService.PopulateIngredientList(ctx, recipe.Ingredients)

		if err != nil {
			return recipe, err
		}
	}

	return recipe, nil
}

func (s *RecipeService) GetRecipesInIDList(ctx context.Context, ids *[]string) ([]types.Recipe, error) {
	return s.dao.GetListInclIDS(ctx, ids)
}

func (s *RecipeService) GetPopulatedRecipesInIDList(ctx context.Context, ids *[]string) ([]types.Recipe, error) {
	recipes, err := s.dao.GetListInclIDS(ctx, ids)

	if err != nil {
		return nil, err
	}

	for i, _ := range recipes {
		recipes[i].PopulatedIngredients, err = s.ingredientService.PopulateIngredientList(ctx, recipes[i].Ingredients)
	}

	return recipes, err
}

func (s *RecipeService) GetRecipesNotInIDList(ctx context.Context, ids *[]string) ([]types.Recipe, error) {
	return s.dao.GetListExclIDS(ctx, ids)
}

func (s *RecipeService) AddIngredientToRecipe(ctx context.Context, recipe types.Recipe, ingredientID string) error {
	if s.ingredientService.ingredientIDInRecipeIngredients(ingredientID, recipe.Ingredients) {
		return fmt.Errorf("error adding ingredient: not unique")
	}

	ingredient := types.IngredientInRecipe{
		IngredientID: ingredientID,
	}

	return s.dao.AddIngredientToRecipe(ctx, recipe.ID, ingredient)
}

func (s *RecipeService) UpdateRecipeIngredient(ctx context.Context, recipeID string, ingredient types.IngredientInRecipe) error {
	return s.dao.UpdateRecipeIngredient(ctx, recipeID, ingredient)
}

func (s *RecipeService) RemoveRecipe(ctx context.Context, recipeID string) error {
	return s.dao.DeleteByID(ctx, recipeID)
}

func (s *RecipeService) RemoveIngredientFromRecipe(ctx context.Context, recipeID, ingredientID string) error {

	return s.dao.DeleteIngredientFromRecipe(ctx, recipeID, ingredientID)
}

func (s *RecipeService) AddNewIngredientToRecipe(ctx context.Context, recipeID string, ingredient types.Ingredient) error {
	id, err := s.ingredientService.Create(ctx, ingredient)

	if err != nil {
		return err
	}
	ingredientInRecipe := types.IngredientInRecipe{
		IngredientID: id,
	}

	return s.dao.AddIngredientToRecipe(ctx, recipeID, ingredientInRecipe)
}

func (s *RecipeService) MapRecipes(recipes *[]types.Recipe) *map[string]types.Recipe {
	out := make(map[string]types.Recipe)

	fmt.Println(recipes)

	for _, recipe := range *recipes {
		if _, ok := out[recipe.ID]; !ok {
			out[recipe.ID] = recipe
		}
	}

	return &out
}
