package recipe_handler

import (
	"fmt"
	"net/http"
	"path"
	"strconv"

	"github.com/go-chi/chi/v5"
	handler_util "github.com/unexpectedtoken/recipes/handler/common"
	"github.com/unexpectedtoken/recipes/types"
	recipe_view "github.com/unexpectedtoken/recipes/view/recipe"
)

func (h *RecipeHandler) HandleUpdateIngredient(w http.ResponseWriter, r *http.Request) {
	recipeID, err := handler_util.ObjectIDFromR(r, "id")

	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	ingredientID, err := handler_util.ObjectIDFromR(r, "ingID")
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = r.ParseForm()
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

	http.Redirect(w, r, path.Join("/recipes", recipeID.Hex(), "ingredients"), http.StatusSeeOther)
	fmt.Println(recipeID, ingredientID, amount)
}

func (h *RecipeHandler) HandleRemoveIngredientFromRecipe(w http.ResponseWriter, r *http.Request) {
	recipeID, err := handler_util.ObjectIDFromR(r, "id")

	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	ingredientID, err := handler_util.ObjectIDFromR(r, "ingID")
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = h.recipeService.RemoveIngredientFromRecipe(r.Context(), recipeID, ingredientID)

	if err != nil {
		handler_util.LogErrorWithMessage(r, "error updating recipe", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, path.Join("/recipes", recipeID.Hex(), "ingredients"), http.StatusSeeOther)
}

func (h *RecipeHandler) HandleAddToGroceryList(w http.ResponseWriter, r *http.Request) {
	recipeID, err := handler_util.ObjectIDFromR(r, "id")

	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	statusString := chi.URLParam(r, "status")

	status, err := strconv.ParseBool(statusString)

	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = h.recipeService.UpdateGroceryListStatus(r.Context(), recipeID, status)

	if err != nil {
		handler_util.LogErrorWithMessage(r, "error updating grocery list status", err)
		w.WriteHeader(500)
		return
	}

	recipe_view.InGroceryListButton(status, recipeID).Render(r.Context(), w)
}
