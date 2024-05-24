package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/unexpectedtoken/recipes/mealplan"
	"github.com/unexpectedtoken/recipes/shopping"
	"github.com/unexpectedtoken/recipes/types"
)

type MealplanRepo interface {
	DAO[mealplan.MealplanEntryV2]
	EntriesForDateWithType(ctx context.Context, fromDate time.Time, toDate time.Time, mealType mealplan.MealType) (*mealplan.MealplanRecipeDataEntries, error)
	EntriesForDate(ctx context.Context, fromDate time.Time, toDate time.Time) (*mealplan.MealplanRecipeDataEntries, error)
	RemoveRecipeEntries(ctx context.Context, recipeID string) error
}

type MealplanService struct {
	repo MealplanRepo
	genericCrudService[mealplan.MealplanEntryV2]
	recipeService *RecipeService
}

func NewMealplanService(repo MealplanRepo, recipeService *RecipeService) *MealplanService {
	return &MealplanService{
		repo: repo,
		genericCrudService: genericCrudService[mealplan.MealplanEntryV2]{
			dao: repo,
		},
		recipeService: recipeService,
	}
}

func startOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func endOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
}

func (s *MealplanService) GetPeriod() mealplan.Period {
	period := mealplan.Period{}
	today := time.Now()

	for i := 0; i < 7; i++ {
		multiplier := 24 * i
		currentIterationDay := today.Add(time.Hour * time.Duration(multiplier))

		period = append(period, currentIterationDay)
	}

	return period
}

func (s *MealplanService) MealtypeEntriesForDate(ctx context.Context, date time.Time, mealType mealplan.MealType) (*mealplan.MealplanRecipeDataEntries, error) {
	return s.repo.EntriesForDateWithType(ctx, startOfDay(date), endOfDay(date), mealType)
}

func (s *MealplanService) allEntriesForDates(ctx context.Context, dr *mealplan.DateRange) (*mealplan.MealplanRecipeDataEntries, error) {
	return s.repo.EntriesForDate(ctx, startOfDay(dr.FromDate), endOfDay(dr.ToDate))
}

func parseDates(from string, to string) (*mealplan.DateRange, error) {
	fromDate, err := mealplan.ParseMealplanDate(from)

	if err != nil {
		return nil, fmt.Errorf("error parsing %s as from date: %w", from, err)
	}
	toDate, err := mealplan.ParseMealplanDate(to)

	if err != nil {
		return nil, fmt.Errorf("error parsing %s as to date: %w", to, err)
	}

	dr := &mealplan.DateRange{
		FromDate: fromDate,
		ToDate:   toDate,
	}

	return dr, nil
}

func (s *MealplanService) ParseDateRange(dateRange string) (*mealplan.DateRange, error) {
	if len(dateRange) != 23 {
		return nil, fmt.Errorf("dateRange %s too short", dateRange)
	}

	rangeSplit := strings.Split(dateRange, "...")

	if len(rangeSplit) != 2 {
		return nil, fmt.Errorf("dateRange %s using wrong seperator", dateRange)
	}

	output, err := parseDates(rangeSplit[0], rangeSplit[1])

	if err != nil {
		return nil, fmt.Errorf("error parsing from and to date: %w", err)
	}

	return output, nil
}

func (s *MealplanService) ShoppingListForDateRange(ctx context.Context, dateRange *mealplan.DateRange) (*shopping.ShoppingList, error) {
	entries, err := s.allEntriesForDates(ctx, dateRange)

	if err != nil {
		return nil, fmt.Errorf("error getting entries: %w", err)
	}

	recipes, err := s.recipeService.GetPopulatedRecipesInIDList(ctx, entries.IDS())

	if err != nil {
		return nil, fmt.Errorf("error getting and populating recipes: %w", err)
	}

	recipeMap := s.recipeService.MapRecipes(&recipes)
	shoppingList := &shopping.ShoppingList{
		DateRange:   dateRange,
		Ingredients: make(map[string]types.PopulatedIngredientInRecipe),
	}

	for _, entry := range *entries {
		if recipe, ok := (*recipeMap)[entry.RecipeID]; ok {
			for _, ingredient := range *recipe.PopulatedIngredients {
				found := false
				for key, val := range shoppingList.Ingredients {
					if key == ingredient.IngredientID {
						found = true
						val.Amount += ingredient.Amount
						shoppingList.Ingredients[ingredient.IngredientID] = val
						break
					}
				}

				if !found {
					shoppingList.Ingredients[ingredient.IngredientID] = ingredient
				}
			}
		}

	}

	return shoppingList, nil
}

func (s *MealplanService) RemoveAllEntriesForRecipe(ctx context.Context, recipeID string) error {
	return s.repo.RemoveRecipeEntries(ctx, recipeID)
}
