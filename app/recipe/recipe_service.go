package recipe

import (
	"strconv"
	"strings"

	"github.com/kilianmandscharo/lethimcook/errutil"
	"github.com/kilianmandscharo/lethimcook/types"
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

func (rs *recipeService) createRecipe(recipe *types.Recipe) errutil.RecipeError {
	return rs.db.createRecipe(recipe)
}

func (rs *recipeService) readRecipe(id uint) (types.Recipe, errutil.RecipeError) {
	return rs.db.readRecipe(id)
}

func (rs *recipeService) readAllRecipes() ([]types.Recipe, errutil.RecipeError) {
	return rs.db.readAllRecipes()
}

func (rs *recipeService) deleteRecipe(id uint) errutil.RecipeError {
	return rs.db.deleteRecipe(id)
}

func (rs *recipeService) updateRecipe(recipe *types.Recipe) errutil.RecipeError {
	return rs.db.updateRecipe(recipe)
}

func (rs *recipeService) getFilteredRecipes(query string) ([]types.Recipe, errutil.RecipeError) {
	recipes, err := rs.readAllRecipes()
	if err != nil {
		return []types.Recipe{}, err
	}

	query = strings.TrimSpace(query)

	if len(query) == 0 {
		return recipes, err
	}

	var filteredRecipes []types.Recipe
	for _, recipe := range recipes {
		if recipe.ContainsQuery(query) {
			filteredRecipes = append(filteredRecipes, recipe)
		}
	}

	return filteredRecipes, nil
}

func (rs *recipeService) getPathId(c echo.Context) (uint, errutil.RecipeError) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return 0, errutil.RecipeErrorInvalidParam
	}

	return uint(id), nil
}

func (rs *recipeService) getRecipeById(c echo.Context) (types.Recipe, errutil.RecipeError) {
	var recipe types.Recipe

	id, err := rs.getPathId(c)
	if err != nil {
		return recipe, errutil.RecipeErrorInvalidParam
	}

	return rs.readRecipe(id)
}

func (rs *recipeService) parseFormData(c echo.Context, withID bool) (types.Recipe, errutil.RecipeError) {
	var recipe types.Recipe

	if err := c.Request().ParseForm(); err != nil {
		return recipe, errutil.RecipeErrorInvalidFormData
	}

	recipe.Title = strings.TrimSpace(c.Request().FormValue("title"))
	recipe.Description = strings.TrimSpace(c.Request().FormValue("description"))
	recipe.Tags = strings.TrimSpace(c.Request().FormValue("tags"))
	recipe.Ingredients = strings.TrimSpace(c.Request().FormValue("ingredients"))
	recipe.Instructions = strings.TrimSpace(c.Request().FormValue("instructions"))

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
