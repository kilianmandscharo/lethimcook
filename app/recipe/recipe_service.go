package recipe

import (
	"errors"
	"fmt"
	"net/http"
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

func (rs *recipeService) createRecipe(recipe *types.Recipe) error {
	return rs.db.createRecipe(recipe)
}

func (rs *recipeService) readRecipe(id uint) (types.Recipe, error) {
	return rs.db.readRecipe(id)
}

func (rs *recipeService) readAllRecipes(isAdmin bool) ([]types.Recipe, error) {
	var recipes []types.Recipe
	var err error

	if isAdmin {
		recipes, err = rs.db.readAllRecipesWithPending()
	} else {
		recipes, err = rs.db.readAllRecipes()
	}

	if err != nil {
		errutil.AddMessageToAppError(
			err,
			fmt.Sprintf("failed at readAllRecipes() with isAdmin = %t", isAdmin),
		)
	}

	return recipes, nil
}

func (rs *recipeService) deleteRecipe(id uint) error {
	return rs.db.deleteRecipe(id)
}

func (rs *recipeService) updateRecipe(recipe *types.Recipe) error {
	return rs.db.updateRecipe(recipe)
}

func (rs *recipeService) updatePending(id uint, pending bool) error {
	return rs.db.updatePending(id, pending)
}

func (rs *recipeService) getFilteredRecipes(query string, isAdmin bool) ([]types.Recipe, error) {
	recipes, err := rs.readAllRecipes(isAdmin)
	if err != nil {
		return []types.Recipe{}, errutil.AddMessageToAppError(
			err,
			fmt.Sprintf("failed at getFilteredRecipes() with query '%s'", query),
		)
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

func (rs *recipeService) getPathId(c echo.Context) (uint, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return 0, &errutil.AppError{
			UserMessage: "Ungültiges Pfadparameter",
			Err: fmt.Errorf(
				"failed at getPathId() with parameter %s: %w",
				c.Param("id"),
				err,
			),
			StatusCode: http.StatusBadRequest,
		}
	}

	return uint(id), nil
}

func (rs *recipeService) getPathPending(c echo.Context) (bool, error) {
	pending := c.Param("pending")
	if len(pending) == 0 {
		return false, &errutil.AppError{
			UserMessage: "Fehlendes Pfadparameter",
			Err:         errors.New("failed at getPathPending()"),
			StatusCode:  http.StatusBadRequest,
		}
	}
	if pending == "true" {
		return true, nil
	}
	if pending == "false" {
		return false, nil
	}
	return false, &errutil.AppError{
		UserMessage: "Ungültiges Pfadparameter",
		Err:         fmt.Errorf("failed at getPathPending() with param %s", pending),
		StatusCode:  http.StatusBadRequest,
	}
}

func (rs *recipeService) getRecipeById(c echo.Context) (types.Recipe, error) {
	var recipe types.Recipe

	id, err := rs.getPathId(c)
	if err != nil {
		return recipe, errutil.AddMessageToAppError(
			err,
			"failed at getRecipeById()",
		)
	}

	return rs.readRecipe(id)
}

func (rs *recipeService) parseFormData(c echo.Context, withID bool) (types.Recipe, map[string]error, error) {
	var recipe types.Recipe

	formErrors := make(map[string]error)

	if err := c.Request().ParseForm(); err != nil {
		return recipe, formErrors, &errutil.AppError{
			UserMessage: "Fehlerhaftes Formular",
			Err:         fmt.Errorf("failed at parseFormData(): %w", err),
			StatusCode:  http.StatusBadRequest,
		}
	}

	recipe.Author = strings.TrimSpace(c.Request().FormValue("author"))
	recipe.Source = strings.TrimSpace(c.Request().FormValue("source"))
	recipe.Title = strings.TrimSpace(c.Request().FormValue("title"))
	recipe.Description = strings.TrimSpace(c.Request().FormValue("description"))
	recipe.Tags = strings.TrimSpace(c.Request().FormValue("tags"))
	recipe.Ingredients = strings.TrimSpace(c.Request().FormValue("ingredients"))
	recipe.Instructions = strings.TrimSpace(c.Request().FormValue("instructions"))

	if len(recipe.Title) == 0 {
		formErrors["title"] = errutil.FormErrorNoTitle
	}
	if len(recipe.Description) == 0 {
		formErrors["description"] = errutil.FormErrorNoDescription
	}
	if len(recipe.Ingredients) == 0 {
		formErrors["ingredients"] = errutil.FormErrorNoIngredients
	}
	if len(recipe.Instructions) == 0 {
		formErrors["instructions"] = errutil.FormErrorNoInstructions
	}

	duration, err := strconv.Atoi(c.Request().FormValue("duration"))

	if err != nil {
		formErrors["duration"] = errutil.FormErrorNoDuration
		recipe.Duration = 0
	} else {
		recipe.Duration = duration
	}

	if withID {
		id, err := rs.getPathId(c)
		if err != nil {
			return recipe, formErrors, err
		}
		recipe.ID = id
	}

	return recipe, formErrors, nil
}

func (rs *recipeService) createRecipeForm(recipe types.Recipe, formErrors map[string]error) []types.FormElement {
	duration := ""
	if recipe.Duration != 0 {
		duration = fmt.Sprintf("%d", recipe.Duration)
	}

	return []types.FormElement{
		{
			Type:      types.FormElementInput,
			Name:      "title",
			Err:       formErrors["title"],
			Value:     recipe.Title,
			InputType: "text",
			Label:     "Titel",
			Required:  true,
		},
		{
			Type:      types.FormElementInput,
			Name:      "description",
			Err:       formErrors["description"],
			Value:     recipe.Description,
			InputType: "text",
			Label:     "Beschreibung",
			Required:  true,
		},
		{
			Type:      types.FormElementInput,
			Name:      "duration",
			Err:       formErrors["duration"],
			Value:     duration,
			InputType: "number",
			Label:     "Zubereitungszeit (Minuten)",
			Required:  true,
		},
		{
			Type:      types.FormElementInput,
			Name:      "author",
			Err:       formErrors["author"],
			Value:     recipe.Author,
			InputType: "text",
			Label:     "Autor",
		},
		{
			Type:      types.FormElementInput,
			Name:      "source",
			Err:       formErrors["source"],
			Value:     recipe.Source,
			InputType: "text",
			Label:     "Quelle",
		},
		{
			Type:      types.FormElementInput,
			Name:      "tags",
			Err:       formErrors["tags"],
			Value:     recipe.Tags,
			InputType: "text",
			Label:     "Tags (getrennt durch Kommas)",
		},
		{
			Type:     types.FormElementTextArea,
			Name:     "ingredients",
			Err:      formErrors["ingredients"],
			Value:    recipe.Ingredients,
			Label:    "Zutaten",
			Required: true,
		},
		{
			Type:     types.FormElementTextArea,
			Name:     "instructions",
			Err:      formErrors["instructions"],
			Value:    recipe.Instructions,
			Label:    "Anleitung",
			Required: true,
		},
	}
}
