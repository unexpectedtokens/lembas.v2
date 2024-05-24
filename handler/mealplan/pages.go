package mealplan_handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/unexpectedtoken/recipes/handler/base"
	services "github.com/unexpectedtoken/recipes/service"
	mealplan_view "github.com/unexpectedtoken/recipes/view/mealplan"
)

type MealplanHandler struct {
	mealplanService *services.MealplanService
	recipeService   *services.RecipeService
	*base.BaseHandler
}

func NewMealplanHandler(mealplanService *services.MealplanService, recipeService *services.RecipeService, baseHandler *base.BaseHandler) *MealplanHandler {
	return &MealplanHandler{
		mealplanService: mealplanService,
		recipeService:   recipeService,
		BaseHandler:     baseHandler,
	}
}

func (h *MealplanHandler) HandleViewMealplanPage(w http.ResponseWriter, r *http.Request) {
	props := &mealplan_view.MealplanOverviewProps{
		Period: h.mealplanService.GetPeriod(),
	}

	h.RenderHTMXWithLayout(w, r, mealplan_view.MealplanOverviewPage(props))
}

func (h *MealplanHandler) HandleViewShoppingListDatePicker(w http.ResponseWriter, r *http.Request) {
	h.RenderHTMXWithLayout(w, r, mealplan_view.MealplanDatePickerPage())
}

func (h *MealplanHandler) HandleViewShoppingListForDateRange(w http.ResponseWriter, r *http.Request) {
	dateRangeString := chi.URLParam(r, "daterange")

	dateRange, err := h.mealplanService.ParseDateRange(dateRangeString)

	if err != nil {
		h.HandleServerError(w, r, err)
		return
	}

	shoppingList, err := h.mealplanService.ShoppingListForDateRange(r.Context(), dateRange)

	if err != nil {
		h.HandleServerError(w, r, err)
		return
	}

	h.RenderHTMXWithLayout(w, r,
		mealplan_view.ShoppingListPage(shoppingList),
	)
}
