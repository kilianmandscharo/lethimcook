package recipe

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type RecipeRequestHandler struct {
	db RecipeDatabase
}

func NewRecipeRequestHandler() RecipeRequestHandler {
	return RecipeRequestHandler{
		db: NewRecipeDatabase(),
	}
}

func (r *RecipeRequestHandler) RenderRecipeListPage(c echo.Context) error {
	recipes, err := r.db.ReadAllRecipes()

	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to read recipes")
	}

	data := struct{ Recipes Recipes }{Recipes: recipes}

	if isHxRequest(c) {
		return c.Render(http.StatusOK, fragmentName(templateNameRecipeList), data)
	}

	return c.Render(http.StatusOK, pageName(templateNameRecipeList), data)
}

func (r *RecipeRequestHandler) RenderNewRecipePage(c echo.Context) error {
	if isHxRequest(c) {
		return c.Render(http.StatusOK, fragmentName(templateNameRecipeNew), nil)
	}

	return c.Render(http.StatusOK, pageName(templateNameRecipeNew), nil)
}

func (r *RecipeRequestHandler) RenderEditRecipePage(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid path parameter")
	}

	recipe, err := r.db.ReadRecipe(uint(id))
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to read recipe")
	}

	data := struct {
		ID                           uint
		Title, Description, Markdown string
	}{
		ID:          recipe.ID,
		Title:       recipe.Title,
		Description: recipe.Description,
	}

	if isHxRequest(c) {
		return c.Render(http.StatusOK, fragmentName(templateNameRecipeEdit), data)
	}

	return c.Render(http.StatusOK, pageName(templateNameRecipeEdit), data)
}

func (r *RecipeRequestHandler) RenderRecipePage(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid path parameter")
	}

	recipe, err := r.db.ReadRecipe(uint(id))
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to read recipe")
	}

	data := struct {
		ID                           uint
		Title, Description, Markdown string
	}{
		ID:          recipe.ID,
		Title:       recipe.Title,
		Description: recipe.Description,
	}

	if isHxRequest(c) {
		return c.Render(http.StatusOK, fragmentName(templateNameRecipe), data)
	}

	return c.Render(http.StatusOK, pageName(templateNameRecipe), data)
}

func (r *RecipeRequestHandler) HandleCreateRecipe(c echo.Context) error {
	if err := c.Request().ParseForm(); err != nil {
		return c.String(http.StatusBadRequest, "failed to parse form data")
	}

	title := c.Request().FormValue("title")
	description := c.Request().FormValue("description")

	newRecipe := Recipe{Title: title, Description: description}

	if err := r.db.CreateRecipe(&newRecipe); err != nil {
		return c.String(http.StatusInternalServerError, "failed to create new recipe")
	}

	return c.String(http.StatusOK, "created recipe")
}

func (r *RecipeRequestHandler) HandleReadAllRecipes(c echo.Context) error {
	recipes, err := r.db.ReadAllRecipes()

	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to read recipes")
	}

	data := struct{ Recipes Recipes }{Recipes: recipes}

	return c.Render(http.StatusOK, "recipe-list.html", data)
}

func (r *RecipeRequestHandler) HandleUpdateRecipe(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid path parameter")
	}

	if err := c.Request().ParseForm(); err != nil {
		return c.String(http.StatusBadRequest, "failed to parse form data")
	}

	title := c.Request().FormValue("title")
	description := c.Request().FormValue("description")

	updatedRecipe := Recipe{ID: uint(id), Title: title, Description: description}

	if err := r.db.UpdateRecipe(&updatedRecipe); err != nil {
		return c.String(http.StatusInternalServerError, "failed to update recipe")
	}

	return c.String(http.StatusOK, "created recipe")
}

func (r *RecipeRequestHandler) HandleDeleteRecipe(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.String(http.StatusBadRequest, "invalid path parameter")
	}

	err = r.db.DeleteRecipe(uint(id))

	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to delete recipe")
	}

	return c.String(http.StatusOK, "")
}

func isHxRequest(c echo.Context) bool {
	hxRequestEntry := c.Request().Header["Hx-Request"]
	return len(hxRequestEntry) > 0 && hxRequestEntry[0] == "true"
}
