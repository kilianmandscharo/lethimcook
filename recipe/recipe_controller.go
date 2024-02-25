package recipe

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/kilianmandscharo/lethimcook/errutil"
	"github.com/kilianmandscharo/lethimcook/servutil"
	"github.com/kilianmandscharo/lethimcook/templutil"
	"github.com/labstack/echo/v4"
)

type recipeTemplateData struct {
	Recipe  recipe
	IsAdmin bool
}

type recipeTemplateListItemData struct {
	recipe
	IsAdmin  bool
	Deleting bool
}

type recipeTemplateListData struct {
	Recipes []recipeTemplateListItemData
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

	isAdmin := servutil.IsAuthorized(c)
	recipeData := make([]recipeTemplateListItemData, len(recipes))

	for i, recipe := range recipes {
		recipeData[i] = recipeTemplateListItemData{
			recipe:  recipe,
			IsAdmin: isAdmin,
		}
	}

	return servutil.RenderTemplate(
		c,
		templutil.PageRecipeList,
		recipeTemplateListData{Recipes: recipeData, IsAdmin: isAdmin},
	)
}

func (rc *RecipeController) RenderRecipeNewPage(c echo.Context) error {
	if !servutil.IsAuthorized(c) {
		return servutil.RenderError(c, errutil.AuthErrorNotAuthorized)
	}

	return servutil.RenderTemplate(
		c,
		templutil.PageRecipeNew,
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
		templutil.PageRecipeEdit,
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

	c.Response().Header().Set("HX-Push-Url", fmt.Sprintf("/recipe/%d", recipe.ID))

	return servutil.RenderTemplate(
		c,
		templutil.PageRecipe,
		recipeTemplateData{Recipe: recipe, IsAdmin: servutil.IsAuthorized(c)},
	)
}

func (rc *RecipeController) HandleSearchRecipe(c echo.Context) error {
	if err := c.Request().ParseForm(); err != nil {
		servutil.RenderError(c, err)
	}
	query := strings.ToLower(c.FormValue("query"))

	recipes, err := rc.recipeService.readAllRecipes()
	if err != nil {
		return servutil.RenderError(c, err)
	}

	var filteredRecipes []recipe
	for _, recipe := range recipes {
		if strings.Contains(strings.ToLower(recipe.Title), query) {
			filteredRecipes = append(filteredRecipes, recipe)
		}
	}

	isAdmin := servutil.IsAuthorized(c)
	recipeData := make([]recipeTemplateListItemData, len(filteredRecipes))

	for i, recipe := range filteredRecipes {
		recipeData[i] = recipeTemplateListItemData{
			recipe:  recipe,
			IsAdmin: isAdmin,
		}
	}

	return servutil.RenderTemplateComponent(
		c,
		templutil.ComponentRecipeList,
		recipeTemplateListData{Recipes: recipeData, IsAdmin: isAdmin},
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

		return c.String(http.StatusOK, "")
	}

	recipe, err := rc.recipeService.readRecipe(id)
	if err != nil {
		return servutil.RenderError(c, err)
	}

	var deleting = true

	if c.QueryParam("cancel") == "true" {
		deleting = false
	}

	return servutil.RenderTemplateComponent(
		c,
		templutil.ComponentRecipeCard,
		recipeTemplateListItemData{
			recipe:   recipe,
			IsAdmin:  servutil.IsAuthorized(c),
			Deleting: deleting,
		})
}
