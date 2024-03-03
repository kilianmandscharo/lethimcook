package auth

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kilianmandscharo/lethimcook/servutil"
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

func newTestCookie(t *testing.T, authController *AuthController) http.Cookie {
	token, err := authController.authService.createToken()
	assert.NoError(t, err)
	return authController.authService.newTokenCookie(token, time.Now().Add(60*time.Minute))
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
		if options.authController.authService.validCookieToken(&options.cookie) {
			c.Set("authorized", true)
		} else {
			c.Set("authorized", false)
		}
	}

	options.handlerFunc(c)

	assert.Equal(t, options.statusWant, w.Code)

	return w, c
}

func TestRenderAdminPage(t *testing.T) {
	authController := newTestAuthController(controllerOptions{})

	t.Run("valid request", func(t *testing.T) {
		assertRequest(
			t,
			requestOptions{
				authController: &authController,
				handlerFunc:    authController.RenderAdminPage,
				method:         http.MethodGet,
				route:          "/auth/admin",
				statusWant:     http.StatusOK,
			},
		)
	})
}

func TestHandleLogin(t *testing.T) {
	authController := newTestAuthController(controllerOptions{})

	t.Run("no admin", func(t *testing.T) {
		w, c := assertRequest(
			t,
			requestOptions{
				authController: &authController,
				handlerFunc:    authController.HandleLogin,
				method:         http.MethodPost,
				route:          "/auth/login",
				statusWant:     http.StatusNotFound,
				withFormData:   true,
				formData:       "password=" + testPassword,
			},
		)
		assert.Equal(t, 0, len(w.Result().Cookies()))
		assert.False(t, servutil.IsAuthorized(c))
	})

	authController = newTestAuthController(controllerOptions{withAdmin: true})

	t.Run("invalid form key", func(t *testing.T) {
		w, c := assertRequest(
			t,
			requestOptions{
				authController: &authController,
				handlerFunc:    authController.HandleLogin,
				method:         http.MethodPost,
				route:          "/auth/login",
				statusWant:     http.StatusUnauthorized,
				withFormData:   true,
				formData:       "invalidKey=" + testPassword,
			},
		)
		assert.Equal(t, 0, len(w.Result().Cookies()))
		assert.False(t, servutil.IsAuthorized(c))
	})

	t.Run("invalid password", func(t *testing.T) {
		w, c := assertRequest(
			t,
			requestOptions{
				authController: &authController,
				handlerFunc:    authController.HandleLogin,
				method:         http.MethodPost,
				route:          "/auth/login",
				statusWant:     http.StatusUnauthorized,
				withFormData:   true,
				formData:       "password=invalid_password",
			},
		)
		assert.Equal(t, 0, len(w.Result().Cookies()))
		assert.False(t, servutil.IsAuthorized(c))
	})

	t.Run("valid password", func(t *testing.T) {
		w, c := assertRequest(
			t,
			requestOptions{
				authController: &authController,
				handlerFunc:    authController.HandleLogin,
				method:         http.MethodPost,
				route:          "/auth/login",
				statusWant:     http.StatusOK,
				withFormData:   true,
				formData:       "password=" + testPassword,
			},
		)
		assert.Equal(t, 1, len(w.Result().Cookies()))
		assert.True(t, servutil.IsAuthorized(c))
	})
}

func TestHandleLogout(t *testing.T) {
	authController := newTestAuthController(controllerOptions{withAdmin: true})

	t.Run("no cookie", func(t *testing.T) {
		assertRequest(
			t,
			requestOptions{
				authController: &authController,
				handlerFunc:    authController.HandleLogout,
				method:         http.MethodPost,
				route:          "/auth/logout",
				statusWant:     http.StatusUnauthorized,
			},
		)
	})

	t.Run("successful logout", func(t *testing.T) {
		_, c := assertRequest(
			t,
			requestOptions{
				authController: &authController,
				handlerFunc:    authController.HandleLogout,
				method:         http.MethodPost,
				route:          "/auth/logout",
				statusWant:     http.StatusOK,
				withCookie:     true,
				cookie:         newTestCookie(t, &authController),
			},
		)
		assert.False(t, servutil.IsAuthorized(c))
	})
}

func TestHandleUpdatePassword(t *testing.T) {
	authController := newTestAuthController(controllerOptions{})

	t.Run("no admin", func(t *testing.T) {
		assertRequest(
			t,
			requestOptions{
				authController: &authController,
				handlerFunc:    authController.HandleUpdatePassword,
				method:         http.MethodPut,
				route:          "/auth/password",
				statusWant:     http.StatusNotFound,
				withFormData:   true,
				formData:       "oldPassword=test&newPassword=test",
			},
		)
	})

	authController = newTestAuthController(controllerOptions{withAdmin: true})

	t.Run("wrong old password", func(t *testing.T) {
		assertRequest(
			t,
			requestOptions{
				authController: &authController,
				handlerFunc:    authController.HandleUpdatePassword,
				method:         http.MethodPut,
				route:          "/auth/password",
				statusWant:     http.StatusUnauthorized,
				withFormData:   true,
				formData:       "oldPassword=invalid_password&newPassword=test",
			},
		)
	})

	t.Run("wrong old password key", func(t *testing.T) {
		assertRequest(
			t,
			requestOptions{
				authController: &authController,
				handlerFunc:    authController.HandleUpdatePassword,
				method:         http.MethodPut,
				route:          "/auth/password",
				statusWant:     http.StatusUnauthorized,
				withFormData:   true,
				formData:       fmt.Sprintf("wrongKey=%s&newPassword=test", testPassword),
			},
		)
	})

	t.Run("wrong new password key", func(t *testing.T) {
		assertRequest(
			t,
			requestOptions{
				authController: &authController,
				handlerFunc:    authController.HandleUpdatePassword,
				method:         http.MethodPut,
				route:          "/auth/password",
				statusWant:     http.StatusBadRequest,
				withFormData:   true,
				formData:       fmt.Sprintf("oldPassword=%s&wrongKey=test", testPassword),
			},
		)
	})

	t.Run("new password too short", func(t *testing.T) {
		assertRequest(
			t,
			requestOptions{
				authController: &authController,
				handlerFunc:    authController.HandleUpdatePassword,
				method:         http.MethodPut,
				route:          "/auth/password",
				statusWant:     http.StatusBadRequest,
				withFormData:   true,
				formData:       fmt.Sprintf("oldPassword=%s&newPassword=new", testPassword),
			},
		)
	})

	t.Run("succesful password update", func(t *testing.T) {
		assertRequest(
			t,
			requestOptions{
				authController: &authController,
				handlerFunc:    authController.HandleUpdatePassword,
				method:         http.MethodPut,
				route:          "/auth/password",
				statusWant:     http.StatusOK,
				withFormData:   true,
				formData:       fmt.Sprintf("oldPassword=%s&newPassword=updated", testPassword),
			},
		)
	})
}
