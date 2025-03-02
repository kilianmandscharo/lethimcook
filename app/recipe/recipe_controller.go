package recipe

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/kilianmandscharo/lethimcook/components"
	"github.com/kilianmandscharo/lethimcook/errutil"
	"github.com/kilianmandscharo/lethimcook/logging"
	"github.com/kilianmandscharo/lethimcook/render"
	"github.com/kilianmandscharo/lethimcook/servutil"
	"github.com/kilianmandscharo/lethimcook/types"
	"github.com/labstack/echo/v4"
)

type RecipeController struct {
	recipeService *recipeService
	logger        *logging.Logger
	renderer      *render.Renderer
}

func NewRecipeController(recipeService *recipeService, logger *logging.Logger, renderer *render.Renderer) *RecipeController {
	return &RecipeController{
		recipeService: recipeService,
		logger:        logger,
		renderer:      renderer,
	}
}

func (rc *RecipeController) AttachHandlerFunctions(e *echo.Echo) {
	// Pages
	e.GET("/", rc.RenderRecipeListPage)
	e.GET("/recipe/:id/edit", rc.RenderRecipeEditPage)
	e.GET("/recipe/new", rc.RenderRecipeNewPage)
	e.GET("/recipe/:id", rc.RenderRecipePage)

	// Actions
	e.GET("/recipe/:id/json", rc.HandleDownloadRecipeAsJson)
	e.POST("/recipe", rc.HandleCreateRecipe)
	e.PUT("/recipe/:id", rc.HandleUpdateRecipe)
	e.PUT("/recipe/:id/pending/:pending", rc.HandleUpdatePending)
	e.DELETE("/recipe/:id", rc.HandleDeleteRecipe)
	e.GET("/recipe/link", rc.HandleGetRecipeLinks)
	e.POST("/recipe/preview", rc.HandlePostRecipePreview)
}

func (rc *RecipeController) RenderRecipeListPage(c echo.Context) error {
	if c.Request().Header.Get("Hx-Target") == "recipe-list" {
		return rc.HandleGetPaginatedRecipes(c)
	}
	return rc.renderRecipeListPageHelper(c, "")
}

func (rc *RecipeController) renderRecipeListPageHelper(c echo.Context, message string) error {
	recipes, paginationInfo, err := rc.recipeService.readRecipes(
		rc.recipeService.getReadRecipeOptionsFromRequest(c),
	)
	if err != nil {
		return rc.renderer.RenderError(
			c,
			errutil.AddMessageToAppError(err, "failed at renderRecipeListPageHelper()"),
		)
	}
	return rc.renderer.RenderComponent(render.RenderComponentOptions{
		Context: c,
		Component: components.RecipesPage(
			servutil.IsAuthorized(c),
			recipes,
			paginationInfo,
		),
		Message: message,
	})
}

func (rc *RecipeController) RenderRecipeNewPage(c echo.Context) error {
	formElements := rc.recipeService.createRecipeForm(types.Recipe{}, make(map[string]error))
	return rc.renderer.RenderComponent(render.RenderComponentOptions{
		Context:   c,
		Component: components.RecipeNewPage(formElements, servutil.IsAuthorized(c)),
	})
}

func (rc *RecipeController) RenderRecipeEditPage(c echo.Context) error {
	isAdmin := servutil.IsAuthorized(c)
	if !isAdmin {
		return rc.renderer.RenderError(
			c,
			errutil.NewAppErrorNotAuthorized("RenderRecipeEditPage()"),
		)
	}

	recipe, err := rc.recipeService.getRecipeById(c)
	if err != nil {
		return rc.renderer.RenderError(
			c,
			errutil.AddMessageToAppError(err, "failed at RenderRecipeNewPage()"),
		)
	}

	formElements := rc.recipeService.createRecipeForm(recipe, make(map[string]error))
	return rc.renderer.RenderComponent(render.RenderComponentOptions{
		Context:   c,
		Component: components.RecipeEditPage(isAdmin, recipe.ID, formElements),
	})
}

func (rc *RecipeController) RenderRecipePage(c echo.Context) error {
	createError := func(err error) error {
		return rc.renderer.RenderError(
			c,
			errutil.AddMessageToAppError(err, "failed at RenderRecipePage()"),
		)
	}
	recipe, err := rc.recipeService.getRecipeById(c)
	if err != nil {
		return createError(err)
	}
	if err := recipe.RenderMarkdown(); err != nil {
		return createError(err)
	}
	return rc.renderer.RenderComponent(render.RenderComponentOptions{
		Context:   c,
		Component: components.RecipePage(servutil.IsAuthorized(c), recipe, recipe.ParseTags()),
	})
}

func (rc *RecipeController) HandleGetPaginatedRecipes(c echo.Context) error {
	recipes, paginationInfo, err := rc.recipeService.readRecipes(
		rc.recipeService.getReadRecipeOptionsFromRequest(c),
	)
	if err != nil {
		return rc.renderer.RenderError(
			c,
			errutil.AddMessageToAppError(err, "failed at HandleGetPaginatedRecipes()"),
		)
	}
	return rc.renderer.RenderComponent(render.RenderComponentOptions{
		Context: c,
		Component: components.Joiner(
			components.RecipeCount(paginationInfo.TotalRecipes, true),
			components.RecipeList(servutil.IsAuthorized(c), recipes, paginationInfo),
			components.PageControl(paginationInfo, true),
		),
	})
}

func (rc *RecipeController) HandleCreateRecipe(c echo.Context) error {
	pending := c.QueryParam("pending") == "true"
	isAdmin := servutil.IsAuthorized(c)

	if !isAdmin && !pending {
		return rc.renderer.RenderError(
			c,
			errutil.NewAppErrorNotAuthorized("HandleCreateRecipe()"),
		)
	}

	var recipe types.Recipe

	formErrors, err := rc.recipeService.updateRecipeWithFormData(c, &recipe)
	if err != nil {
		return rc.renderer.RenderError(
			c,
			errutil.AddMessageToAppError(err, "failed at HandleCreateRecipe()"),
		)
	}

	if len(formErrors) > 0 {
		formElements := rc.recipeService.createRecipeForm(recipe, formErrors)
		return rc.renderer.RenderComponent(render.RenderComponentOptions{
			Context:   c,
			Component: components.RecipeNewPage(formElements, servutil.IsAuthorized(c)),
			Err: &errutil.AppError{
				UserMessage: "Fehlerhaftes Formular",
				StatusCode:  http.StatusBadRequest,
				Err:         fmt.Errorf("failed at HandleCreateRecipe(), invalid form: %v", formErrors),
			},
		})
	}

	recipe.Pending = pending
	recipe.CreatedAt = time.Now().Format(time.RFC3339)

	if err := rc.recipeService.createRecipe(&recipe); err != nil {
		return rc.renderer.RenderError(
			c,
			errutil.AddMessageToAppError(err, "failed at HandleCreateRecipe()"),
		)
	}

	if pending {
		rc.logger.Info("created pending recipe", recipe.ID)
		return rc.renderRecipeListPageHelper(c, "Rezept eingereicht")
	}
	rc.logger.Info("created recipe", recipe.ID)
	return rc.renderRecipeListPageHelper(c, "Rezept erstellt")
}

func (rc *RecipeController) HandleUpdatePending(c echo.Context) error {
	isAdmin := servutil.IsAuthorized(c)
	if !isAdmin {
		return rc.renderer.RenderError(
			c,
			errutil.NewAppErrorNotAuthorized("HandleUpdatePending()"),
		)
	}

	createError := func(err error) error {
		return rc.renderer.RenderError(
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
		rc.logger.Infof("set recipe %d to pending", id)
		return rc.renderRecipeListPageHelper(c, "Rezept auf 'ausstehend' gesetzt")
	} else {
		rc.logger.Infof("set recipe %d to not pending", id)
		return rc.renderRecipeListPageHelper(c, "Rezept angenommen")
	}
}

func (rc *RecipeController) HandleUpdateRecipe(c echo.Context) error {
	isAdmin := servutil.IsAuthorized(c)
	if !isAdmin {
		return rc.renderer.RenderError(
			c, errutil.NewAppErrorNotAuthorized("HandleCreateRecipe()"),
		)
	}

	createError := func(err error) error {
		return rc.renderer.RenderError(
			c,
			errutil.AddMessageToAppError(err, "failed at HandleUpdateRecipe()"),
		)
	}

	id, err := rc.recipeService.getPathId(c)
	if err != nil {
		return createError(err)
	}

	recipe, err := rc.recipeService.readRecipe(id)
	if err != nil {
		return createError(err)
	}
	rc.logger.Info("old recipe:", recipe.String())

	formErrors, err := rc.recipeService.updateRecipeWithFormData(c, &recipe)
	if err != nil {
		return createError(err)
	}

	if len(formErrors) > 0 {
		formElements := rc.recipeService.createRecipeForm(recipe, formErrors)
		return rc.renderer.RenderComponent(render.RenderComponentOptions{
			Context:   c,
			Component: components.RecipeEditPage(isAdmin, recipe.ID, formElements),
			Err: &errutil.AppError{
				UserMessage: "Fehlerhaftes Formular",
				StatusCode:  http.StatusBadRequest,
				Err:         fmt.Errorf("failed at HandleUpdateRecipe(), invalid form: %v", formErrors),
			},
		})
	}

	recipe.LastModifiedAt = time.Now().Format(time.RFC3339)
	if err := rc.recipeService.updateRecipe(&recipe); err != nil {
		return createError(err)
	}

	if err := recipe.RenderMarkdown(); err != nil {
		return rc.renderer.RenderError(
			c,
			errutil.AddMessageToAppError(err, "failed at RenderRecipePageHelper()"),
		)
	}

	return rc.renderer.RenderComponent(render.RenderComponentOptions{
		Context:   c,
		Component: components.RecipePage(servutil.IsAuthorized(c), recipe, recipe.ParseTags()),
		Message:   "Rezept aktualisiert",
	})
}

func (rc *RecipeController) HandleDeleteRecipe(c echo.Context) error {
	isAdmin := servutil.IsAuthorized(c)
	if !isAdmin {
		return rc.renderer.RenderError(
			c,
			errutil.NewAppErrorNotAuthorized("HandleDeleteRecipe()"),
		)
	}

	createError := func(err error) error {
		return rc.renderer.RenderError(
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

	recipes, paginationInfo, err := rc.recipeService.readRecipes(
		rc.recipeService.getReadRecipeOptionsFromRequest(c),
	)
	if err != nil {
		return createError(err)
	}

	rc.logger.Info("deleted recipe", id)
	return rc.renderer.RenderComponent(render.RenderComponentOptions{
		Context:   c,
		Component: components.RecipesPage(isAdmin, recipes, paginationInfo),
		Message:   "Rezept entfernt",
	})
}

func (rc *RecipeController) HandleDownloadRecipeAsJson(c echo.Context) error {
	createError := func(err error) error {
		return rc.renderer.RenderError(
			c,
			errutil.AddMessageToAppError(
				err,
				"failed at HandleDownloadRecipeAsJson()",
			),
		)
	}
	id, err := rc.recipeService.getPathId(c)
	if err != nil {
		return createError(err)
	}
	jsonRecipe, err := rc.recipeService.getRecipeAsJson(id)
	if err != nil {
		return createError(err)
	}
	c.Response().Header().Set(
		echo.HeaderContentDisposition,
		fmt.Sprintf("attachment; filename=reicpe_%d.json", id),
	)
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().Header().Set(echo.HeaderContentLength, strconv.Itoa(len(jsonRecipe)))
	return c.Blob(http.StatusOK, echo.MIMEApplicationJSON, jsonRecipe)
}

func (rc *RecipeController) HandleGetRecipeLinks(c echo.Context) error {
	recipes, err := rc.recipeService.getRecipeLinks(
		servutil.IsAuthorized(c),
		c.QueryParam("query"),
	)
	if err != nil {
		return rc.renderer.RenderError(
			c,
			errutil.AddMessageToAppError(err, "failed at HandleGetRecipeLink()"),
		)
	}
	if len(recipes) == 1 {
		payload, err := json.Marshal(recipes[0])
		if err != nil {
			return rc.renderer.RenderError(
				c,
				&errutil.AppError{
					UserMessage: "Fehler bei der Datenverarbeitung",
					Err:         err,
					StatusCode:  http.StatusInternalServerError,
				},
			)
		}
		return c.String(http.StatusOK, string(payload))
	}
	return rc.renderer.RenderComponent(render.RenderComponentOptions{
		Context:       c,
		Component:     components.SelectDialog(recipes),
		OnlyComponent: true,
	})
}

func (rc *RecipeController) HandlePostRecipePreview(c echo.Context) error {
	createError := func(err error) error {
		return rc.renderer.RenderError(
			c,
			errutil.AddMessageToAppError(err, "failed at HandlePostRecipePreview()"),
		)
	}

	key, value, err := rc.recipeService.extractFirstFormEntry(c)
	if err != nil {
		return createError(err)
	}

	html, err := rc.recipeService.renderMarkdown(value)
	if err != nil {
		return createError(err)
	}

	title := ""
	if key == "ingredients" {
		title = "Zutaten"
	} else if key == "instructions" {
		title = "Anleitung"
	}

	return rc.renderer.RenderComponent(render.RenderComponentOptions{
		Context:       c,
		Component:     components.PreviewModal(title, html),
		OnlyComponent: true,
	})
}
