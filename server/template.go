package server

import (
	"io"
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
	e.Renderer = &templateRegistry{
		templates: template.Must(template.ParseGlob("./templates/*.html")),
	}
}
