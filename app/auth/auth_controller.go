package auth

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/kilianmandscharo/lethimcook/components"
	"github.com/kilianmandscharo/lethimcook/errutil"
	"github.com/kilianmandscharo/lethimcook/logging"
	"github.com/kilianmandscharo/lethimcook/render"
	"github.com/kilianmandscharo/lethimcook/servutil"
	"github.com/labstack/echo/v4"
)

type AuthController struct {
	authService AuthService
	logger      *logging.Logger
	renderer    *render.Renderer
}

func NewAuthController(authService AuthService, logger *logging.Logger, renderer *render.Renderer) AuthController {
	return AuthController{
		authService: authService,
		logger:      logger,
		renderer:    renderer,
	}
}

func (ac *AuthController) AttachHandlerFunctions(e *echo.Echo) {
	// Pages
	e.GET("/admin", ac.RenderAdminPage)

	// Actions
	e.POST("/auth/login", ac.HandleLogin)
	e.POST("/auth/logout", ac.HandleLogout)
	e.PUT("/auth/password", ac.HandleUpdatePassword)
}

func (ac *AuthController) RenderAdminPage(c echo.Context) error {
	return ac.renderAdminPage(renderAdminPageOptions{
		c:            c,
		isAuthorized: servutil.IsAuthorized(c),
	})
}

func (ac *AuthController) HandleLogin(c echo.Context) error {
	if err := c.Request().ParseForm(); err != nil {
		return ac.renderer.RenderError(c, &errutil.AppError{
			UserMessage: "Fehlerhaftes Formular",
			Err: fmt.Errorf(
				"failed at HandleLogin(), invalid form: %w",
				err,
			),
			StatusCode: http.StatusBadRequest,
		})
	}

	err := ac.authService.validatePassword(c.Request().FormValue("password"))
	if err != nil {
		appError := errutil.AddMessageToAppError(
			err,
			"failed at HandleLogin()",
		)
		return ac.renderAdminPage(renderAdminPageOptions{
			c:              c,
			isAuthorized:   servutil.IsAuthorized(c),
			loginFormError: errutil.FormErrorInvalidPassword,
			err:            appError,
		})
	}

	token, err := ac.authService.createToken()
	if err != nil {
		return ac.renderer.RenderError(
			c,
			errutil.AddMessageToAppError(err, "failed at HandleLogin()"),
		)
	}

	cookie := ac.authService.newTokenCookie(token, time.Now().Add(60*time.Minute))
	c.SetCookie(&cookie)

	c.Set("authorized", true)

	ac.logger.Info("admin login successful")

	return ac.renderAdminPage(renderAdminPageOptions{
		c:            c,
		isAuthorized: servutil.IsAuthorized(c),
		message:      "Angemeldet",
	})
}

func (ac *AuthController) HandleLogout(c echo.Context) error {
	if !servutil.IsAuthorized(c) {
		return ac.renderer.RenderError(
			c,
			errutil.NewAppErrorNotAuthorized("HandleLogout()"),
		)
	}

	cookie := ac.authService.newTokenCookie("", time.Unix(0, 0))
	c.SetCookie(&cookie)

	c.Set("authorized", false)

	ac.logger.Info("admin logout successful")

	return ac.renderAdminPage(renderAdminPageOptions{
		c:            c,
		isAuthorized: servutil.IsAuthorized(c),
		message:      "Abgemeldet",
	})
}

func (ac *AuthController) HandleUpdatePassword(c echo.Context) error {
	if err := c.Request().ParseForm(); err != nil {
		return ac.renderer.RenderError(c, &errutil.AppError{
			UserMessage: "Fehlerhaftes Formular",
			Err: fmt.Errorf(
				"failed at HandleUpdatePassword(), invalid form: %w",
				err,
			),
			StatusCode: http.StatusBadRequest,
		})
	}

	err := ac.authService.validatePassword(c.Request().FormValue("old-password"))
	if err != nil {
		appError := errutil.AddMessageToAppError(
			err,
			"failed at HandleUpdatePassword()",
		)
		return ac.renderAdminPage(renderAdminPageOptions{
			c:                c,
			isAuthorized:     servutil.IsAuthorized(c),
			err:              appError,
			oldPasswordError: errutil.FormErrorInvalidPassword,
		})
	}

	err = ac.authService.updateAdminPasswordHash(c.Request().FormValue("new-password"))
	if err != nil {
		formError := errors.Unwrap(err)
		appError := errutil.AddMessageToAppError(
			err,
			"failed at HandleUpdatePassword()",
		)
		return ac.renderAdminPage(renderAdminPageOptions{
			c:                c,
			isAuthorized:     servutil.IsAuthorized(c),
			err:              appError,
			newPasswordError: formError,
		})
	}

	ac.logger.Info("admin password updated successfully")

	return ac.renderAdminPage(renderAdminPageOptions{
		c:            c,
		isAuthorized: servutil.IsAuthorized(c),
		message:      "Passwort aktualisiert",
	})
}

func (ac *AuthController) ValidateTokenMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("token")
		if err != nil || !ac.authService.validCookieToken(cookie) {
			c.Set("authorized", false)
			return next(c)
		}

		c.Set("authorized", true)

		return next(c)
	}
}

type renderAdminPageOptions struct {
	c                echo.Context
	isAuthorized     bool
	loginFormError   error
	message          string
	err              error
	oldPasswordError error
	newPasswordError error
}

func (ac *AuthController) renderAdminPage(options renderAdminPageOptions) error {
	return ac.renderer.RenderComponent(render.RenderComponentOptions{
		Context: options.c,
		Component: components.AdminPage(
			options.isAuthorized,
			ac.authService.createLoginForm(options.isAuthorized, options.loginFormError),
			ac.authService.createNewPasswordForm(options.oldPasswordError, options.newPasswordError),
		),
		Message: options.message,
		Err:     options.err,
	})
}
