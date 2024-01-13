package shopping_handler

import (
	"net/http"

	handler_util "github.com/unexpectedtoken/recipes/handler/common"
	services "github.com/unexpectedtoken/recipes/service"
	shopping_list_view "github.com/unexpectedtoken/recipes/view/shopping"
)

type ShoppingListHandler struct {
	recipeService *services.RecipeService // Create shopping list service
}

func NewShoppingListHandler(recipeService *services.RecipeService) *ShoppingListHandler {
	return &ShoppingListHandler{
		recipeService: recipeService,
	}
}

func (h *ShoppingListHandler) HandleViewRecipesInShoppingList(w http.ResponseWriter, r *http.Request) {
	recipes, err := h.recipeService.RecipesInShoppingList(r.Context())

	if err != nil {
		handler_util.LogErrorWithMessage(r, "error getting grocery list recipes", err)
		w.WriteHeader(500)
		return
	}

	shopping_list_view.ShoppingListOverview(recipes).Render(r.Context(), w)
}
