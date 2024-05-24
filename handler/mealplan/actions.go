package mealplan_handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/unexpectedtoken/recipes/mealplan"
)

func (h *MealplanHandler) HandleAddEntry(w http.ResponseWriter, r *http.Request) {

	day := chi.URLParam(r, "day")

	dayConv, err := time.Parse("02-01-2006", day)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	mealType := chi.URLParam(r, "mealtype")

	recipeId := chi.URLParam(r, "recipeID")

	recipe, err := h.recipeService.GetByID(r.Context(), recipeId)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	newEntry := mealplan.MealplanEntryV2{
		RecipeID:    recipeId,
		RecipeTitle: recipe.Title,
		Amount:      1,
		Date:        dayConv,
		MealType:    mealplan.MealType(mealType),
	}

	_, err = h.mealplanService.Create(r.Context(), newEntry)

	if err != nil {
		h.HandleServerError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/mealplan/%s/%s", day, mealType), http.StatusSeeOther)
}

func (h *MealplanHandler) HandleRemoveEntry(w http.ResponseWriter, r *http.Request) {
	entryID := chi.URLParam(r, "entryID")

	entry, err := h.mealplanService.GetByID(r.Context(), entryID)

	if err != nil {
		h.HandleServerError(w, r, err)
		return
	}

	err = h.mealplanService.DeleteByID(r.Context(), entryID)

	if err != nil {
		h.HandleServerError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/mealplan/%s/%s", mealplan.FormatDateToIdentifierFormat(entry.Date), entry.MealType), http.StatusSeeOther)
}
