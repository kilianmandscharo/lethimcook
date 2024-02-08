package server

import (
	"io"
	"log"
	"path/filepath"
	"text/template"

	"github.com/labstack/echo/v4"
)

type templateRegistry struct {
	templates *template.Template
}

func (t *templateRegistry) Render(w io.Writer, name string, data any, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func attachTemplates(e *echo.Echo) {
	files, err := filepath.Glob("./templates/*.html")

  if err != nil {
    log.Fatal("failed to find template files")
  }

	pages, err := filepath.Glob("./templates/**/*.html")

  if err != nil {
    log.Fatal("failed to find template files")
  }

  files = append(files, pages...)

	log.Println(files)

	e.Renderer = &templateRegistry{
		templates: template.Must(template.ParseFiles(files...)),
	}
}
