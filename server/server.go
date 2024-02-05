package server

import (
	"text/template"

	"github.com/kilianmandscharo/lethimcook/database"
	"github.com/labstack/echo/v4"
)

type Server struct {
	e *echo.Echo
}

func New(db *database.Database) Server {
	e := echo.New()
	r := newRequestHandler(db)

  e.Static("/static", "static")

	e.Renderer = &templateRegistry{
		templates: template.Must(template.ParseGlob("./templates/*.html")),
	}

	e.GET("/", r.handleHome)
	e.POST("/recipe", r.handleCreateRecipe)
	e.GET("/recipes", r.handleReadAllRecipes)
	e.PUT("/recipe", r.handleUpdateRecipe)
	e.DELETE("/recipe/:id", r.handleDeleteRecipe)

	return Server{
		e: e,
	}
}

func (s *Server) Start() {
	s.e.Logger.Fatal(s.e.Start(":8080"))
}
