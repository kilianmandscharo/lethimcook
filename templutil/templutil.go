package templutil

import (
	"errors"
	"io"
	"path"
	"runtime"
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

func getComponentPaths(components []string, dir string) []string {
	var componentPaths []string

	for _, component := range components {
		componentPaths = append(
			componentPaths,
			path.Join(dir, "templates/components", htmlName(component)),
		)
	}

	return componentPaths
}

func registerPageTemplate(tmap map[string]customTemplate, name string, dir string, components ...string) {
	pageFilePath := path.Join(dir, "templates/pages", htmlName(name))
	componentFilePaths := getComponentPaths(components, dir)

	paths := append(
		[]string{path.Join(dir, "templates/page.html"), pageFilePath},
		componentFilePaths...,
	)

	// Register full page
	tmap[PageName(name)] = newCustomTemplate(
		"page.html",
		template.Must(template.ParseFiles(paths...)),
	)

	paths = append(
		[]string{path.Join(dir, "templates/fragment.html"), pageFilePath},
		componentFilePaths...,
	)

	// Register fragment
	tmap[FragmentName(name)] = newCustomTemplate(
		"fragment.html",
		template.Must(template.ParseFiles(paths...)),
	)
}

func registerComponentTemplate(tmap map[string]customTemplate, name string, dir string, components ...string) {
	fragmentFilePath := path.Join(dir, "templates/components", htmlName(FragmentName(name)))
	componentFilePath := path.Join(dir, "templates/components", htmlName(name))
	componentPaths := getComponentPaths(components, dir)

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

func AttachTemplatesTest(e *echo.Echo) {
	_, file, _, _ := runtime.Caller(0)
	dir := path.Join(file, "../../")
	templates := make(map[string]customTemplate)

	attachPageTemplates(templates, dir)
	attachComponentTemplates(templates, dir)

	e.Renderer = &templateRegistry{
		templates: templates,
	}
}

func AttachTemplates(e *echo.Echo) {
	templates := make(map[string]customTemplate)

	attachPageTemplates(templates, ".")
	attachComponentTemplates(templates, ".")

	e.Renderer = &templateRegistry{
		templates: templates,
	}
}

func attachComponentTemplates(tmap map[string]customTemplate, dir string) {
	registerComponentTemplate(tmap, ComponentRecipeCard, dir)
	registerComponentTemplate(tmap, ComponentRecipeList, dir, ComponentRecipeCard)
}

func attachPageTemplates(tmap map[string]customTemplate, dir string) {
	registerPageTemplate(
		tmap,
		PageImprint,
		dir,
		ComponentHomeButton,
		ComponentAdminButton,
		ComponentHeading,
	)
	registerPageTemplate(
		tmap,
		PagePrivacyNotice,
		dir,
		ComponentHomeButton,
		ComponentAdminButton,
		ComponentHeading,
	)
	registerPageTemplate(
		tmap,
		PageRecipeList,
		dir,
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
		dir,
		ComponentAdminButton,
		ComponentHomeButton,
		ComponentHeading,
	)
	registerPageTemplate(
		tmap,
		PageRecipeNew,
		dir,
		ComponentAdminButton,
		ComponentHomeButton,
		ComponentHeading,
	)
	registerPageTemplate(
		tmap,
		PageRecipeEdit,
		dir,
		ComponentAdminButton,
		ComponentHomeButton,
		ComponentHeading,
	)
	registerPageTemplate(tmap,
		PageAdmin,
		dir,
		ComponentHomeButton,
		ComponentAdminButton,
		ComponentHeading,
	)
}
