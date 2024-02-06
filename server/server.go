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
	recipeDatabase := recipe.NewRecipeDatabase()
	recipeRequestHandler := recipe.NewRecipeRequestHandler(&recipeDatabase)

	authDatabase := auth.NewAuthDatabase()
	authRequestHandler := auth.NewAuthRequestHandler(&authDatabase)

	e := echo.New()
	e.Static("/static", "static")
	attachTemplates(e)

	e.GET("/", recipeRequestHandler.HandleHome)
	e.GET("/recipe/new", recipeRequestHandler.HandleNewRecipe)
	e.POST("/recipe", recipeRequestHandler.HandleCreateRecipe)
	e.GET("/recipe/all", recipeRequestHandler.HandleReadAllRecipes)
	e.PUT("/recipe", recipeRequestHandler.HandleUpdateRecipe)
	e.DELETE("/recipe/:id", recipeRequestHandler.HandleDeleteRecipe)

	e.POST("/auth/login", authRequestHandler.HandleLogin)
	e.PUT("/auth/password", authRequestHandler.HandleUpdatePassword)

	return Server{
		e: e,
	}
}

func (s *Server) Start() {
	s.e.Logger.Fatal(s.e.Start(":8080"))
}
