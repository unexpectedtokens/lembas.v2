package mealplan_handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	handler_util "github.com/unexpectedtoken/recipes/handler/common"
	"github.com/unexpectedtoken/recipes/logging"
	"github.com/unexpectedtoken/recipes/types"
	mealplan_view "github.com/unexpectedtoken/recipes/view/mealplan"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h *MealplanHandler) HandleViewMealplanEntryCreateForm(w http.ResponseWriter, r *http.Request) {
	datestring := chi.URLParam(r, "datestring")
	// TODO: manage this logic with AlpineJS, in the frontend
	mealplan_view.EntryCreateForm(datestring).Render(r.Context(), w)
}

func (h *MealplanHandler) HandleViewMealplanEntry(w http.ResponseWriter, r *http.Request) {
	id, err := handler_util.ObjectIDFromR(r, "id")

	if err != nil {
		w.WriteHeader(400)
		logging.LogRequestError(r, err)
		return
	}

	entry, err := h.mealplanService.GetByID(r.Context(), id)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			logging.LogRequestError(r, err)
		}

		return
	}

	entryData := mealplan_view.EntryData{
		Entry: entry,
		Date:  types.FormatMealplanDate(entry.Date),
	}

	mealplan_view.Entry(&entryData).Render(r.Context(), w)
}
