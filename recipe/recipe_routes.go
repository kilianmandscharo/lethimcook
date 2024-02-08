package recipe

import (
	"errors"
	"io"
	"path"
	"text/template"

	"github.com/labstack/echo/v4"
)

type templateName = string

const (
	templateNameRecipe     templateName = "recipe"
	templateNameRecipeList templateName = "recipe-list"
	templateNameRecipeEdit templateName = "recipe-edit"
	templateNameRecipeNew  templateName = "recipe-new"
)

func pageName(name templateName) string {
	return name + "-page"
}

func fragmentName(name templateName) string {
	return name + "-fragment"
}

func htmlName(name templateName) string {
	return name + ".html"
}

type templateRegistry struct {
	templates map[string]customTemplate
}

type customTemplate struct {
	baseName  string
	templates *template.Template
}

func newCustomTemplate(baseName string, templates *template.Template) customTemplate {
	return customTemplate{
		baseName:  baseName,
		templates: templates,
	}
}

func registerTemplate(tmap map[string]customTemplate, name templateName) {
	pathToFragment := path.Join("templates/pages", htmlName(name))

	// Register full page
	tmap[pageName(name)] = newCustomTemplate(
		"page.html",
		template.Must(template.ParseFiles("templates/page.html", pathToFragment)),
	)

	// Register fragment
	tmap[fragmentName(name)] = newCustomTemplate(
		"fragment.html",
		template.Must(template.ParseFiles("templates/fragment.html", pathToFragment)),
	)
}

func (t *templateRegistry) Render(w io.Writer, name string, data any, c echo.Context) error {
	if tmpl, ok := t.templates[name]; ok {
		return tmpl.templates.ExecuteTemplate(w, tmpl.baseName, data)
	}

	return errors.New("Template not found -> " + name)
}

func AttachTemplates(e *echo.Echo) {
	templates := make(map[string]customTemplate)

	registerTemplate(templates, templateNameRecipeList)
	registerTemplate(templates, templateNameRecipe)
	registerTemplate(templates, templateNameRecipeNew)
	registerTemplate(templates, templateNameRecipeEdit)

	e.Renderer = &templateRegistry{
		templates: templates,
	}
}
