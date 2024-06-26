package recipe

import (
	"fmt"
	"strings"

	"github.com/kilianmandscharo/lethimcook/components"
	"github.com/kilianmandscharo/lethimcook/errutil"
	"github.com/kilianmandscharo/lethimcook/servutil"
	"github.com/kilianmandscharo/lethimcook/types"
	"github.com/labstack/echo/v4"
)

type RecipeController struct {
	recipeService recipeService
}

func NewRecipeController() RecipeController {
	return RecipeController{
		recipeService: newRecipeService(),
	}
}

func (rc *RecipeController) AttachHandlerFunctions(e *echo.Echo) {
	// Pages
	e.GET("/", rc.RenderRecipeListPage)
	e.GET("/recipe/edit/:id", rc.RenderRecipeEditPage)
	e.GET("/recipe/new", rc.RenderRecipeNewPage)
	e.GET("/recipe/:id", rc.RenderRecipePage)

	// Actions
	e.POST("/search", rc.HandleSearchRecipe)
	e.POST("/recipe", rc.HandleCreateRecipe)
	e.PUT("/recipe/:id", rc.HandleUpdateRecipe)
	e.DELETE("/recipe/:id", rc.HandleDeleteRecipe)
}

func (rc *RecipeController) RenderRecipeListPage(c echo.Context) error {
	return rc.renderRecipeListPageHelper(c, "")
}

func (rc *RecipeController) renderRecipeListPageHelper(c echo.Context, message string) error {
	recipes, err := rc.recipeService.readAllRecipes()
	if err != nil {
		return servutil.RenderError(c, err)
	}

	return servutil.RenderComponent(servutil.RenderComponentOptions{
		Context:   c,
		Component: components.RecipesPage(servutil.IsAuthorized(c), false, recipes),
		Message:   message,
	})
}

func (rc *RecipeController) RenderRecipeNewPage(c echo.Context) error {
	if !servutil.IsAuthorized(c) {
		return servutil.RenderError(c, errutil.AuthErrorNotAuthorized)
	}

	formElements := rc.recipeService.createRecipeForm(types.Recipe{}, make(map[string]error))
	return servutil.RenderComponent(servutil.RenderComponentOptions{
		Context:   c,
		Component: components.RecipeNewPage(formElements),
	})
}

func (rc *RecipeController) RenderRecipeEditPage(c echo.Context) error {
	if !servutil.IsAuthorized(c) {
		return servutil.RenderError(c, errutil.AuthErrorNotAuthorized)
	}

	recipe, err := rc.recipeService.getRecipeById(c)
	if err != nil {
		return servutil.RenderError(c, err)
	}

	formElements := rc.recipeService.createRecipeForm(recipe, make(map[string]error))
	return servutil.RenderComponent(servutil.RenderComponentOptions{
		Context:   c,
		Component: components.RecipeEditPage(recipe.ID, formElements),
	})
}

func (rc *RecipeController) RenderRecipePage(c echo.Context) error {
	return rc.RenderRecipePageHelper(c, "")
}

func (rc *RecipeController) RenderRecipePageHelper(c echo.Context, message string) error {
	recipe, err := rc.recipeService.getRecipeById(c)
	if err != nil {
		return servutil.RenderError(c, err)
	}

	if err := recipe.RenderMarkdown(); err != nil {
		return servutil.RenderError(c, err)
	}

	c.Response().Header().Set("HX-Push-Url", fmt.Sprintf("/recipe/%d", recipe.ID))

	return servutil.RenderComponent(servutil.RenderComponentOptions{
		Context:   c,
		Component: components.RecipePage(servutil.IsAuthorized(c), recipe),
		Message:   message,
	})
}

func (rc *RecipeController) HandleSearchRecipe(c echo.Context) error {
	if err := c.Request().ParseForm(); err != nil {
		return servutil.RenderError(c, err)
	}
	query := strings.ToLower(c.FormValue("query"))

	filteredRecipes, err := rc.recipeService.getFilteredRecipes(query)
	if err != nil {
		return servutil.RenderError(c, err)
	}

	return servutil.RenderComponent(servutil.RenderComponentOptions{
		Context:   c,
		Component: components.RecipeList(servutil.IsAuthorized(c), false, filteredRecipes),
	})
}

func (rc *RecipeController) HandleCreateRecipe(c echo.Context) error {
	if !servutil.IsAuthorized(c) {
		return servutil.RenderError(c, errutil.AuthErrorNotAuthorized)
	}

	newRecipe, formErrors, err := rc.recipeService.parseFormData(c, false)
	if err != nil {
		return servutil.RenderError(c, err)
	}

	if len(formErrors) > 0 {
		formElements := rc.recipeService.createRecipeForm(newRecipe, formErrors)
		return servutil.RenderComponent(servutil.RenderComponentOptions{
			Context:   c,
			Component: components.RecipeNewPage(formElements),
			Message:   "Fehlerhaftes Formular",
			IsError:   true,
		})
	}

	if err := rc.recipeService.createRecipe(&newRecipe); err != nil {
		return servutil.RenderError(c, err)
	}

	return rc.renderRecipeListPageHelper(c, "Rezept erstellt")
}

func (rc *RecipeController) HandleUpdateRecipe(c echo.Context) error {
	if !servutil.IsAuthorized(c) {
		return servutil.RenderError(c, errutil.AuthErrorNotAuthorized)
	}

	updatedRecipe, formErrors, err := rc.recipeService.parseFormData(c, true)
	if err != nil {
		return servutil.RenderError(c, err)
	}

	if len(formErrors) > 0 {
		formElements := rc.recipeService.createRecipeForm(updatedRecipe, formErrors)
		return servutil.RenderComponent(servutil.RenderComponentOptions{
			Context:   c,
			Component: components.RecipeEditPage(updatedRecipe.ID, formElements),
			Message:   "Fehlerhaftes Formular",
			IsError:   true,
		})
	}

	if err := rc.recipeService.updateRecipe(&updatedRecipe); err != nil {
		return servutil.RenderError(c, err)
	}

	return rc.RenderRecipePageHelper(c, "Rezept aktualisiert")
}

func (rc *RecipeController) HandleDeleteRecipe(c echo.Context) error {
	if !servutil.IsAuthorized(c) {
		return servutil.RenderError(c, errutil.AuthErrorNotAuthorized)
	}

	id, err := rc.recipeService.getPathId(c)
	if err != nil {
		return servutil.RenderError(c, err)
	}

	if c.QueryParam("force") == "true" {
		err = rc.recipeService.deleteRecipe(id)
		if err != nil {
			return servutil.RenderError(c, err)
		}

		recipes, err := rc.recipeService.readAllRecipes()
		if err != nil {
			return servutil.RenderError(c, err)
		}

		return servutil.RenderComponent(servutil.RenderComponentOptions{
			Context:   c,
			Component: components.RecipeList(servutil.IsAuthorized(c), false, recipes),
			Message:   "Rezept gelöscht",
		})
	}

	var deleting = true

	if c.QueryParam("cancel") == "true" {
		deleting = false
	}

	recipe, err := rc.recipeService.readRecipe(id)
	if err != nil {
		return servutil.RenderError(c, err)
	}

	return servutil.RenderComponent(servutil.RenderComponentOptions{
		Context:   c,
		Component: components.RecipeCard(servutil.IsAuthorized(c), deleting, recipe),
	})
}
