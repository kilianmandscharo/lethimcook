package recipe

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/kilianmandscharo/lethimcook/cache"
	"github.com/kilianmandscharo/lethimcook/components"
	"github.com/kilianmandscharo/lethimcook/errutil"
	"github.com/kilianmandscharo/lethimcook/logging"
	"github.com/kilianmandscharo/lethimcook/servutil"
	"github.com/kilianmandscharo/lethimcook/types"
	"github.com/labstack/echo/v4"
	"github.com/yuin/goldmark"
)

type recipeService struct {
	db          *recipeDatabase
	logger      *logging.Logger
	recipeCache *cache.RecipeCache
}

func NewRecipeService(db *recipeDatabase, logger *logging.Logger) *recipeService {
	return &recipeService{
		db:          db,
		logger:      logger,
		recipeCache: cache.NewRecipeCache(logger),
	}
}

func (rs *recipeService) createRecipe(recipe *types.Recipe) error {
	rs.recipeCache.Invalidate()
	return rs.db.createRecipe(recipe)
}

func (rs *recipeService) readRecipe(id uint) (types.Recipe, error) {
	return rs.db.readRecipe(id)
}

func (rs *recipeService) getReadRecipeOptionsFromRequest(c echo.Context) readRecipesOptions {
	isAdmin := servutil.IsAuthorized(c)
	query := c.QueryParam("search")
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		page = 1
	}
	pageSize, err := strconv.Atoi(c.QueryParam("pageSize"))
	if err != nil {
		pageSize = 10
	}
	return readRecipesOptions{
		isAdmin:  isAdmin,
		query:    query,
		page:     page,
		pageSize: pageSize,
	}
}

type readRecipesOptions struct {
	query          string
	page, pageSize int
	isAdmin        bool
}

func (rs *recipeService) readRecipes(options readRecipesOptions) ([]types.Recipe, types.PaginationInfo, error) {
	recipes := []types.Recipe{}
	paginationInfo := types.PaginationInfo{}
	var err error

	if options.page <= 0 || options.pageSize <= 0 {
		return recipes, paginationInfo, nil
	}
	paginationInfo.CurrentPage = options.page

	recipes, err = rs.readAllRecipes(options.isAdmin)
	if err != nil {
		return recipes, paginationInfo, errutil.AddMessageToAppError(
			err,
			fmt.Sprintf("failed at readRecipes() with options %v", options),
		)
	}

	if len(options.query) > 0 {
		recipes = rs.filterRecipes(recipes, options.query)
	}
	paginationInfo.TotalRecipes = len(recipes)

	numberOfPages := int(math.Ceil(float64(len(recipes)) / float64(options.pageSize)))
	if numberOfPages == 0 {
		return recipes, paginationInfo, nil
	}
	paginationInfo.TotalPages = numberOfPages

	start := (options.page - 1) * options.pageSize
	if start >= len(recipes) {
		return []types.Recipe{}, paginationInfo, nil
	}

	end := options.page * options.pageSize
	if end > len(recipes) {
		end = len(recipes)
	}

	return recipes[start:end], paginationInfo, nil
}

func (rs *recipeService) readAllRecipes(isAdmin bool) ([]types.Recipe, error) {
	if rs.recipeCache.IsValid() {
		recipes := *rs.recipeCache.Get(isAdmin)
		return recipes, nil
	}
	recipes, err := rs.db.readAllRecipesWithPending()
	if err != nil {
		errutil.AddMessageToAppError(
			err,
			fmt.Sprintf("failed at readAllRecipes() with isAdmin = %t", isAdmin),
		)
	}
	rs.recipeCache.Set(recipes)
	return *rs.recipeCache.Get(isAdmin), nil
}

func (rs *recipeService) deleteRecipe(id uint) error {
	rs.recipeCache.Invalidate()
	return rs.db.deleteRecipe(id)
}

func (rs *recipeService) updateRecipe(recipe *types.Recipe) error {
	rs.recipeCache.Invalidate()
	return rs.db.updateRecipe(recipe)
}

func (rs *recipeService) updatePending(id uint, pending bool) error {
	rs.recipeCache.Invalidate()
	return rs.db.updatePending(id, pending)
}

func (rs *recipeService) filterRecipes(recipes []types.Recipe, query string) []types.Recipe {
	query = strings.ToLower(strings.TrimSpace(query))
	if len(query) == 0 {
		return recipes
	}
	var filteredRecipes []types.Recipe
	for _, recipe := range recipes {
		if recipe.ContainsQuery(query) {
			filteredRecipes = append(filteredRecipes, recipe)
		}
	}
	return filteredRecipes
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

func (rs *recipeService) updateRecipeWithFormData(c echo.Context, recipe *types.Recipe) (map[string]error, error) {
	formErrors := make(map[string]error)

	if err := rs.parseForm(c); err != nil {
		return formErrors, errutil.AddMessageToAppError(
			err,
			"failed at updateRecipeWithFormData()",
		)
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
	} else {
		recipe.Duration = duration
	}

	return formErrors, nil
}

func (rs *recipeService) getRecipeAsJson(id uint) ([]byte, error) {
	recipe, err := rs.db.readRecipe(id)
	if err != nil {
		return []byte{}, errutil.AddMessageToAppError(err, "failed at getRecipeAsJson()")
	}
	jsonRecipe, err := json.Marshal(recipe)
	if err != nil {
		return []byte{}, &errutil.AppError{
			UserMessage: "Serverfehler",
			Err:         fmt.Errorf("failed at getRecipeAsJson() for id %d: %w", id, err),
			StatusCode:  http.StatusInternalServerError,
		}
	}
	return jsonRecipe, nil
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
			Type:        types.FormElementInput,
			Name:        "duration",
			Err:         formErrors["duration"],
			Value:       duration,
			InputType:   "number",
			Label:       "Zubereitungszeit (Minuten)",
			Placeholder: "Zubereitungszeit",
			Required:    true,
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
			Type:        types.FormElementInput,
			Name:        "tags",
			Err:         formErrors["tags"],
			Value:       recipe.Tags,
			InputType:   "text",
			Label:       "Tags (getrennt durch Kommas)",
			Placeholder: "Tags",
		},
		{
			Type:           types.FormElementTextArea,
			Name:           "ingredients",
			Err:            formErrors["ingredients"],
			Value:          recipe.Ingredients,
			Label:          "Zutaten (Markdown)",
			Placeholder:    "Zutaten",
			Required:       true,
			LabelComponent: components.PreviewButton("ingredients"),
		},
		{
			Type:           types.FormElementTextArea,
			Name:           "instructions",
			Err:            formErrors["instructions"],
			Value:          recipe.Instructions,
			Label:          "Anleitung (Markdown)",
			Placeholder:    "Anleitung",
			Required:       true,
			LabelComponent: components.PreviewButton("instructions"),
		},
	}
}

func (rs *recipeService) getRecipeLinks(isAdmin bool, query string) ([]types.RecipeLinkData, error) {
	links := []types.RecipeLinkData{}
	recipes, err := rs.readAllRecipes(isAdmin)
	if err != nil {
		return links, errutil.AddMessageToAppError(err, "failed at getRecipeLinks()")
	}
	if len(query) > 0 {
		recipes = rs.filterRecipes(recipes, query)
	}
	for _, recipe := range recipes {
		if recipe.Pending {
			continue
		}
		links = append(links, types.RecipeLinkData{
			ID:    recipe.ID,
			Title: recipe.Title,
		})
	}
	return links, nil
}

func (rs *recipeService) renderMarkdown(markdown string) (string, error) {
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(markdown), &buf); err != nil {
		return "", &errutil.AppError{
			UserMessage: "Fehler beim Markdownparsing",
			StatusCode:  http.StatusInternalServerError,
			Err:         fmt.Errorf("failed at renderMarkdown(): %w", err),
		}
	}
	return buf.String(), nil
}

func (rs *recipeService) parseForm(c echo.Context) error {
	if err := c.Request().ParseForm(); err != nil {
		return &errutil.AppError{
			UserMessage: "Fehlerhaftes Formular",
			Err:         fmt.Errorf("failed at parseForm(): %w", err),
			StatusCode:  http.StatusBadRequest,
		}
	}
	return nil
}

func (rs *recipeService) extractFirstFormEntry(c echo.Context) (string, string, error) {
	if err := rs.parseForm(c); err != nil {
		return "", "", errutil.AddMessageToAppError(
			err,
			"failed at extractFirstFormEntry()",
		)
	}

	if len(c.Request().Form) != 1 {
		return "", "", &errutil.AppError{
			UserMessage: "Fehlerhaftes Formular",
			Err: fmt.Errorf(
				"Failed at extractFirstFormEntry(), got %d form fields, want 1",
				len(c.Request().Form),
			),
			StatusCode: http.StatusBadRequest,
		}
	}

	key := ""
	value := ""

	for k, v := range c.Request().Form {
		key = k
		if len(v) > 0 {
			value = v[0]
		}
	}

	switch key {
	case "ingredients":
		return "Zutaten", value, nil
	case "instructions":
		return "Anleitung", value, nil
	}

	return "", "", &errutil.AppError{
		UserMessage: "Fehlerhaftes Formular",
		Err: fmt.Errorf(
			"Failed at extractFirstFormEntry() with key %s",
			key,
		),
		StatusCode: http.StatusBadRequest,
	}
}
