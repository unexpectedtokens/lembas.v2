package recipe_handler

import (
	"fmt"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	handler_util "github.com/unexpectedtoken/recipes/handler/common"
	"github.com/unexpectedtoken/recipes/types"
	recipe_view "github.com/unexpectedtoken/recipes/view/recipe"
)

func (h *RecipeHandler) HandleUpdateIngredient(w http.ResponseWriter, r *http.Request) {
	recipeID := chi.URLParam(r, "id")

	ingredientID := chi.URLParam(r, "ingID")

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	amount := r.FormValue("amount")

	amountFl, err := strconv.ParseFloat(amount, 32)

	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	ingredient := types.IngredientInRecipe{
		IngredientID: ingredientID,
		Amount:       float32(amountFl),
	}

	err = h.recipeService.UpdateRecipeIngredient(r.Context(), recipeID, ingredient)

	if err != nil {
		handler_util.LogErrorWithMessage(r, "error updating recipe", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, path.Join("/recipes", recipeID, "ingredients"), http.StatusSeeOther)
	fmt.Println(recipeID, ingredientID, amount)
}

func (h *RecipeHandler) HandleRemoveIngredientFromRecipe(w http.ResponseWriter, r *http.Request) {
	recipeID := chi.URLParam(r, "id")

	ingredientID := chi.URLParam(r, "ingID")

	err := h.recipeService.RemoveIngredientFromRecipe(r.Context(), recipeID, ingredientID)

	if err != nil {
		handler_util.LogErrorWithMessage(r, "error updating recipe", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, path.Join("/recipes", recipeID, "ingredients"), http.StatusSeeOther)
}

func (h *RecipeHandler) HandleShowRecipeList(w http.ResponseWriter, r *http.Request) {
	logger := handler_util.GetLoggerFromReqContext(r)

	concattedIds := r.URL.Query().Get("ids")
	var recipes []types.Recipe
	var err error
	if len(concattedIds) > 0 {
		ids := strings.Split(concattedIds, ",")

		recipes, err = h.recipeService.GetRecipesInIDList(r.Context(), &ids)
	} else {
		recipes, err = h.recipeService.GetList(r.Context())
	}

	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(500)
		return
	}

	recipe_view.RecipeList(recipes).Render(r.Context(), w)
}
