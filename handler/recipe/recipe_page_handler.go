package recipe_handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	handler_util "github.com/unexpectedtoken/recipes/handler/common"
	recipe_view "github.com/unexpectedtoken/recipes/view/recipe"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h RecipeHandler) HandleViewRecipeDetail(w http.ResponseWriter, r *http.Request) {
	logger := handler_util.GetLoggerFromReqContext(r)
	id := chi.URLParam(r, "id")

	recipe, err := h.recipeService.GetRecipeByID(r.Context(), id, true)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			w.WriteHeader(404)
			return
		}

		logger.Error("non-404 error fetching recipe by id", "err", err, "id", id)
		return

	}

	h.RenderHTMXWithLayout(w, r, recipe_view.RecipeDetail(&recipe_view.RecipeDetailProps{
		Recipe: recipe,
	}))
}

func (h RecipeHandler) HandleViewRecipeOverview(w http.ResponseWriter, r *http.Request) {
	recipes, err := h.recipeService.GetList(r.Context())

	if err != nil {
		h.HandleServerError(w, r, err)
		return
	}

	h.RenderHTMXWithLayout(w, r, recipe_view.RecipeOverview(recipes))
}
