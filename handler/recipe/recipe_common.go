package recipe_handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/unexpectedtoken/recipes/types"
)

func (h *RecipeHandler) recipeFromReqOrHandleError(w http.ResponseWriter, r *http.Request) (*types.Recipe, error) {
	id := chi.URLParam(r, "id")

	recipe, err := h.recipeService.GetRecipeByID(r.Context(), id, false)

	// TODO: Create more generic view handlers like 400, 404, 500
	if err != nil {
		w.WriteHeader(404)
		return nil, err
	}

	return &recipe, nil
}

func (h *RecipeHandler) richRecipeFromReqOrHandleError(w http.ResponseWriter, r *http.Request) (*types.Recipe, error) {
	recipe, err := h.recipeFromReqOrHandleError(w, r)
	if err != nil {
		return nil, err
	}

	popIngredients, err := h.ingredientService.PopulateIngredientList(r.Context(), recipe.Ingredients)

	if err != nil {
		w.WriteHeader(400)
		return nil, err
	}

	recipe.PopulatedIngredients = popIngredients

	return recipe, nil
}
