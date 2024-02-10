package recipe

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type RecipeRequestHandler struct {
	db recipeDatabase
}

func NewRecipeRequestHandler() RecipeRequestHandler {
	return RecipeRequestHandler{
		db: newRecipeDatabase(),
	}
}

type recipeTemplateData struct {
	recipe  recipe
	isAdmin bool
}

type recipeTemplateListData struct {
	recipes recipes
	isAdmin bool
}

type recipeError = error

var (
	recipeErrorInvalidParam    recipeError = errors.New("invalid path parameter")
	recipeErrorInvalidFormData recipeError = errors.New("invalid form data")
	recipeErrorNotFound        recipeError = errors.New("no recipe found")
	recipeErrorDatabaseFailure recipeError = errors.New("database error")
)

func (r *RecipeRequestHandler) RenderRecipeListPage(c echo.Context) error {
	recipes, err := r.db.readAllRecipes()
	if err != nil {
		return r.renderError(c, err)
	}

	return r.renderTemplate(c, templateNameRecipeList, recipeTemplateListData{recipes: recipes})
}

func (r *RecipeRequestHandler) RenderNewRecipePage(c echo.Context) error {
	return r.renderTemplate(c, templateNameRecipeNew, nil)
}

func (r *RecipeRequestHandler) RenderEditRecipePage(c echo.Context) error {
	recipe, err := r.getRecipeById(c)
	if err != nil {
		return r.renderError(c, err)
	}

	return r.renderTemplate(c, templateNameRecipeEdit, recipeTemplateData{recipe: recipe})
}

func (r *RecipeRequestHandler) RenderRecipePage(c echo.Context) error {
	recipe, err := r.getRecipeById(c)
	if err != nil {
		return r.renderError(c, err)
	}

	return r.renderTemplate(c, templateNameRecipe, recipeTemplateData{recipe: recipe})
}

func (r *RecipeRequestHandler) HandleCreateRecipe(c echo.Context) error {
	if err := c.Request().ParseForm(); err != nil {
		return r.renderError(c, err)
	}

	title := c.Request().FormValue("title")
	description := c.Request().FormValue("description")

	newRecipe := recipe{Title: title, Description: description}

	if err := r.db.createRecipe(&newRecipe); err != nil {
		return r.renderError(c, err)
	}

	return r.RenderRecipeListPage(c)
}

func (r *RecipeRequestHandler) HandleUpdateRecipe(c echo.Context) error {
	updatedRecipe, err := r.parseFormData(c)
	if err != nil {
		return r.renderError(c, err)
	}

	if err := r.db.updateRecipe(&updatedRecipe); err != nil {
		return r.renderError(c, err)
	}

	return r.RenderRecipePage(c)
}

func (r *RecipeRequestHandler) HandleDeleteRecipe(c echo.Context) error {
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

func (r *RecipeRequestHandler) renderTemplate(c echo.Context, templateName string, data any) error {
	if r.isHxRequest(c) {
		return c.Render(http.StatusOK, fragmentName(templateName), data)
	}

	return c.Render(http.StatusOK, pageName(templateName), data)
}

func (r *RecipeRequestHandler) isHxRequest(c echo.Context) bool {
	hxRequestEntry := c.Request().Header["Hx-Request"]
	return len(hxRequestEntry) > 0 && hxRequestEntry[0] == "true"
}

func (r *RecipeRequestHandler) renderError(c echo.Context, err recipeError) error {
	if errors.Is(err, recipeErrorInvalidParam) {
		return c.String(http.StatusBadRequest, err.Error())
	}
	if errors.Is(err, recipeErrorInvalidFormData) {
		return c.String(http.StatusBadRequest, err.Error())
	}
	if errors.Is(err, recipeErrorNotFound) {
		return c.String(http.StatusNotFound, err.Error())
	}
	if errors.Is(err, recipeErrorDatabaseFailure) {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusInternalServerError, err.Error())
}

func (r *RecipeRequestHandler) getPathId(c echo.Context) (uint, recipeError) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return 0, recipeErrorInvalidParam
	}

	return uint(id), nil
}

func (r *RecipeRequestHandler) getRecipeById(c echo.Context) (recipe, recipeError) {
	var recipe recipe

	id, err := r.getPathId(c)
	if err != nil {
		return recipe, recipeErrorInvalidParam
	}

	return r.db.readRecipe(id)
}

func (r *RecipeRequestHandler) parseFormData(c echo.Context) (recipe, recipeError) {
	var recipe recipe

	id, err := r.getPathId(c)
	if err != nil {
		return recipe, err
	}

	if err := c.Request().ParseForm(); err != nil {
		return recipe, recipeErrorInvalidFormData
	}

	title := c.Request().FormValue("title")
	description := c.Request().FormValue("description")

	recipe.ID = id
	recipe.Title = title
	recipe.Description = description

	return recipe, nil
}
