package servutil

import (
	"log"

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
	return RenderComponent(RenderComponentOptions{
		Context:   c,
		Component: components.Imprint(IsAuthorized(c)),
	})
}

func renderPrivacyNotice(c echo.Context) error {
	return RenderComponent(RenderComponentOptions{
		Context:   c,
		Component: components.PrivacyNotice(IsAuthorized(c)),
	})
}

func RenderError(c echo.Context, err error) error {
	c.Response().Header().Set("HX-Retarget", "#notification-container")
	c.Response().Header().Set("HX-Reswap", "beforeend:#notification-container")
	c.Response().WriteHeader(
		errutil.GetAppErrorStatusCode(err),
	)
	log.Println("error in RenderError():", err)
	return components.Notification(errutil.GetAppErrorUserMessage(err), true).Render(
		c.Request().Context(),
		c.Response().Writer,
	)
}

type RenderComponentOptions struct {
	Context   echo.Context
	Component templ.Component
	Message   string
	Err       error
}

func RenderComponent(options RenderComponentOptions) error {
	options.Context.Response().Header().Set(
		echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8,
	)

	if options.Err != nil {
		options.Context.Response().WriteHeader(
			errutil.GetAppErrorStatusCode(options.Err),
		)
		log.Println("error in RenderComponent():", options.Err)
	}

	var message string
	if options.Err != nil {
		message = errutil.GetAppErrorUserMessage(options.Err)
	} else {
		message = options.Message
	}

	if isHxRequest(options.Context) {
		if len(message) > 0 {
			return components.Joiner(options.Component, components.NotificationWithSwap(message, options.Err != nil)).Render(
				options.Context.Request().Context(),
				options.Context.Response().Writer,
			)
		}
		return options.Component.Render(
			options.Context.Request().Context(),
			options.Context.Response().Writer,
		)
	}
	return components.Page(options.Component).Render(
		options.Context.Request().Context(),
		options.Context.Response().Writer,
	)
}
