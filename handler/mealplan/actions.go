package mealplan_handler

import (
	"net/http"
	"path"

	"github.com/go-chi/chi/v5"
	"github.com/unexpectedtoken/recipes/logging"
	"github.com/unexpectedtoken/recipes/types"
)

func (h *MealplanHandler) HandleAddEntry(w http.ResponseWriter, r *http.Request) {
	datestring := chi.URLParam(r, "datestring")

	date, err := types.ParseMealplanDate(datestring)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	entry := types.MealplanEntry{
		Date: date,
	}

	id, err := h.mealplanService.Create(r.Context(), entry)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logging.LogRequestError(r, err)
		return
	}

	http.Redirect(w, r, path.Join("/mealplan/entries", id.Hex()), http.StatusSeeOther)
}
