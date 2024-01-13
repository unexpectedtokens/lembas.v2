package ingredient_handler

import (
	"net/http"

	handler_util "github.com/unexpectedtoken/recipes/handler/common"
	services "github.com/unexpectedtoken/recipes/service"
	"github.com/unexpectedtoken/recipes/types"
	ingredient_view "github.com/unexpectedtoken/recipes/view/ingredient"
)

type IngredientHandler struct {
	// TODO: Change to recipe service interface
	ingredientService *services.IngredientService
}

func NewIngredientHandler(serv *services.IngredientService) *IngredientHandler {
	return &IngredientHandler{
		ingredientService: serv,
	}
}

func (h *IngredientHandler) HandleViewIngredients(w http.ResponseWriter, r *http.Request) {
	logger := handler_util.GetLoggerFromReqContext(r)
	ingredients, err := h.ingredientService.GetList(r.Context())

	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(500)
		return
	}

	ingredient_view.IngredientsList(ingredients).Render(r.Context(), w)
}

func (h *IngredientHandler) HandlePostIngredient(w http.ResponseWriter, r *http.Request) {

	var ingredient types.Ingredient

	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(500)
		return
	}

	err = handler_util.Decoder.Decode(&ingredient, r.PostForm)
	if err != nil {
		handler_util.LogErrorWithMessage(r, "error decoding into struct", err)
		w.WriteHeader(422)
		return
	}

	id, err := h.ingredientService.Create(r.Context(), ingredient)

	if err != nil {
		w.WriteHeader(500)
		return
	}

	ingredient.ID = id

	ingredient_view.IngredientCard(ingredient).Render(r.Context(), w)
}
