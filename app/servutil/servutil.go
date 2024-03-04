package servutil

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/kilianmandscharo/lethimcook/components"
	"github.com/kilianmandscharo/lethimcook/errutil"
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
	return RenderComponent(c, components.Imprint())
}

func renderPrivacyNotice(c echo.Context) error {
	return RenderComponent(c, components.PrivacyNotice())
}

// func RenderTemplate(c echo.Context, templateName string, data any) error {
// 	if isHxRequest(c) {
// 		return c.Render(http.StatusOK, templutil.FragmentName(templateName), data)
// 	}
//
// 	return c.Render(http.StatusOK, templutil.PageName(templateName), data)
// }
//
// func RenderTemplateComponent(c echo.Context, templateName string, data any) error {
// 	return c.Render(http.StatusOK, templutil.FragmentName(templateName), data)
// }

func RenderError(c echo.Context, err error) error {
	return c.String(errutil.ErrorHttpCodes[err], err.Error())
}

func RenderComponent(c echo.Context, component templ.Component) error {
	c.Response().Writer.WriteHeader(http.StatusOK)
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	if isHxRequest(c) {
		return component.Render(c.Request().Context(), c.Response().Writer)
	}

	return components.Page(component).Render(c.Request().Context(), c.Response().Writer)
}
