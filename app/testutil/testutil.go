package testutil

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/kilianmandscharo/lethimcook/servutil"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type RequestOptions struct {
	HandlerFunc         func(c echo.Context) error
	Method              string
	Route               string
	StatusWant          int
	Authorized          bool
	WithFormData        bool
	FormData            string
	WithCookie          bool
	Cookie              http.Cookie
	WithPathParam       bool
	PathParamName       string
	PathParamValue      string
	PathParamNames      []string
	PathParamValues     []string
	WithQueryParam      bool
	QueryParam          string
	HeaderErrorCodeWant int
	AssertMessage       bool
	MessageWant         string
}

func AssertRequest(t *testing.T, options RequestOptions) (*httptest.ResponseRecorder, echo.Context) {
	e := echo.New()
	rr := httptest.NewRecorder()

	var body io.Reader

	if options.WithFormData {
		body = bytes.NewBufferString(options.FormData)
	} else {
		body = nil
	}

	if options.WithQueryParam {
		options.Route += options.QueryParam
	}

	req, err := http.NewRequest(options.Method, options.Route, body)
	assert.NoError(t, err)

	if options.WithFormData {
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	}

	c := e.NewContext(req, rr)

	if options.Authorized {
		c.Set("authorized", true)
	} else {
		c.Set("authorized", false)
	}

	if options.WithPathParam {
		c.SetParamNames(options.PathParamName)
		c.SetParamValues(options.PathParamValue)
		if len(options.PathParamNames) > 0 && len(options.PathParamValues) > 0 {
			c.SetParamNames(options.PathParamNames...)
			c.SetParamValues(options.PathParamValues...)
		}
	}

	assert.NoError(t, options.HandlerFunc(c))

	assert.Equal(t, options.StatusWant, rr.Code)

	if options.HeaderErrorCodeWant != 0 {
		errorCodes := rr.Header()["Errorcode"]
		assert.Equal(t, 1, len(errorCodes))
		assert.Equal(t, strconv.Itoa(options.HeaderErrorCodeWant), errorCodes[0])
	}

	if options.AssertMessage {
		headerValue := rr.Header().Get("HX-Trigger")
		var payload servutil.TriggerPayload
		err := json.Unmarshal([]byte(headerValue), &payload)
		var message servutil.ResponseMessage
		err = json.Unmarshal([]byte(payload.Message), &message)
		assert.NoError(t, err)
		assert.Equal(t, options.MessageWant, message.Value)
	}

	return rr, c
}

type TestFormDataStringOptions struct {
	TitleEmpty        bool
	DescriptionEmpty  bool
	DurationEmpty     bool
	InvalidDuration   bool
	TagsEmpty         bool
	IngredientsEmpty  bool
	InstructionsEmpty bool
}

func ConstructTestFormDataString(options TestFormDataStringOptions) string {
	formData := ""

	if !options.TitleEmpty {
		formData += "title=title"
	}
	if !options.DescriptionEmpty {
		if len(formData) != 0 {
			formData += "&"
		}
		formData += "description=description"
	}
	if !options.DurationEmpty {
		if len(formData) != 0 {
			formData += "&"
		}
		if options.InvalidDuration {
			formData += "duration=duration"
		} else {
			formData += "duration=10"
		}
	}
	if !options.TagsEmpty {
		if len(formData) != 0 {
			formData += "&"
		}
		formData += "tags=tags"
	}
	if !options.IngredientsEmpty {
		if len(formData) != 0 {
			formData += "&"
		}
		formData += "ingredients=ingredients"
	}
	if !options.InstructionsEmpty {
		if len(formData) != 0 {
			formData += "&"
		}
		formData += "instructions=instructions"
	}

	return formData
}
