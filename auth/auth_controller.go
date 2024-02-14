package auth

import (
	"net/http"
	"time"

	"github.com/kilianmandscharo/lethimcook/errutil"
	"github.com/kilianmandscharo/lethimcook/routes"
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

func (a *AuthController) AttachHandlerFunctions(e *echo.Echo) {
	// Pages
	e.GET("/auth/admin", a.RenderAdminPage)

	// Actions
	e.POST("/auth/login", a.HandleLogin)
	e.POST("/auth/logout", a.HandleLogout)
	e.PUT("/auth/password", a.HandleUpdatePassword)
}

func (r *AuthController) RenderAdminPage(c echo.Context) error {
	return r.renderTemplate(c, routes.TemplateNameAdmin, r.isAdmin(c))
}

func (r *AuthController) HandleLogin(c echo.Context) error {
	if err := c.Request().ParseForm(); err != nil {
		return r.renderError(c, errutil.AuthErrorInvalidForm)
	}

	err := r.authService.validatePassword(c.Request().FormValue("password"))
	if err != nil {
		return r.renderError(c, err)
	}

	token, err := r.authService.createToken()
	if err != nil {
		return r.renderError(c, err)
	}

	cookie := r.authService.newTokenCookie(token, time.Now().Add(60*time.Minute))
	c.SetCookie(&cookie)

	c.Set("authorized", true)

	return r.RenderAdminPage(c)
}

func (r *AuthController) HandleLogout(c echo.Context) error {
	if !r.isAdmin(c) {
		return c.String(http.StatusUnauthorized, "not authorized")
	}

	cookie := r.authService.newTokenCookie("", time.Unix(0, 0))
	c.SetCookie(&cookie)

	c.Set("authorized", false)

	return r.RenderAdminPage(c)
}

func (r *AuthController) HandleUpdatePassword(c echo.Context) error {
	if err := c.Request().ParseForm(); err != nil {
		return r.renderError(c, errutil.AuthErrorInvalidForm)
	}

	err := r.authService.validatePassword(c.Request().FormValue("oldPassword"))
	if err != nil {
		return r.renderError(c, err)
	}

	err = r.authService.updateAdminPasswordHash(c.Request().FormValue("newPassword"))
	if err != nil {
		return r.renderError(c, err)
	}

	return c.String(http.StatusOK, "Passwort erfolgreich geändert")
}

func (r *AuthController) ValidateTokenMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("token")
		if err != nil || !r.authService.validCookieToken(cookie) {
			c.Set("authorized", false)
			return next(c)
		}

		c.Set("authorized", true)

		return next(c)
	}
}

func (r *AuthController) renderError(c echo.Context, err errutil.AuthError) error {
	return c.String(errutil.AuthErrorHttpCodes[err], err.Error())
}

func (r *AuthController) renderTemplate(c echo.Context, templateName string, data any) error {
	if r.isHxRequest(c) {
		return c.Render(http.StatusOK, routes.FragmentName(templateName), data)
	}

	return c.Render(http.StatusOK, routes.PageName(templateName), data)
}

func (r *AuthController) isHxRequest(c echo.Context) bool {
	hxRequestEntry := c.Request().Header["Hx-Request"]
	return len(hxRequestEntry) > 0 && hxRequestEntry[0] == "true"
}

func (r *AuthController) isAdmin(c echo.Context) bool {
	authorized := c.Get("authorized")
	if isAdmin, ok := authorized.(bool); ok {
		return isAdmin
	}
	return false
}
