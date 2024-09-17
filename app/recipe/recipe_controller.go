package recipe

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/kilianmandscharo/lethimcook/components"
	"github.com/kilianmandscharo/lethimcook/errutil"
	"github.com/kilianmandscharo/lethimcook/servutil"
	"github.com/kilianmandscharo/lethimcook/types"
	"github.com/labstack/echo/v4"
)

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
	e.GET("/recipe/:id/edit", rc.RenderRecipeEditPage)
	e.GET("/recipe/new", rc.RenderRecipeNewPage)
	e.GET("/recipe/:id", rc.RenderRecipePage)

	// Actions
	e.POST("/search", rc.HandleSearchRecipe)
	e.POST("/recipe", rc.HandleCreateRecipe)
	e.PUT("/recipe/:id", rc.HandleUpdateRecipe)
	e.PUT("/recipe/:id/pending/:pending", rc.HandleUpdatePending)
	e.DELETE("/recipe/:id", rc.HandleDeleteRecipe)
}

func (rc *RecipeController) RenderRecipeListPage(c echo.Context) error {
	return rc.renderRecipeListPageHelper(c, "")
}

func (rc *RecipeController) renderRecipeListPageHelper(c echo.Context, message string) error {
	isAdmin := servutil.IsAuthorized(c)

	recipes, err := rc.recipeService.readAllRecipes(isAdmin)
	if err != nil {
		return servutil.RenderError(
			c,
			errutil.AddMessageToAppError(err, "failed at renderRecipeListPageHelper()"),
		)
	}

	return servutil.RenderComponent(servutil.RenderComponentOptions{
		Context:   c,
		Component: components.RecipesPage(isAdmin, recipes),
		Message:   message,
	})
}

func (rc *RecipeController) RenderRecipeNewPage(c echo.Context) error {
	formElements := rc.recipeService.createRecipeForm(types.Recipe{}, make(map[string]error))
	return servutil.RenderComponent(servutil.RenderComponentOptions{
		Context:   c,
		Component: components.RecipeNewPage(formElements, servutil.IsAuthorized(c)),
	})
}

func (rc *RecipeController) RenderRecipeEditPage(c echo.Context) error {
	isAdmin := servutil.IsAuthorized(c)
	if !isAdmin {
		return servutil.RenderError(c, &errutil.AppError{
			UserMessage: "Nicht authorisiert",
			Err:         errors.New("failed at RenderRecipeEditPage(), not authorized"),
			StatusCode:  http.StatusUnauthorized,
		})
	}

	recipe, err := rc.recipeService.getRecipeById(c)
	if err != nil {
		return servutil.RenderError(
			c,
			errutil.AddMessageToAppError(err, "failed at RenderRecipeNewPage()"),
		)
	}

	formElements := rc.recipeService.createRecipeForm(recipe, make(map[string]error))
	return servutil.RenderComponent(servutil.RenderComponentOptions{
		Context:   c,
		Component: components.RecipeEditPage(isAdmin, recipe.ID, formElements),
	})
}

func (rc *RecipeController) RenderRecipePage(c echo.Context) error {
	return rc.RenderRecipePageHelper(c, "")
}

func (rc *RecipeController) RenderRecipePageHelper(c echo.Context, message string) error {
	recipe, err := rc.recipeService.getRecipeById(c)
	if err != nil {
		return servutil.RenderError(
			c,
			errutil.AddMessageToAppError(err, "failed at RenderRecipePageHelper()"),
		)
	}

	if err := recipe.RenderMarkdown(); err != nil {
		return servutil.RenderError(
			c,
			errutil.AddMessageToAppError(err, "failed at RenderRecipePageHelper()"),
		)
	}

	c.Response().Header().Set("HX-Push-Url", fmt.Sprintf("/recipe/%d", recipe.ID))

	return servutil.RenderComponent(servutil.RenderComponentOptions{
		Context:   c,
		Component: components.RecipePage(servutil.IsAuthorized(c), recipe, recipe.ParseTags()),
		Message:   message,
	})
}

func (rc *RecipeController) HandleSearchRecipe(c echo.Context) error {
	isAdmin := servutil.IsAuthorized(c)

	if err := c.Request().ParseForm(); err != nil {
		return servutil.RenderError(
			c,
			errutil.AddMessageToAppError(err, "failed at HandleSearchRecipe()"),
		)
	}
	query := strings.ToLower(c.FormValue("query"))

	filteredRecipes, err := rc.recipeService.getFilteredRecipes(query, isAdmin)
	if err != nil {
		return servutil.RenderError(
			c,
			errutil.AddMessageToAppError(err, "failed at HandleSearchRecipe()"),
		)
	}

	return servutil.RenderComponent(servutil.RenderComponentOptions{
		Context:   c,
		Component: components.RecipeList(isAdmin, filteredRecipes),
	})
}

func (rc *RecipeController) HandleCreateRecipe(c echo.Context) error {
	pending := c.QueryParam("pending") == "true"
	isAdmin := servutil.IsAuthorized(c)

	if !isAdmin && !pending {
		return servutil.RenderError(c, &errutil.AppError{
			UserMessage: "Nicht authorisiert",
			Err:         errors.New("failed at HandleCreateRecipe(), not authorized"),
			StatusCode:  http.StatusUnauthorized,
		})
	}

	newRecipe, formErrors, err := rc.recipeService.parseFormData(c, false)
	if err != nil {
		return servutil.RenderError(
			c,
			errutil.AddMessageToAppError(err, "failed at HandleCreateRecipe()"),
		)
	}
	newRecipe.Pending = pending

	if len(formErrors) > 0 {
		formElements := rc.recipeService.createRecipeForm(newRecipe, formErrors)
		return servutil.RenderComponent(servutil.RenderComponentOptions{
			Context:   c,
			Component: components.RecipeNewPage(formElements, servutil.IsAuthorized(c)),
			Err: &errutil.AppError{
				UserMessage: "Fehlerhaftes Formular",
				StatusCode:  http.StatusBadRequest,
				Err:         fmt.Errorf("failed at HandleCreateRecipe(), invalid form: %v", formErrors),
			},
		})
	}

	if err := rc.recipeService.createRecipe(&newRecipe); err != nil {
		return servutil.RenderError(
			c,
			errutil.AddMessageToAppError(err, "failed at HandleCreateRecipe()"),
		)
	}

	if pending {
		return servutil.RenderComponent(servutil.RenderComponentOptions{
			Context:   c,
			Component: components.RecipeCreationSuccess(),
			Message:   "Rezept eingereicht",
		})
	}
	return rc.renderRecipeListPageHelper(c, "Rezept erstellt")
}

func (rc *RecipeController) HandleUpdatePending(c echo.Context) error {
	isAdmin := servutil.IsAuthorized(c)
	if !isAdmin {
		return servutil.RenderError(c, &errutil.AppError{
			UserMessage: "Nicht authorisiert",
			Err:         errors.New("failed at HandleToggleRecipe(), not authorized"),
			StatusCode:  http.StatusUnauthorized,
		})
	}

	createError := func(err error) error {
		return servutil.RenderError(
			c,
			errutil.AddMessageToAppError(err, "failed at HandleToggleRecipe()"),
		)
	}

	pending, err := rc.recipeService.getPathPending(c)
	if err != nil {
		return createError(err)
	}

	id, err := rc.recipeService.getPathId(c)
	if err != nil {
		return createError(err)
	}

	err = rc.recipeService.updatePending(id, pending)
	if err != nil {
		return createError(err)
	}

	if pending {
		return rc.renderRecipeListPageHelper(c, "Rezept auf 'ausstehend' gesetzt")
	} else {
		return rc.renderRecipeListPageHelper(c, "Rezept angenommen")
	}
}

func (rc *RecipeController) HandleUpdateRecipe(c echo.Context) error {
	isAdmin := servutil.IsAuthorized(c)
	if !isAdmin {
		return servutil.RenderError(c, &errutil.AppError{
			UserMessage: "Nicht authorisiert",
			Err:         errors.New("failed at HandleUpdateRecipe(), not authorized"),
			StatusCode:  http.StatusUnauthorized,
		})
	}

	updatedRecipe, formErrors, err := rc.recipeService.parseFormData(c, true)
	if err != nil {
		return servutil.RenderError(
			c,
			errutil.AddMessageToAppError(err, "failed at HandleUpdateRecipe()"),
		)
	}

	if len(formErrors) > 0 {
		formElements := rc.recipeService.createRecipeForm(updatedRecipe, formErrors)
		return servutil.RenderComponent(servutil.RenderComponentOptions{
			Context:   c,
			Component: components.RecipeEditPage(isAdmin, updatedRecipe.ID, formElements),
			Err: &errutil.AppError{
				UserMessage: "Fehlerhaftes Formular",
				StatusCode:  http.StatusBadRequest,
				Err:         fmt.Errorf("failed at HandleUpdateRecipe(), invalid form: %v", formErrors),
			},
		})
	}

	if err := rc.recipeService.updateRecipe(&updatedRecipe); err != nil {
		return servutil.RenderError(
			c,
			errutil.AddMessageToAppError(err, "failed at HandleUpdateRecipe()"),
		)
	}

	return rc.RenderRecipePageHelper(c, "Rezept aktualisiert")
}

func (rc *RecipeController) HandleDeleteRecipe(c echo.Context) error {
	isAdmin := servutil.IsAuthorized(c)
	if !isAdmin {
		return servutil.RenderError(c, &errutil.AppError{
			UserMessage: "Nicht authorisiert",
			Err:         errors.New("failed at HandleDeleteRecipe(), not authorized"),
			StatusCode:  http.StatusUnauthorized,
		})
	}

	createError := func(err error) error {
		return servutil.RenderError(
			c,
			errutil.AddMessageToAppError(err, "failed at HandleDeleteRecipe()"),
		)
	}

	id, err := rc.recipeService.getPathId(c)
	if err != nil {
		return createError(err)
	}

	err = rc.recipeService.deleteRecipe(id)
	if err != nil {
		return createError(err)
	}

	recipes, err := rc.recipeService.readAllRecipes(isAdmin)
	if err != nil {
		return createError(err)
	}

	return servutil.RenderComponent(servutil.RenderComponentOptions{
		Context:   c,
		Component: components.RecipesPage(isAdmin, recipes),
		Message:   "Rezept entfernt",
	})
}
