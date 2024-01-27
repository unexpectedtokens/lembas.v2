package recipe_handler

import (
	"net/http"
	"path"

	handler_util "github.com/unexpectedtoken/recipes/handler/common"
	"github.com/unexpectedtoken/recipes/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *RecipeHandler) HandleDeleteRecipe(w http.ResponseWriter, r *http.Request) {

	id, err := handler_util.ObjectIDFromR(r, "id")

	if err != nil {
		handler_util.LogErrorWithMessage(r, "error getting id from r", err)
		return
	}

	err = h.recipeService.RemoveRecipe(r.Context(), id)

	if err != nil {
		handler_util.LogErrorWithMessage(r, "error deleting recipe", err)
		return
	}

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusNoContent)
}

func (h *RecipeHandler) HandleAddNewIngredientToRecipe(w http.ResponseWriter, r *http.Request) {
	recipeID, err := handler_util.ObjectIDFromR(r, "id")

	if err != nil {
		w.WriteHeader(400)
		handler_util.HandleClientErr()
		return
	}

	err = r.ParseForm()
	if err != nil {
		handler_util.HandleClientErr()
		return
	}

	ingredient := types.Ingredient{
		Name:          r.FormValue("ingredient-name"),
		MeasuringUnit: r.FormValue("measuring-unit"),
	}

	err = h.recipeService.AddNewIngredientToRecipe(r.Context(), recipeID, ingredient)

	if err != nil {
		handler_util.HandleServerErr()
		return
	}

	http.Redirect(w, r, path.Join("/recipes", recipeID.Hex(), "ingredients"), http.StatusSeeOther)
}

func (h RecipeHandler) HandleAddIngredientToRecipe(w http.ResponseWriter, r *http.Request) {
	recipe, err := h.recipeFromReqOrHandleError(w, r)

	if err != nil {
		return
	}

	ingredientID, err := primitive.ObjectIDFromHex(r.FormValue("ingredient-id"))
	if err != nil {
		handler_util.HandleClientErr()
		return
	}

	err = h.recipeService.AddIngredientToRecipe(r.Context(), *recipe, ingredientID)

	if err != nil {
		handler_util.LogErrorWithMessage(r, "error adding ingredient to recipe", err)
		return
	}

	http.Redirect(w, r, path.Join("/recipes", recipe.ID.Hex(), "ingredients"), http.StatusSeeOther)
}
