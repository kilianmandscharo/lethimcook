package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/kilianmandscharo/lethimcook/auth"
	"github.com/kilianmandscharo/lethimcook/env"
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
	certFilePath := env.Get(env.EnvKeyCertFilePath)
	keyFilePath := env.Get(env.EnvKeyKeyFilePath)

	if len(certFilePath) > 0 && len(keyFilePath) > 0 {
		s.startProd(certFilePath, keyFilePath)
	} else {
		s.startDev()
	}
}

func (s *Server) startDev() {
	if err := s.e.Start(":80"); err != nil {
		s.e.Logger.Fatal("Error starting development server: ", err)
	}
}

func (s *Server) startProd(certFilePath, keyFilePath string) {
	go func() {
		if err := s.e.StartTLS(":443", certFilePath, keyFilePath); err != nil {
			s.e.Logger.Fatal("Error starting server: ", err)
		}
	}()

	go func() {
		if err := http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			http.Redirect(w, req, "https://"+req.Host+req.URL.String(), http.StatusMovedPermanently)
		})); err != nil {
			s.e.Logger.Fatal("Error starting redirect server: ", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.e.Shutdown(ctx); err != nil {
		s.e.Logger.Fatal("Error shutting down server: ", err)
	}
}
