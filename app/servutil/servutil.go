package servutil

import (
	"encoding/json"
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

func RenderError(c echo.Context, appError error) error {
	payload, err := createMessagePayload(
		errutil.GetAppErrorUserMessage(appError),
		true,
	)
	if err != nil {
		log.Printf("failed to create message payload: %v", err)
	} else {
		c.Response().Header().Set(
			"HX-Trigger",
			string(payload),
		)
	}
	log.Println("error in RenderError():", appError)
	return c.String(
		errutil.GetAppErrorStatusCode(appError),
		errutil.GetAppErrorUserMessage(appError),
	)
}

type RenderComponentOptions struct {
	Context   echo.Context
	Component templ.Component
	Message   string
	Err       error
}

type ResponseMessage struct {
	Value   string `json:"value"`
	IsError bool   `json:"isError"`
}

type TriggerPayload struct {
	Message string `json:"message"`
}

func createMessagePayload(message string, isError bool) (string, error) {
	responseMessage, err := json.Marshal(ResponseMessage{
		Value:   message,
		IsError: isError,
	})
	if err != nil {
		return "", err
	}
	payload, err := json.Marshal(TriggerPayload{Message: string(responseMessage)})
	if err != nil {
		return "", err
	}
	return string(payload), nil
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
			return components.Joiner(options.Component, components.Notification(message, options.Err != nil)).Render(
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
