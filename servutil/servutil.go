package servutil

import (
	"net/http"

	"github.com/kilianmandscharo/lethimcook/errutil"
	"github.com/kilianmandscharo/lethimcook/templutil"
	"github.com/labstack/echo/v4"
)

func isHxRequest(c echo.Context) bool {
	hxRequestEntry := c.Request().Header["Hx-Request"]
	return len(hxRequestEntry) > 0 && hxRequestEntry[0] == "true"
}

func IsAuthorized(c echo.Context) bool {
	authorizedValue := c.Get("authorized")
	if authorized, ok := authorizedValue.(bool); ok {
		return authorized
	}
	return false
}

func AttachHandlerFunctions(e *echo.Echo) {
	e.GET("/imprint", renderImprint)
	e.GET("/privacy-notice", renderPrivacyNotice)
}

func renderImprint(c echo.Context) error {
	return RenderTemplate(c, templutil.TemplateNameImprint, nil)
}

func renderPrivacyNotice(c echo.Context) error {
	return RenderTemplate(c, templutil.TemplateNamePrivacyNotice, nil)
}

func RenderTemplate(c echo.Context, templateName string, data any) error {
	if isHxRequest(c) {
		return c.Render(http.StatusOK, templutil.FragmentName(templateName), data)
	}

	return c.Render(http.StatusOK, templutil.PageName(templateName), data)
}

func RenderError(c echo.Context, err error) error {
	return c.String(errutil.ErrorHttpCodes[err], err.Error())
}
