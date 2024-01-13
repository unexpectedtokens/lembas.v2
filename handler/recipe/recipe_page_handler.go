package recipe_handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	handler_util "github.com/unexpectedtoken/recipes/handler/common"
	recipe_view "github.com/unexpectedtoken/recipes/view/recipe"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h RecipeHandler) HandleViewRecipeDetail(w http.ResponseWriter, r *http.Request) {
	logger := handler_util.GetLoggerFromReqContext(r)
	id := chi.URLParam(r, "id")

	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		w.WriteHeader(400)
		logger.Error("error parsing objectId", "err", err)
		return
	}

	recipe, err := h.recipeService.GetRecipeByID(r.Context(), objectId, true)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			w.WriteHeader(404)
			return
		}

		logger.Error("non-404 error fetching recipe by id", "err", err, "id", id)
		return

	}

	ing, err := h.ingredientService.GetIngredientsNotInList(r.Context(), recipe.Ingredients)

	if err != nil {
		w.WriteHeader(500)
		handler_util.LogErrorWithMessage(r, "error getting ingredients not in recipe", err)
		return
	}

	recipe_view.RecipeDetail(recipe, len(ing) > 0).Render(r.Context(), w)
}

func (h RecipeHandler) HandleViewRecipeOverview(w http.ResponseWriter, r *http.Request) {
	logger := handler_util.GetLoggerFromReqContext(r)
	recipes, err := h.recipeService.GetList(r.Context())

	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(500)
		return
	}

	recipe_view.RecipeOverview(recipes).Render(r.Context(), w)
}
