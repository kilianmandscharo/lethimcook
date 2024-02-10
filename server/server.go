package server

import (
	"github.com/kilianmandscharo/lethimcook/auth"
	"github.com/kilianmandscharo/lethimcook/recipe"
	"github.com/kilianmandscharo/lethimcook/routes"
	"github.com/labstack/echo/v4"
)

type Server struct {
	e *echo.Echo
}

func New(authService auth.AuthService) Server {
	recipeController := recipe.NewRecipeController()
	authController := auth.NewAuthController(authService)

	e := echo.New()
	e.Use(authController.ValidateToken)
	e.Static("/static", "static")
	routes.AttachTemplates(e)

	// Pages
	e.GET("/", recipeController.RenderRecipeListPage)
	e.GET("/recipe/edit/:id", recipeController.RenderRecipeEditPage)
	e.GET("/recipe/new", recipeController.RenderRecipeNewPage)
	e.GET("/recipe/:id", recipeController.RenderRecipePage)

	// Actions
	e.POST("/recipe", recipeController.HandleCreateRecipe)
	e.PUT("/recipe/:id", recipeController.HandleUpdateRecipe)
	e.DELETE("/recipe/:id", recipeController.HandleDeleteRecipe)

	// Pages
	e.GET("/auth/admin", authController.RenderAdminPage)

	// Actions
	e.POST("/auth/login", authController.HandleLogin)
	e.POST("/auth/logout", authController.HandleLogout)
	e.PUT("/auth/password", authController.HandleUpdatePassword)

	return Server{
		e: e,
	}
}

func (s *Server) Start() {
	s.e.Logger.Fatal(s.e.Start(":8080"))
}
