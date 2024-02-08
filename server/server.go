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
	attachTemplates(e)

	// e.Use(Log)

	e.GET("/", recipeRequestHandler.HandleHome)
	e.GET("/recipe/new", recipeRequestHandler.HandleNewRecipe)
	e.GET("/recipe/edit/:id", recipeRequestHandler.HandleEditRecipe)
	e.POST("/recipe", recipeRequestHandler.HandleCreateRecipe)
	e.GET("/recipe/all", recipeRequestHandler.HandleReadAllRecipes)
	e.GET("/recipe/:id", recipeRequestHandler.HandleReadRecipe)
	e.PUT("/recipe/:id", recipeRequestHandler.HandleUpdateRecipe)
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

// func Log(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		log.Println("Hx-Request:", c.Request().Header["Hx-Request"])
// 		return next(c)
// 	}
// }
