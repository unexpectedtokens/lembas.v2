package mealplan_handler

import (
	"net/http"
	"time"

	services "github.com/unexpectedtoken/recipes/service"
	"github.com/unexpectedtoken/recipes/types"
	mealplan_view "github.com/unexpectedtoken/recipes/view/mealplan"
)

type MealplanHandler struct {
	mealplanService *services.MealplanService
}

func NewMealplanHandler(mealplanService *services.MealplanService) *MealplanHandler {
	return &MealplanHandler{
		mealplanService: mealplanService,
	}
}

func startOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func endOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
}

func (h *MealplanHandler) HandleViewMealplanWeek(w http.ResponseWriter, r *http.Request) {
	days := []string{}

	today := time.Now()

	fromFilter := startOfDay(today)
	var toFilter time.Time
	days = append(days, types.FormatMealplanDate(today))
	for i := 1; i < 7; i++ {
		multiplier := 24 * i
		day := today.Add(time.Hour * time.Duration(multiplier))

		days = append(days, types.FormatMealplanDate(day))

		if i == 6 {
			toFilter = endOfDay(day)
		}
	}

	h.mealplanService.EntriesForDates(r.Context(), fromFilter, toFilter)
	mealplan_view.MealplanOverviewPage(days).Render(r.Context(), w)
}
