package recipe

import (
	"net/http"
	"testing"

	"github.com/kilianmandscharo/lethimcook/logging"
	"github.com/kilianmandscharo/lethimcook/render"
	"github.com/kilianmandscharo/lethimcook/testutil"
	"github.com/kilianmandscharo/lethimcook/types"
	"github.com/stretchr/testify/assert"
)

func newTestRecipeController() RecipeController {
	logger := logging.New(logging.Debug)
	renderer := render.New(&logger)
	recipeService := newTestRecipeService()
	return NewRecipeController(recipeService, &logger, &renderer)
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

	t.Run("not authorized and not pending", func(t *testing.T) {
		testutil.AssertRequest(
			t,
			testutil.RequestOptions{
				HandlerFunc:    recipeController.HandleCreateRecipe,
				Method:         http.MethodPost,
				Route:          "/recipe",
				StatusWant:     http.StatusUnauthorized,
				WithQueryParam: true,
				QueryParam:     "?pending=false",
			},
		)
	})

	t.Run("not authorized but pending", func(t *testing.T) {
		testutil.AssertRequest(
			t,
			testutil.RequestOptions{
				HandlerFunc:    recipeController.HandleCreateRecipe,
				Method:         http.MethodPost,
				Route:          "/recipe",
				StatusWant:     http.StatusOK,
				WithQueryParam: true,
				QueryParam:     "?pending=true",
				WithFormData:   true,
				FormData:       "title=title&description=description&ingredients=ingredients&instructions=instructions&duration=30",
				AssertMessage:  true,
				MessageWant:    "Rezept eingereicht",
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
				HandlerFunc:   recipeController.HandleCreateRecipe,
				Method:        http.MethodPost,
				Route:         "/recipe",
				Authorized:    true,
				StatusWant:    http.StatusOK,
				WithFormData:  true,
				FormData:      "title=title&description=description&ingredients=ingredients&instructions=instructions&duration=30",
				AssertMessage: true,
				MessageWant:   "Rezept erstellt",
			},
		)
	})
}

func TestHandleUpdatePending(t *testing.T) {
	recipeController := newTestRecipeController()

	t.Run("not authorized", func(t *testing.T) {
		testutil.AssertRequest(
			t,
			testutil.RequestOptions{
				HandlerFunc: recipeController.HandleUpdatePending,
				Method:      http.MethodPut,
				Route:       "/recipe/:id/pending/:pending",
				StatusWant:  http.StatusUnauthorized,
			},
		)
	})

	t.Run("no pending path param", func(t *testing.T) {
		testutil.AssertRequest(
			t,
			testutil.RequestOptions{
				HandlerFunc: recipeController.HandleUpdatePending,
				Method:      http.MethodPut,
				Route:       "/recipe/:id/pending/:pending",
				StatusWant:  http.StatusBadRequest,
				Authorized:  true,
			},
		)
	})

	t.Run("no id path param", func(t *testing.T) {
		testutil.AssertRequest(
			t,
			testutil.RequestOptions{
				HandlerFunc:    recipeController.HandleUpdatePending,
				Method:         http.MethodPut,
				Route:          "/recipe/:id/pending/:pending",
				StatusWant:     http.StatusBadRequest,
				Authorized:     true,
				WithPathParam:  true,
				PathParamName:  "pending",
				PathParamValue: "true",
			},
		)
	})

	t.Run("no recipe found", func(t *testing.T) {
		testutil.AssertRequest(
			t,
			testutil.RequestOptions{
				HandlerFunc:     recipeController.HandleUpdatePending,
				Method:          http.MethodPut,
				Route:           "/recipe/:id/pending/:pending",
				StatusWant:      http.StatusNotFound,
				Authorized:      true,
				WithPathParam:   true,
				PathParamNames:  []string{"pending", "id"},
				PathParamValues: []string{"false", "5"},
			},
		)
	})

	t.Run("valid with pending false", func(t *testing.T) {
		assert.NoError(
			t,
			recipeController.recipeService.createRecipe(&types.Recipe{ID: 1}),
		)
		testutil.AssertRequest(
			t,
			testutil.RequestOptions{
				HandlerFunc:     recipeController.HandleUpdatePending,
				Method:          http.MethodPut,
				Route:           "/recipe/:id/pending/:pending",
				StatusWant:      http.StatusOK,
				Authorized:      true,
				WithPathParam:   true,
				PathParamNames:  []string{"pending", "id"},
				PathParamValues: []string{"false", "1"},
				AssertMessage:   true,
				MessageWant:     "Rezept angenommen",
			},
		)
	})

	t.Run("valid with pending true", func(t *testing.T) {
		testutil.AssertRequest(
			t,
			testutil.RequestOptions{
				HandlerFunc:     recipeController.HandleUpdatePending,
				Method:          http.MethodPut,
				Route:           "/recipe/:id/pending/:pending",
				StatusWant:      http.StatusOK,
				Authorized:      true,
				WithPathParam:   true,
				PathParamNames:  []string{"pending", "id"},
				PathParamValues: []string{"true", "1"},
				AssertMessage:   true,
				MessageWant:     "Rezept auf &#39;ausstehend&#39; gesetzt<",
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
				AssertMessage:  true,
				MessageWant:    "Ungültiges Pfadparameter",
			},
		)
	})

	t.Run("recipe not found", func(t *testing.T) {
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
				StatusWant:     http.StatusNotFound,
				AssertMessage:  true,
				MessageWant:    "Rezept nicht gefunden",
			},
		)
	})

	assert.NoError(
		t,
		recipeController.recipeService.createRecipe(&types.Recipe{}),
	)

	t.Run("no form data", func(t *testing.T) {
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
				StatusWant:     http.StatusBadRequest,
				AssertMessage:  true,
				MessageWant:    "Fehlerhaftes Formular",
			},
		)
	})

	t.Run("invalid form data", func(t *testing.T) {
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
				StatusWant:     http.StatusBadRequest,
				WithFormData:   true,
				FormData:       "title=title&description=description&ingredients=ingredients&&duration=xx",
				AssertMessage:  true,
				MessageWant:    "Fehlerhaftes Formular",
			},
		)
	})

	err := recipeController.recipeService.createRecipe(&types.Recipe{})
	assert.NoError(t, err)

	t.Run("valid request", func(t *testing.T) {
		testutil.AssertRequest(
			t,
			testutil.RequestOptions{
				HandlerFunc:     recipeController.HandleUpdateRecipe,
				Method:          http.MethodPut,
				Route:           "/recipe",
				WithPathParam:   true,
				PathParamNames:  []string{"id", "pending"},
				PathParamValues: []string{"1", "false"},
				Authorized:      true,
				StatusWant:      http.StatusOK,
				WithFormData:    true,
				FormData:        "title=title&description=description&ingredients=ingredients&instructions=instructions&duration=30",
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

	assert.NoError(
		t,
		recipeController.recipeService.createRecipe(&types.Recipe{}),
	)

	t.Run("valid request without force", func(t *testing.T) {
		testutil.AssertRequest(
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
				AssertMessage:  true,
				MessageWant:    "Rezept entfernt",
			},
		)
	})
}

func TestHandleDownloadRecipeAsJson(t *testing.T) {
	recipeController := newTestRecipeController()

	t.Run("invalid path id", func(t *testing.T) {
		testutil.AssertRequest(
			t,
			testutil.RequestOptions{
				HandlerFunc:    recipeController.HandleDownloadRecipeAsJson,
				Method:         http.MethodGet,
				Route:          "/recipe/json",
				WithPathParam:  true,
				PathParamName:  "id",
				PathParamValue: "xx",
				StatusWant:     http.StatusBadRequest,
				AssertMessage:  true,
				MessageWant:    "Ungültiges Pfadparameter",
			},
		)
	})

	t.Run("recipe not found", func(t *testing.T) {
		testutil.AssertRequest(
			t,
			testutil.RequestOptions{
				HandlerFunc:    recipeController.HandleDownloadRecipeAsJson,
				Method:         http.MethodGet,
				Route:          "/recipe/json",
				WithPathParam:  true,
				PathParamName:  "id",
				PathParamValue: "1",
				StatusWant:     http.StatusNotFound,
				AssertMessage:  true,
				MessageWant:    "Rezept nicht gefunden",
			},
		)
	})

	assert.NoError(t, recipeController.recipeService.createRecipe(&types.Recipe{ID: 1}))

	t.Run("valid request", func(t *testing.T) {
		testutil.AssertRequest(
			t,
			testutil.RequestOptions{
				HandlerFunc:    recipeController.HandleDownloadRecipeAsJson,
				Method:         http.MethodGet,
				Route:          "/recipe/json",
				WithPathParam:  true,
				PathParamName:  "id",
				PathParamValue: "1",
				StatusWant:     http.StatusOK,
			},
		)
	})
}

func TestHandleGetPaginatedRecipe(t *testing.T) {
	recipeController := newTestRecipeController()

	t.Run("valid request", func(t *testing.T) {
		testutil.AssertRequest(
			t,
			testutil.RequestOptions{
				HandlerFunc: recipeController.HandleGetPaginatedRecipes,
				Method:      http.MethodGet,
				Route:       "/",
				StatusWant:  http.StatusOK,
			},
		)
	})
}
