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
	"github.com/unexpectedtoken/recipes/dao"
	"github.com/unexpectedtoken/recipes/db"
	ingredient_handler "github.com/unexpectedtoken/recipes/handler/ingredient"
	recipe_handler "github.com/unexpectedtoken/recipes/handler/recipe"
	shopping_handler "github.com/unexpectedtoken/recipes/handler/shopping"
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

	// DAO
	recipeDAO := dao.NewRecipeDAO(database)
	ingredientsDAO := dao.NewMongoDAO[types.Ingredient](database, "ingredients")

	// Services
	ingredientsService := services.NewIngredientService(ingredientsDAO)
	recipesService := services.NewRecipeService(recipeDAO, ingredientsService)

	// Handlers
	recipesHandler := recipe_handler.NewRecipeHandler(recipesService, ingredientsService)
	ingredientsHandler := ingredient_handler.NewIngredientHandler(ingredientsService)
	shoppingListHandler := shopping_handler.NewShoppingListHandler(recipesService)

	// Global middlewares
	logger := httplog.NewLogger("lembas-server", httplog.Options{
		// JSON:             true,
		LogLevel:         slog.LevelDebug,
		Concise:          true,
		RequestHeaders:   true,
		MessageFieldName: "message",
		// TimeFieldFormat: time.RFC850,
		Tags: map[string]string{
			"version": "v1.0-81aa4244d9fc8076a",
			"env":     "dev",
		},
		QuietDownRoutes: []string{
			"/",
			"/ping",
		},
		QuietDownPeriod: 10 * time.Second,
		// SourceFieldName: "source",
	})
	r.Use(httplog.RequestLogger(logger), middleware.Recoverer)

	// Register routes on root
	r.Get("/", recipesHandler.HandleViewRecipeOverview)
	r.Get("/recipes/create", recipesHandler.HandleViewRecipeCreateForm)
	r.Post("/recipes/create", recipesHandler.HandleRecipeSubmit)
	r.Get("/recipes/{id}", recipesHandler.HandleViewRecipeDetail)
	r.Get("/recipes/{id}/add-ingredient", recipesHandler.HandleViewRecipeAddIngredientForm)
	r.Post("/recipes/{id}/add-ingredient", recipesHandler.HandleAddIngredientToRecipe)
	r.Get("/recipes/{id}/ingredients", recipesHandler.HandleViewRecipeIngredientList)
	r.Put("/recipes/{id}/ingredients/{ingID}", recipesHandler.HandleUpdateIngredient)
	r.Put("/recipes/{id}/set-status/{status}", recipesHandler.HandleAddToGroceryList)

	r.Get("/ingredients", ingredientsHandler.HandleViewIngredients)
	r.Post("/ingredients/create", ingredientsHandler.HandlePostIngredient)

	r.Get("/shopping-list/overview", shoppingListHandler.HandleViewRecipesInShoppingList)
	// Serve public assets
	fs := http.FileServer(http.Dir("./assets/public"))
	r.Handle("/public/*", http.StripPrefix("/public/", fs))

	// Register api routes
	// empty

	log.Println("Listening on port http://localhost:4200")
	http.ListenAndServe(":4200", r)
}
