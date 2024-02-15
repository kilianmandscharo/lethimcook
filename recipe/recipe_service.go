package recipe

import (
	"strconv"

	"github.com/kilianmandscharo/lethimcook/errutil"
	"github.com/labstack/echo/v4"
)

type recipeService struct {
	db recipeDatabase
}

func newRecipeService() recipeService {
	return recipeService{
		db: newRecipeDatabase(),
	}
}

func (rs *recipeService) createRecipe(recipe *recipe) errutil.RecipeError {
	return rs.db.createRecipe(recipe)
}

func (rs *recipeService) readRecipe(id uint) (recipe, errutil.RecipeError) {
	return rs.db.readRecipe(id)
}

func (rs *recipeService) readAllRecipes() ([]recipe, errutil.RecipeError) {
	return rs.db.readAllRecipes()
}

func (rs *recipeService) deleteRecipe(id uint) errutil.RecipeError {
	return rs.db.deleteRecipe(id)
}

func (rs *recipeService) updateRecipe(recipe *recipe) errutil.RecipeError {
	return rs.db.updateRecipe(recipe)
}

func (rs *recipeService) getPathId(c echo.Context) (uint, errutil.RecipeError) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return 0, errutil.RecipeErrorInvalidParam
	}

	return uint(id), nil
}

func (rs *recipeService) getRecipeById(c echo.Context) (recipe, errutil.RecipeError) {
	var recipe recipe

	id, err := rs.getPathId(c)
	if err != nil {
		return recipe, errutil.RecipeErrorInvalidParam
	}

	return rs.readRecipe(id)
}

func (rs *recipeService) parseFormData(c echo.Context, withID bool) (recipe, errutil.RecipeError) {
	var recipe recipe

	if err := c.Request().ParseForm(); err != nil {
		return recipe, errutil.RecipeErrorInvalidFormData
	}

	recipe.Title = c.Request().FormValue("title")
	recipe.Description = c.Request().FormValue("description")
	recipe.Ingredients = c.Request().FormValue("ingredients")
	recipe.Instructions = c.Request().FormValue("instructions")

	if len(recipe.Title) == 0 ||
		len(recipe.Description) == 0 ||
		len(recipe.Ingredients) == 0 ||
		len(recipe.Instructions) == 0 {
		return recipe, errutil.RecipeErrorInvalidFormData
	}

	duration, err := strconv.Atoi(c.Request().FormValue("duration"))
	if err != nil {
		return recipe, errutil.RecipeErrorInvalidFormData
	}
	recipe.Duration = duration

	if withID {
		id, err := rs.getPathId(c)
		if err != nil {
			return recipe, err
		}
		recipe.ID = id
	}

	return recipe, nil
}
