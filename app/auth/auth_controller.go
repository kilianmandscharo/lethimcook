package auth

import (
	"time"

	"github.com/kilianmandscharo/lethimcook/errutil"
	"github.com/kilianmandscharo/lethimcook/servutil"
	"github.com/labstack/echo/v4"
)

type AuthController struct {
	authService AuthService
}

func NewAuthController(authService AuthService) AuthController {
	return AuthController{
		authService: authService,
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
	return ac.authService.createAdminPage(createAdminPageOptions{
		c:            c,
		isAuthorized: servutil.IsAuthorized(c),
	})
}

func (ac *AuthController) HandleLogin(c echo.Context) error {
	if err := c.Request().ParseForm(); err != nil {
		return servutil.RenderError(c, errutil.AuthErrorInvalidForm)
	}

	err := ac.authService.validatePassword(c.Request().FormValue("password"))
	if err != nil {
		return ac.authService.createAdminPage(createAdminPageOptions{
			c:              c,
			isAuthorized:   servutil.IsAuthorized(c),
			loginFormError: err,
			err:            err,
		})
	}

	token, err := ac.authService.createToken()
	if err != nil {
		return servutil.RenderError(c, err)
	}

	cookie := ac.authService.newTokenCookie(token, time.Now().Add(60*time.Minute))
	c.SetCookie(&cookie)

	c.Set("authorized", true)

	return ac.authService.createAdminPage(createAdminPageOptions{
		c:            c,
		isAuthorized: servutil.IsAuthorized(c),
		message:      "Angemeldet",
	})
}

func (ac *AuthController) HandleLogout(c echo.Context) error {
	if !servutil.IsAuthorized(c) {
		return servutil.RenderError(c, errutil.AuthErrorNotAuthorized)
	}

	cookie := ac.authService.newTokenCookie("", time.Unix(0, 0))
	c.SetCookie(&cookie)

	c.Set("authorized", false)

	return ac.authService.createAdminPage(createAdminPageOptions{
		c:            c,
		isAuthorized: servutil.IsAuthorized(c),
		message:      "Abgemeldet",
	})
}

func (ac *AuthController) HandleUpdatePassword(c echo.Context) error {
	if err := c.Request().ParseForm(); err != nil {
		return servutil.RenderError(c, errutil.AuthErrorInvalidForm)
	}

	err := ac.authService.validatePassword(c.Request().FormValue("old-password"))
	if err != nil {
		return ac.authService.createAdminPage(createAdminPageOptions{
			c:                c,
			isAuthorized:     servutil.IsAuthorized(c),
			err:              err,
			oldPasswordError: err,
		})
	}

	err = ac.authService.updateAdminPasswordHash(c.Request().FormValue("new-password"))
	if err != nil {
		return ac.authService.createAdminPage(createAdminPageOptions{
			c:                c,
			isAuthorized:     servutil.IsAuthorized(c),
			err:              err,
			newPasswordError: err,
		})
	}

	return ac.authService.createAdminPage(createAdminPageOptions{
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
