package recipe

import (
	"net/http"

	"github.com/kilianmandscharo/lethimcook/errutil"
	"github.com/kilianmandscharo/lethimcook/routes"
	"github.com/kilianmandscharo/lethimcook/servutil"
	"github.com/labstack/echo/v4"
)

type recipeTemplateData struct {
	Recipe  recipe
	IsAdmin bool
}

type recipeTemplateListData struct {
	Recipes recipes
	IsAdmin bool
}

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
	e.POST("/recipe", rc.HandleCreateRecipe)
	e.PUT("/recipe/:id", rc.HandleUpdateRecipe)
	e.DELETE("/recipe/:id", rc.HandleDeleteRecipe)
}

func (rc *RecipeController) RenderRecipeListPage(c echo.Context) error {
	recipes, err := rc.recipeService.readAllRecipes()
	if err != nil {
		return servutil.RenderError(c, err)
	}

	return servutil.RenderTemplate(
		c,
		routes.TemplateNameRecipeList,
		recipeTemplateListData{Recipes: recipes, IsAdmin: servutil.IsAuthorized(c)},
	)
}

func (rc *RecipeController) RenderRecipeNewPage(c echo.Context) error {
	if !servutil.IsAuthorized(c) {
		return servutil.RenderError(c, errutil.AuthErrorNotAuthorized)
	}

	return servutil.RenderTemplate(
		c,
		routes.TemplateNameRecipeNew,
		servutil.IsAuthorized(c),
	)
}

func (rc *RecipeController) RenderRecipeEditPage(c echo.Context) error {
	if !servutil.IsAuthorized(c) {
		return servutil.RenderError(c, errutil.AuthErrorNotAuthorized)
	}

	recipe, err := rc.recipeService.getRecipeById(c)
	if err != nil {
		return servutil.RenderError(c, err)
	}

	return servutil.RenderTemplate(
		c,
		routes.TemplateNameRecipeEdit,
		recipeTemplateData{Recipe: recipe, IsAdmin: servutil.IsAuthorized(c)},
	)
}

func (rc *RecipeController) RenderRecipePage(c echo.Context) error {
	recipe, err := rc.recipeService.getRecipeById(c)
	if err != nil {
		return servutil.RenderError(c, err)
	}

	if err := recipe.renderMarkdown(); err != nil {
		return servutil.RenderError(c, err)
	}

	return servutil.RenderTemplate(
		c,
		routes.TemplateNameRecipe,
		recipeTemplateData{Recipe: recipe, IsAdmin: servutil.IsAuthorized(c)},
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

	err = rc.recipeService.deleteRecipe(id)
	if err != nil {
		return servutil.RenderError(c, err)
	}

	return c.String(http.StatusOK, "")
}
