package servutil

import (
	"encoding/json"
	"fmt"

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
		Component: components.Imprint(),
	})
}

func renderPrivacyNotice(c echo.Context) error {
	return RenderComponent(RenderComponentOptions{
		Context:   c,
		Component: components.PrivacyNotice(),
	})
}

func RenderError(c echo.Context, err error) error {
	return c.String(errutil.ErrorHttpCodes[err], err.Error())
}

type RenderComponentOptions struct {
	Context   echo.Context
	Component templ.Component
	Message   string
	Err       error
}

type responseMessage struct {
	Value   string `json:"value"`
	IsError bool   `json:"isError"`
}

type triggerPayload struct {
	Message string `json:"message"`
}

func RenderComponent(options RenderComponentOptions) error {
	options.Context.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	// Usually you would want to return the status code correctly, however,
	// for this purpose to still render the component with htmx and still
	// somehow mark the response as an error, the code for the error status is
	// set in the header (there might be a better way...)
	if options.Err != nil {
		options.Context.Response().Header().Set(
			"Errorcode",
			fmt.Sprintf("%d", errutil.ErrorHttpCodes[options.Err]),
		)
	}

	var message string
	if options.Err != nil {
		message = options.Err.Error()
	} else {
		message = options.Message
	}

	if len(message) > 0 {
		message, err := json.Marshal(responseMessage{
			Value:   message,
			IsError: options.Err != nil,
		})

		if err == nil {
			payload, err := json.Marshal(triggerPayload{Message: string(message)})

			if err == nil {
				options.Context.Response().Header().Set(
					"HX-Trigger",
					string(payload),
				)
			}
		}

	}

	if isHxRequest(options.Context) {
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
