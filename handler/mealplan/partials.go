package mealplan_handler

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/unexpectedtoken/recipes/mealplan"
	mealplan_view "github.com/unexpectedtoken/recipes/view/mealplan"
)

func (h *MealplanHandler) HandleViewMealplanEntrySection(w http.ResponseWriter, r *http.Request) {
	day := chi.URLParam(r, "day")
	mealType := chi.URLParam(r, "mealtype")

	dayConv, err := time.Parse("02-01-2006", day)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	entries, err := h.mealplanService.MealtypeEntriesForDate(r.Context(), dayConv, mealplan.MealType(mealType))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	if err != nil {
		return
	}

	recipes, err := h.recipeService.GetRecipesNotInIDList(r.Context(), entries.IDS())

	if err != nil {
		h.HandleServerError(w, r, err)
		return
	}

	h.RenderHTMX(w, r, mealplan_view.MealplanSection((*mealplan.MealplanRecipeDataEntries)(entries), day, len(recipes) > 0, mealType))
}

func (h *MealplanHandler) HandleViewDay(w http.ResponseWriter, r *http.Request) {
	day := chi.URLParam(r, "day")

	dayConv, err := time.Parse("02-01-2006", day)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mealplan_view.Day(dayConv.Format("02-01-2006")).Render(r.Context(), w)
}

func (h *MealplanHandler) HandleViewAddEntry(w http.ResponseWriter, r *http.Request) {
	day := chi.URLParam(r, "day")
	mealType := chi.URLParam(r, "mealtype")
	dayConv, err := time.Parse("02-01-2006", day)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	entries, err := h.mealplanService.MealtypeEntriesForDate(r.Context(), dayConv, mealplan.MealType(mealType))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	recipes, err := h.recipeService.GetRecipesNotInIDList(r.Context(), entries.IDS())

	if err != nil {
		return
	}

	mealplan_view.MealOptions(recipes, day, mealplan.MealType(mealType)).Render(r.Context(), w)
}
