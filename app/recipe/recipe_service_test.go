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

func newTestContext(t *testing.T, formData string, pathId string) echo.Context {
	e := echo.New()
	w := httptest.NewRecorder()

	var body io.Reader
	body = bytes.NewBufferString(formData)
	req, err := http.NewRequest(http.MethodPost, "", body)
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

	c := e.NewContext(req, w)
	c.SetParamNames("id")
	c.SetParamValues(pathId)

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
		filteredRecipes, err := recipeService.getFilteredRecipes(test.query)
		assert.NoError(t, err)
		assert.Equal(t, test.hits, len(filteredRecipes))
	}
}

func TestGetPathId(t *testing.T) {
	// Given
	recipeService := newTestRecipeService()
	c := newTestContext(t, "", "xx")

	// When
	_, err := recipeService.getPathId(c)

	// Then
	assert.ErrorIs(t, err, errutil.RecipeErrorInvalidParam)

	// Given
	c = newTestContext(t, "", "1")

	// When
	id, err := recipeService.getPathId(c)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, uint(1), id)
}

func TestParseFormData(t *testing.T) {
	// Given
	recipeService := newTestRecipeService()

	testCases := []struct {
		formData      string
		withPathParam bool
		pathParamId   string
		shouldError   bool
		formErrors    []errutil.FormError
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
			formErrors: []errutil.FormError{errutil.FormErrorNoTitle},
		},
		{
			formData: testutil.ConstructTestFormDataString(
				testutil.TestFormDataStringOptions{
					DescriptionEmpty: true,
				},
			),
			formErrors: []errutil.FormError{errutil.FormErrorNoDescription},
		},
		{
			formData: testutil.ConstructTestFormDataString(
				testutil.TestFormDataStringOptions{
					DurationEmpty: true,
				},
			),
			formErrors: []errutil.FormError{errutil.FormErrorNoDuration},
		},
		{
			formData: testutil.ConstructTestFormDataString(
				testutil.TestFormDataStringOptions{
					InvalidDuration: true,
				},
			),
			formErrors: []errutil.FormError{errutil.FormErrorNoDuration},
		},
		{
			formData: testutil.ConstructTestFormDataString(
				testutil.TestFormDataStringOptions{
					IngredientsEmpty: true,
				},
			),
			formErrors: []errutil.FormError{errutil.FormErrorNoIngredients},
		},
		{
			formData: testutil.ConstructTestFormDataString(
				testutil.TestFormDataStringOptions{
					InstructionsEmpty: true,
				},
			),
			formErrors: []errutil.FormError{errutil.FormErrorNoInstructions},
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
			formErrors: []errutil.FormError{
				errutil.FormErrorNoTitle,
				errutil.FormErrorNoDescription,
				errutil.FormErrorNoDuration,
				errutil.FormErrorNoIngredients,
				errutil.FormErrorNoInstructions,
			},
		},
	}

	for _, test := range testCases {
		c := newTestContext(t, test.formData, test.pathParamId)
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
	c := newTestContext(t, formData, "1")

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
