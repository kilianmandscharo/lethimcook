package templutil

import (
	"errors"
	"io"
	"path"
	"text/template"

	"github.com/labstack/echo/v4"
)

type templateName = string

const (
	TemplateNameImprint       templateName = "imprint"
	TemplateNamePrivacyNotice templateName = "privacy-notice"
)

const (
	TemplateNameAdmin templateName = "admin"
)

const (
	TemplateNameRecipe     templateName = "recipe"
	TemplateNameRecipeList templateName = "recipe-list"
	TemplateNameRecipeEdit templateName = "recipe-edit"
	TemplateNameRecipeNew  templateName = "recipe-new"
)

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

func getComponentPaths(components []string) []string {
	var componentPaths []string

	for _, component := range components {
		componentPaths = append(
			componentPaths,
			path.Join("./templates/components", htmlName(component)),
		)
	}

	return componentPaths
}

func registerPageTemplate(tmap map[string]customTemplate, name templateName, components ...string) {
	pathToFragment := path.Join("./templates/pages", htmlName(name))
	componentPaths := getComponentPaths(components)

	paths := append([]string{"./templates/page.html", pathToFragment}, componentPaths...)

	// Register full page
	tmap[PageName(name)] = newCustomTemplate(
		"page.html",
		template.Must(template.ParseFiles(paths...)),
	)

	paths = append([]string{"./templates/fragment.html", pathToFragment}, componentPaths...)

	// Register fragment
	tmap[FragmentName(name)] = newCustomTemplate(
		"fragment.html",
		template.Must(template.ParseFiles(paths...)),
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

	// General
	registerPageTemplate(templates, TemplateNameImprint, "home-button", "admin-button")
	registerPageTemplate(templates, TemplateNamePrivacyNotice, "home-button", "admin-button")

	// Recipe
	registerPageTemplate(templates, TemplateNameRecipeList, "admin-button", "new-recipe-button", "home-button")
	registerPageTemplate(templates, TemplateNameRecipe, "admin-button", "home-button")
	registerPageTemplate(templates, TemplateNameRecipeNew, "admin-button", "home-button")
	registerPageTemplate(templates, TemplateNameRecipeEdit, "admin-button", "home-button")

	// Auth
	registerPageTemplate(templates, TemplateNameAdmin, "home-button", "logout-button", "admin-button")

	e.Renderer = &templateRegistry{
		templates: templates,
	}
}
