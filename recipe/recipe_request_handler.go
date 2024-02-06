package recipe

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type RecipeRequestHandler struct {
	db *RecipeDatabase
}

func NewRecipeRequestHandler(db *RecipeDatabase) RecipeRequestHandler {
	return RecipeRequestHandler{
		db: db,
	}
}

func (r *RecipeRequestHandler) HandleHome(c echo.Context) error {
	return c.Render(http.StatusOK, "page.html", nil)
}

func (r *RecipeRequestHandler) HandleNewRecipe(c echo.Context) error {
  return c.Render(http.StatusOK, "new-recipe.html", nil);
}

func (r *RecipeRequestHandler) HandleCreateRecipe(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func (r *RecipeRequestHandler) HandleReadAllRecipes(c echo.Context) error {
	recipes, err := r.db.ReadAllRecipes()

	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to read recipes")
	}

	return c.Render(http.StatusOK, "recipe-list.html", struct{ Recipes Recipes }{Recipes: recipes})
}

func (r *RecipeRequestHandler) HandleUpdateRecipe(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
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
