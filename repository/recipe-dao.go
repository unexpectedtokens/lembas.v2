package repository

import (
	"context"
	"fmt"

	"github.com/unexpectedtoken/recipes/types"
	"go.mongodb.org/mongo-driver/bson"
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

func (d *RecipeDAO) AddIngredientToRecipe(ctx context.Context, ID string, ingredient types.IngredientInRecipe) error {
	recipeID, err := parseId(ID)

	if err != nil {
		return err
	}

	_, err = d.col.UpdateOne(ctx, bson.M{"_id": recipeID}, bson.M{"$push": bson.M{"ingredients": ingredient}})

	return err
}

func (d *RecipeDAO) UpdateRecipeIngredient(ctx context.Context, recipeID string, ingredient types.IngredientInRecipe) error {
	recipeObjID, err := parseId(recipeID)

	if err != nil {
		return err
	}

	// Define the filter to match the specific recipe by ID and the ingredient by ID
	filter := bson.M{
		"_id":                       recipeObjID,
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

func (d *RecipeDAO) DeleteRecipe(ctx context.Context, ID string) error {
	res, err := d.col.DeleteOne(ctx, bson.M{"_id": ID})

	if err != nil {
		return fmt.Errorf("an error occured deleting recipe from repository: %w", err)
	}

	if res.DeletedCount == 1 {
		return nil
	}

	return fmt.Errorf("no recipe deleted: deletedCount 0")
}

func (d *RecipeDAO) DeleteIngredientFromRecipe(ctx context.Context, ID string, ingredientID string) error {
	// Define the filter to match the specific recipe by ID and the ingredient by ID
	filter := bson.M{
		"_id": ID,
	}

	// Define the update to remove the matched ingredient from the array
	update := bson.M{
		"$pull": bson.M{"ingredients": bson.M{"ingredient_id": ingredientID}},
	}

	// Perform the update
	_, err := d.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return fmt.Errorf("no recipe deleted: deletedCount 0")
}
