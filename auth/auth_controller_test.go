package auth

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

const (
	testPassword = "test_password"
)

type controllerOptions struct {
	withAdmin bool
}

func newTestAuthController(options controllerOptions) AuthController {
	authService := newTestAuthService()

	if options.withAdmin {
		authService.createAdmin(testPassword)
	}

	return NewAuthController(authService)
}

func newTestCookie(t *testing.T, a *AuthController) http.Cookie {
	token, err := a.authService.createToken()
	assert.NoError(t, err)
	return a.authService.newTokenCookie(token, time.Now().Add(60*time.Minute))
}

type requestOptions struct {
	authController *AuthController
	handlerFunc    func(c echo.Context) error
	method         string
	route          string
	statusWant     int
	withFormData   bool
	formData       string
	withCookie     bool
	cookie         http.Cookie
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
	if options.withCookie {
		c.SetCookie(&options.cookie)
	}

	options.authController.ValidateToken(options.handlerFunc)(c)

	assert.Equal(t, options.statusWant, w.Code)

	return w, c
}

func isAuthorized(c echo.Context) bool {
	valueInterface := c.Get("authorized")
	if value, ok := valueInterface.(bool); ok {
		return value
	}
	return false
}

func TestRenderAdminPage(t *testing.T) {
	a := newTestAuthController(controllerOptions{})

	t.Run("valid request", func(t *testing.T) {
		assertRequest(
			t,
			requestOptions{
				authController: &a,
				handlerFunc:    a.RenderAdminPage,
				method:         http.MethodGet,
				route:          "/auth/admin",
				statusWant:     http.StatusOK,
			},
		)
	})
}

func TestHandleLogin(t *testing.T) {
	a := newTestAuthController(controllerOptions{})

	t.Run("no admin", func(t *testing.T) {
		w, c := assertRequest(
			t,
			requestOptions{
				authController: &a,
				handlerFunc:    a.HandleLogin,
				method:         http.MethodPost,
				route:          "/auth/login",
				statusWant:     http.StatusNotFound,
				withFormData:   true,
				formData:       "password=" + testPassword,
			},
		)
		assert.Equal(t, 0, len(w.Result().Cookies()))
		assert.False(t, isAuthorized(c))
	})

	a = newTestAuthController(controllerOptions{withAdmin: true})

	t.Run("invalid form key", func(t *testing.T) {
		w, c := assertRequest(
			t,
			requestOptions{
				authController: &a,
				handlerFunc:    a.HandleLogin,
				method:         http.MethodPost,
				route:          "/auth/login",
				statusWant:     http.StatusUnauthorized,
				withFormData:   true,
				formData:       "invalidKey=" + testPassword,
			},
		)
		assert.Equal(t, 0, len(w.Result().Cookies()))
		assert.False(t, isAuthorized(c))
	})

	t.Run("invalid password", func(t *testing.T) {
		w, c := assertRequest(
			t,
			requestOptions{
				authController: &a,
				handlerFunc:    a.HandleLogin,
				method:         http.MethodPost,
				route:          "/auth/login",
				statusWant:     http.StatusUnauthorized,
				withFormData:   true,
				formData:       "password=invalid_password",
			},
		)
		assert.Equal(t, 0, len(w.Result().Cookies()))
		assert.False(t, isAuthorized(c))
	})

	t.Run("valid password", func(t *testing.T) {
		w, c := assertRequest(
			t,
			requestOptions{
				authController: &a,
				handlerFunc:    a.HandleLogin,
				method:         http.MethodPost,
				route:          "/auth/login",
				statusWant:     http.StatusOK,
				withFormData:   true,
				formData:       "password=" + testPassword,
			},
		)
		assert.Equal(t, 1, len(w.Result().Cookies()))
		assert.True(t, isAuthorized(c))
	})
}

// func TestHandleLogout(t *testing.T) {
// 	a := newTestAuthController(controllerOptions{withAdmin: true})
//
// 	t.Run("no cookie", func(t *testing.T) {
// 		assertRequest(
// 			t,
// 			requestOptions{
// 				authController: &a,
// 				handlerFunc:    a.HandleLogout,
// 				method:         http.MethodPost,
// 				route:          "/auth/logout",
// 				statusWant:     http.StatusUnauthorized,
// 			},
// 		)
// 	})
//
// 	t.Run("successful logout", func(t *testing.T) {
// 		assertRequest(
// 			t,
// 			requestOptions{
// 				authController: &a,
// 				handlerFunc:    a.HandleLogout,
// 				method:         http.MethodPost,
// 				route:          "/auth/logout",
// 				statusWant:     http.StatusOK,
// 				withCookie:     true,
// 				cookie:         newTestCookie(t, &a),
// 			},
// 		)
// 	})
// }
