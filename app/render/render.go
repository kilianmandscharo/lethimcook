package render

import (
	"github.com/a-h/templ"
	"github.com/kilianmandscharo/lethimcook/components"
	"github.com/kilianmandscharo/lethimcook/errutil"
	"github.com/kilianmandscharo/lethimcook/logging"
	"github.com/kilianmandscharo/lethimcook/servutil"
	"github.com/labstack/echo/v4"
)

type Renderer struct {
	logger *logging.Logger
}

func New(logger *logging.Logger) Renderer {
	return Renderer{logger: logger}
}

func (r *Renderer) RenderError(c echo.Context, err error) error {
	r.logger.Error(err)
	userMessage := errutil.GetAppErrorUserMessage(err)
	statusCode := errutil.GetAppErrorStatusCode(err)
	if servutil.IsHxRequest(c) {
		c.Response().Header().Set("HX-Retarget", "#notification-container")
		c.Response().Header().Set("HX-Reswap", "beforeend:#notification-container")
		c.Response().WriteHeader(statusCode)
		return components.Notification(userMessage, true).Render(
			c.Request().Context(),
			c.Response().Writer,
		)
	}
	return c.String(statusCode, userMessage)
}

type RenderComponentOptions struct {
	Context   echo.Context
	Component templ.Component
	Message   string
	Err       error
}

func (r *Renderer) RenderComponent(options RenderComponentOptions) error {
	options.Context.Response().Header().Set(
		echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8,
	)

	if options.Err != nil {
		options.Context.Response().WriteHeader(
			errutil.GetAppErrorStatusCode(options.Err),
		)
		r.logger.Error(options.Err)
	}

	message := getRenderComponentOptionsMessage(options)
	component := getRenderComponentOptionsComponent(options, message)

	if servutil.IsHxRequest(options.Context) {
		return component.Render(
			options.Context.Request().Context(),
			options.Context.Response().Writer,
		)
	}
	return components.Page(component).Render(
		options.Context.Request().Context(),
		options.Context.Response().Writer,
	)
}

func getRenderComponentOptionsComponent(options RenderComponentOptions, message string) templ.Component {
	if len(message) > 0 {
		return components.Joiner(
			options.Component,
			components.NotificationWithSwap(message, options.Err != nil),
		)
	}
	return options.Component
}

func getRenderComponentOptionsMessage(options RenderComponentOptions) string {
	if options.Err != nil {
		return errutil.GetAppErrorUserMessage(options.Err)
	}
	return options.Message
}
