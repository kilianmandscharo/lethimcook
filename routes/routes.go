package routes

import (
	"errors"
	"io"
	"path"
	"text/template"

	"github.com/labstack/echo/v4"
)

type templateName = string

func PageName(name templateName) string {
	return name + "-page"
}

func FragmentName(name templateName) string {
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

func registerPageTemplate(tmap map[string]customTemplate, name templateName) {
	pathToFragment := path.Join("templates/pages", htmlName(name))

	// Register full page
	tmap[PageName(name)] = newCustomTemplate(
		"page.html",
		template.Must(template.ParseFiles("templates/page.html", pathToFragment)),
	)

	// Register fragment
	tmap[FragmentName(name)] = newCustomTemplate(
		"fragment.html",
		template.Must(template.ParseFiles("templates/fragment.html", pathToFragment)),
	)
}

func registerComponentTemplate(tmap map[string]customTemplate, name templateName) {
	pathToFragment := path.Join("templates/components", htmlName(name))

	tmap[FragmentName(name)] = newCustomTemplate(
		htmlName(name),
		template.Must(template.ParseFiles(pathToFragment)),
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

	// Recipe
	registerPageTemplate(templates, TemplateNameRecipeList)
	registerPageTemplate(templates, TemplateNameRecipe)
	registerPageTemplate(templates, TemplateNameRecipeNew)
	registerPageTemplate(templates, TemplateNameRecipeEdit)

	// Auth
	registerPageTemplate(templates, TemplateNameAdmin)
	registerComponentTemplate(templates, TemplateNameLoginSuccessful)

	e.Renderer = &templateRegistry{
		templates: templates,
	}
}
