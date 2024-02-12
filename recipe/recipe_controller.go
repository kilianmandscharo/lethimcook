package recipe

import (
	"net/http"
	"strconv"

	"github.com/kilianmandscharo/lethimcook/errutil"
	"github.com/kilianmandscharo/lethimcook/routes"
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
	db recipeDatabase
}

func NewRecipeController() RecipeController {
	return RecipeController{
		db: newRecipeDatabase(),
	}
}

func (r *RecipeController) AttachHandlerFunctions(e *echo.Echo) {
	// Pages
	e.GET("/", r.RenderRecipeListPage)
	e.GET("/recipe/edit/:id", r.RenderRecipeEditPage)
	e.GET("/recipe/new", r.RenderRecipeNewPage)
	e.GET("/recipe/:id", r.RenderRecipePage)

	// Actions
	e.POST("/recipe", r.HandleCreateRecipe)
	e.PUT("/recipe/:id", r.HandleUpdateRecipe)
	e.DELETE("/recipe/:id", r.HandleDeleteRecipe)
}

func (r *RecipeController) RenderRecipeListPage(c echo.Context) error {
	recipes, err := r.db.readAllRecipes()
	if err != nil {
		return r.renderError(c, err)
	}

	return r.renderTemplate(c, routes.TemplateNameRecipeList, recipeTemplateListData{Recipes: recipes, IsAdmin: r.isAdmin(c)})
}

func (r *RecipeController) RenderRecipeNewPage(c echo.Context) error {
	return r.renderTemplate(c, routes.TemplateNameRecipeNew, r.isAdmin(c))
}

func (r *RecipeController) RenderRecipeEditPage(c echo.Context) error {
	recipe, err := r.getRecipeById(c)
	if err != nil {
		return r.renderError(c, err)
	}

	return r.renderTemplate(c, routes.TemplateNameRecipeEdit, recipeTemplateData{Recipe: recipe, IsAdmin: r.isAdmin(c)})
}

func (r *RecipeController) RenderRecipePage(c echo.Context) error {
	recipe, err := r.getRecipeById(c)
	if err != nil {
		return r.renderError(c, err)
	}

	if err := recipe.renderMarkdown(); err != nil {
		return r.renderError(c, err)
	}

	return r.renderTemplate(c, routes.TemplateNameRecipe, recipeTemplateData{Recipe: recipe, IsAdmin: r.isAdmin(c)})
}

func (r *RecipeController) HandleCreateRecipe(c echo.Context) error {
	if !r.isAdmin(c) {
		return c.String(http.StatusUnauthorized, "not authorized")
	}

	newRecipe, err := r.parseFormData(c, false)
	if err != nil {
		return r.renderError(c, err)
	}

	if err := r.db.createRecipe(&newRecipe); err != nil {
		return r.renderError(c, err)
	}

	return r.RenderRecipeListPage(c)
}

func (r *RecipeController) HandleUpdateRecipe(c echo.Context) error {
	if !r.isAdmin(c) {
		return c.String(http.StatusUnauthorized, "not authorized")
	}

	updatedRecipe, err := r.parseFormData(c, true)
	if err != nil {
		return r.renderError(c, err)
	}

	if err := r.db.updateRecipe(&updatedRecipe); err != nil {
		return r.renderError(c, err)
	}

	return r.RenderRecipePage(c)
}

func (r *RecipeController) HandleDeleteRecipe(c echo.Context) error {
	if !r.isAdmin(c) {
		return c.String(http.StatusUnauthorized, "not authorized")
	}

	id, err := r.getPathId(c)
	if err != nil {
		return r.renderError(c, err)
	}

	err = r.db.deleteRecipe(id)
	if err != nil {
		return r.renderError(c, err)
	}

	return c.String(http.StatusOK, "")
}

func (r *RecipeController) renderTemplate(c echo.Context, templateName string, data any) error {
	if r.isHxRequest(c) {
		return c.Render(http.StatusOK, routes.FragmentName(templateName), data)
	}

	return c.Render(http.StatusOK, routes.PageName(templateName), data)
}

func (r *RecipeController) isHxRequest(c echo.Context) bool {
	hxRequestEntry := c.Request().Header["Hx-Request"]
	return len(hxRequestEntry) > 0 && hxRequestEntry[0] == "true"
}

func (r *RecipeController) isAdmin(c echo.Context) bool {
	authorized := c.Get("authorized")
	if isAdmin, ok := authorized.(bool); ok {
		return isAdmin
	}
	return false
}

func (r *RecipeController) renderError(c echo.Context, err errutil.RecipeError) error {
	return c.String(errutil.RecipeErrorHttpCodes[err], err.Error())
}

func (r *RecipeController) getPathId(c echo.Context) (uint, errutil.RecipeError) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return 0, errutil.RecipeErrorInvalidParam
	}

	return uint(id), nil
}

func (r *RecipeController) getRecipeById(c echo.Context) (recipe, errutil.RecipeError) {
	var recipe recipe

	id, err := r.getPathId(c)
	if err != nil {
		return recipe, errutil.RecipeErrorInvalidParam
	}

	return r.db.readRecipe(id)
}

func (r *RecipeController) parseFormData(c echo.Context, withID bool) (recipe, errutil.RecipeError) {
	var recipe recipe

	if err := c.Request().ParseForm(); err != nil {
		return recipe, errutil.RecipeErrorInvalidFormData
	}

	recipe.Title = c.Request().FormValue("title")
	recipe.Description = c.Request().FormValue("description")
	recipe.Ingredients = c.Request().FormValue("ingredients")
	recipe.Instructions = c.Request().FormValue("instructions")

	duration, err := strconv.Atoi(c.Request().FormValue("duration"))
	if err != nil {
		return recipe, errutil.RecipeErrorInvalidFormData
	}
	recipe.Duration = duration

	if withID {
		id, err := r.getPathId(c)
		if err != nil {
			return recipe, err
		}
		recipe.ID = id
	}

	return recipe, nil
}
