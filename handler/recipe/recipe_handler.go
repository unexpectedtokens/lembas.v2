package recipe_handler

import (
	"net/http"
	"path"

	handler_util "github.com/unexpectedtoken/recipes/handler/common"
	services "github.com/unexpectedtoken/recipes/service"
	"github.com/unexpectedtoken/recipes/types"
	recipe_view "github.com/unexpectedtoken/recipes/view/recipe"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RecipeHandler struct {
	// TODO: Change to recipe service interface
	recipeService     *services.RecipeService
	ingredientService *services.IngredientService
}

func NewRecipeHandler(recServ *services.RecipeService, ingServ *services.IngredientService) *RecipeHandler {
	return &RecipeHandler{
		recipeService:     recServ,
		ingredientService: ingServ,
	}
}

func (h RecipeHandler) HandleViewRecipeIngredientList(w http.ResponseWriter, r *http.Request) {
	recipe, err := h.richRecipeFromReqOrHandleError(w, r)

	if err != nil {
		return
	}

	ing, err := h.ingredientService.GetIngredientsNotInList(r.Context(), recipe.Ingredients)
	if err != nil {
		w.WriteHeader(500)
		handler_util.LogErrorWithMessage(r, "error getting ingredients not in recipe", err)
		return
	}
	recipe_view.RecipeIngredientList(recipe.ID, recipe.PopulatedIngredients, len(ing) > 0).Render(r.Context(), w)
}

func (h RecipeHandler) HandleViewRecipeAddIngredientForm(w http.ResponseWriter, r *http.Request) {
	recipe, err := h.recipeFromReqOrHandleError(w, r)
	if err != nil {

		return
	}

	ingredients, err := h.ingredientService.GetIngredientsNotInList(r.Context(), recipe.Ingredients)

	if err != nil {
		handler_util.LogErrorWithMessage(r, "error getting ingredient options", err)
		return
	}

	recipe_view.AddIngredientToRecipeForm(recipe.ID, ingredients).Render(r.Context(), w)
}

func (h RecipeHandler) HandleAddIngredientToRecipe(w http.ResponseWriter, r *http.Request) {
	recipe, err := h.recipeFromReqOrHandleError(w, r)

	if err != nil {
		return
	}

	ingredientID, err := primitive.ObjectIDFromHex(r.FormValue("ingredient-id"))
	if err != nil {
		w.WriteHeader(400)
		return
	}

	ingredient := types.IngredientInRecipe{
		IngredientID: ingredientID,
	}

	err = h.recipeService.AddIngredientToRecipe(r.Context(), *recipe, ingredient)

	if err != nil {
		handler_util.LogErrorWithMessage(r, "error adding ingredient to recipe", err)
		return
	}

	http.Redirect(w, r, path.Join("/recipes", recipe.ID.Hex(), "ingredients"), http.StatusSeeOther)
}

func (h RecipeHandler) HandleViewRecipeCreateForm(w http.ResponseWriter, r *http.Request) {
	recipe_view.RecipeCreateForm().Render(r.Context(), w)
}

func (h RecipeHandler) HandleRecipeSubmit(w http.ResponseWriter, r *http.Request) {
	logger := handler_util.GetLoggerFromReqContext(r)
	err := r.ParseForm()
	if err != nil {
		handler_util.LogErrorWithMessage(r, "error parsing form", err)
		w.WriteHeader(400)
		return
	}

	var recipe types.Recipe

	err = handler_util.Decoder.Decode(&recipe, r.PostForm)
	if err != nil {
		handler_util.LogErrorWithMessage(r, "error decoding into struct", err)
		w.WriteHeader(400)
		return
	}
	recipe.Ingredients = []types.IngredientInRecipe{}
	id, err := h.recipeService.Create(r.Context(), recipe)

	if err != nil {
		handler_util.LogErrorWithMessage(r, "error saving recipe", err)
	}

	logger.Info("recipe created.", "id", id.Hex())

	http.Redirect(w, r, path.Join("/recipes", id.Hex()), http.StatusFound)
}
