package recipe_handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	handler_util "github.com/unexpectedtoken/recipes/handler/common"
	"github.com/unexpectedtoken/recipes/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *RecipeHandler) recipeFromReqOrHandleError(w http.ResponseWriter, r *http.Request) (*types.Recipe, error) {
	id := chi.URLParam(r, "id")

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		w.WriteHeader(400)
		return nil, err
	}

	if err != nil {
		w.WriteHeader(400)
		handler_util.LogErrorWithMessage(r, "error parsing ObjectID", err)
		return nil, err
	}

	recipe, err := h.recipeService.GetRecipeByID(r.Context(), objectId, false)

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
		if err != nil {
			w.WriteHeader(400)
			return nil, err
		}
	}

	recipe.PopulatedIngredients = popIngredients

	return recipe, nil
}
