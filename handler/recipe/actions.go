package recipe_handler

import (
	"net/http"
	"path"

	"github.com/go-chi/chi/v5"
	handler_util "github.com/unexpectedtoken/recipes/handler/common"
	"github.com/unexpectedtoken/recipes/types"
)

func (h *RecipeHandler) HandleDeleteRecipe(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	err := h.recipeService.RemoveRecipe(r.Context(), id)

	if err != nil {
		h.HandleServerError(w, r, err)
		return
	}

	err = h.mealplanService.RemoveAllEntriesForRecipe(r.Context(), id)

	if err != nil {
		h.HandleServerError(w, r, err)
		return
	}

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusNoContent)
}

func (h *RecipeHandler) HandleAddNewIngredientToRecipe(w http.ResponseWriter, r *http.Request) {
	recipeID := chi.URLParam(r, "id")

	err := r.ParseForm()
	if err != nil {
		h.HandleClientError(w, r, err)
		return
	}

	ingredient := types.Ingredient{
		Name:          r.FormValue("ingredient-name"),
		MeasuringUnit: r.FormValue("measuring-unit"),
	}

	err = h.recipeService.AddNewIngredientToRecipe(r.Context(), recipeID, ingredient)

	if err != nil {
		h.HandleServerError(w, r, err)
		return
	}

	http.Redirect(w, r, path.Join("/recipes", recipeID, "ingredients"), http.StatusSeeOther)
}

func (h RecipeHandler) HandleAddIngredientToRecipe(w http.ResponseWriter, r *http.Request) {
	recipe, err := h.recipeFromReqOrHandleError(w, r)

	if err != nil {
		return
	}

	ingredientID := r.FormValue("ingredient-id")

	err = h.recipeService.AddIngredientToRecipe(r.Context(), *recipe, ingredientID)

	if err != nil {
		handler_util.LogErrorWithMessage(r, "error adding ingredient to recipe", err)
		return
	}

	http.Redirect(w, r, path.Join("/recipes", recipe.ID, "ingredients"), http.StatusSeeOther)
}
