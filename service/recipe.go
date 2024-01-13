package services

import (
	"context"
	"fmt"

	"github.com/unexpectedtoken/recipes/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RecipeDAO interface {
	DAO[types.Recipe]
	AddIngredientToRecipe(ctx context.Context, ID primitive.ObjectID, ingredient types.IngredientInRecipe) error
	UpdateRecipeIngredient(ctx context.Context, recipeID primitive.ObjectID, ingredient types.IngredientInRecipe) error
	SetInGroceryListStatus(ctx context.Context, ID primitive.ObjectID, status bool) error
	RecipesInGroceryList(ctx context.Context) ([]types.Recipe, error)
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

func (s *RecipeService) GetRecipeByID(ctx context.Context, id primitive.ObjectID, populateIngredients bool) (types.Recipe, error) {
	recipe, err := s.dao.GetSingle(ctx, id)

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

func (s *RecipeService) AddIngredientToRecipe(ctx context.Context, recipe types.Recipe, ingredient types.IngredientInRecipe) error {
	if s.ingredientService.ingredientIDInRecipeIngredients(ingredient.IngredientID, recipe.Ingredients) {
		return fmt.Errorf("error adding ingredient: not unique")
	}

	return s.dao.AddIngredientToRecipe(ctx, recipe.ID, ingredient)
}

func (s *RecipeService) UpdateRecipeIngredient(ctx context.Context, recipeID primitive.ObjectID, ingredient types.IngredientInRecipe) error {
	return s.dao.UpdateRecipeIngredient(ctx, recipeID, ingredient)
}

func (s *RecipeService) UpdateGroceryListStatus(ctx context.Context, recipeID primitive.ObjectID, status bool) error {
	// TODO: Create recipe exists logic to reduce data fetched, we only need to know it exists
	recipe, err := s.dao.GetSingle(ctx, recipeID)

	// No need to do further querying on a non-existing recipe
	if err != nil {
		return err
	}

	if recipe.InGroceryList == status {
		return fmt.Errorf("error updating status: status already assigned")
	}

	return s.dao.SetInGroceryListStatus(ctx, recipeID, status)
}

func (s *RecipeService) RecipesInShoppingList(ctx context.Context) ([]types.Recipe, error) {
	recipes, err := s.dao.RecipesInGroceryList(ctx)

	if err != nil {
		// TODO: Move check to dao
		if err == mongo.ErrNoDocuments {
			return []types.Recipe{}, nil
		}
		return nil, fmt.Errorf("error getting recipes from dao: %w", err)

	}

	return recipes, nil
}
