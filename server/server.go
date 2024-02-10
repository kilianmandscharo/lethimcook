package server

import (
	"github.com/kilianmandscharo/lethimcook/auth"
	"github.com/kilianmandscharo/lethimcook/recipe"
	"github.com/labstack/echo/v4"
)

type Server struct {
	e *echo.Echo
}

func New() Server {
	recipeRequestHandler := recipe.NewRecipeRequestHandler()
	authRequestHandler := auth.NewAuthRequestHandler()

	e := echo.New()
	e.Static("/static", "static")
	recipe.AttachTemplates(e)

	// Pages
	e.GET("/", recipeRequestHandler.RenderRecipeListPage)
	e.GET("/recipe/edit/:id", recipeRequestHandler.RenderEditRecipePage)
	e.GET("/recipe/new", recipeRequestHandler.RenderNewRecipePage)
	e.GET("/recipe/:id", recipeRequestHandler.RenderRecipePage)

	// Actions
	e.POST("/recipe", recipeRequestHandler.HandleCreateRecipe)
	e.PUT("/recipe/:id", recipeRequestHandler.HandleUpdateRecipe)
	e.DELETE("/recipe/:id", recipeRequestHandler.HandleDeleteRecipe)

	// Auth
	e.POST("/auth/login", authRequestHandler.HandleLogin)
	e.PUT("/auth/password", authRequestHandler.HandleUpdatePassword)

	return Server{
		e: e,
	}
}

func (s *Server) Start() {
	s.e.Logger.Fatal(s.e.Start(":8080"))
}
