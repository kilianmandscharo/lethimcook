package recipe

import (
	"net/http"
	"strings"
	"testing"

	"github.com/kilianmandscharo/lethimcook/testutil"
	"github.com/kilianmandscharo/lethimcook/types"
	"github.com/stretchr/testify/assert"
)

func newTestRecipeController() RecipeController {
	return RecipeController{
		recipeService: newTestRecipeService(),
	}
}

func TestRenderRecipeListPage(t *testing.T) {
	recipeController := newTestRecipeController()

	t.Run("valid request", func(t *testing.T) {
		testutil.AssertRequest(
			t,
			testutil.RequestOptions{
				HandlerFunc: recipeController.RenderRecipeListPage,
				Method:      http.MethodGet,
				Route:       "/",
				StatusWant:  http.StatusOK,
			},
		)
	})
}

func TestRenderRecipeNewPage(t *testing.T) {
	recipeController := newTestRecipeController()

	t.Run("not authorized", func(t *testing.T) {
		testutil.AssertRequest(
			t,
			testutil.RequestOptions{
				HandlerFunc: recipeController.RenderRecipeNewPage,
				Method:      http.MethodGet,
				Route:       "/recipe/new",
				StatusWant:  http.StatusUnauthorized,
			},
		)
	})

	t.Run("authorized", func(t *testing.T) {
		testutil.AssertRequest(
			t,
			testutil.RequestOptions{
				HandlerFunc: recipeController.RenderRecipeNewPage,
				Method:      http.MethodGet,
				Route:       "/recipe/new",
				StatusWant:  http.StatusOK,
				Authorized:  true,
			},
		)
	})
}

func TestRenderRecipeEditPage(t *testing.T) {
	recipeController := newTestRecipeController()

	t.Run("not authorized", func(t *testing.T) {
		testutil.AssertRequest(
			t,
			testutil.RequestOptions{
				HandlerFunc:    recipeController.RenderRecipeEditPage,
				Method:         http.MethodGet,
				Route:          "/recipe/edit",
				WithPathParam:  true,
				PathParamName:  "id",
				PathParamValue: "1",
				StatusWant:     http.StatusUnauthorized,
			},
		)
	})

	t.Run("recipe not found", func(t *testing.T) {
		testutil.AssertRequest(
			t,
			testutil.RequestOptions{
				HandlerFunc:    recipeController.RenderRecipeEditPage,
				Method:         http.MethodGet,
				Route:          "/recipe/edit",
				WithPathParam:  true,
				PathParamName:  "id",
				PathParamValue: "1",
				StatusWant:     http.StatusNotFound,
				Authorized:     true,
			},
		)
	})

	err := recipeController.recipeService.createRecipe(&types.Recipe{})
	assert.NoError(t, err)

	t.Run("authorized", func(t *testing.T) {
		testutil.AssertRequest(
			t,
			testutil.RequestOptions{
				HandlerFunc:    recipeController.RenderRecipeEditPage,
				Method:         http.MethodGet,
				Route:          "/recipe/edit",
				WithPathParam:  true,
				PathParamName:  "id",
				PathParamValue: "1",
				StatusWant:     http.StatusOK,
				Authorized:     true,
			},
		)
	})

	t.Run("invalid path param", func(t *testing.T) {
		testutil.AssertRequest(
			t,
			testutil.RequestOptions{
				HandlerFunc:    recipeController.RenderRecipeEditPage,
				Method:         http.MethodGet,
				Route:          "/recipe/edit",
				WithPathParam:  true,
				PathParamName:  "id",
				PathParamValue: "xx",
				StatusWant:     http.StatusBadRequest,
				Authorized:     true,
			},
		)
	})
}

func TestRenderRecipePage(t *testing.T) {
	recipeController := newTestRecipeController()

	t.Run("recipe not found", func(t *testing.T) {
		testutil.AssertRequest(
			t,
			testutil.RequestOptions{
				HandlerFunc:    recipeController.RenderRecipePage,
				Method:         http.MethodGet,
				Route:          "/recipe",
				WithPathParam:  true,
				PathParamName:  "id",
				PathParamValue: "1",
				StatusWant:     http.StatusNotFound,
			},
		)
	})

	t.Run("invalid path param", func(t *testing.T) {
		testutil.AssertRequest(
			t,
			testutil.RequestOptions{
				HandlerFunc:    recipeController.RenderRecipePage,
				Method:         http.MethodGet,
				Route:          "/recipe",
				WithPathParam:  true,
				PathParamName:  "id",
				PathParamValue: "xx",
				StatusWant:     http.StatusBadRequest,
			},
		)
	})

	err := recipeController.recipeService.createRecipe(&types.Recipe{})
	assert.NoError(t, err)

	t.Run("valid request", func(t *testing.T) {
		testutil.AssertRequest(
			t,
			testutil.RequestOptions{
				HandlerFunc:    recipeController.RenderRecipePage,
				Method:         http.MethodGet,
				Route:          "/recipe",
				WithPathParam:  true,
				PathParamName:  "id",
				PathParamValue: "1",
				StatusWant:     http.StatusOK,
			},
		)
	})
}

func TestHandleCreateRecipe(t *testing.T) {
	recipeController := newTestRecipeController()

	t.Run("not authorized", func(t *testing.T) {
		testutil.AssertRequest(
			t,
			testutil.RequestOptions{
				HandlerFunc: recipeController.HandleCreateRecipe,
				Method:      http.MethodPost,
				Route:       "/recipe",
				StatusWant:  http.StatusUnauthorized,
			},
		)
	})

	t.Run("no form data", func(t *testing.T) {
		testutil.AssertRequest(
			t,
			testutil.RequestOptions{
				HandlerFunc: recipeController.HandleCreateRecipe,
				Method:      http.MethodPost,
				Route:       "/recipe",
				Authorized:  true,
				StatusWant:  http.StatusBadRequest,
			},
		)
	})

	t.Run("invalid form data", func(t *testing.T) {
		testutil.AssertRequest(
			t,
			testutil.RequestOptions{
				HandlerFunc:  recipeController.HandleCreateRecipe,
				Method:       http.MethodPost,
				Route:        "/recipe",
				Authorized:   true,
				StatusWant:   http.StatusBadRequest,
				WithFormData: true,
				FormData:     "title=title&description=description&ingredients=ingredients&instructions=instructions&duration=xx",
			},
		)
	})

	t.Run("valid request", func(t *testing.T) {
		testutil.AssertRequest(
			t,
			testutil.RequestOptions{
				HandlerFunc:  recipeController.HandleCreateRecipe,
				Method:       http.MethodPost,
				Route:        "/recipe",
				Authorized:   true,
				StatusWant:   http.StatusOK,
				WithFormData: true,
				FormData:     "title=title&description=description&ingredients=ingredients&instructions=instructions&duration=30",
			},
		)
	})
}

func TestHandleUpdateRecipe(t *testing.T) {
	recipeController := newTestRecipeController()

	t.Run("not authorized", func(t *testing.T) {
		testutil.AssertRequest(
			t,
			testutil.RequestOptions{
				HandlerFunc: recipeController.HandleUpdateRecipe,
				Method:      http.MethodPut,
				Route:       "/recipe",
				StatusWant:  http.StatusUnauthorized,
			},
		)
	})

	t.Run("no form data", func(t *testing.T) {
		testutil.AssertRequest(
			t,
			testutil.RequestOptions{
				HandlerFunc: recipeController.HandleUpdateRecipe,
				Method:      http.MethodPut,
				Route:       "/recipe",
				Authorized:  true,
				StatusWant:  http.StatusBadRequest,
			},
		)
	})

	t.Run("invalid form data", func(t *testing.T) {
		testutil.AssertRequest(
			t,
			testutil.RequestOptions{
				HandlerFunc:  recipeController.HandleUpdateRecipe,
				Method:       http.MethodPut,
				Route:        "/recipe",
				Authorized:   true,
				StatusWant:   http.StatusBadRequest,
				WithFormData: true,
				FormData:     "title=title&description=description&ingredients=ingredients&instructions=instructions&duration=xx",
			},
		)
	})

	t.Run("invalid id", func(t *testing.T) {
		testutil.AssertRequest(
			t,
			testutil.RequestOptions{
				HandlerFunc:    recipeController.HandleUpdateRecipe,
				Method:         http.MethodPut,
				Route:          "/recipe",
				WithPathParam:  true,
				PathParamName:  "id",
				PathParamValue: "xx",
				Authorized:     true,
				StatusWant:     http.StatusBadRequest,
				WithFormData:   true,
				FormData:       "title=title&description=description&ingredients=ingredients&instructions=instructions&duration=30",
			},
		)
	})

	err := recipeController.recipeService.createRecipe(&types.Recipe{})
	assert.NoError(t, err)

	t.Run("valid request", func(t *testing.T) {
		testutil.AssertRequest(
			t,
			testutil.RequestOptions{
				HandlerFunc:    recipeController.HandleUpdateRecipe,
				Method:         http.MethodPut,
				Route:          "/recipe",
				WithPathParam:  true,
				PathParamName:  "id",
				PathParamValue: "1",
				Authorized:     true,
				StatusWant:     http.StatusOK,
				WithFormData:   true,
				FormData:       "title=title&description=description&ingredients=ingredients&instructions=instructions&duration=30",
			},
		)
	})
}

func TestHandleDeleteRecipe(t *testing.T) {
	recipeController := newTestRecipeController()

	t.Run("not authorized", func(t *testing.T) {
		testutil.AssertRequest(
			t,
			testutil.RequestOptions{
				HandlerFunc:    recipeController.HandleDeleteRecipe,
				Method:         http.MethodDelete,
				Route:          "/recipe",
				WithPathParam:  true,
				PathParamName:  "id",
				PathParamValue: "1",
				StatusWant:     http.StatusUnauthorized,
			},
		)
	})

	t.Run("not found", func(t *testing.T) {
		testutil.AssertRequest(
			t,
			testutil.RequestOptions{
				HandlerFunc:    recipeController.HandleDeleteRecipe,
				Method:         http.MethodDelete,
				Route:          "/recipe",
				WithPathParam:  true,
				PathParamName:  "id",
				PathParamValue: "1",
				StatusWant:     http.StatusNotFound,
				Authorized:     true,
			},
		)
	})

	t.Run("invalid path param", func(t *testing.T) {
		testutil.AssertRequest(
			t,
			testutil.RequestOptions{
				HandlerFunc:    recipeController.HandleDeleteRecipe,
				Method:         http.MethodDelete,
				Route:          "/recipe",
				WithPathParam:  true,
				PathParamName:  "id",
				PathParamValue: "xx",
				StatusWant:     http.StatusBadRequest,
				Authorized:     true,
			},
		)
	})

	err := recipeController.recipeService.createRecipe(&types.Recipe{})
	assert.NoError(t, err)

	t.Run("valid request without force", func(t *testing.T) {
		rr, _ := testutil.AssertRequest(
			t,
			testutil.RequestOptions{
				HandlerFunc:    recipeController.HandleDeleteRecipe,
				Method:         http.MethodDelete,
				Route:          "/recipe",
				WithPathParam:  true,
				PathParamName:  "id",
				PathParamValue: "1",
				StatusWant:     http.StatusOK,
				Authorized:     true,
			},
		)
		assert.True(t, strings.Contains(rr.Body.String(), "Löschen bestätigen"))
	})

	t.Run("valid request with cancel", func(t *testing.T) {
		rr, _ := testutil.AssertRequest(
			t,
			testutil.RequestOptions{
				HandlerFunc:    recipeController.HandleDeleteRecipe,
				Method:         http.MethodDelete,
				Route:          "/recipe",
				WithPathParam:  true,
				PathParamName:  "id",
				PathParamValue: "1",
				WithQueryParam: true,
				QueryParam:     "?cancel=true",
				StatusWant:     http.StatusOK,
				Authorized:     true,
			},
		)
		resBody := rr.Body.String()
		assert.False(t, strings.Contains(resBody, "Löschen bestätigen"))
		assert.True(t, len(resBody) != 0)
	})

	t.Run("valid request with force", func(t *testing.T) {
		rr, _ := testutil.AssertRequest(
			t,
			testutil.RequestOptions{
				HandlerFunc:    recipeController.HandleDeleteRecipe,
				Method:         http.MethodDelete,
				Route:          "/recipe",
				WithPathParam:  true,
				PathParamName:  "id",
				PathParamValue: "1",
				WithQueryParam: true,
				QueryParam:     "?force=true",
				StatusWant:     http.StatusOK,
				Authorized:     true,
			},
		)
		resBody := rr.Body.String()
		assert.False(t, strings.Contains(resBody, "Löschen bestätigen"))
		assert.NotEqual(t, 0, len(resBody))
	})
}

func TestHandleSearchRecipe(t *testing.T) {
	recipeController := newTestRecipeController()

	t.Run("valid request", func(t *testing.T) {
		testutil.AssertRequest(
			t,
			testutil.RequestOptions{
				HandlerFunc:  recipeController.HandleSearchRecipe,
				Method:       http.MethodPost,
				Route:        "/search",
				WithFormData: true,
				FormData:     "query=test",
				StatusWant:   http.StatusOK,
			},
		)
	})
}
