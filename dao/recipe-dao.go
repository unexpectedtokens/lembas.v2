package dao

import (
	"context"
	"fmt"

	"github.com/unexpectedtoken/recipes/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RecipeDAO struct {
	col *mongo.Collection
	MongoDAO[types.Recipe]
}

func NewRecipeDAO(db *mongo.Database) *RecipeDAO {
	col := db.Collection("recipes")
	return &RecipeDAO{
		col: col,
		MongoDAO: MongoDAO[types.Recipe]{
			col: col,
		},
	}
}

func (d *RecipeDAO) RecipesInGroceryList(ctx context.Context) ([]types.Recipe, error) {
	return d.GetList(ctx, bson.M{"inGroceryList": true})
}

func (d *RecipeDAO) AddIngredientToRecipe(ctx context.Context, ID primitive.ObjectID, ingredient types.IngredientInRecipe) error {
	_, err := d.col.UpdateOne(ctx, bson.M{"_id": ID}, bson.M{"$push": bson.M{"ingredients": ingredient}})

	return err
}

func (d *RecipeDAO) UpdateRecipeIngredient(ctx context.Context, recipeID primitive.ObjectID, ingredient types.IngredientInRecipe) error {
	// Define the filter to match the specific recipe by ID and the ingredient by ID
	filter := bson.M{
		"_id":                       recipeID,
		"ingredients.ingredient_id": ingredient.IngredientID,
	}

	// Define the update to modify the amount of the matched ingredient
	update := bson.M{
		"$set": bson.M{"ingredients.$.amount": ingredient.Amount},
	}

	result, err := d.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount < 1 {
		return fmt.Errorf("error updating recipe ingredient: modified count 0")
	}

	return nil
}

func (d *RecipeDAO) SetInGroceryListStatus(ctx context.Context, ID primitive.ObjectID, status bool) error {
	_, err := d.col.UpdateOne(ctx, bson.M{"_id": ID}, bson.M{"$set": bson.M{"inGroceryList": status}})

	return err
}
