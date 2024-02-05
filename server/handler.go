package server

import (
	"net/http"
	"strconv"

	"github.com/kilianmandscharo/lethimcook/database"
	"github.com/kilianmandscharo/lethimcook/recipe"
	"github.com/labstack/echo/v4"
)

type RequestHandler struct {
	db *database.Database
}

func newRequestHandler(db *database.Database) RequestHandler {
	return RequestHandler{
		db: db,
	}
}

func (r *RequestHandler) handleHome(c echo.Context) error {
	return c.Render(http.StatusOK, "page.html", nil)
}

func (r *RequestHandler) handleCreateRecipe(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func (r *RequestHandler) handleReadAllRecipes(c echo.Context) error {
	recipes, err := r.db.ReadAllRecipes()

	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to read recipes")
	}

	return c.Render(http.StatusOK, "recipe-list.html", struct{ Recipes recipe.Recipes }{Recipes: recipes})
}

func (r *RequestHandler) handleUpdateRecipe(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func (r *RequestHandler) handleDeleteRecipe(c echo.Context) error {
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
