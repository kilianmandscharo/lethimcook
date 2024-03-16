package recipe

import (
	"fmt"
	"strings"

	"github.com/kilianmandscharo/lethimcook/components"
	"github.com/kilianmandscharo/lethimcook/errutil"
	"github.com/kilianmandscharo/lethimcook/servutil"
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
	recipes, err := rc.recipeService.readAllRecipes()
	if err != nil {
		return servutil.RenderError(c, err)
	}

	return servutil.RenderComponent(
		c,
		components.RecipesPage(servutil.IsAuthorized(c), false, recipes),
	)
}

func (rc *RecipeController) RenderRecipeNewPage(c echo.Context) error {
	if !servutil.IsAuthorized(c) {
		return servutil.RenderError(c, errutil.AuthErrorNotAuthorized)
	}

	return servutil.RenderComponent(c, components.RecipeNewPage())
}

func (rc *RecipeController) RenderRecipeEditPage(c echo.Context) error {
	if !servutil.IsAuthorized(c) {
		return servutil.RenderError(c, errutil.AuthErrorNotAuthorized)
	}

	recipe, err := rc.recipeService.getRecipeById(c)
	if err != nil {
		return servutil.RenderError(c, err)
	}

	return servutil.RenderComponent(c, components.RecipeEditPage(recipe))
}

func (rc *RecipeController) RenderRecipePage(c echo.Context) error {
	recipe, err := rc.recipeService.getRecipeById(c)
	if err != nil {
		return servutil.RenderError(c, err)
	}

	if err := recipe.RenderMarkdown(); err != nil {
		return servutil.RenderError(c, err)
	}

	c.Response().Header().Set("HX-Push-Url", fmt.Sprintf("/recipe/%d", recipe.ID))

	return servutil.RenderComponent(
		c,
		components.RecipePage(servutil.IsAuthorized(c), recipe),
	)
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

	return servutil.RenderComponent(
		c,
		components.RecipeList(servutil.IsAuthorized(c), false, filteredRecipes),
	)
}

func (rc *RecipeController) HandleCreateRecipe(c echo.Context) error {
	if !servutil.IsAuthorized(c) {
		return servutil.RenderError(c, errutil.AuthErrorNotAuthorized)
	}

	newRecipe, err := rc.recipeService.parseFormData(c, false)
	if err != nil {
		return servutil.RenderError(c, err)
	}

	if err := rc.recipeService.createRecipe(&newRecipe); err != nil {
		return servutil.RenderError(c, err)
	}

	return rc.RenderRecipeListPage(c)
}

func (rc *RecipeController) HandleUpdateRecipe(c echo.Context) error {
	if !servutil.IsAuthorized(c) {
		return servutil.RenderError(c, errutil.AuthErrorNotAuthorized)
	}

	updatedRecipe, err := rc.recipeService.parseFormData(c, true)
	if err != nil {
		return servutil.RenderError(c, err)
	}

	if err := rc.recipeService.updateRecipe(&updatedRecipe); err != nil {
		return servutil.RenderError(c, err)
	}

	return rc.RenderRecipePage(c)
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

		return servutil.RenderComponent(
			c, components.RecipeList(servutil.IsAuthorized(c), false, recipes),
		)
	}

	var deleting = true

	if c.QueryParam("cancel") == "true" {
		deleting = false
	}

	recipe, err := rc.recipeService.readRecipe(id)
	if err != nil {
		return servutil.RenderError(c, err)
	}

	return servutil.RenderComponent(
		c,
		components.RecipeCard(servutil.IsAuthorized(c), deleting, recipe),
	)
}
