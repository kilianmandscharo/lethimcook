package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/kilianmandscharo/lethimcook/auth"
	"github.com/kilianmandscharo/lethimcook/env"
	"github.com/kilianmandscharo/lethimcook/recipe"
	"github.com/kilianmandscharo/lethimcook/servutil"
	"github.com/labstack/echo/v4"
)

type Server struct {
	e *echo.Echo
}

func New(authController auth.AuthController, recipeController recipe.RecipeController) Server {
	e := echo.New()

	e.Use(authController.ValidateTokenMiddleware)
	e.Static("/static", "./static")

	servutil.AttachHandlerFunctions(e)
	recipeController.AttachHandlerFunctions(e)
	authController.AttachHandlerFunctions(e)

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
	go func() {
		log.Fatal(s.e.Start(":8080"))
	}()

	s.listenForShutdown()
}

func (s *Server) startProd(certFilePath, keyFilePath string) {
	go func() {
		s.e.Logger.Fatal(s.e.StartTLS(":443", certFilePath, keyFilePath))
	}()

	go func() {
		log.Fatal(http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			http.Redirect(w, req, "https://"+req.Host+req.URL.String(), http.StatusMovedPermanently)
		})))
	}()

	s.listenForShutdown()
}

func (s *Server) listenForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.e.Shutdown(ctx); err != nil {
		log.Fatal("Error shutting down server: ", err)
	}
}
