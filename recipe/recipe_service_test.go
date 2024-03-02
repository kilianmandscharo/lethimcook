package recipe

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kilianmandscharo/lethimcook/errutil"
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
	assert.NoError(t, recipeService.createRecipe(&recipe{
		Title: "Naan",
	}))
	assert.NoError(t, recipeService.createRecipe(&recipe{
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
		shouldBeValid bool
	}{
		{
			formData:      "title=title&description=description&ingredients=ingredients&instructions=instructions&duration=10",
			withPathParam: true,
			pathParamId:   "1",
			shouldBeValid: true,
		},
		{
			formData:      "title=title&description=description&ingredients=ingredients&instructions=instructions&duration=10",
			shouldBeValid: true,
		},
		{
			formData:      "title=title&description=description&ingredients=ingredients&instructions=instructions&duration=10",
			withPathParam: true,
			shouldBeValid: false,
		},
		{
			formData:      "title=title&description=description&ingredients=ingredients&instructions=instructions&duration=xx",
			shouldBeValid: false,
		},
		{
			formData:      "description=description&ingredients=ingredients&instructions=instructions&duration=xx",
			shouldBeValid: false,
		},
		{
			formData:      "title=title&ingredients=ingredients&instructions=instructions&duration=xx",
			shouldBeValid: false,
		},
		{
			formData:      "title=title&description=description&instructions=instructions&duration=xx",
			shouldBeValid: false,
		},
		{
			formData:      "title=title&description=description&ingredients=ingredients&duration=xx",
			shouldBeValid: false,
		},
		{
			formData:      "title=title&description=description&ingredients=ingredients&instructions=instructions",
			shouldBeValid: false,
		},
	}

	for _, test := range testCases {
		c := newTestContext(t, test.formData, test.pathParamId)
		_, err := recipeService.parseFormData(c, test.withPathParam)
		assert.Equal(t, test.shouldBeValid, err == nil)
	}

	// Given
	formData := "title=title&description=description&ingredients=ingredients&instructions=instructions&duration=20"
	c := newTestContext(t, formData, "1")

	// When
	parsedRecipe, err := recipeService.parseFormData(c, true)
	assert.NoError(t, err)
	assertRecipesEqual(t, recipe{
		ID:           uint(1),
		Title:        "title",
		Description:  "description",
		Instructions: "instructions",
		Ingredients:  "ingredients",
		Duration:     20,
	}, parsedRecipe)
}
