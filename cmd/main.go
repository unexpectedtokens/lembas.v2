package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
	"github.com/joho/godotenv"
	"github.com/unexpectedtoken/recipes/db"
	auth_handler "github.com/unexpectedtoken/recipes/handler/auth"
	"github.com/unexpectedtoken/recipes/handler/base"
	ingredient_handler "github.com/unexpectedtoken/recipes/handler/ingredient"
	mealplan_handler "github.com/unexpectedtoken/recipes/handler/mealplan"
	custom_middleware "github.com/unexpectedtoken/recipes/handler/middleware"
	recipe_handler "github.com/unexpectedtoken/recipes/handler/recipe"
	"github.com/unexpectedtoken/recipes/repository"
	services "github.com/unexpectedtoken/recipes/service"
	"github.com/unexpectedtoken/recipes/types"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		panic(err)
	}

	r := chi.NewMux()

	database, err := db.NewDBConn(os.Getenv("MONGO_URI"))

	if err != nil {
		panic(err)
	}

	authCFG := services.AuthServiceCFG{
		Username: os.Getenv("USERNAME"),
		Password: os.Getenv("PASSWORD"),
	}

	// DAO
	recipeDAO := repository.NewRecipeDAO(database)
	ingredientsDAO := repository.NewMongoDAO[types.Ingredient](database, "ingredients")
	mealplanRepo := repository.NewMealplanRepo(database)

	// Services
	ingredientsService := services.NewIngredientService(ingredientsDAO)
	recipesService := services.NewRecipeService(recipeDAO, ingredientsService)
	mealplanService := services.NewMealplanService(mealplanRepo, recipesService)
	authService := services.NewAuthService(authCFG)

	// Handlers
	baseHandler := base.New()
	authHandler := auth_handler.New(baseHandler, authService)
	recipesHandler := recipe_handler.NewRecipeHandler(recipesService, ingredientsService, mealplanService, baseHandler)
	ingredientsHandler := ingredient_handler.NewIngredientHandler(ingredientsService, baseHandler)
	mealPlanHandler := mealplan_handler.NewMealplanHandler(mealplanService, recipesService, baseHandler)

	// Global middlewares
	logger := httplog.NewLogger("lembas-server", httplog.Options{
		LogLevel:         slog.LevelDebug,
		Concise:          true,
		RequestHeaders:   true,
		MessageFieldName: "message",
		Tags: map[string]string{
			"version": "v1.0-81aa4244d9fc8076a",
			"env":     "dev",
		},
		QuietDownRoutes: []string{
			"/",
			"/ping",
		},
		QuietDownPeriod: 10 * time.Second,
	})
	r.Use(httplog.RequestLogger(logger), middleware.Recoverer)

	// Register routes on root
	r.Group(func(r chi.Router) {
		r.Use(
			custom_middleware.AuthMiddleware,
		)

		r.Get("/", recipesHandler.HandleViewRecipeOverview)
		r.Get("/recipes/list", recipesHandler.HandleShowRecipeList)
		r.Get("/recipes/create", recipesHandler.HandleViewRecipeCreateForm)
		r.Post("/recipes/create", recipesHandler.HandleRecipeSubmit)
		r.Get("/recipes/{id}", recipesHandler.HandleViewRecipeDetail)
		r.Delete("/recipes/{id}", recipesHandler.HandleDeleteRecipe)
		r.Get("/recipes/{id}/add-ingredient", recipesHandler.HandleViewRecipeAddIngredientForm)
		r.Post("/recipes/{id}/add-ingredient", recipesHandler.HandleAddIngredientToRecipe)
		r.Post("/recipes/{id}/add-new-ingredient", recipesHandler.HandleAddNewIngredientToRecipe)
		r.Get("/recipes/{id}/ingredients", recipesHandler.HandleViewRecipeIngredientList)
		r.Put("/recipes/{id}/ingredients/{ingID}", recipesHandler.HandleUpdateIngredient)
		r.Delete("/recipes/{id}/ingredients/{ingID}", recipesHandler.HandleRemoveIngredientFromRecipe)

		r.Get("/ingredients", ingredientsHandler.HandleViewIngredients)
		r.Post("/ingredients/create", ingredientsHandler.HandlePostIngredient)

		r.Get("/mealplan", mealPlanHandler.HandleViewMealplanPage)
		r.Get("/mealplan/{day}", mealPlanHandler.HandleViewDay)
		r.Get("/mealplan/{day}/{mealtype}", mealPlanHandler.HandleViewMealplanEntrySection)
		r.Get("/mealplan/{day}/{mealtype}/new", mealPlanHandler.HandleViewAddEntry)
		r.Post("/mealplan/{day}/{mealtype}/new/{recipeID}", mealPlanHandler.HandleAddEntry)
		r.Delete("/mealplan/entries/{entryID}", mealPlanHandler.HandleRemoveEntry)
		r.Get("/mealplan/shopping-list/datepicker", mealPlanHandler.HandleViewShoppingListDatePicker)
		r.Get("/mealplan/shopping-list/{daterange}", mealPlanHandler.HandleViewShoppingListForDateRange)
	})

	r.Get("/auth", authHandler.HandleViewLoginPage)
	r.Post("/auth", authHandler.HandleLoginAttempt)

	// Serve public assets
	fs := http.FileServer(http.Dir("./assets/public"))
	r.Handle("/public/*", http.StripPrefix("/public/", fs))

	// Register api routes
	// empty

	log.Println("Listening on port http://localhost:4200")
	panic(http.ListenAndServe("localhost:4200", r))
}
