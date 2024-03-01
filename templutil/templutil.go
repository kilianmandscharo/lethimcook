package templutil

import (
	"errors"
	"io"
	"path"
	"text/template"

	"github.com/labstack/echo/v4"
)

const (
	PageImprint       = "imprint"
	PagePrivacyNotice = "privacy-notice"
)

const (
	PageAdmin = "admin"
)

const (
	PageRecipe     = "recipe"
	PageRecipeList = "recipes"
	PageRecipeEdit = "recipe-edit"
	PageRecipeNew  = "recipe-new"
)

const (
	ComponentHeading         = "heading"
	ComponentAdminButton     = "admin-button"
	ComponentHomeButton      = "home-button"
	ComponentNewRecipeButton = "new-recipe-button"
	ComponentRecipeCard      = "recipe-card"
	ComponentRecipeList      = "recipe-list"
)

func PageName(name string) string {
	return name + "-page"
}

func FragmentName(name string) string {
	return name + "-fragment"
}

func htmlName(name string) string {
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

func registerPageTemplate(tmap map[string]customTemplate, name string, components ...string) {
	pageFilePath := path.Join("./templates/pages", htmlName(name))
	componentFilePaths := getComponentPaths(components)

	paths := append(
		[]string{"./templates/page.html", pageFilePath},
		componentFilePaths...,
	)

	// Register full page
	tmap[PageName(name)] = newCustomTemplate(
		"page.html",
		template.Must(template.ParseFiles(paths...)),
	)

	paths = append(
		[]string{"./templates/fragment.html", pageFilePath},
		componentFilePaths...,
	)

	// Register fragment
	tmap[FragmentName(name)] = newCustomTemplate(
		"fragment.html",
		template.Must(template.ParseFiles(paths...)),
	)
}

func registerComponentTemplate(tmap map[string]customTemplate, name string, components ...string) {
	fragmentFilePath := path.Join("./templates/components", htmlName(FragmentName(name)))
	componentFilePath := path.Join("./templates/components", htmlName(name))
	componentPaths := getComponentPaths(components)

	paths := append(
		[]string{fragmentFilePath, componentFilePath},
		componentPaths...,
	)

	tmap[FragmentName(name)] = newCustomTemplate(
		htmlName(FragmentName(name)),
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

	attachPageTemplates(templates)
	attachComponentTemplates(templates)

	e.Renderer = &templateRegistry{
		templates: templates,
	}
}

func attachComponentTemplates(tmap map[string]customTemplate) {
	registerComponentTemplate(tmap, ComponentRecipeCard)
	registerComponentTemplate(tmap, ComponentRecipeList, ComponentRecipeCard)
}

func attachPageTemplates(tmap map[string]customTemplate) {
	registerPageTemplate(
		tmap,
		PageImprint,
		ComponentHomeButton,
		ComponentAdminButton,
		ComponentHeading,
	)
	registerPageTemplate(
		tmap,
		PagePrivacyNotice,
		ComponentHomeButton,
		ComponentAdminButton,
		ComponentHeading,
	)
	registerPageTemplate(
		tmap,
		PageRecipeList,
		ComponentRecipeList,
		ComponentRecipeCard,
		ComponentAdminButton,
		ComponentNewRecipeButton,
		ComponentHomeButton,
		ComponentHeading,
	)
	registerPageTemplate(
		tmap,
		PageRecipe,
		ComponentAdminButton,
		ComponentHomeButton,
		ComponentHeading,
	)
	registerPageTemplate(
		tmap,
		PageRecipeNew,
		ComponentAdminButton,
		ComponentHomeButton,
		ComponentHeading,
	)
	registerPageTemplate(
		tmap,
		PageRecipeEdit,
		ComponentAdminButton,
		ComponentHomeButton,
		ComponentHeading,
	)
	registerPageTemplate(tmap,
		PageAdmin,
		ComponentHomeButton,
		ComponentAdminButton,
		ComponentHeading,
	)
}
