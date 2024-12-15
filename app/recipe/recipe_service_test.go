package recipe

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kilianmandscharo/lethimcook/cache"
	"github.com/kilianmandscharo/lethimcook/errutil"
	"github.com/kilianmandscharo/lethimcook/logging"
	"github.com/kilianmandscharo/lethimcook/testutil"
	"github.com/kilianmandscharo/lethimcook/types"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func newTestRecipeService() *recipeService {
	logger := logging.New(logging.Debug)
	return &recipeService{
		db:          newTestRecipeDatabase(),
		logger:      logger,
		recipeCache: cache.NewRecipeCache(logger),
	}
}

type newTestContextOptions struct {
	formData    string
	pathId      string
	pathPending string
}

func newTestContext(t *testing.T, options newTestContextOptions) echo.Context {
	e := echo.New()
	w := httptest.NewRecorder()

	var body io.Reader
	body = bytes.NewBufferString(options.formData)
	req, err := http.NewRequest(http.MethodPost, "", body)
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

	c := e.NewContext(req, w)
	c.SetParamNames("id", "pending")
	c.SetParamValues(options.pathId, options.pathPending)

	return c
}

func TestGetFilteredRecipes(t *testing.T) {
	recipeService := newTestRecipeService()
	recipes := []types.Recipe{
		{Description: "Italienische Knoblauchnudeln"},
		{
			Title: "Naan",
			Tags:  "indisch, Beilage",
		},
	}

	tests := []struct {
		query string
		hits  int
	}{
		{"naan", 1},
		{"Naan", 1},
		{"xx", 0},
		{"Italienische Knoblauchnudeln", 1},
		{"a", 2},
		{"indisch", 1},
		{"Beilage", 1},
	}

	for _, tt := range tests {
		filteredRecipes := recipeService.filterRecipes(recipes, tt.query)
		assert.Equal(t, tt.hits, len(filteredRecipes))
	}
}

func TestGetPathId(t *testing.T) {
	// Given
	recipeService := newTestRecipeService()
	c := newTestContext(t, newTestContextOptions{
		pathId: "xx",
	})

	// When
	_, err := recipeService.getPathId(c)

	// Then
	assert.Error(t, err)
	appError, ok := err.(*errutil.AppError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, appError.StatusCode)
	assert.Equal(t, "Ungültiges Pfadparameter", appError.UserMessage)

	// Given
	c = newTestContext(t, newTestContextOptions{
		pathId: "1",
	})

	// When
	id, err := recipeService.getPathId(c)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, uint(1), id)
}

func TestGetPathPending(t *testing.T) {
	// Given
	recipeService := newTestRecipeService()
	c := newTestContext(t, newTestContextOptions{})

	// When
	_, err := recipeService.getPathPending(c)

	// Then
	assert.Error(t, err)
	appError, ok := err.(*errutil.AppError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, appError.StatusCode)
	assert.Equal(t, "Fehlendes Pfadparameter", appError.UserMessage)

	// Given
	c = newTestContext(t, newTestContextOptions{
		pathPending: "xx",
	})

	// When
	_, err = recipeService.getPathPending(c)

	// Then
	assert.Error(t, err)
	appError, ok = err.(*errutil.AppError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, appError.StatusCode)
	assert.Equal(t, "Ungültiges Pfadparameter", appError.UserMessage)

	// Given
	c = newTestContext(t, newTestContextOptions{
		pathPending: "true",
	})

	// When
	pending, err := recipeService.getPathPending(c)

	// Then
	assert.NoError(t, err)
	assert.True(t, pending)

	// Given
	c = newTestContext(t, newTestContextOptions{
		pathPending: "false",
	})

	// When
	pending, err = recipeService.getPathPending(c)

	// Then
	assert.NoError(t, err)
	assert.False(t, pending)
}

func TestUpdateRecipeWithFormData(t *testing.T) {
	// Given
	recipeService := newTestRecipeService()

	testCases := []struct {
		formData   string
		formErrors []error
	}{
		{
			formData: testutil.ConstructTestFormDataString(
				testutil.TestFormDataStringOptions{
					TitleEmpty: true,
				},
			),
			formErrors: []error{errutil.FormErrorNoTitle},
		},
		{
			formData: testutil.ConstructTestFormDataString(
				testutil.TestFormDataStringOptions{
					DescriptionEmpty: true,
				},
			),
			formErrors: []error{errutil.FormErrorNoDescription},
		},
		{
			formData: testutil.ConstructTestFormDataString(
				testutil.TestFormDataStringOptions{
					DurationEmpty: true,
				},
			),
			formErrors: []error{errutil.FormErrorNoDuration},
		},
		{
			formData: testutil.ConstructTestFormDataString(
				testutil.TestFormDataStringOptions{
					InvalidDuration: true,
				},
			),
			formErrors: []error{errutil.FormErrorNoDuration},
		},
		{
			formData: testutil.ConstructTestFormDataString(
				testutil.TestFormDataStringOptions{
					IngredientsEmpty: true,
				},
			),
			formErrors: []error{errutil.FormErrorNoIngredients},
		},
		{
			formData: testutil.ConstructTestFormDataString(
				testutil.TestFormDataStringOptions{
					InstructionsEmpty: true,
				},
			),
			formErrors: []error{errutil.FormErrorNoInstructions},
		},
		{
			formData: testutil.ConstructTestFormDataString(
				testutil.TestFormDataStringOptions{
					TitleEmpty:        true,
					DescriptionEmpty:  true,
					DurationEmpty:     true,
					IngredientsEmpty:  true,
					InstructionsEmpty: true,
				},
			),
			formErrors: []error{
				errutil.FormErrorNoTitle,
				errutil.FormErrorNoDescription,
				errutil.FormErrorNoDuration,
				errutil.FormErrorNoIngredients,
				errutil.FormErrorNoInstructions,
			},
		},
	}

	for _, test := range testCases {
		c := newTestContext(t, newTestContextOptions{
			formData: test.formData,
		})
		formErrors, err := recipeService.updateRecipeWithFormData(c, &types.Recipe{})

		assert.NoError(t, err)
		assert.Equal(t, len(test.formErrors), len(formErrors))

		for _, err := range test.formErrors {
			errExists := false
			for _, e := range formErrors {
				if err == e {
					errExists = true
					break
				}
			}
			assert.True(t, errExists)
		}
	}

	// Given
	formData := "title=title&description=description&ingredients=ingredients&instructions=instructions&duration=20&tags=vegan"
	c := newTestContext(t, newTestContextOptions{
		formData: formData,
		pathId:   "1",
	})

	// When
	var recipe types.Recipe
	_, err := recipeService.updateRecipeWithFormData(c, &recipe)

	// Then
	assert.NoError(t, err)
	assert.True(t, recipe == types.Recipe{
		Title:        "title",
		Description:  "description",
		Instructions: "instructions",
		Ingredients:  "ingredients",
		Duration:     20,
		Tags:         "vegan",
	})
}

func TestReadAllRecipesService(t *testing.T) {
	// Given
	recipeService := newTestRecipeService()
	assert.NoError(t, recipeService.db.createRecipe(&types.Recipe{
		Pending: true,
	}))
	assert.NoError(t, recipeService.db.createRecipe(&types.Recipe{
		Pending: false,
	}))

	// When
	recipes, err := recipeService.readAllRecipes(false)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, 1, len(recipes))
	assert.False(t, recipes[0].Pending)

	// When
	recipes, err = recipeService.readAllRecipes(true)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, 2, len(recipes))
}

func TestGetRecipeAsJson(t *testing.T) {
	// Given
	recipeService := newTestRecipeService()

	// When
	_, err := recipeService.getRecipeAsJson(1)

	// Then
	assert.Error(t, err)

	// When
	assert.NoError(t, recipeService.createRecipe(&types.Recipe{ID: 1, Title: "Naan", Author: "Phillip Jeffries"}))
	recipeJson, err := recipeService.getRecipeAsJson(1)

	// Then
	assert.NoError(t, err)
	assert.Equal(
		t,
		[]byte("{\"id\":1,\"author\":\"Phillip Jeffries\",\"source\":\"\",\"title\":\"Naan\",\"description\":\"\",\"duration\":0,\"ingredients\":\"\",\"instructions\":\"\",\"tags\":\"\",\"createdAt\":\"\"}"),
		recipeJson,
	)
}

func TestReadRecipes(t *testing.T) {
	tests := []struct {
		options            readRecipesOptions
		wantRecipes        []types.Recipe
		wantPaginationInfo types.PaginationInfo
	}{
		{
			options: readRecipesOptions{
				page:     1,
				pageSize: 0,
			},
			wantRecipes: []types.Recipe{},
		},
		{
			options: readRecipesOptions{
				page:     0,
				pageSize: 1,
			},
			wantRecipes: []types.Recipe{},
		},
		{
			options: readRecipesOptions{
				page:     1,
				pageSize: 5,
			},
			wantRecipes: []types.Recipe{
				{ID: 10},
				{ID: 9},
				{ID: 8},
				{ID: 7},
				{ID: 6},
			},
			wantPaginationInfo: types.PaginationInfo{
				TotalRecipes: 10,
				TotalPages:   2,
				CurrentPage:  1,
			},
		},
		{
			options: readRecipesOptions{
				page:     1,
				pageSize: 10,
			},
			wantRecipes: []types.Recipe{
				{ID: 10},
				{ID: 9},
				{ID: 8},
				{ID: 7},
				{ID: 6},
				{ID: 5},
				{ID: 4},
				{ID: 3},
				{ID: 2},
				{ID: 1},
			},
			wantPaginationInfo: types.PaginationInfo{
				TotalRecipes: 10,
				TotalPages:   1,
				CurrentPage:  1,
			},
		},
		{
			options: readRecipesOptions{
				page:     1,
				pageSize: 9,
			},
			wantRecipes: []types.Recipe{
				{ID: 10},
				{ID: 9},
				{ID: 8},
				{ID: 7},
				{ID: 6},
				{ID: 5},
				{ID: 4},
				{ID: 3},
				{ID: 2},
			},
			wantPaginationInfo: types.PaginationInfo{
				TotalRecipes: 10,
				TotalPages:   2,
				CurrentPage:  1,
			},
		},
		{
			options: readRecipesOptions{
				page:     1,
				pageSize: 11,
			},
			wantRecipes: []types.Recipe{
				{ID: 10},
				{ID: 9},
				{ID: 8},
				{ID: 7},
				{ID: 6},
				{ID: 5},
				{ID: 4},
				{ID: 3},
				{ID: 2},
				{ID: 1},
			},
			wantPaginationInfo: types.PaginationInfo{
				TotalRecipes: 10,
				TotalPages:   1,
				CurrentPage:  1,
			},
		},
		{
			options: readRecipesOptions{
				page:     2,
				pageSize: 10,
			},
			wantRecipes: []types.Recipe{},
			wantPaginationInfo: types.PaginationInfo{
				TotalRecipes: 10,
				TotalPages:   1,
				CurrentPage:  2,
			},
		},
		{
			options: readRecipesOptions{
				page:     2,
				pageSize: 5,
			},
			wantRecipes: []types.Recipe{
				{ID: 5},
				{ID: 4},
				{ID: 3},
				{ID: 2},
				{ID: 1},
			},
			wantPaginationInfo: types.PaginationInfo{
				TotalRecipes: 10,
				TotalPages:   2,
				CurrentPage:  2,
			},
		},
		{
			options: readRecipesOptions{
				page:     3,
				pageSize: 5,
			},
			wantRecipes: []types.Recipe{},
			wantPaginationInfo: types.PaginationInfo{
				TotalRecipes: 10,
				TotalPages:   2,
				CurrentPage:  3,
			},
		},
		{
			options: readRecipesOptions{
				page:     2,
				pageSize: 3,
			},
			wantRecipes: []types.Recipe{
				{ID: 7},
				{ID: 6},
				{ID: 5},
			},
			wantPaginationInfo: types.PaginationInfo{
				TotalRecipes: 10,
				TotalPages:   4,
				CurrentPage:  2,
			},
		},
		{
			options: readRecipesOptions{
				page:     3,
				pageSize: 3,
			},
			wantRecipes: []types.Recipe{
				{ID: 4},
				{ID: 3},
				{ID: 2},
			},
			wantPaginationInfo: types.PaginationInfo{
				TotalRecipes: 10,
				TotalPages:   4,
				CurrentPage:  3,
			},
		},
	}

	recipeService := newTestRecipeService()
	for range 10 {
		assert.NoError(t, recipeService.createRecipe(&types.Recipe{}))
	}

	for _, tt := range tests {
		recipes, paginationInfo, err := recipeService.readRecipes(tt.options)
		assert.NoError(t, err)
		assert.Equal(t, tt.wantRecipes, recipes)
		assert.Equal(t, tt.wantPaginationInfo, paginationInfo)
	}
}

func TestGetRecipeLinks(t *testing.T) {
	// Given
	recipeService := newTestRecipeService()
	assert.NoError(t, recipeService.createRecipe(&types.Recipe{
		Title: "Naan",
	}))
	assert.NoError(t, recipeService.createRecipe(&types.Recipe{
		Title: "Pita",
	}))
	assert.NoError(t, recipeService.createRecipe(&types.Recipe{
		Title: "Muffins",
		Pending: true,
	}))

	// When
	links, err := recipeService.getRecipeLinks(false, "Naan")

	// Then
	assert.NoError(t, err)
	assert.Equal(t, []recipeLinkData{{ID: 1, Title: "Naan"}}, links)

	// When
	links, err = recipeService.getRecipeLinks(false, "")

	// Then
	assert.NoError(t, err)
	assert.Equal(t, []recipeLinkData{
		{ID: 2, Title: "Pita"},
		{ID: 1, Title: "Naan"},
	}, links)

	// When
	links, err = recipeService.getRecipeLinks(false, "nope")

	// Then
	assert.NoError(t, err)
	assert.Equal(t, []recipeLinkData{}, links)
}

func TestGetRecipeLinksPayload(t *testing.T) {
	// Given
	recipeService := newTestRecipeService()
	assert.NoError(t, recipeService.createRecipe(&types.Recipe{
		Title: "Naan",
	}))

	// When
	links, err := recipeService.getRecipeLinksPayload(false, "Naan")

	// Then
	assert.NoError(t, err)
	assert.Equal(t, "[{\"id\":1,\"title\":\"Naan\"}]", links)
}
