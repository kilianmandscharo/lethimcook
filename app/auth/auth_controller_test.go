package auth

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/kilianmandscharo/lethimcook/servutil"
	"github.com/kilianmandscharo/lethimcook/testutil"
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

func TestRenderAdminPage(t *testing.T) {
	authController := newTestAuthController(controllerOptions{})

	t.Run("valid request", func(t *testing.T) {
		testutil.AssertRequest(
			t,
			testutil.RequestOptions{
				HandlerFunc: authController.RenderAdminPage,
				Method:      http.MethodGet,
				Route:       "/auth/admin",
				StatusWant:  http.StatusOK,
			},
		)
	})
}

func TestHandleLogin(t *testing.T) {
	authController := newTestAuthController(controllerOptions{})

	t.Run("no admin", func(t *testing.T) {
		w, c := testutil.AssertRequest(
			t,
			testutil.RequestOptions{
				HandlerFunc:  authController.HandleLogin,
				Method:       http.MethodPost,
				Route:        "/auth/login",
				StatusWant:   http.StatusNotFound,
				WithFormData: true,
				FormData:     "password=" + testPassword,
			},
		)
		assert.Equal(t, 0, len(w.Result().Cookies()))
		assert.False(t, servutil.IsAuthorized(c))
	})

	authController = newTestAuthController(controllerOptions{withAdmin: true})

	t.Run("invalid form key", func(t *testing.T) {
		w, c := testutil.AssertRequest(
			t,
			testutil.RequestOptions{
				HandlerFunc:  authController.HandleLogin,
				Method:       http.MethodPost,
				Route:        "/auth/login",
				StatusWant:   http.StatusUnauthorized,
				WithFormData: true,
				FormData:     "invalidKey=" + testPassword,
			},
		)
		assert.Equal(t, 0, len(w.Result().Cookies()))
		assert.False(t, servutil.IsAuthorized(c))
	})

	t.Run("invalid password", func(t *testing.T) {
		w, c := testutil.AssertRequest(
			t,
			testutil.RequestOptions{
				HandlerFunc:  authController.HandleLogin,
				Method:       http.MethodPost,
				Route:        "/auth/login",
				StatusWant:   http.StatusUnauthorized,
				WithFormData: true,
				FormData:     "password=invalid_password",
			},
		)
		assert.Equal(t, 0, len(w.Result().Cookies()))
		assert.False(t, servutil.IsAuthorized(c))
	})

	t.Run("valid password", func(t *testing.T) {
		w, c := testutil.AssertRequest(
			t,
			testutil.RequestOptions{
				HandlerFunc:  authController.HandleLogin,
				Method:       http.MethodPost,
				Route:        "/auth/login",
				StatusWant:   http.StatusOK,
				WithFormData: true,
				FormData:     "password=" + testPassword,
			},
		)
		assert.Equal(t, 1, len(w.Result().Cookies()))
		assert.True(t, servutil.IsAuthorized(c))
	})
}

func TestHandleLogout(t *testing.T) {
	authController := newTestAuthController(controllerOptions{withAdmin: true})

	t.Run("no cookie", func(t *testing.T) {
		testutil.AssertRequest(
			t,
			testutil.RequestOptions{
				HandlerFunc: authController.HandleLogout,
				Method:      http.MethodPost,
				Route:       "/auth/logout",
				StatusWant:  http.StatusUnauthorized,
			},
		)
	})

	t.Run("successful logout", func(t *testing.T) {
		_, c := testutil.AssertRequest(
			t,
			testutil.RequestOptions{
				HandlerFunc: authController.HandleLogout,
				Method:      http.MethodPost,
				Route:       "/auth/logout",
				StatusWant:  http.StatusOK,
				WithCookie:  true,
				Cookie:      newTestCookie(t, &authController),
				Authorized:  true,
			},
		)
		assert.False(t, servutil.IsAuthorized(c))
	})
}

func TestHandleUpdatePassword(t *testing.T) {
	authController := newTestAuthController(controllerOptions{})

	t.Run("no admin", func(t *testing.T) {
		testutil.AssertRequest(
			t,
			testutil.RequestOptions{
				HandlerFunc:  authController.HandleUpdatePassword,
				Method:       http.MethodPut,
				Route:        "/auth/password",
				StatusWant:   http.StatusNotFound,
				WithFormData: true,
				FormData:     "oldPassword=test&newPassword=test",
			},
		)
	})

	authController = newTestAuthController(controllerOptions{withAdmin: true})

	t.Run("wrong old password", func(t *testing.T) {
		testutil.AssertRequest(
			t,
			testutil.RequestOptions{
				HandlerFunc:  authController.HandleUpdatePassword,
				Method:       http.MethodPut,
				Route:        "/auth/password",
				StatusWant:   http.StatusUnauthorized,
				WithFormData: true,
				FormData:     "oldPassword=invalid_password&newPassword=test",
			},
		)
	})

	t.Run("wrong old password key", func(t *testing.T) {
		testutil.AssertRequest(
			t,
			testutil.RequestOptions{
				HandlerFunc:  authController.HandleUpdatePassword,
				Method:       http.MethodPut,
				Route:        "/auth/password",
				StatusWant:   http.StatusUnauthorized,
				WithFormData: true,
				FormData:     fmt.Sprintf("wrongKey=%s&newPassword=test", testPassword),
			},
		)
	})

	t.Run("wrong new password key", func(t *testing.T) {
		testutil.AssertRequest(
			t,
			testutil.RequestOptions{
				HandlerFunc:  authController.HandleUpdatePassword,
				Method:       http.MethodPut,
				Route:        "/auth/password",
				StatusWant:   http.StatusBadRequest,
				WithFormData: true,
				FormData:     fmt.Sprintf("oldPassword=%s&wrongKey=test", testPassword),
			},
		)
	})

	t.Run("new password too short", func(t *testing.T) {
		testutil.AssertRequest(
			t,
			testutil.RequestOptions{
				HandlerFunc:  authController.HandleUpdatePassword,
				Method:       http.MethodPut,
				Route:        "/auth/password",
				StatusWant:   http.StatusBadRequest,
				WithFormData: true,
				FormData:     fmt.Sprintf("oldPassword=%s&newPassword=new", testPassword),
			},
		)
	})

	t.Run("succesful password update", func(t *testing.T) {
		testutil.AssertRequest(
			t,
			testutil.RequestOptions{
				HandlerFunc:  authController.HandleUpdatePassword,
				Method:       http.MethodPut,
				Route:        "/auth/password",
				StatusWant:   http.StatusOK,
				WithFormData: true,
				FormData:     fmt.Sprintf("oldPassword=%s&newPassword=updated", testPassword),
			},
		)
	})
}
