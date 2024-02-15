package recipe

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func newTestRecipeController() RecipeController {
	return RecipeController{
		recipeService: newTestRecipeService(),
	}
}

type requestOptions struct {
	RecipeController *RecipeController
	handlerFunc      func(c echo.Context) error
	method           string
	route            string
	statusWant       int
	withFormData     bool
	formData         string
	authorized       bool
	withPathParam    bool
	pathParamId      string
}

func assertRequest(t *testing.T, options requestOptions) (*httptest.ResponseRecorder, echo.Context) {
	e := echo.New()
	w := httptest.NewRecorder()

	var body io.Reader

	if options.withFormData {
		body = bytes.NewBufferString(options.formData)
	} else {
		body = nil
	}

	req, err := http.NewRequest(options.method, options.route, body)
	assert.NoError(t, err)

	if options.withFormData {
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	}

	c := e.NewContext(req, w)
	if options.authorized {
		c.Set("authorized", true)
	} else {
		c.Set("authorized", false)
	}

	if options.withPathParam {
		c.SetParamNames("id")
		c.SetParamValues(options.pathParamId)
	}

	options.handlerFunc(c)

	assert.Equal(t, options.statusWant, w.Code)

	return w, c
}

func TestRenderRecipeListPage(t *testing.T) {
	recipeController := newTestRecipeController()

	t.Run("valid request", func(t *testing.T) {
		assertRequest(
			t,
			requestOptions{
				RecipeController: &recipeController,
				handlerFunc:      recipeController.RenderRecipeListPage,
				method:           http.MethodGet,
				route:            "/",
				statusWant:       http.StatusOK,
			},
		)
	})
}

func TestRenderRecipeNewPage(t *testing.T) {
	recipeController := newTestRecipeController()

	t.Run("not authorized", func(t *testing.T) {
		assertRequest(
			t,
			requestOptions{
				RecipeController: &recipeController,
				handlerFunc:      recipeController.RenderRecipeNewPage,
				method:           http.MethodGet,
				route:            "/recipe/new",
				statusWant:       http.StatusUnauthorized,
			},
		)
	})

	t.Run("authorized", func(t *testing.T) {
		assertRequest(
			t,
			requestOptions{
				RecipeController: &recipeController,
				handlerFunc:      recipeController.RenderRecipeNewPage,
				method:           http.MethodGet,
				route:            "/recipe/new",
				statusWant:       http.StatusOK,
				authorized:       true,
			},
		)
	})
}

func TestRenderRecipeEditPage(t *testing.T) {
	recipeController := newTestRecipeController()

	t.Run("not authorized", func(t *testing.T) {
		assertRequest(
			t,
			requestOptions{
				RecipeController: &recipeController,
				handlerFunc:      recipeController.RenderRecipeEditPage,
				method:           http.MethodGet,
				route:            "/recipe/edit",
				withPathParam:    true,
				pathParamId:      "1",
				statusWant:       http.StatusUnauthorized,
			},
		)
	})

	t.Run("recipe not found", func(t *testing.T) {
		assertRequest(
			t,
			requestOptions{
				RecipeController: &recipeController,
				handlerFunc:      recipeController.RenderRecipeEditPage,
				method:           http.MethodGet,
				route:            "/recipe/edit",
				withPathParam:    true,
				pathParamId:      "1",
				statusWant:       http.StatusNotFound,
				authorized:       true,
			},
		)
	})

	err := recipeController.recipeService.createRecipe(&recipe{})
	assert.NoError(t, err)

	t.Run("authorized", func(t *testing.T) {
		assertRequest(
			t,
			requestOptions{
				RecipeController: &recipeController,
				handlerFunc:      recipeController.RenderRecipeEditPage,
				method:           http.MethodGet,
				route:            "/recipe/edit",
				withPathParam:    true,
				pathParamId:      "1",
				statusWant:       http.StatusOK,
				authorized:       true,
			},
		)
	})

	t.Run("invalid path param", func(t *testing.T) {
		assertRequest(
			t,
			requestOptions{
				RecipeController: &recipeController,
				handlerFunc:      recipeController.RenderRecipeEditPage,
				method:           http.MethodGet,
				route:            "/recipe/edit",
				withPathParam:    true,
				pathParamId:      "xx",
				statusWant:       http.StatusBadRequest,
				authorized:       true,
			},
		)
	})
}

func TestRenderRecipePage(t *testing.T) {
	recipeController := newTestRecipeController()

	t.Run("recipe not found", func(t *testing.T) {
		assertRequest(
			t,
			requestOptions{
				RecipeController: &recipeController,
				handlerFunc:      recipeController.RenderRecipePage,
				method:           http.MethodGet,
				route:            "/recipe",
				withPathParam:    true,
				pathParamId:      "1",
				statusWant:       http.StatusNotFound,
			},
		)
	})

	t.Run("invalid path param", func(t *testing.T) {
		assertRequest(
			t,
			requestOptions{
				RecipeController: &recipeController,
				handlerFunc:      recipeController.RenderRecipePage,
				method:           http.MethodGet,
				route:            "/recipe",
				withPathParam:    true,
				pathParamId:      "xx",
				statusWant:       http.StatusBadRequest,
			},
		)
	})

	err := recipeController.recipeService.createRecipe(&recipe{})
	assert.NoError(t, err)

	t.Run("valid request", func(t *testing.T) {
		assertRequest(
			t,
			requestOptions{
				RecipeController: &recipeController,
				handlerFunc:      recipeController.RenderRecipePage,
				method:           http.MethodGet,
				route:            "/recipe",
				withPathParam:    true,
				pathParamId:      "1",
				statusWant:       http.StatusOK,
			},
		)
	})
}

func TestHandleCreateRecipe(t *testing.T) {
	recipeController := newTestRecipeController()

	t.Run("not authorized", func(t *testing.T) {
		assertRequest(
			t,
			requestOptions{
				RecipeController: &recipeController,
				handlerFunc:      recipeController.HandleCreateRecipe,
				method:           http.MethodPost,
				route:            "/recipe",
				statusWant:       http.StatusUnauthorized,
			},
		)
	})

	t.Run("no form data", func(t *testing.T) {
		assertRequest(
			t,
			requestOptions{
				RecipeController: &recipeController,
				handlerFunc:      recipeController.HandleCreateRecipe,
				method:           http.MethodPost,
				route:            "/recipe",
				authorized:       true,
				statusWant:       http.StatusBadRequest,
			},
		)
	})

	t.Run("invalid form data", func(t *testing.T) {
		assertRequest(
			t,
			requestOptions{
				RecipeController: &recipeController,
				handlerFunc:      recipeController.HandleCreateRecipe,
				method:           http.MethodPost,
				route:            "/recipe",
				authorized:       true,
				statusWant:       http.StatusBadRequest,
				withFormData:     true,
				formData:         "title=title&description=description&ingredients=ingredients&instructions=instructions&duration=xx",
			},
		)
	})

	t.Run("valid request", func(t *testing.T) {
		assertRequest(
			t,
			requestOptions{
				RecipeController: &recipeController,
				handlerFunc:      recipeController.HandleCreateRecipe,
				method:           http.MethodPost,
				route:            "/recipe",
				authorized:       true,
				statusWant:       http.StatusOK,
				withFormData:     true,
				formData:         "title=title&description=description&ingredients=ingredients&instructions=instructions&duration=30",
			},
		)
	})
}

func TestHandleUpdateRecipe(t *testing.T) {
	recipeController := newTestRecipeController()

	t.Run("not authorized", func(t *testing.T) {
		assertRequest(
			t,
			requestOptions{
				RecipeController: &recipeController,
				handlerFunc:      recipeController.HandleUpdateRecipe,
				method:           http.MethodPut,
				route:            "/recipe",
				statusWant:       http.StatusUnauthorized,
			},
		)
	})

	t.Run("no form data", func(t *testing.T) {
		assertRequest(
			t,
			requestOptions{
				RecipeController: &recipeController,
				handlerFunc:      recipeController.HandleUpdateRecipe,
				method:           http.MethodPut,
				route:            "/recipe",
				authorized:       true,
				statusWant:       http.StatusBadRequest,
			},
		)
	})

	t.Run("invalid form data", func(t *testing.T) {
		assertRequest(
			t,
			requestOptions{
				RecipeController: &recipeController,
				handlerFunc:      recipeController.HandleUpdateRecipe,
				method:           http.MethodPut,
				route:            "/recipe",
				authorized:       true,
				statusWant:       http.StatusBadRequest,
				withFormData:     true,
				formData:         "title=title&description=description&ingredients=ingredients&instructions=instructions&duration=xx",
			},
		)
	})

	t.Run("invalid id", func(t *testing.T) {
		assertRequest(
			t,
			requestOptions{
				RecipeController: &recipeController,
				handlerFunc:      recipeController.HandleUpdateRecipe,
				method:           http.MethodPut,
				route:            "/recipe",
				withPathParam:    true,
				pathParamId:      "xx",
				authorized:       true,
				statusWant:       http.StatusBadRequest,
				withFormData:     true,
				formData:         "title=title&description=description&ingredients=ingredients&instructions=instructions&duration=30",
			},
		)
	})

	err := recipeController.recipeService.createRecipe(&recipe{})
	assert.NoError(t, err)

	t.Run("valid request", func(t *testing.T) {
		assertRequest(
			t,
			requestOptions{
				RecipeController: &recipeController,
				handlerFunc:      recipeController.HandleUpdateRecipe,
				method:           http.MethodPut,
				route:            "/recipe",
				withPathParam:    true,
				pathParamId:      "1",
				authorized:       true,
				statusWant:       http.StatusOK,
				withFormData:     true,
				formData:         "title=title&description=description&ingredients=ingredients&instructions=instructions&duration=30",
			},
		)
	})
}

func TestHandleDeleteRecipe(t *testing.T) {
	recipeController := newTestRecipeController()

	t.Run("not authorized", func(t *testing.T) {
		assertRequest(
			t,
			requestOptions{
				RecipeController: &recipeController,
				handlerFunc:      recipeController.HandleDeleteRecipe,
				method:           http.MethodDelete,
				route:            "/recipe",
				withPathParam:    true,
				pathParamId:      "1",
				statusWant:       http.StatusUnauthorized,
			},
		)
	})

	t.Run("not found", func(t *testing.T) {
		assertRequest(
			t,
			requestOptions{
				RecipeController: &recipeController,
				handlerFunc:      recipeController.HandleDeleteRecipe,
				method:           http.MethodDelete,
				route:            "/recipe",
				withPathParam:    true,
				pathParamId:      "1",
				statusWant:       http.StatusNotFound,
				authorized:       true,
			},
		)
	})

	t.Run("invalid path param", func(t *testing.T) {
		assertRequest(
			t,
			requestOptions{
				RecipeController: &recipeController,
				handlerFunc:      recipeController.HandleDeleteRecipe,
				method:           http.MethodDelete,
				route:            "/recipe",
				withPathParam:    true,
				pathParamId:      "xx",
				statusWant:       http.StatusBadRequest,
				authorized:       true,
			},
		)
	})

	err := recipeController.recipeService.createRecipe(&recipe{})
	assert.NoError(t, err)

	t.Run("valid request", func(t *testing.T) {
		assertRequest(
			t,
			requestOptions{
				RecipeController: &recipeController,
				handlerFunc:      recipeController.HandleDeleteRecipe,
				method:           http.MethodDelete,
				route:            "/recipe",
				withPathParam:    true,
				pathParamId:      "1",
				statusWant:       http.StatusOK,
				authorized:       true,
			},
		)
	})
}
