package mealplan

import (
	"time"
)

type MealType string

const MealTypeBreakfast MealType = "breakfast"
const MealTypeLunch MealType = "lunch"
const MealTypeDinner MealType = "dinner"
const MealTypeSnacks MealType = "snacks"

type MealplanRecipeDataEntries []MealplanEntryV2

type MealplanRecipeData struct {
	Amount      int    `bson:"amount"`
	RecipeID    string `bson:"recipe_id"`
	RecipeTitle string `bson:"recipe_title"`
}

type MealplanEntryV2 struct {
	ID          string    `bson:"_id,omitempty"`
	Date        time.Time `bson:"date"`
	Amount      int       `bson:"amount"`
	RecipeID    string    `bson:"recipe_id"`
	RecipeTitle string    `bson:"recipe_title"`
	MealType    MealType
}

type MealplanEntry struct {
	ID        string               `bson:"_id,omitempty"`
	Date      time.Time            `bson:"date"`
	Breakfast []MealplanRecipeData `bson:"breakfast"`
	Lunch     []MealplanRecipeData `bson:"lunch"`
	Dinner    []MealplanRecipeData `bson:"dinner"`
	Snacks    []MealplanRecipeData `bson:"snacks"`
}

type MealplanDay struct {
	Date  time.Time
	Entry *MealplanEntry
}

type Period []time.Time

type MealplanWeekOverview []MealplanDay

const EntryDateIdentifierFormat = "02-01-2006"

func FormatDateToIdentifierFormat(t time.Time) string {
	return t.Format(EntryDateIdentifierFormat)
}

func ParseMealplanDate(date string) (time.Time, error) {
	return time.Parse(EntryDateIdentifierFormat, date)
}

func (e *MealplanRecipeDataEntries) IDS() *[]string {
	ids := []string{}

	for _, recipeEntry := range *e {
		ids = append(ids, recipeEntry.RecipeID)
	}

	return &ids
}

type DateRange struct {
	FromDate time.Time
	ToDate   time.Time
}
