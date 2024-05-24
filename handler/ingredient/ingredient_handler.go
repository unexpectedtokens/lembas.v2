package ingredient_handler

import (
	"net/http"

	"github.com/unexpectedtoken/recipes/handler/base"
	handler_util "github.com/unexpectedtoken/recipes/handler/common"
	services "github.com/unexpectedtoken/recipes/service"
	"github.com/unexpectedtoken/recipes/types"
	ingredient_view "github.com/unexpectedtoken/recipes/view/ingredient"
)

type IngredientHandler struct {
	// TODO: Change to recipe service interface
	ingredientService *services.IngredientService
	*base.BaseHandler
}

func NewIngredientHandler(serv *services.IngredientService, baseHandler *base.BaseHandler) *IngredientHandler {
	return &IngredientHandler{
		ingredientService: serv,
		BaseHandler:       baseHandler,
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

	h.RenderHTMXWithLayout(w, r, ingredient_view.IngredientsList(ingredients))
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

	h.RenderHTMX(w, r, ingredient_view.IngredientCard(ingredient))
}
