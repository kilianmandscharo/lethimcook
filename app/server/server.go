package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/kilianmandscharo/lethimcook/auth"
	"github.com/kilianmandscharo/lethimcook/components"
	"github.com/kilianmandscharo/lethimcook/env"
	"github.com/kilianmandscharo/lethimcook/logging"
	"github.com/kilianmandscharo/lethimcook/recipe"
	"github.com/kilianmandscharo/lethimcook/render"
	"github.com/kilianmandscharo/lethimcook/servutil"
	"github.com/labstack/echo/v4"
)

type Server struct {
	e        *echo.Echo
	logger   *logging.Logger
	renderer *render.Renderer
	isProd   bool
}

func New(
	authController *auth.AuthController,
	recipeController *recipe.RecipeController,
	logger *logging.Logger,
	renderer *render.Renderer,
	isProd bool,
) Server {
	e := echo.New()

	e.Use(authController.ValidateTokenMiddleware)
	e.Static("/static", "./static")

	e.GET("/imprint", func(c echo.Context) error {
		return renderer.RenderComponent(render.RenderComponentOptions{
			Context:   c,
			Component: components.Imprint(servutil.IsAuthorized(c)),
		})
	})
	e.GET("/privacy-notice", func(c echo.Context) error {
		return renderer.RenderComponent(render.RenderComponentOptions{
			Context:   c,
			Component: components.PrivacyNotice(servutil.IsAuthorized(c)),
		})
	})
	e.GET("/info", func(c echo.Context) error {
		return renderer.RenderComponent(render.RenderComponentOptions{
			Context:   c,
			Component: components.Information(servutil.IsAuthorized(c)),
		})
	})

	recipeController.AttachHandlerFunctions(e)
	authController.AttachHandlerFunctions(e)

	return Server{
		e:        e,
		logger:   logger,
		renderer: renderer,
		isProd:   isProd,
	}
}

func (s *Server) Start() {
	certFilePath := env.Get(env.EnvKeyCertFilePath)
	keyFilePath := env.Get(env.EnvKeyKeyFilePath)

	if s.isProd {
		if len(certFilePath) == 0 || len(keyFilePath) == 0 {
			s.logger.Fatal("env variables certFilePath and keyFilePath need to be defined when running in production mode")
		}
		s.startProd(certFilePath, keyFilePath)
	} else {
		s.startDev()
	}
}

func (s *Server) startDev() {
	s.logger.Info("Starting server and listening on :8080")
	go func() {
		s.logger.Fatal(s.e.Start(":8080"))
	}()
	s.listenForShutdown()
}

func (s *Server) startProd(certFilePath, keyFilePath string) {
	s.logger.Info("Starting server and listening on :443/:80")
	go func() {
		s.logger.Fatal(s.e.StartTLS(":443", certFilePath, keyFilePath))
	}()
	go func() {
		s.logger.Fatal(http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
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
