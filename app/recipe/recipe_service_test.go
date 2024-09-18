package recipe

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kilianmandscharo/lethimcook/errutil"
	"github.com/kilianmandscharo/lethimcook/testutil"
	"github.com/kilianmandscharo/lethimcook/types"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func newTestRecipeService() recipeService {
	return recipeService{
		db: newTestRecipeDatabase(),
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
	// Given
	recipeService := newTestRecipeService()
	assert.NoError(t, recipeService.createRecipe(&types.Recipe{
		Title: "Naan",
		Tags:  "indisch, Beilage",
	}))
	assert.NoError(t, recipeService.createRecipe(&types.Recipe{
		Description: "Italienische Knoblauchnudeln",
	}))

	testCases := []struct {
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

	for _, test := range testCases {
		filteredRecipes, err := recipeService.getFilteredRecipes(test.query, false)
		assert.NoError(t, err)
		assert.Equal(t, test.hits, len(filteredRecipes))
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

func TestParseFormData(t *testing.T) {
	// Given
	recipeService := newTestRecipeService()

	testCases := []struct {
		formData      string
		withPathParam bool
		pathParamId   string
		shouldError   bool
		formErrors    []error
	}{
		{
			formData: testutil.ConstructTestFormDataString(
				testutil.TestFormDataStringOptions{},
			),
			withPathParam: true,
			pathParamId:   "1",
			shouldError:   false,
		},
		{
			formData: testutil.ConstructTestFormDataString(
				testutil.TestFormDataStringOptions{},
			),
			withPathParam: true,
			shouldError:   true,
		},
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
			pathId:   test.pathParamId,
		})
		_, formErrors, err := recipeService.parseFormData(c, test.withPathParam)

		assert.Equal(t, test.shouldError, err != nil)
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
	formData := "title=title&description=description&ingredients=ingredients&instructions=instructions&duration=20"
	c := newTestContext(t, newTestContextOptions{
		formData: formData,
		pathId:   "1",
	})

	// When
	parsedRecipe, _, err := recipeService.parseFormData(c, true)
	assert.NoError(t, err)
	assertRecipesEqual(t, types.Recipe{
		ID:           uint(1),
		Title:        "title",
		Description:  "description",
		Instructions: "instructions",
		Ingredients:  "ingredients",
		Duration:     20,
	}, parsedRecipe)
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
